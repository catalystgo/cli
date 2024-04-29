package component

var _ Component = reviveComponent{}

type (
	reviveComponent struct{}
)

// Content implements Component.
func (d reviveComponent) Content() ([]byte, error) {
	return []byte(
		`ignoreGeneratedHeader = false
severity = "warning"
confidence = 0.8
errorCode = 0
warningCode = 0

[rule.blank-imports]
[rule.context-as-argument]
[rule.context-keys-type]
[rule.dot-imports]
[rule.error-return]
[rule.error-strings]
[rule.error-naming]
[rule.increment-decrement]
[rule.var-naming]
[rule.var-declaration]
[rule.range]
[rule.receiver-naming]
[rule.time-naming]
[rule.unexported-return]
[rule.indent-error-flow]
[rule.errorf]
[rule.empty-block]
[rule.superfluous-else]
[rule.unused-parameter]
[rule.unreachable-code]
[rule.redefines-builtin-id]
[rule.time-equal]
[rule.if-return]
`,
	), nil
}

// Name implements Component.
func (d reviveComponent) Name() string {
	return "revive.toml"
}

// Path implements Component.
func (d reviveComponent) Path() string {
	return "."
}

func NewReviveComponent() Component {
	return reviveComponent{}
}
