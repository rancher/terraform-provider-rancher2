Terraform Provider for Rancher v2
==================================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Features WIP
-------------

[Implemented resources](https://github.com/rancher/terraform-provider-rancher2/blob/master/website/docs/r)

- Auth Config providers
  - ActiveDirectory
  - ADFS
  - AzureAD
  - Github
  - FreeIpa
  - OpenLdap
  - Ping
- Catalogs
- Clusters
  - Amazon EKS
  - Azure AKS
  - Google GKE
  - Imported
  - RKE
    - Cloud providers adding node pools
    - Custom
- Clusters & Projects logging
  - Elasticsearch
  - Fluentd
  - Kafka
  - Splunk
  - Syslog
- Clusters & Projects Role Template Bindings
- Namespaces
  - Resource quota limits (Rancher v2.1.x or higher )
- Node Pools
- Node Driver
- Projects
  - Resource quota limits (Rancher v2.1.x or higher )

TODO

- Node Template


Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.11.x
-	[Docker](https://docs.docker.com/install/) to build the provider plugin
- [Trash](https://github.com/rancher/trash/releases) 0.2.6 (to manage vendor dependencies)

Using the provider
----------------------

- Build or download rancher2 provider binary from [relases](https://github.com/rancher/terraform-provider-rancher2/releases)
- Copy rancher2 provider binary to same folder as terraform binary.
- Create tf file and use it.

```sh
$ terraform init
$ terraform plan
$ terraform apply
```

Building The Provider
---------------------

Clone this repository and run make:

```sh
$ git clone git@github.com:rancher/terraform-provider-rancher2
$ cd terraform-provider-rancher2
```

- Building Linux binary. Released at `bin/` folder

```sh
$ make
```

- Building linux, osx and windoes binaries. Released at `build/bin/` folder

```sh
$ CROSS=1 make
```

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.8+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make`. This will build the provider and put the provider binary in `bin/terraform-provider-rancher2` .

To compile binary on repo path and test on terraform:

```sh
$ make
$ terraform init
$ terraform plan
$ terraform apply
```

Managing vendor dependencies
-----------------------------

Go vendor dependencies are managed with [Trash](https://github.com/rancher/trash) and vendor.conf file.

To update vendor dependencies, edit `vendor.conf` file and execute trash

```sh
$ trash
```

Acceptance tests
----------------

For execute acceptance tests, a running rancher HA system and a rancher API key are needed.

To run acceptance tests, export `RANCHER_URL` with rancher url and `RANCHER_TOKEN_KEY` with bearer token or `RANCHER_ACCESS_KEY` with rancher acces key and `RANCHER_SECRET_KEY` with rancher secret key and execute

```sh
$ export RANCHER_URL=<URL>
$ export RANCHER_TOKEN_KEY=<TOKEN>
$ scripts/testacc
```
