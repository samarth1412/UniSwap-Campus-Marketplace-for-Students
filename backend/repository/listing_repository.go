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
	GetAll(ctx context.Context, params models.ListingListParams) (*models.PaginatedListings, error)
	GetByID(ctx context.Context, id int64) (*models.Listing, error)
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

func (r *PostgresListingRepository) GetAll(ctx context.Context, params models.ListingListParams) (*models.PaginatedListings, error) {
	whereClause := ""
	countArgs := make([]interface{}, 0)
	selectArgs := make([]interface{}, 0)

	if params.Search != "" {
		whereClause = " WHERE title ILIKE $1"
		searchValue := "%" + strings.TrimSpace(params.Search) + "%"
		countArgs = append(countArgs, searchValue)
		selectArgs = append(selectArgs, searchValue)
	}

	countQuery := `SELECT COUNT(*) FROM listings` + whereClause
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return nil, fmt.Errorf("count listings: %w", err)
	}

	offset := (params.Page - 1) * params.Limit
	limitPlaceholder := len(selectArgs) + 1
	offsetPlaceholder := len(selectArgs) + 2
	selectArgs = append(selectArgs, params.Limit, offset)

	selectQuery := `
		SELECT id, seller_id, title, description, price, category, created_at
		FROM listings` + whereClause + fmt.Sprintf(`
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, limitPlaceholder, offsetPlaceholder)

	rows, err := r.db.QueryContext(ctx, selectQuery, selectArgs...)
	if err != nil {
		return nil, fmt.Errorf("get listings: %w", err)
	}
	defer rows.Close()

	items := make([]models.Listing, 0)
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
		items = append(items, listing)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate listings: %w", err)
	}

	totalPages := 0
	if total > 0 {
		totalPages = (total + params.Limit - 1) / params.Limit
	}

	return &models.PaginatedListings{
		Items:      items,
		Page:       params.Page,
		Limit:      params.Limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
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
