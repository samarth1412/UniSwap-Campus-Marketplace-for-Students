package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"uniswap-campus-marketplace/models"
)

var ErrListingNotFound = errors.New("listing not found")

type ListingRepository interface {
	Create(ctx context.Context, listing *models.Listing) (*models.Listing, error)
	GetAll(ctx context.Context, search string) ([]models.Listing, error)
	GetByID(ctx context.Context, id int64) (*models.Listing, error)
}

type PostgresListingRepository struct {
	db *sql.DB
}

func NewPostgresListingRepository(db *sql.DB) *PostgresListingRepository {
	return &PostgresListingRepository{db: db}
}

func (r *PostgresListingRepository) Create(ctx context.Context, listing *models.Listing) (*models.Listing, error) {
	const query = `
		INSERT INTO listings (seller_id, title, description, category, price)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, seller_id, title, description, price, category, created_at
	`

	created := &models.Listing{}
	err := r.db.QueryRowContext(
		ctx,
		query,
		listing.UserID,
		listing.Title,
		listing.Description,
		listing.Category,
		listing.Price,
	).Scan(
		&created.ID,
		&created.UserID,
		&created.Title,
		&created.Description,
		&created.Price,
		&created.Category,
		&created.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create listing: %w", err)
	}

	return created, nil
}

func (r *PostgresListingRepository) GetAll(ctx context.Context, search string) ([]models.Listing, error) {
	base := `
		SELECT id, seller_id, title, description, price, category, created_at
		FROM listings
	`

	var (
		args  []interface{}
		query string
	)

	if strings.TrimSpace(search) != "" {
		query = base + ` WHERE title ILIKE $1 ORDER BY created_at DESC`
		args = append(args, "%"+strings.TrimSpace(search)+"%")
	} else {
		query = base + ` ORDER BY created_at DESC`
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get listings: %w", err)
	}
	defer rows.Close()

	listings := make([]models.Listing, 0)
	for rows.Next() {
		var listing models.Listing
		if err := rows.Scan(
			&listing.ID,
			&listing.UserID,
			&listing.Title,
			&listing.Description,
			&listing.Price,
			&listing.Category,
			&listing.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan listing: %w", err)
		}
		listings = append(listings, listing)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate listings: %w", err)
	}

	return listings, nil
}

func (r *PostgresListingRepository) GetByID(ctx context.Context, id int64) (*models.Listing, error) {
	const query = `
		SELECT id, seller_id, title, description, price, category, created_at
		FROM listings
		WHERE id = $1
	`

	listing := &models.Listing{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&listing.ID,
		&listing.UserID,
		&listing.Title,
		&listing.Description,
		&listing.Price,
		&listing.Category,
		&listing.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrListingNotFound
		}
		return nil, fmt.Errorf("get listing by id: %w", err)
	}

	return listing, nil
}
