Terraform Provider for Rancher v2
==================================

[![Go Report Card](https://goreportcard.com/badge/github.com/terraform-providers/terraform-provider-rancher2)](https://goreportcard.com/report/github.com/terraform-providers/terraform-provider-rancher2)

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/public/img/logo-hashicorp.svg" width="600px">

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) >= 0.11.x
- [Go](https://golang.org/doc/install) 1.13 to build the provider plugin
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

Using the Provider
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

In order to run the full suite of Acceptance tests, a running rancher system, a rancher API key and a working k8s cluster imported are needed. Also, acceptance tests are covering Rancher server upgrade, v2.3.6 and v2.4.2

To run acceptance tests, you can simply run `make testacc`. `scripts/gotestacc.sh` will be run, deploying all needed requirements, running acceptance tests and cleanup.

```sh
$ make testacc
```

Due to [network limitation](https://docs.docker.com/docker-for-mac/networking/#known-limitations-use-cases-and-workarounds) on docker for osx and/or windows, there is a way to run dockerized acceptance test.

```sh
$ EXPOSE_HOST_PORTS=true make docker-testacc
```

Releasing the Provider
---------------------------

* Create a draft of the [release](https://github.com/rancher/terraform-provider-rancher2/releases) and select create new tag for the version you are releasing
* Create release notes by clicking `Generate release notes`
* Copy the release notes to the CHANGELOG and update to the following format

```
# <tag version> (Month Day, Year)
FEATURES:
ENHANCEMENTS:
BUG FIXES:
```

* Create a PR to update CHANGELOG
* Copy the updated notes back to the draft release (DO NOT release with just the generated notes. Those are just a template to help you)
* Publish the release (using the CL: git tag v0.0.0; git push -u origin v0.0.0)
* Create an [EIO issue](https://github.com/rancherlabs/eio) for Hashicorp to publish the release