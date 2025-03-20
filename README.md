# cli

[![wakatime](https://wakatime.com/badge/user/965e81db-2a88-4564-b236-537c4a901130/project/4cfc2a67-bfe6-432b-a9b7-abf550e6be1c.svg)](https://wakatime.com/badge/user/965e81db-2a88-4564-b236-537c4a901130/project/4cfc2a67-bfe6-432b-a9b7-abf550e6be1c)
![Build Status](https://github.com/catalystgo/cli/actions/workflows/ci.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/catalystgo/cli)](https://goreportcard.com/report/github.com/catalystgo/cli)
[![codecov](https://codecov.io/gh/catalystgo/cli/graph/badge.svg?token=KN3G1NL58M)](https://codecov.io/gh/catalystgo/cli)

[![GitHub issues](https://img.shields.io/github/issues/catalystgo/cli.svg)](https://github.com/catalystgo/cli/issues)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/catalystgo/cli.svg)](https://github.com/catalystgo/cli/pulls)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

code generation CLI tool for catalystgo projects

## Installation 🏗

### Using Go 🐹

```bash
go install github.com/catalystgo/cli/cmd/catalystgo@latest
```

### Using Docker 🐳

```bash
docker pull catalystgo/cli:latest
```

## Usage 🚀

### Commands 📜

| Command     | Short | Description                            | 
|-------------|-------|----------------------------------------|
| `init`      | i     | Initialize the project files           |
| `implement` | impl  | Generate the gRPC code for the project |
| `version`   | ver   | Print the version of the tool          |
| `help`      |       | Print the help message                 |

## Example

### Prerequisites

* [Go](https://go.dev/doc/install)
* [Task](https://taskfile.dev/installation/)

### Steps

1) Initialize the project

    ```bash
    catalystgo init github.com/catalystgo/awesome-app
    ```

2) Create proto file

   ```bash
   mkdir -p api/user
   touch api/user/user.proto
   ```

3) Insert proto file content 

    ```protobuf
    syntax = "proto3";
    
    package user_pb;
    
    import "google/api/annotations.proto";
    import "google/api/http.proto";
    import "protoc-gen-openapiv2/options/annotations.proto";
    import "protoc-gen-openapiv2/options/openapiv2.proto";
    
    option go_package = "github.com/catalystgo/awesome-app/pkg/user";
    
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

4) Generate the gRPC code

    ```bash
    task generate
    ```

5) Modify the content of `internal/api/user/authenticate.go` to

   ```go
   package user
   
   import (
      "context"
      desc "github.com/catalystgo/awesome-app/pkg/user"
      "google.golang.org/grpc/codes"
      "google.golang.org/grpc/status"
   )

   func (i *Implementation) Authenticate(ctx context.Context, req *desc.AuthenticateRequest) (*desc.AuthenticateResponse, error) {
       if req.GetUsername() == "admin" && req.GetPassword() == "admin" {
           return &desc.AuthenticateResponse{Token: "admin-token"}, nil
       }
       return nil, status.Error(codes.Unauthenticated, "invalid credentials")
   }
   ```
   
6) Modify the content of `cmd/awesome-app/main.go` to be

     ```go
    package main
    
    import (
      "context"
    
      "github.com/catalystgo/catalystgo"
      "github.com/catalystgo/logger/logger"
      "github.com/catalystgo/awesome-app/internal/api/user"
    )
    
    func main() {
      ctx := context.Background()
    
      app, err := catalystgo.New()
      if err != nil {
        logger.Fatalf(ctx, "create app: %v", err)
      }
      
      srv := user.NewUserService()
    
      if err := app.Run(srv); err != nil {
        logger.Fatalf(ctx, "run app: %v", err)
      }
    }
    ```

7) Run the app

    ```bash
    task run
    ```     

8) Test the app
    ```bash
    # using http
    curl -s -X POST -d '{"username": "admin", "password": "admin"}' http://localhost:7000/user/authenticate | jq
   
    # using grpc
    grpcurl -plaintext -d '{"username": "admin", "password": "admin"}' localhost:8000 user_pb.UserService/Authenticate | jq
    ```

Congratz!!! You have built a golang app in seconds 🥳

## Contributing 🤝

![Alt](https://repobeats.axiom.co/api/embed/e373563a2553ccafa09a9dc1a6e18d31c761fc4b.svg "Repobeats analytics image")
