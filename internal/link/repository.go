package link

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pkab00/shortenit/internal/apperr"
	"github.com/pkab00/shortenit/pkg/encode"
)

type Repository interface {
	Create(ctx context.Context, url string) (*Link, error)
	Delete(ctx context.Context, code string) (*Link, error)
	All(ctx context.Context) ([]Link, error)
	ByCode(ctx context.Context, code string) (*Link, error)
	ByURL(ctx context.Context, url string) (*Link, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, url string) (*Link, error) {
	var res Link
	var query string
	var err error

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("create link: %w", err)
	}
	defer tx.Rollback()

	var id int
	query = "SELECT nextval('links_link_id_seq')"
	err = tx.
		QueryRowContext(ctx, query).
		Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("create link: %w", err)
	}

	code := encode.NewHashEncoder().Encode(uint64(id))
	err = tx.
		QueryRowContext(ctx, `
    	INSERT INTO links (link_id, link_code, url)
		VALUES ($1, $2, $3)
		RETURNING *
	`, id, code, url).
		Scan(&res.ID, &res.Code, &res.URL, &res.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create link: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("create link: %w", err)
	}

	return &res, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, code string) (*Link, error) {
	var res Link

	query := "DELETE FROM links WHERE link_code = $1 RETURNING *"
	err := r.db.
		QueryRowContext(ctx, query, code).
		Scan(&res.ID, &res.Code, &res.URL, &res.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get by id: %w", apperr.ErrorLinkNotFound)
		}
		return nil, fmt.Errorf("delete link: %w", err)
	}
	return &res, nil
}

func (r *PostgresRepository) All(ctx context.Context) ([]Link, error) {
	var links []Link

	query := "SELECT * FROM links ORDER BY created_at DESC"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get all links: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var l Link
		if err := rows.Scan(&l.ID, &l.Code, &l.URL, &l.CreatedAt); err != nil {
			return nil, fmt.Errorf("get all links: %w", err)
		}
		links = append(links, l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get all links: %w", err)
	}

	return links, nil
}

func (r *PostgresRepository) ByCode(ctx context.Context, code string) (*Link, error) {
	var res Link

	query := "SELECT * FROM links WHERE link_code = $1 LIMIT 1"
	err := r.db.
		QueryRowContext(ctx, query, code).
		Scan(&res.ID, &res.Code, &res.URL, &res.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get by code: %w", apperr.ErrorLinkNotFound)
		}
		return nil, fmt.Errorf("get by code: %w", err)
	}
	return &res, nil
}

func (r *PostgresRepository) ByURL(ctx context.Context, url string) (*Link, error) {
	var res Link

	query := "SELECT * FROM links WHERE url = $1 LIMIT 1"
	err := r.db.
		QueryRowContext(ctx, query, url).
		Scan(&res.ID, &res.URL, &res.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get by url: %w", apperr.ErrorLinkNotFound)
		}
		return nil, fmt.Errorf("get by url: %w", err)
	}
	return &res, nil
}
