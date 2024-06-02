FROM golang:1.22 AS development
WORKDIR /go/src/{{.Module}}
COPY . .
RUN go mod download
RUN go install github.com/cespare/reflex@latest
CMD reflex -sr '\.go$' go run ./cmd/.

FROM golang:alpine AS builder
WORKDIR /go/src/{{.Module}}
COPY . .
RUN go build -o /go/bin/{{.AppName}} ./cmd/.

FROM alpine:3.19 AS production
COPY --from=builder /go/bin/{{.AppName}} /go/bin/{{.AppName}}
# COPY ./{{.AppName}}/migrations /migrations
COPY ./{{.AppName}}/config.yml /{{.AppName}}/config.yml
ENTRYPOINT ["/go/bin/{{.AppName}}"]
