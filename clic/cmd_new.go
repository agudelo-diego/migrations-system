package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/agudelo-diego/migrations-system/migration"
)

var templates = map[migration.Type]string{
	migration.TypeTable: `-- ══════════════════════════════════════════════════════
-- Módulo:  %s
-- Tipo:    table
-- Nombre:  %s
-- Fecha:   %s
-- ══════════════════════════════════════════════════════

CREATE TABLE %s (
    id         VARCHAR(26)  PRIMARY KEY DEFAULT generate_ulid(),
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);
`,

	migration.TypeSP: `-- ══════════════════════════════════════════════════════
-- Módulo:  %s
-- Tipo:    sp
-- Nombre:  %s
-- Fecha:   %s
-- ══════════════════════════════════════════════════════

CREATE OR REPLACE FUNCTION %s(

)
RETURNS VOID
LANGUAGE plpgsql
AS $$
DECLARE

BEGIN

END;
$$;
`,

	migration.TypeSeed: `-- ══════════════════════════════════════════════════════
-- Módulo:  %s
-- Tipo:    seed
-- Nombre:  %s
-- Fecha:   %s
-- ══════════════════════════════════════════════════════

INSERT INTO %s () VALUES ()
ON CONFLICT DO NOTHING;
`,
}

func runNew(basePath, module, migType, name string) {
	t := migration.Type(migType)

	if !t.IsValid() {
		printError(fmt.Sprintf(
			"tipo inválido '%s' — válidos: table | sp | seed", migType,
		))
		return
	}

	template := templates[t]
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s_%s.sql", timestamp, migType, name)
	dir := filepath.Join(basePath, module)
	path := filepath.Join(dir, filename)

	if err := os.MkdirAll(dir, 0755); err != nil {
		printError(fmt.Sprintf("no se pudo crear el directorio: %s", err))
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	content := fmt.Sprintf(template, module, name, now, name)

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		printError(fmt.Sprintf("no se pudo crear el archivo: %s", err))
		return
	}

	printSection("Nueva migración creada")
	fmt.Printf("  %s│%s  📦 Módulo:  %s\n", Cyan, Reset, module)
	fmt.Printf("  %s│%s  📋 Tipo:    %s (%s)\n", Cyan, Reset, migType, t.Label())
	fmt.Printf("  %s│%s  📄 Archivo: %s\n", Cyan, Reset, path)
	printFooter("Edita el archivo y ejecuta: migrate up")
}
