package cli

import (
	"context"
	"fmt"

	"github.com/agudelo-diego/migrations-system/migration"
)

func runUp(ctx context.Context, runner *migration.Runner) {
	printSection("Aplicando migraciones")

	results, err := runner.Up(ctx)

	applied := 0
	skipped := 0
	errCount := 0

	for _, r := range results {
		switch {
		case r.Error != nil:
			errCount++
			printRow("❌", r.Migration.Module, string(r.Migration.Type),
				r.Migration.Filename, "error")
			fmt.Printf("  %s│%s       └─ %s%s%s\n",
				Cyan, Reset, Red, r.Error.Error(), Reset)

		case r.Skipped:
			skipped++
			printRow("⏭️ ", r.Migration.Module, string(r.Migration.Type),
				r.Migration.Filename, "omitida")

		case r.Applied:
			applied++
			printRow("✅", r.Migration.Module, string(r.Migration.Type),
				r.Migration.Filename,
				fmt.Sprintf("%dms", r.Duration.Milliseconds()),
			)
		}
	}

	fmt.Printf("  %s│%s\n", Cyan, Reset)
	printFooter(fmt.Sprintf(
		"Total: %d | %s✅ Aplicadas: %d%s | ⏭️  Omitidas: %d | %s❌ Errores: %d%s",
		len(results),
		Green, applied, Reset,
		skipped,
		Red, errCount, Reset,
	))

	if err != nil {
		printError(err.Error())
		return
	}

	printSuccess("Migraciones completadas")
}
