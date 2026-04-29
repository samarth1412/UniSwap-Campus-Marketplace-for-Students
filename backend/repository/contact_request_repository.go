package repository

import (
	"context"
	"database/sql"
	"fmt"

	"uniswap-campus-marketplace/models"
)

type ContactRequestRepository interface {
	Create(ctx context.Context, request *models.ContactRequest) (*models.ContactRequest, error)
	ListBySellerID(ctx context.Context, sellerID int64) ([]models.ContactRequest, error)
}

type PostgresContactRequestRepository struct {
	db *sql.DB
}

func NewPostgresContactRequestRepository(db *sql.DB) *PostgresContactRequestRepository {
	return &PostgresContactRequestRepository{db: db}
}

func (r *PostgresContactRequestRepository) Create(ctx context.Context, request *models.ContactRequest) (*models.ContactRequest, error) {
	const query = `
		INSERT INTO contact_requests (listing_id, buyer_id, seller_id, message, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, listing_id, buyer_id, seller_id, message, status, created_at
	`

	created := &models.ContactRequest{}
	err := r.db.QueryRowContext(
		ctx,
		query,
		request.ListingID,
		request.BuyerID,
		request.SellerID,
		request.Message,
		request.Status,
	).Scan(
		&created.ID,
		&created.ListingID,
		&created.BuyerID,
		&created.SellerID,
		&created.Message,
		&created.Status,
		&created.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create contact request: %w", err)
	}

	return created, nil
}

func (r *PostgresContactRequestRepository) ListBySellerID(ctx context.Context, sellerID int64) ([]models.ContactRequest, error) {
	const query = `
		SELECT id, listing_id, buyer_id, seller_id, message, status, created_at
		FROM contact_requests
		WHERE seller_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, sellerID)
	if err != nil {
		return nil, fmt.Errorf("list contact requests by seller id: %w", err)
	}
	defer rows.Close()

	requests := make([]models.ContactRequest, 0)
	for rows.Next() {
		var request models.ContactRequest
		if err := rows.Scan(
			&request.ID,
			&request.ListingID,
			&request.BuyerID,
			&request.SellerID,
			&request.Message,
			&request.Status,
			&request.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan contact request: %w", err)
		}
		requests = append(requests, request)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate contact requests: %w", err)
	}

	return requests, nil
}
