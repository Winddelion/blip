package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Check struct {
	ID            int64
	URL           string
	Status_code   sql.NullInt32
	Response_time sql.NullInt32
	CheckedAt     time.Time
}

func InsertCheck(ctx context.Context, pool *pgxpool.Pool, url string, status, response_time int32) (int64, error) {
	var id int64
	err := pool.QueryRow(ctx,
		`INSERT INTO checks (url, status_code, response_time)
		VALUES ($1, $2, $3)
		RETURNING id`, url, status, response_time).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("insert check: %w", err)
	}
	return id, nil
}

func GetChecksByURL(ctx context.Context, pool *pgxpool.Pool, url string, limit int) ([]Check, error) {
	rows, err := pool.Query(ctx,
		`SELECT id, url, status_code, response_time, checked_at
		FROM checks
		WHERE url = $1
		ORDER BY checked_at DESC
		LIMIT $2`,
		url, limit,
	)
	if err != nil {
		return nil, fmt.Errorf("query checks: %w", err)
	}
	defer rows.Close()

	checks, err := pgx.CollectRows(rows, pgx.RowToStructByName[Check])
	if err != nil {
		return nil, fmt.Errorf("collect checks: %w", err)
	}
	return checks, nil
}

func GetCheckByID(ctx context.Context, pool *pgxpool.Pool, id int64) (*Check, error) {
	var c Check
	err := pool.QueryRow(ctx,
		`SELECT id, url, status_code, response_time, checked_at
		FROM checks
		WHERE id = $1`,
		id,
	).Scan(&c.ID, &c.URL, &c.Status_code, &c.Response_time, &c.CheckedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}
