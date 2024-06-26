package main

import (
	"path"

	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	protogen.Options{}.Run(func(gen *protogen.Plugin) error {
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generateFile(gen, f)
		}
		return nil
	})
}

func generateFile(gen *protogen.Plugin, file *protogen.File) {
	if len(file.Services) == 0 || file.GoPackageName == "" {
		return // skip files without services or no package
	}

	filename := file.GeneratedFilenamePrefix + "_catalystgo.pb.go"
	g := gen.NewGeneratedFile(filename, file.GoImportPath)
	g.P("// Code generated by protoc-gen-catalystgo. DO NOT EDIT.")
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()

	// add imports
	g.P("import (")
	g.P(" _ \"embed\"")
	g.P("	context \"context\"")
	g.P()
	g.P("	go_grpc_middleware \"github.com/grpc-ecosystem/go-grpc-middleware\"")
	g.P("	\"github.com/grpc-ecosystem/grpc-gateway/v2/runtime\"")
	g.P("	\"google.golang.org/grpc\"")
	g.P(")")
	g.P()

	g.P("//go:embed " + path.Base(file.GeneratedFilenamePrefix) + ".swagger.json")
	g.P("var swaggerJSON []byte")
	g.P()

	// loop over all services and generate the service descriptor definition for each one
	for _, s := range file.Services {
		g.P("type ", s.GoName+"ServiceDesc", " struct {")
		g.P("	svc ", s.GoName, "Server")
		g.P("	i   grpc.UnaryServerInterceptor")
		g.P("}")
		g.P()

		g.P("func New", s.GoName, "ServiceDesc(svc ", s.GoName, "Server) *", s.GoName, "ServiceDesc {")
		g.P("	return &", s.GoName, "ServiceDesc{")
		g.P("		svc: svc,")
		g.P("	}")
		g.P("}")
		g.P()

		descSignature := "func (d *" + s.GoName + "ServiceDesc) "

		// Add GRPC service registration
		g.P(descSignature + "RegisterGRPC(s *grpc.Server) {")
		g.P("	Register", s.GoName, "Server(s, d.svc)")
		g.P("}")
		g.P()

		// Add HTTP gateway registration
		g.P(descSignature + "RegisterHTTP(ctx context.Context, mux *runtime.ServeMux) error {")
		g.P("	if d.i == nil {")
		g.P("		return Register", s.GoName, "HandlerServer(ctx, mux, d.svc)")
		g.P("	}")
		g.P()
		g.P("	return Register", s.GoName, "HandlerServer(ctx, mux, &proxy", s.GoName, "Server{")
		g.P("		", s.GoName, "Server: d.svc,")
		g.P("		interceptor: d.i,")
		g.P("	})")
		g.P("}")
		g.P()

		// Add swagger methods
		g.P(descSignature + "SwaggerJSON()[]byte {")
		g.P("	return swaggerJSON")
		g.P("}")
		g.P()

		// Add HTTP interceptors
		g.P("//WithHTTPUnaryInterceptor adds GRPC Server interceptors for the HTTP unary endpoints. Call again to add more.")
		g.P(descSignature + "WithHTTPUnaryInterceptor(i grpc.UnaryServerInterceptor) {")
		g.P("	if d.i == nil {")
		g.P("		d.i = i")
		g.P("	} else {")
		g.P("		d.i = go_grpc_middleware.ChainUnaryServer(d.i, i)")
		g.P("	}")
		g.P("}")

		g.P("type proxy", s.GoName, "Server struct {")
		g.P("	", s.GoName, "Server")
		g.P("	interceptor grpc.UnaryServerInterceptor")
		g.P("}")
		g.P()

		for _, m := range s.Methods {
			// TODO: research if streaming is supported in grpc-gateway (probably not)
			if m.Desc.IsStreamingClient() || m.Desc.IsStreamingServer() {
				continue
			}

			g.P("func (p *proxy", s.GoName, "Server) ", m.GoName, "(ctx context.Context, req *", m.Input.GoIdent, ") (*", m.Output.GoIdent, ", error) {")
			g.P("	info := &grpc.UnaryServerInfo{")
			g.P("		Server: p.", s.GoName, "Server,")
			g.P("		FullMethod: \"/", file.GoPackageName, ".", s.GoName, "/", m.GoName, "\",")
			g.P("	}")
			g.P()
			g.P("	handler := func(ctx context.Context, req interface{}) (interface{}, error) {")
			g.P("		return p.", s.GoName, "Server.", m.GoName, "(ctx, req.(*", m.Input.GoIdent, "))")
			g.P("	}")
			g.P()
			g.P("	resp, err := p.interceptor(ctx, req, info, handler)")
			g.P("	if err != nil || resp == nil {")
			g.P("		return nil, err")
			g.P("	}")
			g.P()
			g.P("	return resp.(*", m.Output.GoIdent, "), nil")
			g.P("}")
			g.P()
		}
	}
}
