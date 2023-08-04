FROM golang:1.19

USER root

RUN mkdir -p /.cache && chmod -R 777 /.cache

RUN mkdir -p $GOPATH/pkg/mod && chmod -R 777 $GOPATH/pkg/mod

RUN chown -R root:root $GOPATH/pkg/mod && chmod -R g+rwx $GOPATH/pkg/mod

ENV WORKSPACE ${GOPATH}/src/github.com/rancher/terraform-provider-rancher2/tfp-automation

WORKDIR $WORKSPACE

COPY [".", "$WORKSPACE"]

RUN go mod download && \
    go install gotest.tools/gotestsum@latest

COPY . .

ARG TERRAFORM_VERSION=0.13.7
RUN wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && apt-get update && apt-get install unzip &&  unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip && rm terraform_${TERRAFORM_VERSION}_linux_amd64.zip && chmod u+x terraform && mv terraform /usr/bin/terraform

ARG CONFIG_FILE
COPY ${CONFIG_FILE} /config.yml

SHELL ["/bin/bash", "-c"] 