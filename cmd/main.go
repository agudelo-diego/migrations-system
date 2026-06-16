package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/agudelo-diego/migrations-system/cli"
	"github.com/agudelo-diego/migrations-system/migration"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var scanner = bufio.NewScanner(os.Stdin)

func prompt(label string) string {
	fmt.Printf("  \033[36m%s:\033[0m ", label)
	scanner.Scan()
	return strings.TrimSpace(scanner.Text())
}

func promptNew(c *cli.CLI) {
	fmt.Println()

	// Módulo — reintenta si está vacío
	var module string
	for {
		module = prompt("Módulo (ej: sales, banks)")
		if module != "" {
			break
		}
		fmt.Println("  ⚠️   Módulo requerido, intenta de nuevo")
	}

	// Tipo — reintenta si es inválido
	var migType string
	for {
		fmt.Println("  Tipos: [1] table  [2] sp  [3] seed")
		typeInput := prompt("Tipo")
		switch typeInput {
		case "1", "table":
			migType = "table"
		case "2", "sp":
			migType = "sp"
		case "3", "seed":
			migType = "seed"
		default:
			fmt.Println("  ⚠️   Tipo inválido, intenta de nuevo")
			continue
		}
		break
	}

	// Nombre — reintenta si está vacío
	var name string
	for {
		if migType == "sp" {
			spName := prompt("Nombre SP (ej: get_all_v1, create_v1)")
			if spName == "" {
				fmt.Println("  ⚠️   Nombre requerido, intenta de nuevo")
				continue
			}
			name = module + "__" + spName
		} else {
			name = prompt("Nombre (ej: create_bank_accounts)")
			if name == "" {
				fmt.Println("  ⚠️   Nombre requerido, intenta de nuevo")
				continue
			}
		}
		break
	}

	fmt.Println()
	c.Run([]string{"migrate", "new", module, migType, name})
}

func main() {
	_ = godotenv.Load()

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

	// Modo directo con argumentos
	if len(os.Args) > 1 {
		c.Run(os.Args)
		return
	}

	// Modo interactivo
	c.Run([]string{"migrate"})

	for {
		fmt.Print("  \033[36m❯\033[0m ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		switch input {
		case "exit", "quit", "q":
			fmt.Println("\n  👋  Hasta luego\n")
			return

		case "new":
			promptNew(c)

		case "up", "status", "validate", "help":
			args := append([]string{"migrate"}, strings.Fields(input)...)
			c.Run(args)

		default:
			fmt.Printf("  ⚠️   Comando desconocido: '%s' — escribe %shelp%s para ver los comandos\n",
				input, "\033[36m", "\033[0m")
		}
	}
}
