FROM golang:alpine AS builder
WORKDIR /go/src/{{.Module}}
COPY . .
RUN go mod download
RUN go build -o /go/bin/{{.AppName}} ./cmd/.

FROM alpine:3.19 AS production
COPY --from=builder /go/bin/{{.AppName}} /go/bin/{{.AppName}}
COPY /go/bin/{{.AppName}} /go/bin/{{.AppName}}
COPY ./.cataystgo/ /go/bin/.catalystgo/
ENTRYPOINT ["/go/bin/{{.AppName}}"]
