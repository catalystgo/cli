package component

import (
	_ "embed"
	"strings"
	"text/template"
)

var (
	//go:embed template/gomod.txt
	gomodContent []byte

	//go:embed template/gitignore.txt
	gitignoreContent []byte

	//go:embed template/revive.toml
	reviveConfig []byte

	//go:embed template/Taskfile.yml
	taskfileContent []byte

	//go:embed template/precommit/pre-commit-config.yml
	preCommitConfig []byte

	//go:embed template/precommit/commitlint.config.js
	commitLintConfig []byte
)

var (
	gomodTemplate template.Template
)

//////////////////
// GOMOD COMPONENT
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
	return executeTemplate(&gomodTemplate, d)
}

func (d gomodComponent) Name() string {
	return "go.mod"
}

func (d gomodComponent) Path() string {
	return "."
}

//////////////////
// GITIGNORE COMPONENT
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
// REVIVE COMPONENT
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
// TASKFILE COMPONENT
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
	// Replace the placeholder with the app name, could be done with go-template
	// because the Taskfile has templates inside.
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

// ////////////////
// PRE-COMMIT COMPONENT
// ////////////////

type preCommitComponent struct{}

func NewPreCommitComponent() Component {
	return preCommitComponent{}
}

func (p preCommitComponent) Content() ([]byte, error) {
	return preCommitConfig, nil
}

func (p preCommitComponent) Name() string {
	return ".pre-commit-config.yml"
}

func (p preCommitComponent) Path() string {
	return "."
}

//////////////////
// COMMIT-LINT CONFIG COMPONENT
//////////////////

type commitLintConfigComponent struct{}

func NewcommitLintConfigComponent() Component {
	return commitLintConfigComponent{}
}

func (c commitLintConfigComponent) Content() ([]byte, error) {
	return commitLintConfig, nil
}

func (c commitLintConfigComponent) Name() string {
	return "commitlint.config.js"
}

func (c commitLintConfigComponent) Path() string {
	return "."
}
