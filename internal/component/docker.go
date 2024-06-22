package component

import (
	_ "embed"

	"github.com/catalystgo/cli/internal/domain"
)

var (
	//go:embed template/docker/Dockerfile
	dockerfileContent []byte

	//go:embed template/docker/docker-compose.txt
	dockerComposeContent []byte
)

var (
	dockerTemplate        = domain.MustParseTemplate("Dockerfile", dockerfileContent)
	dockerComposeTemplate = domain.MustParseTemplate("docker-compose.yml", dockerComposeContent)
)

//////////////////
// DOCKER COMPONENT
//////////////////

type dockerComponent struct {
	Module  string
	AppName string
}

func NewDockerComponent(module string) Component {
	return dockerComponent{
		Module:  module,
		AppName: getAppNameFromModule(module)}
}

func (d dockerComponent) Content() ([]byte, error) {
	return domain.ExecuteTemplate(dockerTemplate, d)
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
	Module  string
}

func NewDockerComposeComponent(module string) Component {
	return dockerComposeComponent{
		AppName: getAppNameFromModule(module),
		Module:  module,
	}
}

func (d dockerComposeComponent) Content() ([]byte, error) {
	return domain.ExecuteTemplate(dockerComposeTemplate, d)
}

func (d dockerComposeComponent) Name() string {
	return "docker-compose.yml"
}

func (d dockerComposeComponent) Path() string {
	return ".deployment"
}
