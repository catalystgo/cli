package service

import (
	"fmt"
	"go/ast"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/catalystgo/cli/internal/component"
	"github.com/catalystgo/cli/internal/domain"
	"github.com/catalystgo/helpers"
	log "github.com/catalystgo/logger/cli"
)

type Service struct{}

func New() *Service { return &Service{} }

// Init implements CommandService.
func (s *Service) Init(components []component.Component, override bool) {
	for _, c := range components {
		name := filepath.Join(c.Path(), c.Name())

		log.Debugf("file = %s | exist = %t | override = %t", name, fileExists(name), override)

		// Override files content only if override flag is passed
		if fileExists(name) {
			if override {
				log.Warnf("override file (%s)", name)
			} else {
				log.Warnf("skip file (%s) => already exist", name)
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
		err = helpers.SaveFile(name, b, &helpers.SaveFileOpt{Override: override})
		if err != nil {
			continue
		}
	}
}

func (s *Service) Implement(input string, output string) error {
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
		err = helpers.SaveFile(file, data, &helpers.SaveFileOpt{Override: false})
		if err != nil {
			continue
		}
	}

	return nil
}

func (s *Service) implement(nodes map[string]*ast.File, input string, output string) (map[string][]byte, error) {
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
		serviceContent, err := s.buildContructor(module, file, packageName, serviceName)
		if err != nil {
			log.Errorf("build service (%s) => %v", serviceName, err)
			continue
		}

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
			methodContent, err := s.buildStructMethod(module, file, packageName, method)
			if err != nil {
				log.Errorf("build method (%s) => %v", method.Name.Name, err)
				continue
			}

			files[methodFile] = methodContent
		}
	}

	return files, nil
}

func (s *Service) buildContructor(module string, file string, packageName string, serviceName string) ([]byte, error) {
	// Get the service name without the unimplemented prefix and suffix
	serviceName = strings.TrimSuffix(serviceName, unimplementedStructSuffix)
	serviceName = strings.TrimPrefix(serviceName, unimplementedStructPrefix)

	data := serviceTemplateData{
		PackageName:       packageName,
		ProtoImportPath:   path.Dir(path.Join(module, file)),
		ServiceStructName: implementationStructName,
		ServiceName:       serviceName,
	}

	return domain.ExecuteTemplate(serviceTemplate, data)
}

func (s *Service) buildStructMethod(module string, file string, packageName string, method *ast.FuncDecl) ([]byte, error) {
	var (
		tmpl *template.Template
		data = methodTemplateData{
			PackageName:       packageName,
			ProtoImportPath:   path.Dir(path.Join(module, file)),
			ServiceStructName: implementationStructName,
			MethodName:        method.Name.Name,
		}
	)

	// Get the method type
	mt := getMethodType(method)
	switch mt {
	case unaryMethod:
		tmpl = methodUnaryTemplate
		data.Request = method.Type.Params.List[1].Type.(*ast.StarExpr).X.(*ast.Ident).Name   // request
		data.Response = method.Type.Results.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name // response
	case serverStreamMethod:
		tmpl = methodServerStreamTemplate
		data.Request = method.Type.Params.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).Name // request
		data.Stream = method.Type.Params.List[1].Type.(*ast.Ident).Name                    // stream
	case bidirectionalOrClientStreamMethod:
		tmpl = methodBidirectionalOrClientStreamTemplate
		data.Stream = method.Type.Params.List[0].Type.(*ast.Ident).Name // stream
	default:
		log.Warnf("unsupported method type (%s) in file (%s)", mt, file)
		return nil, nil
	}

	return domain.ExecuteTemplate(tmpl, data)
}
