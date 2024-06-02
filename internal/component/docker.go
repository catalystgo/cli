package component

import (
	"bytes"
	_ "embed"
	"text/template"
)

var (
	//go:embed template/docker/Dockerfile
	dockerfileContent []byte

	//go:embed template/docker/docker-compose.txt
	dockerComposeContent []byte
)

var (
	dockerTemplate        template.Template
	dockerComposeTemplate template.Template
)

//////////////////
// DOCKER COMPONENT
//////////////////

type dockerComponent struct {
	AppName string
}

func NewDockerComponent(module string) Component {
	return dockerComponent{
		AppName: getAppNameFromModule(module)}
}

func (d dockerComponent) Content() ([]byte, error) {
	var b bytes.Buffer

	err := dockerTemplate.Execute(&b, d)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// Name implements Component.
func (d dockerComponent) Name() string {
	return "Dockerfile"
}

// Path implements Component.
func (d dockerComponent) Path() string {
	return ".deployment"
}

//////////////////
// DOCKER COMPOSE COMPONENT
//////////////////

type dockerComposeComponent struct {
	AppName string
}

func NewDockerComposeComponent(module string) Component {
	return dockerComposeComponent{
		AppName: getAppNameFromModule(module)}
}

func (d dockerComposeComponent) Content() ([]byte, error) {
	var b bytes.Buffer

	err := dockerComposeTemplate.Execute(&b, d)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (d dockerComposeComponent) Name() string {
	return "docker-compose.yml"
}

func (d dockerComposeComponent) Path() string {
	return ".deployment"
}
