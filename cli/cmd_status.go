package cli

import (
	"context"
	"fmt"

	"github.com/agudelo-diego/migrations-system/migration"
)

func runStatus(ctx context.Context, runner *migration.Runner) {
	printSection("Estado de migraciones")

	results, err := runner.Status(ctx)
	if err != nil {
		printError(err.Error())
		return
	}

	printDivider()
	fmt.Printf("  %s│%s  %-4s  %-20s  %-6s  %-48s  %s\n",
		Cyan, Reset, "", "MÓDULO", "TIPO", "ARCHIVO", "APLICADO")
	printDivider()

	applied := 0
	pending := 0

	for _, r := range results {
		if r.Pending {
			pending++
			printRow("⏳", r.Module, r.Type, r.Filename, "pendiente")
		} else {
			applied++
			printRow("✅", r.Module, r.Type, r.Filename,
				r.AppliedAt.Format("2006-01-02 15:04"),
			)
		}
	}

	fmt.Printf("  %s│%s\n", Cyan, Reset)
	printFooter(fmt.Sprintf(
		"Total: %d | %s✅ Aplicadas: %d%s | %s⏳ Pendientes: %d%s",
		len(results),
		Green, applied, Reset,
		Yellow, pending, Reset,
	))

	if err != nil {
		printError(err.Error())
	}
}
