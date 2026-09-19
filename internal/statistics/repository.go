package statistics

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pkab00/shortenit/internal/apperr"
)

type Repository interface {
	Get(ctx context.Context, code string) (*Statistics, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Get(ctx context.Context, code string) (*Statistics, error) {
	var err error
	var res Statistics

	query := "SELECT links.link_id, link_code, url, created_at, redirect_counter " +
		"FROM links LEFT JOIN redirects ON links.link_id = redirects.link_id " +
		"WHERE links.link_code = $1"

	err = r.db.
		QueryRowContext(ctx, query, code).
		Scan(&res.ID, &res.Code, &res.URL, &res.CreatedAt, &res.Redirects)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get by statistics: %w", apperr.ErrorLinkNotFound)
		}
		return nil, fmt.Errorf("get statistics: %w", err)
	}
	return &res, nil
}
