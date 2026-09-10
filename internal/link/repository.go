package link

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/pkab00/shortenit/internal/apperr"
)

type Repository interface {
	Create(ctx context.Context, url string) (*Link, error)
	Delete(ctx context.Context, id int) (*Link, error)
	All(ctx context.Context) ([]Link, error)
	ByID(ctx context.Context, id int) (*Link, error)
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

	query := "INSERT INTO links (url) VALUES ($1) RETURNING *"
	err := r.db.
		QueryRowContext(ctx, query, url).
		Scan(&res.ID, &res.URL, &res.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("create link:", err)
	}
	return &res, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id int) (*Link, error) {
	var res Link

	query := "DELETE FROM links WHERE link_id = $1 RETURNING *"
	err := r.db.
		QueryRowContext(ctx, query, id).
		Scan(&res.ID, &res.URL, &res.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("delete links:", err)
	}
	return &res, nil
}

func (r *PostgresRepository) All(ctx context.Context) ([]Link, error) {
	var links []Link

	query := "SELECT * FROM links ORDER BY created_at DESC"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get all links:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var l Link
		if err := rows.Scan(&l.ID, &l.URL, &l.CreatedAt); err != nil {
			return nil, fmt.Errorf("get all links:", err)
		}
		links = append(links, l)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get all links:", err)
	}

	return links, nil
}

func (r *PostgresRepository) ByID(ctx context.Context, id int) (*Link, error) {
	var res Link

	query := "SELECT * FROM links WHERE link_id = $1 LIMIT 1"
	err := r.db.
		QueryRowContext(ctx, query, id).
		Scan(&res.ID, &res.URL, &res.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("get by id:", apperr.ErrorLinkNotFound)
		}
		return nil, fmt.Errorf("get by id:", err)
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
			return nil, fmt.Errorf("get by url:", apperr.ErrorLinkNotFound)
		}
		return nil, fmt.Errorf("get by url:", err)
	}
	return &res, nil
}
