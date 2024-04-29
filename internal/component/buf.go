package component

var _ Component = bufComponent{}

type (
	bufComponent struct{}
)

// Content implements Component.
func (f bufComponent) Content() ([]byte, error) {
	return []byte(`version: v1

lint:
    use:
        - DEFAULT
    rpc_allow_same_request_response: false
    rpc_allow_google_protobuf_empty_requests: true
    rpc_allow_google_protobuf_empty_responses: true
    allow_comment_ignores: true
`), nil
}

// Name implements Component.
func (f bufComponent) Name() string {
	return "buf.yml"
}

// Path implements Component.
func (f bufComponent) Path() string {
	return "."
}

func NewBufComponent() Component {
	return bufComponent{}
}
