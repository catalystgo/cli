package component

import (
	"bytes"
	"text/template"
)

var _ Component = dockerComponent{}

var dockerComponentTemplate *template.Template

func init() {
	var err error

	dockerComponentTemplate, err = template.New("docker-compose").Parse(
		`FROM golang:1.22 AS development
WORKDIR /go/src/{{.Module}}
COPY . .
RUN go mod download
RUN go install github.com/cespare/reflex@latest
CMD reflex -sr '\.go$' go run ./cmd/.

FROM golang:alpine AS builder
WORKDIR /go/src/{{.Module}}
COPY . .
RUN go build -o /go/bin/{{.AppName}} ./cmd/.

FROM alpine:3.19 AS production
COPY --from=builder /go/bin/{{.AppName}} /go/bin/{{.AppName}}
# COPY ./{{.AppName}}/migrations /migrations
COPY ./{{.AppName}}/config.yml /{{.AppName}}/config.yml
ENTRYPOINT ["/go/bin/{{.AppName}}"]
`)

	if err != nil {
		panic(err)
	}
}

type (
	dockerComponent struct {
		AppName string
		Module  string
	}
)

// Content implements Component.
func (d dockerComponent) Content() ([]byte, error) {
	var (
		b   bytes.Buffer
		err error
	)

	err = dockerComponentTemplate.Execute(&b, d)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// Name implements Component.
func (d dockerComponent) Name() string {
	return "Dockerfile"
}

// Path implements Component.
func (d dockerComponent) Path() string {
	return ".deployment"
}

func NewDockerComponent(module string) Component {
	return dockerComponent{
		Module:  module,
		AppName: getAppNameFromModule(module)}
}
