package service

import "go/ast"

const (
	gprcFileSuffix = "_grpc.pb.go"
)

const (
	unimplementedStructPrefix = "Unimplemented"
	unimplementedStructSuffix = "Server"
	implementationStructName  = "Implementation"
)

type methodType string

const (
	unknownMethod                   methodType = "unknown"
	unaryMethod                     methodType = "unary"
	serverStreamMethod              methodType = "serverStream"
	bidectionalOrClientStreamMethod methodType = "bidectionalOrClientStream"
)

func getMethodType(method *ast.FuncDecl) methodType {
	var (
		paramsCount  = method.Type.Params.NumFields()
		resultsCount = method.Type.Results.NumFields()
	)

	switch {
	// Example: func (i *Implementation) MethodName(ctx context.Context, req *Request) (*Response, error)
	case paramsCount == 2 && resultsCount == 2:
		return unaryMethod

	// Example: func (i *Implementation) MethodName(req *Request, stream desc.ServiceName_MethodNameServer) error
	case paramsCount == 2 && resultsCount == 1:
		return serverStreamMethod

	// Example: func (i *Implementation) MethodName(stream desc.ServiceName_MethodNameServer) error
	case paramsCount == 1 && resultsCount == 1:
		return bidectionalOrClientStreamMethod
	default:
		return unknownMethod
	}
}
