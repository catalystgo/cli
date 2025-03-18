FROM golang:alpine AS builder

ARG VERSION=${VERSION:-dev}
ARG VERSION_PATH="github.com/catalystgo/cli/internal/build.Version"

WORKDIR /go/src/github.com/catalystgo/cli
COPY . .
RUN go mod download
RUN go build -o bin/catalystgo -ldflags "-X '${VERSION_PATH}=${VERSION}'" cmd/catalystgo/main.go

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/catalystgo/cli/bin/catalystgo .
ENTRYPOINT ["./catalystgo"]
