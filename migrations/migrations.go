package migrations

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func checkIfMigrationApplied(ctx context.Context, pool *pgxpool.Pool, name string) (bool, error) {
	var exists bool
	err := pool.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM Migrations WHERE name = $1)", name).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func RunInitMigration(ctx context.Context, pool *pgxpool.Pool) error {
	exists, err := checkIfMigrationApplied(ctx, pool, "001_init.sql")
	if err != nil {
		return err
	}

	if exists {
		log.Printf("Migration 001_init.sql already applied, skipping.")
		return nil
	}

	sql, err := os.ReadFile("migrations/sql/001_init.sql")
	if err != nil {
		panic(err)
	}

	log.Default().Println(string(sql))
	_, err = pool.Exec(ctx, string(sql))
	if err != nil {
		return err
	}

	_, err = pool.Exec(ctx, "INSERT INTO migrations (name) VALUES ($1)", "001_init.sql")
	return nil
}

func RunCustomMigration(ctx context.Context, pool *pgxpool.Pool, migrationFile string) error {
	exists, err := checkIfMigrationApplied(ctx, pool, migrationFile)
	if err != nil {
		return err
	}

	if exists {
		log.Printf("Migration %s already applied, skipping.", migrationFile)
		return nil
	}

	sql, err := os.ReadFile("migrations/sql/" + migrationFile)
	if err != nil {
		panic(err)
	}

	_, err = pool.Exec(ctx, string(sql))
	if err != nil {
		return err
	}

	_, err = pool.Exec(ctx, "INSERT INTO migrations (name) VALUES ($1)", migrationFile)
	return nil
}

func CreateMigrationsTable(ctx context.Context, pool *pgxpool.Pool) error {
	var exists bool
	err := pool.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'migrations')").Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		log.Println("Migrations table already exists, skipping creation.")
		return nil
	}

	sql, err := os.ReadFile("migrations/sql/000_migration.sql")
	if err != nil {
		panic(err)
	}

	_, err = pool.Exec(ctx, string(sql))
	return err
}
