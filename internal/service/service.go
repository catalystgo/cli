package service

import (
	"fmt"
	"go/ast"
	"path"
	"path/filepath"
	"strings"

	"github.com/catalystgo/cli/internal/component"
	"github.com/catalystgo/cli/internal/log"
)

var _ Service = &service{}

type (
	Service interface {
		Init(components []component.Component, override bool)
		Implement(input string, output string) error
	}

	WriterSrv interface {
		WriteFile(file string, data []byte, opts ...WriteOption) error
	}
)

type service struct {
	w WriterSrv
}

func New() Service {
	return &service{
		w: NewWriter(),
	}
}

// Init implements CommandService.
func (s *service) Init(components []component.Component, override bool) {
	for _, c := range components {
		name := filepath.Join(c.Path(), c.Name())

		log.Debugf("file = %s | exist = %t | override = %t", name, fileExists(name), override)

		// Override files content only if override flag is passed
		if fileExists(name) {
			if override {
				log.Warnf("overriding file (%s)", name)
			} else {
				log.Warnf("skipping file (%s) => already exist", name)
				continue
			}
		}

		// Get content
		b, err := c.Content()
		if err != nil {
			log.Errorf("read content file (%s) => %v", name, err)
			continue
		}

		// Write content
		err = s.w.WriteFile(name, b, WithOverride(override))
		if err != nil {
			continue
		}
	}
}

func (s *service) Implement(input string, output string) error {
	files, err := getGrpcGoFiles(input)
	if err != nil {
		return err
	}

	if len(files) == 0 {
		return fmt.Errorf("no grpc go files found in directory (%s)", input)
	}

	nodes, err := parseGoFiles(files)
	if err != nil {
		return err
	}

	content, err := s.implement(nodes, input, output)
	if err != nil {
		return err
	}

	for file, data := range content {
		err = s.w.WriteFile(file, data)
		if err != nil {
			continue
		}
	}

	return nil
}

func (s *service) implement(nodes map[string]*ast.File, input string, output string) (map[string][]byte, error) {
	module, err := getCurrentModule()
	if err != nil {
		return nil, err
	}

	// Clean the input and output paths for further processing
	input = path.Clean(input)
	output = path.Clean(output)

	files := make(map[string][]byte, len(nodes))
	for file, node := range nodes {
		file = path.Clean(file)

		packageName := node.Name.Name
		serviceName := getStructName(node, unimplementedStructPrefix, unimplementedStructSuffix)
		if serviceName == "" {
			log.Warnf("no unimplemented struct found in file (%s)", file)
			continue
		}

		log.Infof("implementing service (%s) in file (%s)", serviceName, file)

		// Get the directory of the file  in output directory relative to the input directory
		// Example:
		//  input   = /path/to/input
		//  output  = /path/to/output
		//  file    = /path/to/input/level_one/level_two/service_grpc.pb.go
		//  fileDir = /path/to/output/level_one/level_two
		fileDir := path.Join(output, path.Dir(strings.TrimPrefix(file, input)))

		log.Debugf("implementing service (%s) in directory (%s)", serviceName, fileDir)

		// Build the service.go file
		serviceFile := path.Join(fileDir, "service.go")
		serviceContent := s.buildContructor(module, file, packageName, serviceName)
		files[serviceFile] = serviceContent

		// Build the method files
		serviceMethods := getStructMethods(node, serviceName)
		for _, method := range serviceMethods {
			// Skip unexported methods
			if !method.Name.IsExported() {
				continue
			}

			// Build the method file name
			// Example:
			//  fileDir     = /path/to/output/level_one/level_two
			// 	methodName  = MethodName
			// 	methodFile  = /path/to/output/level_one/level_two/method_name.go
			methodFile := path.Join(fileDir, fmt.Sprintf("%s.go", toSnakeCase(method.Name.Name)))
			methodContent := s.buildStructMethod(module, file, packageName, method)
			files[methodFile] = methodContent
		}
	}

	return files, nil
}

func (s *service) buildContructor(module string, file string, packageName string, serviceName string) []byte {
	b := strings.Builder{}

	// Get the service name without the unimplemented prefix and suffix
	serviceName = strings.TrimSuffix(serviceName, unimplementedStructSuffix)
	serviceName = strings.TrimPrefix(serviceName, unimplementedStructPrefix)

	// Write package and imports
	b.WriteString(fmt.Sprintf("package %s", packageName))
	b.WriteString("\n\n")
	b.WriteString("import (")
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("\t%q", "github.com/catalystgo/catalystgo"))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("\tdesc %q", path.Dir(path.Join(module, file))))
	b.WriteString("\n")
	b.WriteString(")")
	b.WriteString("\n\n")

	// Write struct object
	b.WriteString(fmt.Sprintf("type %s struct {", implementationStructName))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("\tdesc.Unimplemented%sServer", serviceName))
	b.WriteString("\n")
	b.WriteString("}")
	b.WriteString("\n\n")

	// Write constructor method
	b.WriteString(fmt.Sprintf("func New%s() *%s {", serviceName, implementationStructName))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("\treturn &%s{}", implementationStructName))
	b.WriteString("\n")
	b.WriteString("}")
	b.WriteString("\n\n")

	// Write GetDescription method
	b.WriteString(fmt.Sprintf("func (i *%s) GetDescription() catalystgo.ServiceDesc {", implementationStructName))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("\treturn desc.New%sServiceDesc(i)", serviceName))
	b.WriteString("\n")
	b.WriteString("}")
	b.WriteString("\n\n")

	return []byte(b.String())
}

func (s *service) buildStructMethod(module string, file string, packageName string, method *ast.FuncDecl) []byte {
	b := &strings.Builder{}

	// Get the method type
	methodType := getMethodType(method)
	if methodType == unknownMethod {
		log.Warnf("unsupported method type for method %q", method.Name.Name)
		return nil
	}

	// Write package and imports
	b.WriteString(fmt.Sprintf("package %s", packageName))
	b.WriteString("\n\n")
	b.WriteString("import (")
	b.WriteString("\n")

	// Only unary methods require context import
	if methodType == unaryMethod {
		b.WriteString(fmt.Sprintf("\t%q", "context"))
		b.WriteString("\n\n")
	}

	b.WriteString(fmt.Sprintf("\tdesc %q", path.Dir(path.Join(module, file))))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("\t%q", "google.golang.org/grpc/codes"))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("\t%q", "google.golang.org/grpc/status"))
	b.WriteString("\n")
	b.WriteString(")\n\n")

	// Write method signature
	s.buildStructMethodSignature(b, method, methodType)

	return []byte(b.String())
}

// buildStructMethodSignature writes the method signature to the string builder
// Output example:
// Unary RPC									: func (i *Implementation) MethodName(ctx context.Context, req *Request) (*Response, error) {
// Server Streaming RPC				: func (i *Implementation) MethodName(req *Request, stream desc.ServiceName_MethodNameServer) error {
// Client Streaming RPC				: func (i *Implementation) MethodName(stream desc.ServiceName_MethodNameServer) error {
// Bidirectional Streaming RPC: func (i *Implementation) MethodName(stream desc.ServiceName_MethodNameServer) error {
func (s *service) buildStructMethodSignature(b *strings.Builder, method *ast.FuncDecl, methodType methodType) {
	b.WriteString(fmt.Sprintf("func (i *%s) %s(", implementationStructName, method.Name.Name))

	switch methodType {
	// Unary RPC
	case unaryMethod:
		b.WriteString(fmt.Sprintf("ctx context.Context, req *desc.%s) (*desc.%s, error) {",
			method.Type.Params.List[1].Type.(*ast.StarExpr).X.(*ast.Ident).Name,  // request
			method.Type.Results.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name, // response
		))
		// Write the body
		b.WriteString(fmt.Sprintf("\n\treturn nil, status.Error(codes.Unimplemented, `method %q not implemented`)", method.Name.Name))

	// Server Streaming RPC
	case serverStreamMethod:
		b.WriteString(fmt.Sprintf("req *desc.%s, stream desc.%s) error {",
			method.Type.Params.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name, // request
			method.Type.Params.List[1].Type.(*ast.Ident).Name,                   // stream
		))
		// Write the body
		b.WriteString(fmt.Sprintf("\n\treturn status.Error(codes.Unimplemented, `method %q not implemented`)", method.Name.Name))

	// Bidirectional or Client Streaming RPC (same signature)
	case bidectionalOrClientStreamMethod:
		b.WriteString(fmt.Sprintf("stream desc.%s) error {",
			method.Type.Params.List[0].Type.(*ast.Ident).Name, // stream
		))
		// Write the body
		b.WriteString(fmt.Sprintf("\n\treturn status.Error(codes.Unimplemented, `method %q not implemented`)", method.Name.Name))

	default:
		log.Warnf("unsupported method signature ()=> %s", method.Name.Name)
	}

	b.WriteString("\n}\n")
}
