package component

var _ Component = bufworkComponent{}

type (
	bufworkComponent struct{}
)

// Content implements Component.
func (f bufworkComponent) Content() ([]byte, error) {
	return []byte(
		`version: v1

directories:
  - "api"
`), nil
}

// Name implements Component.
func (f bufworkComponent) Name() string {
	return "buf.work.yaml"
}

// Path implements Component.
func (f bufworkComponent) Path() string {
	return "."
}

func NewBufWorkComponent() Component {
	return bufworkComponent{}
}
