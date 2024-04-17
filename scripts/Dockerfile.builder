FROM golang:1.22.2-alpine3.19
RUN apk -U add bash git gcc musl-dev make docker-cli curl ca-certificates
WORKDIR /go/src/github.com/terraform-providers/terraform-provider-rancher2
