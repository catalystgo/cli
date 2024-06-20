package component

import (
	"path/filepath"
	"strings"
)

type Component interface {
	Name() string
	Path() string
	Content() ([]byte, error)
}

func getAppNameFromModule(module string) string {
	return filepath.Base(module)
}

func getUserFromModule(module string) string {
	return strings.Split(module, "/")[1]
}
