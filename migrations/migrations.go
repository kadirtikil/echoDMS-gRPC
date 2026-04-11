package migrations

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunInitMigration(ctx context.Context, pool *pgxpool.Pool) error {
	sql, err := os.ReadFile("migrations/sql/001_init.sql")
	if err != nil {
		panic(err)
	}

	log.Default().Println(string(sql))
	_, err = pool.Exec(ctx, string(sql))
	if err != nil {
		return err
	}

	return nil
}

func RunCustomMigration(ctx context.Context, pool *pgxpool.Pool, migrationFile string) error {
	sql, err := os.ReadFile("migrations/sql/" + migrationFile)
	if err != nil {
		panic(err)
	}

	_, err = pool.Exec(ctx, string(sql))
	if err != nil {
		return err
	}

	return nil
}
