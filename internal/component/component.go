package component

import (
	"bytes"
	"path/filepath"
	"text/template"

	"github.com/catalystgo/cli/internal/log"
)

func init() {
	// Go
	loadTemplate("go.mod", gomodContent, &gomodTemplate)

	// Docker
	loadTemplate("docker-compose.yml", dockerComposeContent, &dockerComposeTemplate)
	loadTemplate("Dockerfile", dockerfileContent, &dockerTemplate)
}

type Component interface {
	Name() string
	Path() string
	Content() ([]byte, error)
}

func getAppNameFromModule(module string) string {
	return filepath.Base(module)
}

func loadTemplate(name string, content []byte, t *template.Template) {
	gotT, err := template.New(name).Parse(string(content))
	if err != nil {
		log.Errorf("parse template %s: %v", name, err)
		return
	}
	*t = *gotT
}

func executeTemplate(t *template.Template, data interface{}) ([]byte, error) {
	var b bytes.Buffer
	err := t.Execute(&b, data)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
