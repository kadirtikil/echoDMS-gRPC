package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 10
	config.MinConns = 1

	pool, err := pgxpool.New(ctx, config.ConnString())
	if err != nil {
		return nil, err
	}

	return pool, nil

}
