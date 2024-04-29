package component

import (
	"strings"
)

var _ Component = taskfileComponent{}

type (
	taskfileComponent struct {
		AppName string
	}
)

// Content implements Component.
func (t taskfileComponent) Content() ([]byte, error) {
	b := `version: '3'

vars:
  # GENERAL
  GOBIN: ./bin

  # DB
  DB_DSN: "user=postgres dbname={{.AppName}} password=postgres sslmode=disable"
  MIGRATION_DIR: "./migrations"

tasks:
  run:
    cmds:
      - |
        export GOLANG_PROTOBUF_REGISTRATION_CONFLICT=warn
        go run ./cmd/.

  mock:
    cmds:
      - echo "mocking!"

  generate:
    cmds:
      - "{{.GOBIN}}/buf generate"
      - go mod tidy

  format:
    cmds:
      - task: go_files
        vars: { COMMAND: "gofmt -w  {} +"}
      - task: go_files
        vars: { COMMAND: "{{.GOBIN}}/goimports -w  {} +"}

  test:
    cmds:
      - go test -v -race -cover ./...

  lint:
    cmds:
      - "{{.GOBIN}}/revive
        -config revive.toml
        -formatter friendly
        -exclude ./mock
        ./..."

  deps:
    cmds:
      - |
        export GOBIN="$(pwd)/bin"

        # GRPC

        go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.16.0
        go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.16.0
        go install github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto@v0.6.0
        go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
        go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2.0
        go install github.com/bufbuild/buf/cmd/buf@v1.31.0

        # OTHERS

        go install github.com/pressly/goose/v3/cmd/goose@v3.20.0
        go install golang.org/x/tools/cmd/goimports@v0.19.0
        go install github.com/mgechev/revive@v1.3.7

        # Uncomment to use pre-commit hooks (pip must be installed)
        # pip install pre-commit

## DOCKER

  docker-up:
    cmds:
      - docker compose up -f .deployment/docker-compose.yml

  docker-down:
    cmds:
      - docker compose down -f .deployment/docker-compose.yml

# MIGRATIONS

  goose-add:
    cmds:
    - "{{.GOBIN}}/goose -dir {{.MIGRATION_DIR}} create fix_me sql"

  goose-status:
    cmds:
      - "{{.GOBIN}}/goose -dir {{.MIGRATION_DIR}} postgres \"{{.DB_DSN}}\" status"

  goose-up:
    cmds:
    - "{{.GOBIN}}/goose -dir {{.MIGRATION_DIR}} postgres \"{{.DB_DSN}}\" up"

  goose-down:
    cmds:
      - "{{.GOBIN}}/goose -dir {{.MIGRATION_DIR}} postgres \"{{.DB_DSN}}\" down"

## INTERNAL COMMANDS

  go_files:
    desc: "Return all .go files and run .COMMAND on them"
    internal: true
    cmds:
     - find .
        -name "*.go"
        -not -path ./mock
        -exec {{.COMMAND}};
`

	// can't use gotemplate since the file has template syntax inside
	return []byte(strings.ReplaceAll(b, "{{.AppName}}", t.AppName)), nil
}

// Name implements Component.
func (t taskfileComponent) Name() string {
	return "Taskfile.yml"
}

// Path implements Component.
func (t taskfileComponent) Path() string {
	return "."
}

func NewTaskfileComponent(module string) Component {
	return taskfileComponent{
		AppName: getAppNameFromModule(module),
	}
}
