package main

import (
	"context"
	"os"

	"cloud-notes/internal/migrator"
)

func main() {
	ctx := context.Background()
	dir := os.Getenv("MIGRATIONS_DIR")
	url := os.Getenv("POSTGRES_URL")

	migrator.MustMigrate(ctx, dir, url)
}
