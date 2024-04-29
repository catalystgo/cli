FROM golang:alpine AS builder
WORKDIR /go/src/github.com/catalystgo/cli
COPY . .
RUN go mod download
RUN go build -o bin/catalystgo cmd/catalystgo/main.go

FROM alpine:3.19
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/catalystgo/cli/bin/catalystgo .
ENTRYPOINT ["./catalystgo"]
