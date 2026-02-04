package main

import (
	"context"
	"log"

	"github.com/aryeko/modkit/examples/hello-mysql/internal/platform/config"
	"github.com/aryeko/modkit/examples/hello-mysql/internal/platform/mysql"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	db, err := mysql.Open(cfg.MySQLDSN)
	if err != nil {
		log.Fatalf("open db failed: %v", err)
	}
	defer db.Close()

	if err := mysql.ApplyMigrations(ctx, db, "migrations"); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}
}
