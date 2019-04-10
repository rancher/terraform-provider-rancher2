---
layout: "rancher2"
page_title: "Rancher2: rancher2_cloud_credential"
sidebar_current: "docs-rancher2-resource-cloud_credential"
description: |-
  Provides a Rancher v2.2.x Cloud Credential resource. This can be used to create Cloud Credential for rancher v2.2 node templates and retrieve their information.
---

# rancher2\_cloud\_credential

Provides a Rancher v2 Cloud Credential resource. This can be used to create Cloud Credential for rancher v2.2.x and retrieve their information. 

The following credentials config support typed parameters.

- amazonec2
- azure
- digitalocean
- openstack
- vsphere

Other credentials config can specify parameter as key/value pair using `generic_credential_config`.

## Example Usage

```hcl
# Create a new rancher2 Cloud Credential
resource "rancher2_cloud_credential" "foo" {
  name = "foo"
  description = "foo test"
  amazonec2_credential_config {
    access_key = "<AWS_ACCESS_KEY>"
    secret_key = "<AWS_SECRET_KEY>"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Cloud Credential (string)
* `amazonec2_credential_config` - (Optional) AWS config for the Cloud Credential (list maxitems:1)
* `azure_credential_config` - (Optional) Azure config for the Cloud Credential (list maxitems:1)
* `description` - (Optional) Description for the Cloud Credential (string)
* `digitalocean_credential_config` - (Optional) Digitalocean config for the Cloud Credential (list maxitems:1)
* `openstack_credential_config` - (Optional) Openstack config for the Cloud Credential (list maxitems:1)
* `vsphere_credential_config` - (Optional) vSphere config for the Cloud Credential (list maxitems:1)
* `annotations` - (Optional) Annotations for Cloud Credential object (map)
* `labels` - (Optional/Computed) Labels for Cloud Credential object (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `driver` - (Computed) The driver of the Cloud Credential (string)

## Nested blocks

### `amazonec2_credential_config`

#### Arguments

* `access_key` - (Required/Sensitive) AWS access key (string)
* `secret_key` - (Required/Sensitive) AWS secret key (string)

### `azure_credential_config`

#### Arguments

* `client_id` - (Required/Sensitive) Azure Service Principal Account ID (string)
* `client_secret` - (Required/Sensitive) Azure Service Principal Account password (string)
* `subscription_id` - (Required/Sensitive) Azure Subscription ID (string)

### `digitalocean_credential_config`

#### Arguments

* `access_token` - (Required/Sensitive) Digital Ocean access token (string)

### `openstack_credential_config`

#### Arguments

* `password` - (Required/Sensitive) Openstack password (string)

### `vsphere_credential_config`

#### Arguments

* `password` - (Required/Sensitive) vSphere password (string)
* `username` - (Required) vSphere username (string)
* `vcenter` - (Required) vSphere IP/hostname for vCenter (string)
* `vcenter_port` - (Optional) vSphere Port for vCenter. Default `443` (string)

### `generic_credential_config`

#### Arguments

* `driver` - (Required) The ID of the node driver 
* `config` - (Required) The parameters used by node driver (map)

#### Example Usage of `generic_credential_config`

```hcl
# with builtin driver
resource "rancher2_cloud_credential" "example" {
  name        = "example"
  description = "cloud credential with builtin driver"

  generic_config {
    driver = "rackspace"
    config {
      username = "XXXXXXXXXXXXXXXXXXXX"
      apiKey   = "XXXXXXXXXXXXXXXXXXXX"
    }
  }  
}
```

```hcl
# with custom driver
resource "rancher2_node_driver" "example" {
  active            = true
  builtin           = false
  checksum          = "xxx"
  name              = "example"
  ui_url            = "https://www.example.com/ui-driver-example/component.js"
  url               = "https://www.example.com/ui-driver-example/docker-machine-driver-example_linux-amd64.zip"
  whitelist_domains = ["www.example.com"]
}

resource "rancher2_cloud_credential" "example" {
  name        = "example-credential"
  description = "cloud credential with custom driver"

  generic_config {
    driver = "${rancher2_node_driver.example.id}"
    config {
      username = "XXXXXXXXXXXXXXXXXXXX"
      apiKey   = "XXXXXXXXXXXXXXXXXXXX"
    }
  }  
}

```

## Timeouts

`rancher2_cloud_credential` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating cloud credentials.
- `update` - (Default `10 minutes`) Used for cloud credential modifications.
- `delete` - (Default `10 minutes`) Used for deleting cloud credentials.

