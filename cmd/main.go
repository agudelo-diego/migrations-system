package main

import (
	"context"
	"fmt"
	"os"

	"github.com/agudelo-diego/migrations-system/cli"
	"github.com/agudelo-diego/migrations-system/migration"
	"github.com/jackc/pgx/v5/pgxpool"
)

// kk
func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		fmt.Println("  ❌  DATABASE_URL no configurado")
		os.Exit(1)
	}

	basePath := os.Getenv("MIGRATIONS_PATH")
	if basePath == "" {
		fmt.Println("  ❌  MIGRATIONS_PATH no configurado")
		os.Exit(1)
	}

	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Printf("  ❌  Error conectando DB: %s\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(context.Background()); err != nil {
		fmt.Printf("  ❌  No se pudo conectar a PostgreSQL: %s\n", err)
		os.Exit(1)
	}

	registry := migration.NewRegistry()
	if err := registry.DiscoverFromPath(basePath); err != nil {
		fmt.Printf("  ❌  %s\n", err)
		os.Exit(1)
	}

	c := cli.New(db, registry, basePath)
	c.Run(os.Args)
}
