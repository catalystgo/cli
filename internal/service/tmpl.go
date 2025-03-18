package service

import (
	_ "embed"

	"github.com/catalystgo/cli/internal/domain"
)

var (
	//go:embed template/implement/service.go.tmpl
	serviceContent []byte

	//go:embed template/implement/method_unary.go.tmpl
	methodUnaryContent []byte

	//go:embed template/implement/method_server_stream.go.tmpl
	methodServerStreamContent []byte

	//go:embed template/implement/method_bidirectional_or_client_stream.go.tmpl
	methodBidirectionalOrClientStreamContent []byte
)

var (
	serviceTemplate                           = domain.MustParseTemplate("service.go", serviceContent)
	methodUnaryTemplate                       = domain.MustParseTemplate("method_unary.go", methodUnaryContent)
	methodServerStreamTemplate                = domain.MustParseTemplate("method_server_stream.go", methodServerStreamContent)
	methodBidirectionalOrClientStreamTemplate = domain.MustParseTemplate("method_bidirectional_or_client_stream.go", methodBidirectionalOrClientStreamContent)
)

type (
	serviceTemplateData struct {
		PackageName       string
		ProtoImportPath   string
		ServiceStructName string
		ServiceName       string
	}

	methodTemplateData struct {
		PackageName       string
		ProtoImportPath   string
		ServiceStructName string
		MethodName        string
		Request           string
		Response          string
		Stream            string
	}
)
