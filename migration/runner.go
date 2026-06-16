package migration

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io/fs"
	"sort"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const createTable = `
CREATE TABLE IF NOT EXISTS schema_migrations (
    id         SERIAL       PRIMARY KEY,
    module     VARCHAR(100) NOT NULL,
    version    BIGINT       NOT NULL,
    filename   VARCHAR(255) NOT NULL,
    type       VARCHAR(20)  NOT NULL,
    checksum   VARCHAR(64)  NOT NULL,
    applied_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (module, version)
)`

type Runner struct {
	db       *pgxpool.Pool
	registry *Registry
}

func NewRunner(db *pgxpool.Pool, registry *Registry) *Runner {
	return &Runner{db: db, registry: registry}
}

func (r *Runner) ensureTable(ctx context.Context) error {
	_, err := r.db.Exec(ctx, createTable)
	return err
}

func (r *Runner) Collect() ([]Migration, error) {
	var all []Migration

	for _, mod := range r.registry.Modules() {
		files, err := fs.Glob(mod.FS, "*.sql")
		if err != nil {
			return nil, fmt.Errorf("error leyendo módulo '%s': %w", mod.Name, err)
		}

		for _, file := range files {
			version, migType, name, err := ParseFilename(file)
			if err != nil {
				// ← warning y continúa, no rompe
				fmt.Printf("  ⚠️   [%s] ignorando '%s': nombre inválido\n", mod.Name, file)
				continue
			}

			content, err := fs.ReadFile(mod.FS, file)
			if err != nil {
				return nil, fmt.Errorf("[%s] error leyendo %s: %w", mod.Name, file, err)
			}

			all = append(all, Migration{
				Module:   mod.Name,
				Version:  version,
				Name:     name,
				Type:     migType,
				Filename: file,
				Content:  content,
			})
		}
	}

	sort.Slice(all, func(i, j int) bool {
		if all[i].Version == all[j].Version {
			return all[i].Type.Order() < all[j].Type.Order()
		}
		return all[i].Version < all[j].Version
	})

	return all, nil
}

func (r *Runner) Up(ctx context.Context) ([]RunResult, error) {
	if err := r.ensureTable(ctx); err != nil {
		return nil, fmt.Errorf("error creando tabla de migraciones: %w", err)
	}

	migrations, err := r.Collect()
	if err != nil {
		return nil, err
	}

	var results []RunResult

	for _, m := range migrations {
		var exists bool
		r.db.QueryRow(ctx,
			`SELECT EXISTS(
				SELECT 1 FROM schema_migrations
				WHERE module=$1 AND version=$2
			)`,
			m.Module, m.Version,
		).Scan(&exists)

		if exists {
			results = append(results, RunResult{Migration: m, Skipped: true})
			continue
		}

		start := time.Now()

		tx, err := r.db.Begin(ctx)
		if err != nil {
			return results, err
		}

		if _, err := tx.Exec(ctx, string(m.Content)); err != nil {
			tx.Rollback(ctx)
			results = append(results, RunResult{
				Migration: m,
				Error:     err,
				Duration:  time.Since(start),
			})
			// ← no retorna error, continúa con las demás
			continue
		}

		if _, err := tx.Exec(ctx,
			`INSERT INTO schema_migrations
			 (module, version, filename, type, checksum)
			 VALUES ($1,$2,$3,$4,$5)`,
			m.Module, m.Version, m.Filename,
			string(m.Type), checksum(m.Content),
		); err != nil {
			tx.Rollback(ctx)
			results = append(results, RunResult{
				Migration: m,
				Error:     err,
				Duration:  time.Since(start),
			})
			continue
		}

		if err := tx.Commit(ctx); err != nil {
			results = append(results, RunResult{
				Migration: m,
				Error:     err,
				Duration:  time.Since(start),
			})
			continue
		}

		results = append(results, RunResult{
			Migration: m,
			Applied:   true,
			Duration:  time.Since(start),
		})
	}

	return results, nil
}

func (r *Runner) Status(ctx context.Context) ([]StatusResult, error) {
	if err := r.ensureTable(ctx); err != nil {
		return nil, err
	}

	migrations, err := r.Collect()
	if err != nil {
		return nil, err
	}

	applied := map[string]time.Time{}
	rows, err := r.db.Query(ctx,
		`SELECT module, version, applied_at FROM schema_migrations`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var module string
		var version int64
		var appliedAt time.Time
		rows.Scan(&module, &version, &appliedAt)
		applied[migKey(module, version)] = appliedAt
	}

	var results []StatusResult
	for _, m := range migrations {
		if t, ok := applied[migKey(m.Module, m.Version)]; ok {
			results = append(results, StatusResult{
				Module:    m.Module,
				Version:   m.Version,
				Filename:  m.Filename,
				Type:      string(m.Type),
				AppliedAt: &t,
			})
		} else {
			results = append(results, StatusResult{
				Module:   m.Module,
				Version:  m.Version,
				Filename: m.Filename,
				Type:     string(m.Type),
				Pending:  true,
			})
		}
	}

	return results, nil
}

func (r *Runner) Validate(ctx context.Context) ([]ValidateResult, error) {
	rows, err := r.db.Query(ctx,
		`SELECT module, filename, checksum
		 FROM schema_migrations
		 ORDER BY version`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []ValidateResult

	for rows.Next() {
		var module, filename, storedChecksum string
		rows.Scan(&module, &filename, &storedChecksum)

		var content []byte
		for _, mod := range r.registry.Modules() {
			if mod.Name == module {
				content, _ = fs.ReadFile(mod.FS, filename)
				break
			}
		}

		if content == nil {
			results = append(results, ValidateResult{
				Module:   module,
				Filename: filename,
				Valid:    false,
				Error:    "archivo no encontrado en el sistema",
			})
			continue
		}

		if checksum(content) != storedChecksum {
			results = append(results, ValidateResult{
				Module:   module,
				Filename: filename,
				Valid:    false,
				Error:    "checksum inválido — migración modificada después de aplicarse",
			})
		} else {
			results = append(results, ValidateResult{
				Module:   module,
				Filename: filename,
				Valid:    true,
			})
		}
	}

	return results, nil
}

func checksum(content []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(content))
}

func migKey(module string, version int64) string {
	return fmt.Sprintf("%s_%d", module, version)
}
