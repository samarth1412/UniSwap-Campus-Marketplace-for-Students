package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"uniswap-campus-marketplace/models"
)

var ErrWishlistAlreadyExists = errors.New("wishlist item already exists")
var ErrWishlistNotFound = errors.New("wishlist item not found")

type WishlistRepository interface {
	Create(ctx context.Context, item *models.WishlistItem) (*models.WishlistItem, error)
	Delete(ctx context.Context, wishlistID, userID int64) error
	ListByUser(ctx context.Context, userID int64) ([]models.WishlistListing, error)
}

type PostgresWishlistRepository struct {
	db *sql.DB
}

func NewPostgresWishlistRepository(db *sql.DB) *PostgresWishlistRepository {
	return &PostgresWishlistRepository{db: db}
}

func (r *PostgresWishlistRepository) Create(ctx context.Context, item *models.WishlistItem) (*models.WishlistItem, error) {
	const query = `
		INSERT INTO wishlist (user_id, listing_id)
		VALUES ($1, $2)
		RETURNING id, user_id, listing_id, created_at
	`

	created := &models.WishlistItem{}
	err := r.db.QueryRowContext(ctx, query, item.UserID, item.ListingID).Scan(
		&created.ID,
		&created.UserID,
		&created.ListingID,
		&created.CreatedAt,
	)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate key") {
			return nil, ErrWishlistAlreadyExists
		}
		return nil, fmt.Errorf("create wishlist item: %w", err)
	}

	return created, nil
}

func (r *PostgresWishlistRepository) Delete(ctx context.Context, wishlistID, userID int64) error {
	const query = `
		DELETE FROM wishlist
		WHERE id = $1 AND user_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, wishlistID, userID)
	if err != nil {
		return fmt.Errorf("delete wishlist item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("wishlist rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return ErrWishlistNotFound
	}

	return nil
}

func (r *PostgresWishlistRepository) ListByUser(ctx context.Context, userID int64) ([]models.WishlistListing, error) {
	const query = `
		SELECT
			w.id,
			w.created_at,
			l.id,
			l.seller_id,
			l.title,
			l.description,
			l.price,
			l.category,
			COALESCE((
				SELECT li.image_url
				FROM listing_images li
				WHERE li.listing_id = l.id
				ORDER BY li.is_primary DESC, li.created_at ASC
				LIMIT 1
			), ''),
			l.created_at
		FROM wishlist w
		JOIN listings l ON l.id = w.listing_id
		WHERE w.user_id = $1
		ORDER BY w.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list wishlist items: %w", err)
	}
	defer rows.Close()

	items := make([]models.WishlistListing, 0)
	for rows.Next() {
		var item models.WishlistListing
		if err := rows.Scan(
			&item.WishlistID,
			&item.SavedAt,
			&item.Listing.ID,
			&item.Listing.UserID,
			&item.Listing.Title,
			&item.Listing.Description,
			&item.Listing.Price,
			&item.Listing.Category,
			&item.Listing.PrimaryImageURL,
			&item.Listing.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan wishlist item: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate wishlist items: %w", err)
	}

	return items, nil
}
