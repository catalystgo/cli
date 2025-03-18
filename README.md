# cli

[![wakatime](https://wakatime.com/badge/user/965e81db-2a88-4564-b236-537c4a901130/project/4cfc2a67-bfe6-432b-a9b7-abf550e6be1c.svg)](https://wakatime.com/badge/user/965e81db-2a88-4564-b236-537c4a901130/project/4cfc2a67-bfe6-432b-a9b7-abf550e6be1c)
![Build Status](https://github.com/catalystgo/cli/actions/workflows/ci.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/catalystgo/cli)](https://goreportcard.com/report/github.com/catalystgo/cli)
[![codecov](https://codecov.io/gh/catalystgo/cli/graph/badge.svg?token=KN3G1NL58M)](https://codecov.io/gh/catalystgo/cli)

[![GitHub issues](https://img.shields.io/github/issues/catalystgo/cli.svg)](https://github.com/catalystgo/cli/issues)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/catalystgo/cli.svg)](https://github.com/catalystgo/cli/pulls)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

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

## Features 🎯

- [ ] Write tests for service.Service methods in generated code
