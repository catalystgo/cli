package component

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/catalystgo/cli/internal/domain"
)

var (
	//go:embed template/gomod.txt
	gomodContent []byte

	//go:embed template/main.go.txt
	goMainContent []byte

	//go:embed template/gitignore.txt
	gitignoreContent []byte

	//go:embed template/revive.toml
	reviveConfig []byte

	//go:embed template/Taskfile.yml
	taskfileContent []byte
)

var (
	gomodTemplate = domain.MustParseTemplate("go.mod", gomodContent)
)

//////////////////
// gomod component
//////////////////

type gomodComponent struct {
	Module  string
	Version string
}

func NewGomodComponent(module string, version string) Component {
	return gomodComponent{
		Module:  module,
		Version: version,
	}
}

func (d gomodComponent) Content() ([]byte, error) {
	return domain.ExecuteTemplate(gomodTemplate, d)
}

func (d gomodComponent) Name() string {
	return "go.mod"
}

func (d gomodComponent) Path() string {
	return "."
}

//////////////////
// gitignore component
//////////////////

type gitignoreComponent struct{}

func NewGitignoreComponent() Component {
	return gitignoreComponent{}
}

func (d gitignoreComponent) Content() ([]byte, error) {
	return gitignoreContent, nil
}

func (d gitignoreComponent) Name() string {
	return ".gitignore"
}

func (d gitignoreComponent) Path() string {
	return "."
}

//////////////////
// revive component
//////////////////

type reviveComponent struct{}

func NewReviveComponent() Component {
	return reviveComponent{}
}

func (d reviveComponent) Content() ([]byte, error) {
	return reviveConfig, nil
}

func (d reviveComponent) Name() string {
	return "revive.toml"
}

func (d reviveComponent) Path() string {
	return "."
}

//////////////////
// taskfile component
//////////////////

type taskfileComponent struct {
	AppName  string
	Username string
}

func NewTaskfileComponent(module string) Component {
	return taskfileComponent{
		AppName:  getAppNameFromModule(module),
		Username: getUserFromModule(module),
	}
}

func (t taskfileComponent) Content() ([]byte, error) {
	// Here we use `strings.ReplaceAll` instead of go templates because the Taskfile has templates inside.
	file := string(taskfileContent)
	file = strings.ReplaceAll(file, "{{.AppName}}", t.AppName)
	file = strings.ReplaceAll(file, "{{.Username}}", t.Username)

	return []byte(file), nil
}

func (t taskfileComponent) Name() string {
	return "Taskfile.yml"
}

func (t taskfileComponent) Path() string {
	return "."
}

//////////////////
// main.go component
//////////////////

type goMainComponent struct {
	AppName string
}

func NewGoMainComponent(module string) Component {
	return goMainComponent{
		AppName: getAppNameFromModule(module),
	}
}

func (g goMainComponent) Content() ([]byte, error) {
	return goMainContent, nil
}

func (g goMainComponent) Name() string {
	return "main.go"
}

func (g goMainComponent) Path() string {
	return fmt.Sprintf("./cmd/%s", g.AppName)
}
