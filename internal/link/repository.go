package link

import (
	"context"
	"database/sql"
)

type Repository interface {
	Create(ctx context.Context, url string) (*Link, error)
	Delete(ctx context.Context, url string) (*Link, error)
	All(ctx context.Context) ([]Link, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, url string) (*Link, error) {
	var res Link

	query := "INSERT INTO links (link_body) VALUES ($1) RETURNING *"
	err := r.db.
		QueryRowContext(ctx, query, url).
		Scan(&res.ID, &res.Body, &res.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, url string) (*Link, error) {
	var res Link

	query := "DELETE FROM links WHERE link_body = $1 RETURNING *"
	err := r.db.
		QueryRowContext(ctx, query, url).
		Scan(&res.ID, &res.Body, &res.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *PostgresRepository) All(ctx context.Context) ([]Link, error) {
	var links []Link

	query := "SELECT * FROM links ORDER BY created_at DESC"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var l Link
		if err := rows.Scan(&l.ID, &l.Body, &l.CreatedAt); err != nil {
			return nil, err
		}
		links = append(links, l)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return links, nil
}
