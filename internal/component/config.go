package component

import (
	_ "embed"
	"text/template"

	"github.com/catalystgo/cli/internal/domain"
)

var (
	//go:embed template/config/config-production.yml
	configProductionContent []byte
	//go:embed template/config/config-staging.yml
	configStagingContent []byte
	//go:embed template/config/config-development.yml
	configDevelopmentContent []byte
)

var (
	configProduction  = domain.MustParseTemplate("config-production.yml", configProductionContent)
	configStaging     = domain.MustParseTemplate("config-staging.yml", configStagingContent)
	configDevelopment = domain.MustParseTemplate("config-development.yml", configDevelopmentContent)
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
			file:    "config-production.yml",
			tmpl:    configProduction,
			AppName: getAppNameFromModule(module),
		},
		configComponent{
			file:    "config-staging.yml",
			tmpl:    configStaging,
			AppName: getAppNameFromModule(module),
		},
		configComponent{
			file:    "config-development.yml",
			tmpl:    configDevelopment,
			AppName: getAppNameFromModule(module),
		},
	}
}
