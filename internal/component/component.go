package component

import "path/filepath"

type Component interface {
	Name() string
	Path() string
	Content() ([]byte, error)
}

func getAppNameFromModule(module string) string {
	return filepath.Base(module)
}
