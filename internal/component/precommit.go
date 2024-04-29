package component

var _ Component = preCommitComponent{}

type (
	preCommitComponent struct{}
)

// Content implements Component.
func (p preCommitComponent) Content() ([]byte, error) {
	return []byte(
		`repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.5.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-added-large-files
      - id: detect-private-key

      - id: check-yaml
      - id: check-toml
      - id: sort-simple-yaml

  - repo: https://github.com/dnephin/pre-commit-golang
    rev: v0.5.1
    hooks:
      - id: go-fmt
        stages: [push, commit]
      - id: go-imports
        stages: [push, commit]
      - id: go-unit-tests
        stages: [push, commit]

  - repo: https://github.com/TekWizely/pre-commit-golang
    rev: v1.0.0-rc.1
    hooks:
      - id: go-revive
        stages: [push, commit]

  - repo: https://github.com/alessandrojcm/commitlint-pre-commit-hook
    rev: v9.13.0
    hooks:
      - id: commitlint # avaliable types are: build, ci, chore, docs, feat, fix, perf, refactor, revert, style, test
        stages: [push, commit]
        additional_dependencies: ["@commitlint/config-conventional"]
`), nil
}

// Name implements Component.
func (p preCommitComponent) Name() string {
	return ".pre-commit-config.yml"
}

// Path implements Component.
func (p preCommitComponent) Path() string {
	return "."
}

func NewPreCommitComponent() Component {
	return preCommitComponent{}
}
