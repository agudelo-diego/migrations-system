package cli

import (
	"context"
	"fmt"

	"github.com/hftech/migrate/internal/migration"
)

func runValidate(ctx context.Context, runner *migration.Runner) {
	printSection("Validando checksums")

	results, err := runner.Validate(ctx)
	if err != nil {
		printError(err.Error())
		return
	}

	errCount := 0

	for _, r := range results {
		if r.Valid {
			printRow("✅", r.Module, "", r.Filename, "ok")
		} else {
			errCount++
			printRow("❌", r.Module, "", r.Filename, "inválido")
			fmt.Printf("  %s│%s       └─ %s%s%s\n",
				Cyan, Reset, Red, r.Error, Reset)
		}
	}

	fmt.Printf("  %s│%s\n", Cyan, Reset)
	printFooter(fmt.Sprintf("Total: %d | Errores: %d", len(results), errCount))

	if errCount > 0 {
		printError(fmt.Sprintf("%d migración(es) con checksum inválido", errCount))
		return
	}

	printSuccess("Todos los checksums son válidos")
}
