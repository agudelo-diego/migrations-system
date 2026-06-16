package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/hftech/migrate/internal/migration"
	"github.com/jackc/pgx/v5/pgxpool"
)

const version = "1.0.0"

const banner = `
 ███╗   ███╗██╗ ██████╗ ██████╗  █████╗ ████████╗███████╗
 ████╗ ████║██║██╔════╝ ██╔══██╗██╔══██╗╚══██╔══╝██╔════╝
 ██╔████╔██║██║██║  ███╗██████╔╝███████║   ██║   █████╗
 ██║╚██╔╝██║██║██║   ██║██╔══██╗██╔══██║   ██║   ██╔══╝
 ██║ ╚═╝ ██║██║╚██████╔╝██║  ██║██║  ██║   ██║   ███████╗
 ╚═╝     ╚═╝╚═╝ ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝  ╚═╝   ╚══════╝`

type CLI struct {
	runner   *migration.Runner
	basePath string
}

func New(db *pgxpool.Pool, registry *migration.Registry, basePath string) *CLI {
	return &CLI{
		runner:   migration.NewRunner(db, registry),
		basePath: basePath,
	}
}

func (c *CLI) Run(args []string) {
	printBanner(banner, version)

	if len(args) < 2 {
		printHelp(version)
		return
	}

	ctx := context.Background()

	switch args[1] {

	case "up":
		runUp(ctx, c.runner)

	case "status":
		runStatus(ctx, c.runner)

	case "validate":
		runValidate(ctx, c.runner)

	case "new":
		if len(args) < 5 {
			printError("Uso: migrate new <modulo> <tipo> <nombre>")
			fmt.Println("  Ejemplos:")
			fmt.Println("    migrate new treasury table  create_bank_accounts")
			fmt.Println("    migrate new treasury sp     bank_accounts_create_v1")
			fmt.Println("    migrate new treasury seed   payment_methods\n")
			os.Exit(1)
		}
		runNew(c.basePath, args[2], args[3], args[4])

	case "help", "--help", "-h":
		printHelp(version)

	default:
		printError(fmt.Sprintf("comando desconocido: '%s'", args[1]))
		printHelp(version)
		os.Exit(1)
	}
}
