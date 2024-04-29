package component

import (
	"bytes"
	"text/template"
)

var _ Component = dockerComposeComponent{}

var dockerComposeComponentTemplate *template.Template

func init() {
	var err error

	dockerComposeComponentTemplate, err = template.New("docker-compose").Parse(
		`version: "3.8"

services:
  {{.AppName}}:
    container_name: "{{.AppName}}"
    restart: unless-stopped
    volumes:
      - .:/go/src/{{.Module}}
    build:
      context: .
      dockerfile: ./development/Dockerfile
      target: development
    # depends_on:
    #   srv1:
    #     condition: service_healthy
    #   srv2:
    #     condition: service_started
`)

	if err != nil {
		panic(err)
	}
}

type (
	dockerComposeComponent struct {
		AppName string
		Module  string
	}
)

// Content implements Component.
func (d dockerComposeComponent) Content() ([]byte, error) {
	var (
		b   bytes.Buffer
		err error
	)

	err = dockerComposeComponentTemplate.Execute(&b, d)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// Name implements Component.
func (d dockerComposeComponent) Name() string {
	return "docker-compose.yml"
}

// Path implements Component.
func (d dockerComposeComponent) Path() string {
	return ".deployment"
}

func NewDockerComposeComponent(module string) Component {
	return dockerComposeComponent{
		Module:  module,
		AppName: getAppNameFromModule(module)}
}
