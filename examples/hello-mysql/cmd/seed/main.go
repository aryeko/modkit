package main

import (
	"context"
	"log"

	"github.com/aryeko/modkit/examples/hello-mysql/internal/platform/config"
	"github.com/aryeko/modkit/examples/hello-mysql/internal/platform/mysql"
	"github.com/aryeko/modkit/examples/hello-mysql/internal/seed"
)

func main() {
	cfg := config.Load()
	db, err := mysql.Open(cfg.MySQLDSN)
	if err != nil {
		log.Fatalf("open db failed: %v", err)
	}
	defer db.Close()

	if err := seed.Seed(context.Background(), db); err != nil {
		log.Fatalf("seed failed: %v", err)
	}

	log.Printf("seed complete")
}
