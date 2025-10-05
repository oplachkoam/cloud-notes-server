package migrator

import (
	"context"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
)

func Migrate(ctx context.Context, dir string, url string) error {
	conn, err := pgx.Connect(ctx, url)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	defer conn.Close(ctx) // nolint

	query := `CREATE TABLE IF NOT EXISTS migrations (
		    	id       SERIAL PRIMARY KEY,
		        file     TEXT NOT NULL)`
	_, err = conn.Exec(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to create migration table: %w", err)
	}

	applied := map[string]bool{}
	query = `SELECT file FROM migrations`
	rows, err := conn.Query(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var file string
		if err := rows.Scan(&file); err != nil {
			return fmt.Errorf("failed to scan row: %w", err)
		}
		applied[file] = true
	}

	d, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read migrations dir: %w", err)
	}

	var files []string
	for _, file := range d {
		if strings.HasSuffix(file.Name(), ".sql") {
			files = append(files, file.Name())
		}
	}
	sort.Strings(files)

	if len(files) == 0 {
		fmt.Printf("migrations dir is empty\n")
		return nil
	}

	for _, file := range files {
		if applied[file] {
			fmt.Printf("migration %s already applied\n", file)
			continue
		}

		sql, err := os.ReadFile(path.Join(dir, file))
		if err != nil {
			panic(fmt.Errorf("failed to read migrations file: %w", err))
		}

		fmt.Printf("applying migration %s\n", file)

		tx, err := conn.Begin(ctx)
		if err != nil {
			panic(fmt.Errorf("failed to start transaction: %w", err))
		}

		_, err = tx.Exec(ctx, string(sql))
		if err != nil {
			_ = tx.Rollback(ctx)
			panic(fmt.Errorf("failed to apply migration %s", file))
		}

		query = `INSERT INTO migrations (file) VALUES ($1)`
		_, err = tx.Exec(ctx, query, file)
		if err != nil {
			_ = tx.Rollback(ctx)
			panic(fmt.Errorf("failed to record migration %s", file))
		}

		err = tx.Commit(ctx)
		if err != nil {
			_ = tx.Rollback(ctx)
			panic(fmt.Errorf("failed to commit migration: %w", err))
		}

		fmt.Printf("migration %s applied\n", file)
	}

	fmt.Printf("migrations applied successfully\n")

	return nil
}

func MustMigrate(ctx context.Context, dir string, url string) {
	err := Migrate(ctx, dir, url)
	if err != nil {
		panic(err)
	}
}
