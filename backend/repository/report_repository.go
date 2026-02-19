package repository

import (
	"context"
	"database/sql"
	"fmt"

	"uniswap-campus-marketplace/models"
)

type ReportRepository interface {
	Create(ctx context.Context, report *models.Report) (*models.Report, error)
}

type PostgresReportRepository struct {
	db *sql.DB
}

func NewPostgresReportRepository(db *sql.DB) *PostgresReportRepository {
	return &PostgresReportRepository{db: db}
}

func (r *PostgresReportRepository) Create(ctx context.Context, report *models.Report) (*models.Report, error) {
	const query = `
		INSERT INTO reports (listing_id, reporter_id, reason)
		VALUES ($1, $2, $3)
		RETURNING id, listing_id, reporter_id, reason, created_at
	`

	created := &models.Report{}
	err := r.db.QueryRowContext(
		ctx,
		query,
		report.ListingID,
		report.ReporterUserID,
		report.Reason,
	).Scan(
		&created.ID,
		&created.ListingID,
		&created.ReporterUserID,
		&created.Reason,
		&created.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("create report: %w", err)
	}

	return created, nil
}
