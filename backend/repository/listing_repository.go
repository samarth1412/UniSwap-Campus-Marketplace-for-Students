package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"

	"uniswap-campus-marketplace/models"
)

var ErrListingNotFound = errors.New("listing not found")

type ListingRepository interface {
	Create(ctx context.Context, listing *models.Listing) (*models.Listing, error)
	GetAll(ctx context.Context, query models.ListingQuery) (*models.PaginatedListings, error)
	GetByID(ctx context.Context, id int64) (*models.Listing, error)
	GetByUserID(ctx context.Context, userID int64) ([]models.Listing, error)
	UpdateByID(ctx context.Context, id int64, listing *models.Listing) (*models.Listing, error)
	DeleteByID(ctx context.Context, id int64) (int64, error)
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

func (r *PostgresListingRepository) GetAll(ctx context.Context, query models.ListingQuery) (*models.PaginatedListings, error) {
	where, args := listingFilters(query)

	countQuery := `SELECT COUNT(*) FROM listings l` + where
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, fmt.Errorf("get listings count: %w", err)
	}

	page := query.Page
	limit := query.Limit
	offset := (page - 1) * limit
	listArgs := append(args, limit, offset)
	limitArg := len(args) + 1
	offsetArg := len(args) + 2

	listQuery := fmt.Sprintf(`
		SELECT
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
		FROM listings l
		%s
		ORDER BY l.created_at DESC
		LIMIT $%d OFFSET $%d
	`, where, limitArg, offsetArg)

	rows, err := r.db.QueryContext(ctx, listQuery, listArgs...)
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
			&listing.PrimaryImageURL,
			&listing.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan listing: %w", err)
		}
		listings = append(listings, listing)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate listings: %w", err)
	}

	totalPages := 1
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	return &models.PaginatedListings{
		Items:      listings,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func listingFilters(query models.ListingQuery) (string, []interface{}) {
	conditions := make([]string, 0, 4)
	args := make([]interface{}, 0, 4)

	if keyword := strings.TrimSpace(query.Keyword); keyword != "" {
		args = append(args, "%"+keyword+"%")
		conditions = append(conditions, fmt.Sprintf("(l.title ILIKE $%d OR l.description ILIKE $%d)", len(args), len(args)))
	}
	if category := strings.TrimSpace(query.Category); category != "" {
		args = append(args, category)
		conditions = append(conditions, fmt.Sprintf("LOWER(l.category) = LOWER($%d)", len(args)))
	}
	if query.MinPrice != nil {
		args = append(args, *query.MinPrice)
		conditions = append(conditions, fmt.Sprintf("l.price >= $%d", len(args)))
	}
	if query.MaxPrice != nil {
		args = append(args, *query.MaxPrice)
		conditions = append(conditions, fmt.Sprintf("l.price <= $%d", len(args)))
	}

	if len(conditions) == 0 {
		return "", args
	}

	return " WHERE " + strings.Join(conditions, " AND "), args
}

func (r *PostgresListingRepository) GetByID(ctx context.Context, id int64) (*models.Listing, error) {
	const query = `
		SELECT
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
		FROM listings l
		WHERE l.id = $1
	`

	listing := &models.Listing{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&listing.ID,
		&listing.UserID,
		&listing.Title,
		&listing.Description,
		&listing.Price,
		&listing.Category,
		&listing.PrimaryImageURL,
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

func (r *PostgresListingRepository) GetByUserID(ctx context.Context, userID int64) ([]models.Listing, error) {
	const query = `
		SELECT
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
		FROM listings l
		WHERE l.seller_id = $1
		ORDER BY l.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get listings by user id: %w", err)
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
			&listing.PrimaryImageURL,
			&listing.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan listing by user id: %w", err)
		}
		listings = append(listings, listing)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate listings by user id: %w", err)
	}

	return listings, nil
}

func (r *PostgresListingRepository) UpdateByID(ctx context.Context, id int64, listing *models.Listing) (*models.Listing, error) {
	const query = `
		UPDATE listings
		SET title = $1,
			description = $2,
			category = $3,
			price = $4,
			updated_at = NOW()
		WHERE id = $5
		RETURNING id, seller_id, title, description, price, category, created_at
	`

	updated := &models.Listing{}
	err := r.db.QueryRowContext(
		ctx,
		query,
		listing.Title,
		listing.Description,
		listing.Category,
		listing.Price,
		id,
	).Scan(
		&updated.ID,
		&updated.UserID,
		&updated.Title,
		&updated.Description,
		&updated.Price,
		&updated.Category,
		&updated.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrListingNotFound
		}
		return nil, fmt.Errorf("update listing by id: %w", err)
	}

	return updated, nil
}

func (r *PostgresListingRepository) DeleteByID(ctx context.Context, id int64) (int64, error) {
	const query = `
		DELETE FROM listings
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return 0, fmt.Errorf("delete listing by id: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("delete listing by id rows affected: %w", err)
	}

	return rowsAffected, nil
}
