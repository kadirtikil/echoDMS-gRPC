package db_utils

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ReseedTestDB(ctx context.Context, pool *pgxpool.Pool) error {
	resetSql, err := os.ReadFile("../../db/sql/reset_test_db.sql")
	if err != nil {
		return err
	}
	_, err = pool.Exec(ctx, string(resetSql))
	if err != nil {
		return err
	}

	seedSql, err := os.ReadFile("../../db/sql/002_seed_test_db.sql")
	if err != nil {
		return err
	}
	_, err = pool.Exec(ctx, string(seedSql))

	return err
}
