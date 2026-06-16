package cli

import (
	"fmt"
	"strings"
	"time"
)

const (
	Reset  = "\033[0m"
	Bold   = "\033[1m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Cyan   = "\033[36m"
	Gray   = "\033[90m"
)

func printBanner(banner, version string) {
	fmt.Println(Cyan + banner + Reset)
	fmt.Printf("\n  %sMigration CLI%s  v%s\n", Bold, Reset, version)
	fmt.Printf("  %s%s%s\n\n", Gray, time.Now().Format("2006-01-02 15:04:05"), Reset)
}

func printSection(title string) {
	fmt.Printf("  %s┌─%s %s\n", Cyan, Reset, title)
	fmt.Printf("  %s│%s\n", Cyan, Reset)
}

func printSuccess(msg string) {
	fmt.Printf("  %s✅  %s%s\n\n", Green, msg, Reset)
}

func printError(msg string) {
	fmt.Printf("  %s❌  %s%s\n\n", Red, msg, Reset)
}

func printRow(icon, module, migType, filename, extra string) {
	fmt.Printf("  %s│%s  %s  %s[%-20s]%s  %s[%-6s]%s  %-48s  %s%s%s\n",
		Cyan, Reset,
		icon,
		Gray, module, Reset,
		Gray, migType, Reset,
		filename,
		Gray, extra, Reset,
	)
}

func printDivider() {
	fmt.Printf("  %s│%s  %s\n", Cyan, Reset, strings.Repeat("─", 98))
}

func printFooter(msg string) {
	fmt.Printf("  %s└─%s %s\n\n", Cyan, Reset, msg)
}

func printHelp(version string) {
	fmt.Printf("  %sComandos:%s\n\n", Bold, Reset)

	cmds := [][]string{
		{"up", "Aplicar todas las migraciones pendientes"},
		{"status", "Ver estado de todas las migraciones"},
		{"validate", "Validar checksums de migraciones aplicadas"},
		{"new <modulo> <tipo> <nombre>", "Crear nueva migración"},
		{"help", "Mostrar esta ayuda"},
	}

	for _, cmd := range cmds {
		fmt.Printf("    %s%-35s%s %s\n", Cyan, cmd[0], Reset, cmd[1])
	}

	fmt.Printf("\n  %sTipos de migración:%s\n\n", Bold, Reset)
	fmt.Printf("    %s%-10s%s %s\n", Green, "table", Reset, "CREATE TABLE, ALTER TABLE, índices")
	fmt.Printf("    %s%-10s%s %s\n", Green, "sp", Reset, "Stored procedures y funciones")
	fmt.Printf("    %s%-10s%s %s\n", Green, "seed", Reset, "Datos iniciales (INSERT)")

	fmt.Printf("\n  %sEjemplos:%s\n\n", Bold, Reset)
	fmt.Printf("    migrate up\n")
	fmt.Printf("    migrate status\n")
	fmt.Printf("    migrate validate\n")
	fmt.Printf("    migrate new treasury table  create_bank_accounts\n")
	fmt.Printf("    migrate new treasury sp     bank_accounts_create_v1\n")
	fmt.Printf("    migrate new treasury seed   payment_methods\n")

	fmt.Printf("\n  %sVariables de entorno:%s\n\n", Bold, Reset)
	fmt.Printf("    %s%-20s%s postgres://user:pass@host:port/dbname\n", Cyan, "DATABASE_URL", Reset)
	fmt.Printf("    %s%-20s%s ruta al directorio de migraciones\n\n", Cyan, "MIGRATIONS_PATH", Reset)
}
