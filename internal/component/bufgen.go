package component

import (
	"bytes"
	"text/template"
)

var _ Component = bufGenComponent{}

var bufGenComponentTemplate *template.Template

func init() {
	var err error

	bufGenComponentTemplate, err = template.New("buf-gen").Parse(
		`version: v1

managed:
  enabled: true
  go_package_prefix:
    default: {{.Module}}

plugins:
  - plugin: go
    out: pkg
    opt: paths=source_relative

  - plugin: go-grpc
    out: pkg
    opt:
      - paths=source_relative

  - plugin: grpc-gateway
    out: pkg
    opt:
      - paths=source_relative
      - generate_unbound_methods=true

  - plugin: vtproto
    out: pkg
    opt:
      - paths=source_relative
      - features=marshal+unmarshal+size

  - plugin: openapiv2
    out: pkg
`)

	if err != nil {
		panic(err)
	}
}

type (
	bufGenComponent struct {
		AppName string
		Module  string
	}
)

// Content implements Component.
func (f bufGenComponent) Content() ([]byte, error) {
	var (
		b   bytes.Buffer
		err error
	)

	err = bufGenComponentTemplate.Execute(&b, f)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// Name implements Component.
func (f bufGenComponent) Name() string {
	return "buf.gen.yml"
}

// Path implements Component.
func (f bufGenComponent) Path() string {
	return "."
}

func NewBufGenComponent(module string) Component {
	return bufGenComponent{
		Module:  module,
		AppName: getAppNameFromModule(module)}
}
