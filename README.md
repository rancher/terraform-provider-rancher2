Terraform Provider for Rancher v2
==================================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) >= 0.11.x
- [Go](https://golang.org/doc/install) 1.12 to build the provider plugin
- [Docker](https://docs.docker.com/install/) 17.03.x to run acceptance tests

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-rancher2`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:terraform-providers/terraform-provider-rancher2
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-rancher2
$ make build
```

Using the provider
----------------------

If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory,  run `terraform init` to initialize it. Documentation about the provider specific configuration options can be found on the [provider's website](https://www.terraform.io/docs/providers/rancher2/index.html).

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.9+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in `$GOPATH/bin` .

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-rancher2
...
```

To just compile provider binary on repo path and test on terraform:

```sh
$ make bin
$ terraform init
$ terraform plan
$ terraform apply
```

Testing the Provider
---------------------------

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, a running rancher system, a rancher API key and a working k8s cluster imported are needed.

To run acceptance tests, you can simply run `make testacc`. `scripts/gotestacc.sh` will be run, deploying all needed requirements, running acceptance tests and cleanup.

```sh
$ make testacc
```

Due to [network limitation](https://docs.docker.com/docker-for-mac/networking/#known-limitations-use-cases-and-workarounds) on docker for osx and/or windows, try to add variable `TESTACC_EXPOSE_HOST_PORTS=true` to run acceptance test ont them.

```sh
$ TESTACC_EXPOSE_HOST_PORTS=true make testacc
```
