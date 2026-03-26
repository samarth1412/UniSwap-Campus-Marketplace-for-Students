package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"uniswap-campus-marketplace/models"
)

type ListingImageRepository interface {
	Create(ctx context.Context, image *models.ListingImage) (*models.ListingImage, error)
	CreateMany(ctx context.Context, listingID int64, imageURLs []string) ([]models.ListingImage, error)
}

type PostgresListingImageRepository struct {
	db *sql.DB
}

func NewPostgresListingImageRepository(db *sql.DB) *PostgresListingImageRepository {
	return &PostgresListingImageRepository{db: db}
}

func (r *PostgresListingImageRepository) Create(ctx context.Context, image *models.ListingImage) (*models.ListingImage, error) {
	const query = `
		INSERT INTO listing_images (listing_id, image_url, is_primary)
		VALUES ($1, $2, $3)
		RETURNING id, listing_id, image_url, is_primary, created_at
	`

	created := &models.ListingImage{}
	err := r.db.QueryRowContext(ctx, query, image.ListingID, image.ImageURL, image.IsPrimary).Scan(
		&created.ID,
		&created.ListingID,
		&created.ImageURL,
		&created.IsPrimary,
		&created.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create listing image: %w", err)
	}

	return created, nil
}

func (r *PostgresListingImageRepository) CreateMany(ctx context.Context, listingID int64, imageURLs []string) ([]models.ListingImage, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin listing images transaction: %w", err)
	}
	defer tx.Rollback()

	const query = `
		INSERT INTO listing_images (listing_id, image_url, is_primary)
		VALUES ($1, $2, $3)
		RETURNING id, listing_id, image_url, is_primary, created_at
	`

	created := make([]models.ListingImage, 0, len(imageURLs))
	for _, imageURL := range imageURLs {
		var image models.ListingImage
		err := tx.QueryRowContext(ctx, query, listingID, strings.TrimSpace(imageURL), false).Scan(
			&image.ID,
			&image.ListingID,
			&image.ImageURL,
			&image.IsPrimary,
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
