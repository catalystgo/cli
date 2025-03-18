# cli

[![Go Reference](https://pkg.go.dev/badge/github.com/catalystgo/cli.svg)](https://pkg.go.dev/github.com/catalystgo/cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/catalystgo/cli)](https://goreportcard.com/report/github.com/catalystgo/cli)
![License](https://img.shields.io/github/license/catalystgo/cli)

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/catalystgo/cli)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/catalystgo/cli)


## Quickstart


Create the directory for the proto file

```bash
mkdir -p api/example 
```

Example proto file

```protobuf
syntax = "proto3";

package example_pb;

import "google/api/annotations.proto";
import "google/api/http.proto";
import "protoc-gen-openapiv2/options/annotations.proto";
import "protoc-gen-openapiv2/options/openapiv2.proto";

option go_package = "github./pkg/example"; // change the go_package

service UserService {
  rpc Authenticate(AuthenticateRequest) returns (AuthenticateResponse) {
    option (google.api.http) = {
      post: "/user/authenticate"
      body: "*"
    };
    option (grpc.gateway.protoc_gen_openapiv2.options.openapiv2_operation) = {
      description: "Authenticate a user"
    };
  }
}

message AuthenticateRequest {
  string username = 1;
  string password = 2;
}

message AuthenticateResponse {
  string token = 1;
}
```

## Future milestones 💎

### Priority (high)

- [ ] Write tests for service.Service methods in generated code
- [ ] Add support to many services in a single proto file/directory

### Priority (medium)

- [ ] Add github actions to run tests, lint, build/push docker images
