Terraform Provider for Rancher v2
==================================

[![Go Report Card](https://goreportcard.com/badge/github.com/terraform-providers/terraform-provider-rancher2)](https://goreportcard.com/report/github.com/terraform-providers/terraform-provider-rancher2)

Requirements
------------

- [Terraform](https://www.terraform.io/downloads.html) >= 1.5.0
- [Go](https://golang.org/doc/install) 1.23 to build the provider plugin

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

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.20+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

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

See [development process](docs/development-process.md) for more details.

Testing the Provider
---------------------------

In order to test the provider, simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, a running rancher system, a rancher API key and a working k8s cluster imported are needed. Acceptance tests cover a Rancher server upgrade, v2.3.6 and v2.4.2.

To run the Acceptance tests, simply run `make testacc`. `scripts/gotestacc.sh` will be run, deploying all needed requirements, running tests and cleanup.

```sh
$ make testacc
```

There is no way to run dockerized acceptance tests.

To run the structure tests, run

```sh
$ go clean --testcache && go test -v ./rancher2
```

See [test process](docs/test-process.md) for details on release testing (_Terraform Maintainers Only_).

Branching the Provider
---------------------------

This provider is branched in correlation with minor versions of Rancher:
* the `release/v5` branch with 5.0.0+ is aligned with Rancher 2.9
* the `release/v6` branch with 6.0.0+ is aligned with Rancher 2.10 
* the `master` branch with 7.0.0+ is aligned with Rancher 2.11

The lifecycle of each major provider version is aligned with the lifecycle of each Rancher minor version.
For example, provider versions 4.x are aligned with Rancher 2.8.x will only be actively maintained until the EOM for Rancher 2.8.x and supported until EOL for Rancher 2.8.x.

See the [Rancher support matrix](https://www.suse.com/lifecycle/#rancher) for details.

Aligning major provider releases with minor Rancher releases means:

* We can follow semver
* We can cut patch/minor versions on an as-needed basis to fix bugs or add new resources
* We have 'out of band' flexibility and are only tied to releasing a new version of the provider when we get a new 2.x Rancher minor version.

See the [compatibility matrix](docs/compatibility-matrix.md) for details.

If you are using Terraform to provision clusters on instances of Rancher 2.8 and 2.9,
  you must have a separate configuration in a separate dir for each provider.
Otherwise, Terraform will overwrite the `.tfstate` file every time you switch versions.

Releasing the Provider
---------------------------

As of v2.0.0 and v3.0.0, the provider is tied to Rancher minor releases but can be released 'out of band' within that minor version.
For example, v7.0.0 will be released 1-2 weeks after Rancher 2.11.x and fixes and features in the 7.0.0 release will be supported for clusters provisioned via Terraform on Rancher 2.11.x.
A critical bug fix can be released immediately as a patch release. For instance, as v7.0.1

To release the provider

* Make sure that the various QA teams have approved the rc versions, see [test process](./docs/test-process.md) for more information.
* Update the `CHANGELOG.md` with the release notes
* Push a tag to the release branch (`release/v2`, `release/v3`, or `master`) which does not have a `-rc` suffix
* The CI will build the provider and generate a release on GitHub
* Make sure to validate that the release is picked up by the Terraform registry, you may need to find the "resync" button to accomplish this.
