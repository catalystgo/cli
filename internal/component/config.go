package component

import (
	_ "embed"
	"text/template"

	"github.com/catalystgo/cli/internal/domain"
)

var (
	//go:embed template/config/config.yml
	configContent []byte
)

var (
	configTemplate = domain.MustParseTemplate("config-prod.yml", configContent)
)

type configComponent struct {
	file    string
	tmpl    *template.Template
	AppName string
}

func (c configComponent) Content() ([]byte, error) {
	return domain.ExecuteTemplate(c.tmpl, c)
}

func (c configComponent) Name() string {
	return c.file
}

func (c configComponent) Path() string {
	return ".catalystgo"
}

func NewConfigComponent(module string) []Component {
	return []Component{
		configComponent{
			file:    "config-prod.yml",
			tmpl:    configTemplate,
			AppName: getAppNameFromModule(module),
		},
		configComponent{
			file:    "config-stage.yml",
			tmpl:    configTemplate,
			AppName: getAppNameFromModule(module),
		},
		configComponent{
			file:    "config-local.yml",
			tmpl:    configTemplate,
			AppName: getAppNameFromModule(module),
		},
	}
}
