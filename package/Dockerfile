FROM alpine:3.8
RUN apk add --no-cache ca-certificates openssh-client
RUN wget https://releases.hashicorp.com/terraform/0.11.10/terraform_0.11.10_linux_amd64.zip -O - | unzip -d /usr/bin/ - && \
    chmod 755 /usr/bin/terraform
COPY terraform-provider-rancher2 /usr/bin/
WORKDIR /mnt
ENTRYPOINT ["terraform"]
CMD  ["-help"]