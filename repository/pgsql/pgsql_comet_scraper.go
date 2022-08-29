package pgsql

import (
	"cometScraper/entity"
	"context"
	"database/sql"
	"fmt"
)

// CometScraperRepository represent the cometScraper's repository contract
type CometScraperRepository interface {
	Create(ctx context.Context, cometScraper *entity.CometScraper) error
	GetByID(ctx context.Context, id string) (entity.CometScraper, error)
	Fetch(ctx context.Context) ([]entity.CometScraper, error)
	Update(ctx context.Context, c *entity.CometScraper) error
	UpdateStatus(ctx context.Context, comet *entity.CometScraper) (err error)
	Delete(ctx context.Context, id string) error
}

type pgsqlCometScraperRepository struct {
	db *sql.DB
}

// NewPgsqlCometScraperRepository NewCometScraperRepository will create new an cometScraperRepository object representation of CometScraperRepository interface
func NewPgsqlCometScraperRepository(db *sql.DB) CometScraperRepository {
	return &pgsqlCometScraperRepository{
		db: db,
	}
}

func (r *pgsqlCometScraperRepository) UpdateStatus(ctx context.Context, comet *entity.CometScraper) (err error) {
	query := "UPDATE comet_scraper SET status = $1, updated_at = $2 WHERE uuid = $3"
	res, err := r.db.ExecContext(ctx, query, comet.Status, comet.UpdatedAt, comet.Uuid)
	if err != nil {
		return
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return
	}

	if affect != 1 {
		err = fmt.Errorf("weird behavior, total affected: %d", affect)
	}

	return
}

func (r *pgsqlCometScraperRepository) Update(ctx context.Context, comet *entity.CometScraper) (err error) {
	query := `UPDATE comet_scraper SET status = $1, applicant = $2,time_taken = $3, updated_at = $4 WHERE uuid = $5`
	res, err := r.db.ExecContext(ctx, query, comet.Status, comet.Applicant, comet.TimeTaken, comet.UpdatedAt, comet.Uuid)
	if err != nil {
		return
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return
	}

	if affect != 1 {
		err = fmt.Errorf("weird behavior, total affected: %d", affect)
	}

	return
}

func (r *pgsqlCometScraperRepository) Create(ctx context.Context, cometScraper *entity.CometScraper) (err error) {
	query := `INSERT INTO comet_scraper (uuid, time_taken, applicant, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = r.db.ExecContext(ctx, query, cometScraper.Uuid, cometScraper.TimeTaken, cometScraper.Applicant, cometScraper.Status, cometScraper.CreatedAt, cometScraper.UpdatedAt)
	return
}

func (r *pgsqlCometScraperRepository) GetByID(ctx context.Context, id string) (cometScraper entity.CometScraper, err error) {
	query := "SELECT uuid, applicant, time_taken, status, created_at, updated_at FROM comet_scraper WHERE uuid = $1"
	err = r.db.QueryRowContext(ctx, query, id).Scan(&cometScraper.Uuid, &cometScraper.Applicant, &cometScraper.TimeTaken, &cometScraper.Status, &cometScraper.CreatedAt, &cometScraper.UpdatedAt)

	return
}

func (r *pgsqlCometScraperRepository) Fetch(ctx context.Context) (cometScrapers []entity.CometScraper, err error) {
	query := "SELECT uuid, time_taken, applicant, status, created_at, updated_at FROM comet_scraper"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return cometScrapers, err
	}

	defer rows.Close()

	for rows.Next() {
		var cometScraper entity.CometScraper
		err := rows.Scan(&cometScraper.Uuid, &cometScraper.TimeTaken, &cometScraper.Applicant, &cometScraper.Status, &cometScraper.CreatedAt, &cometScraper.UpdatedAt)
		if err != nil {
			return cometScrapers, err
		}

		cometScrapers = append(cometScrapers, cometScraper)
	}

	return cometScrapers, nil
}

func (r *pgsqlCometScraperRepository) Delete(ctx context.Context, id string) (err error) {
	query := "DELETE FROM comet_scraper WHERE uuid = $1"
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return
	}

	if affect != 1 {
		err = fmt.Errorf("weird behavior, total affected: %d", affect)
	}

	return
}
