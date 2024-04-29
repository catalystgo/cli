package component

import (
	"bytes"
	"text/template"
)

var _ Component = gomodComponent{}

var gomodComponentTemplate *template.Template

func init() {
	var err error

	gomodComponentTemplate, err = template.New("docker-compose").Parse(
		`module {{.Module}}

go 1.22.0
`)

	if err != nil {
		panic(err)
	}
}

type gomodComponent struct {
	Module  string
	Version string
}

// Content implements Component.
func (d gomodComponent) Content() ([]byte, error) {
	var (
		b   bytes.Buffer
		err error
	)

	err = gomodComponentTemplate.Execute(&b, d)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// Name implements Component.
func (d gomodComponent) Name() string {
	return "go.mod"
}

// Path implements Component.
func (d gomodComponent) Path() string {
	return "."
}

func NewGomodComponent(module string, version string) Component {
	return gomodComponent{Module: module, Version: version}
}
