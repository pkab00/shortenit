package redirect

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	Increment(ctx context.Context, code string) (*Redirect, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Increment(ctx context.Context, code string) (*Redirect, error) {
	var err error
	var res Redirect
	var query string

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("increment redirect counter: %w", err)
	}
	defer tx.Commit()

	var id int
	query = "SELECT link_id FROM links WHERE linK_code = $1"
	err = tx.QueryRowContext(ctx, query, code).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("increment redirect counter: %w", err)
	}

	query = "SELECT * FROM increment_redirect_counter ($1)"
	err = r.db.QueryRowContext(ctx, query, id).
		Scan(&res.ID, &res.LinkID, &res.Counter)
	if err != nil {
		return nil, fmt.Errorf("increment redirect counter: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("increment redirect counter: %w", err)
	}

	return &res, nil
}
