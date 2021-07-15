FROM golang:1.16.5-alpine3.13
RUN apk -U add bash git gcc musl-dev make docker-cli curl ca-certificates
WORKDIR /go/src/github.com/terraform-providers/terraform-provider-rancher2