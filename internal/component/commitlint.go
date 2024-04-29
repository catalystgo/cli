package component

var _ Component = commitLintConfigComponent{}

type (
	commitLintConfigComponent struct{}
)

// Content implements Component.
func (c commitLintConfigComponent) Content() ([]byte, error) {
	return []byte(
		`module.exports = { extends: ["@commitlint/config-conventional"] };\n`,
	), nil
}

// Name implements Component.
func (c commitLintConfigComponent) Name() string {
	return "commitlint.config.js"
}

// Path implements Component.
func (c commitLintConfigComponent) Path() string {
	return "."
}

func NewcommitLintConfigComponent() Component {
	return commitLintConfigComponent{}
}
