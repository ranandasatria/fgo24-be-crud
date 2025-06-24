package utils

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB() (*pgxpool.Conn, error) {
	connectionString := "postgres://postgres:1@localhost:5433/backend"

	pool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}

	return conn, nil

}
