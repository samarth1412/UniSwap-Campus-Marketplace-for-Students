package repository

import (
	"context"
	"database/sql"
	"fmt"

	"uniswap-campus-marketplace/models"
)

type ListingImageRepository interface {
	CreateMany(ctx context.Context, listingID int64, imageURLs []string) ([]models.ListingImage, error)
}

type PostgresListingImageRepository struct {
	db *sql.DB
}

func NewPostgresListingImageRepository(db *sql.DB) *PostgresListingImageRepository {
	return &PostgresListingImageRepository{db: db}
}

func (r *PostgresListingImageRepository) CreateMany(ctx context.Context, listingID int64, imageURLs []string) ([]models.ListingImage, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin listing images transaction: %w", err)
	}
	defer tx.Rollback()

	const query = `
		INSERT INTO listing_images (listing_id, image_url)
		VALUES ($1, $2)
		RETURNING id, listing_id, image_url, created_at
	`

	created := make([]models.ListingImage, 0, len(imageURLs))
	for _, imageURL := range imageURLs {
		var image models.ListingImage
		err := tx.QueryRowContext(ctx, query, listingID, imageURL).Scan(
			&image.ID,
			&image.ListingID,
			&image.ImageURL,
			&image.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("insert listing image: %w", err)
		}
		created = append(created, image)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit listing images transaction: %w", err)
	}

	return created, nil
}
