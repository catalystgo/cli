package component

import (
	_ "embed"

	"github.com/catalystgo/cli/internal/domain"
)

var (
	//go:embed template/docker/Dockerfile
	dockerfileContent []byte
)

var (
	dockerTemplate = domain.MustParseTemplate("Dockerfile", dockerfileContent)
)

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

func (d dockerComponent) Name() string {
	return "Dockerfile"
}

func (d dockerComponent) Path() string {
	return ".catalystgo"
}
