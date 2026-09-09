package redirect

import (
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	Increment(ctx context.Context, linkID int) (*Redirect, error)
}

type PostgresRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Increment(ctx context.Context, linkID int) (*Redirect, error) {
	var err error
	var res Redirect

	query := "SELECT * FROM increment_redirect_counter ($1)"
	err = r.db.QueryRowContext(ctx, query, linkID).
		Scan(&res.ID, &res.LinkID, &res.Counter)
	if err != nil {
		return nil, fmt.Errorf("increment redirect counter: ", err)
	}
	return &res, nil
}
