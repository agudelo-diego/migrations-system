package migration

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type Module struct {
	Name string
	FS   fs.FS
	Path string
}

type Registry struct {
	modules []Module
}

func NewRegistry() *Registry {
	return &Registry{}
}

func (r *Registry) RegisterDir(name, path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("directorio no existe: %s", path)
	}
	r.modules = append(r.modules, Module{
		Name: name,
		FS:   os.DirFS(path),
		Path: path,
	})
	return nil
}

func (r *Registry) DiscoverFromPath(basePath string) error {
	entries, err := os.ReadDir(basePath)
	if err != nil {
		return fmt.Errorf("no se pudo leer '%s': %w", basePath, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		modulePath := filepath.Join(basePath, entry.Name())
		r.modules = append(r.modules, Module{
			Name: entry.Name(),
			FS:   os.DirFS(modulePath),
			Path: modulePath,
		})
	}

	if len(r.modules) == 0 {
		return fmt.Errorf("no se encontraron módulos en '%s'", basePath)
	}

	return nil
}

func (r *Registry) Modules() []Module {
	return r.modules
}
