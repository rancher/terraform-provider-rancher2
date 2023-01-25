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
- [Docker](https://docs.docker.com/install/) 20.10.x to run acceptance tests

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

Branching the Provider
---------------------------

The provider is branched into two release lines that have major version alignment with Rancher 2.6 and 2.7. The `release/v2` branch with 2.0.0+ is aligned with Rancher 2.6 and `master` with 3.0.0+ is aligned with Rancher 2.7. Terraform provider fixes and new features will be available on `master` but only bug fixes will be backported to `release/v2` as necessary.

Aligning major provider releases with minor Rancher releases means,

* We can follow semver
* We can cut patch/minor versions on an as-needed basis to fix bugs or add new resources 
* We have 'out of band' flexibility and are only tied to releasing a new version of the provider when we get a new 2.x Rancher minor version.

See the [compatibility matrix](docs/compatibility-matrix.md) for details.

If you are using Terraform to provision clusters on instances of Rancher 2.6 and 2.7, you must have a separate configuration in a separate dir for each provider. Otherwise, Terraform will overwrite the `.tfstate` file every time you switch versions.

Releasing the Provider
---------------------------

As of Terraform 2.0.0 and 3.0.0, the provider is tied to Rancher minor releases but can be released 'out of band' within that minor version. For example, 3.0.0 will be released a few days after Rancher 2.7.x and fixes and features in the 3.0.0 release will be supported for clusters provisioned via Terraform on Rancher 2.7.x. A critical bug fix can be released 'out of band' as 3.0.1 and backported to `release/v2` as 2.0.1. A new feature can also be released 'out of band' as 3.1.0 but not backported. Terraform 4.0.0 must be released with Rancher 2.8.

The [RKE provider](https://github.com/rancher/terraform-provider-rke) should be released after every RKE or KDM release. For example, if upstream RKE 1.3.15 was released, bump the RKE version to 1.3.15 and release the provider.

To release the provider

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
* Copy the updated notes back to the draft release and save (DO NOT release with just the generated notes. Those are just a template to help you)
* Create an [EIO issue](https://github.com/rancherlabs/eio) for Hashicorp to publish the release