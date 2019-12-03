---
layout: "rancher2"
page_title: "Rancher2: rancher2_node_template"
sidebar_current: "docs-rancher2-resource-node_template"
description: |-
  Provides a Rancher v2 Node Template resource. This can be used to create Node template for Rancher v2 RKE clusters and retrieve their information.
---

# rancher2\_node\_template

Provides a Rancher v2 Node Template resource. This can be used to create Node Template for Rancher v2 and retrieve their information. 

amazonec2, azure, digitalocean, openstack and vsphere drivers are supported for node templates.

**Note** If you are upgrading to Rancher v2.3.3, please take a look to [final section](#Upgrading-to-Rancher-v2.3.3)

## Example Usage

```hcl
# Create a new rancher2 Node Template up to Rancher 2.1.x
resource "rancher2_node_template" "foo" {
  name = "foo"
  description = "foo test"
  amazonec2_config {
    access_key = "AWS_ACCESS_KEY"
    secret_key = "<AWS_SECRET_KEY>"
    ami =  "<AMI_ID>"
    region = "<REGION>"
    security_group = ["<AWS_SECURITY_GROUP>"]
    subnet_id = "<SUBNET_ID>"
    vpc_id = "<VPC_ID>"
    zone = "<ZONE>"
  }
}
```

```hcl
# Create a new rancher2 Node Template from Rancher 2.2.x
resource "rancher2_cloud_credential" "foo" {
  name = "foo"
  description = "foo test"
  amazonec2_credential_config {
    access_key = "<AWS_ACCESS_KEY>"
    secret_key = "<AWS_SECRET_KEY>"
  }
}
resource "rancher2_node_template" "foo" {
  name = "foo"
  description = "foo test"
  cloud_credential_id = "${rancher2_cloud_credential.foo.id}"
  amazonec2_config {
    ami =  "<AMI_ID>"
    region = "<REGION>"
    security_group = ["<AWS_SECURITY_GROUP>"]
    subnet_id = "<SUBNET_ID>"
    vpc_id = "<VPC_ID>"
    zone = "<ZONE>"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Node Template (string)
* `amazonec2_config` - (Optional) AWS config for the Node Template (list maxitems:1)
* `auth_certificate_authority` - (Optional/Sensitive) Auth certificate authority for the Node Template (string)
* `auth_key` - (Optional/Sensitive) Auth key for the Node Template (string)
* `azure_config` - (Optional) Azure config for the Node Template (list maxitems:1)
* `cloud_credential_id` - (Optional) Cloud credential ID for the Node Template. Required from Rancher v2.2.x (string)
* `description` - (Optional) Description for the Node Template (string)
* `digitalocean_config` - (Optional) Digitalocean config for the Node Template (list maxitems:1)
* `engine_env` - (Optional) Engine environment for the node template (string)
* `engine_insecure_registry` - (Optional) Insecure registry for the node template (list)
* `engine_install_url` - (Optional) Docker engine install URL for the node template. Default `https://releases.rancher.com/install-docker/18.09.sh`. Available install docker versions at `https://github.com/rancher/install-docker` (string)
* `engine_label` - (Optional) Engine label for the node template (string)
* `engine_opt` - (Optional) Engine options for the node template (map)
* `engine_registry_mirror` - (Optional) Engine registry mirror for the node template (list)
* `engine_storage_driver` - (Optional) Engine storage driver for the node template (string)
* `openstack_config` - (Optional) Openstack config for the Node Template (list maxitems:1)
* `use_internal_ip_address` - (Optional) Engine storage driver for the node template (bool)
* `vsphere_config` - (Optional) vSphere config for the Node Template (list maxitems:1)
* `annotations` - (Optional) Annotations for Node Template object (map)
* `labels` - (Optional/Computed) Labels for Node Template object (map)

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `driver` - (Computed) The driver of the node template (string)

## Nested blocks

### `amazonec2_config`

#### Arguments

* `ami` - (Required) AWS machine image (string)
* `region` - (Required) AWS region. (string)
* `security_group` - (Required) AWS VPC security group. (list)
* `subnet_id` - (Required) AWS VPC subnet id (string)
* `vpc_id` - (Required) AWS VPC id. (string)
* `zone` - (Required) AWS zone for instance (i.e. a,b,c,d,e) (string)
* `access_key` - (Optional/Sensitive) AWS access key. Required on Rancher v2.0.x and v2.1.x. Use `rancher2_cloud_credential` from Rancher v2.2.x (string)
* `block_duration_minutes` - (Optional) AWS spot instance duration in minutes (60, 120, 180, 240, 300, or 360). Default `0` (string)
* `device_name` - (Optional) AWS root device name. Default `/dev/sda1` (string)
* `endpoint` - (Optional) Optional endpoint URL (hostname only or fully qualified URI) (string)
* `iam_instance_profile` - (Optional) AWS IAM Instance Profile (string)
* `insecure_transport` - (Optional) Disable SSL when sending requests (bool)
* `instance_type` - (Optional) AWS instance type. Default `t2.micro` (string)
* `keypair_name` - (Optional) AWS keypair to use; requires --amazonec2-ssh-keypath (string)
* `monitoring` - (Optional) Set this flag to enable CloudWatch monitoring. Deafult `false` (bool)
* `open_port` - (Optional) Make the specified port number accessible from the Internet. (list)
* `private_address_only` - (Optional) Only use a private IP address. Default `false` (bool)
* `request_spot_instance` - (Optional) Set this flag to request spot instance. Default `false` (bool)
* `retries` - (Optional) Set retry count for recoverable failures (use -1 to disable). Default `5` (string)
* `root_size` - (Optional) AWS root disk size (in GB). Default `16` (string)
* `secret_key` - (Optional/Sensitive) AWS secret key. Required on Rancher v2.0.x and v2.1.x. Use `rancher2_cloud_credential` from Rancher v2.2.x (string)
* `security_group_readonly` - (Optional) Skip adding default rules to security groups (bool)
* `session_token` - (Optional/Sensitive) AWS Session Token (string)
* `spot_price` - (Optional) AWS spot instance bid price (in dollar). Default `0.50` (string)
* `ssh_keypath` - (Optional) SSH Key for Instance (string)
* `ssh_user` - (Optional) Set the name of the ssh user (string)
* `tags` - (Optional) AWS Tags (e.g. key1,value1,key2,value2) (string)
* `use_ebs_optimized_instance` - (Optional) Create an EBS optimized instance. Default `false` (bool)
* `use_private_address` - (Optional) Force the usage of private IP address. Default `false` (bool)
* `userdata` - (Optional) Path to file with cloud-init user data (string)
* `volume_type` - (Optional) Amazon EBS volume type. Default `gp2` (string)

### `azure_config`

#### Arguments

* `client_id` - (Optional/Sensitive) Azure Service Principal Account ID. Mandatory on Rancher v2.0.x and v2.1.x. Use `rancher2_cloud_credential` from Rancher v2.2.x (string)
* `client_secret` - (Optional/Sensitive) Azure Service Principal Account password. Mandatory on Rancher v2.0.x and v2.1.x. Use `rancher2_cloud_credential` from Rancher v2.2.x (string)
* `subscription_id` - (Optional/Sensitive) Azure Subscription ID. Mandatory on Rancher v2.0.x and v2.1.x. Use `rancher2_cloud_credential` from Rancher v2.2.x (string)
* `availability_set` - (Optional) Azure Availability Set to place the virtual machine into. Default `docker-machine` (string)
* `custom_data` - (Optional) Path to file with custom-data (string)
* `disk_size` - (Optional) Disk size if using managed disk. Just for Rancher v2.3.x and above. Default `30` (string)
* `dns` - (Optional) A unique DNS label for the public IP adddress (string)
* `docker_port` - (Optional) Port number for Docker engine. Default `2376` (string)
* `environment` - (Optional) Azure environment (e.g. AzurePublicCloud, AzureChinaCloud). Default `AzurePublicCloud` (string)
`fault_domain_count` - (Optional) Fault domain count to use for availability set. Default `3` (string)
* `image` - (Optional) Azure virtual machine OS image. Default `canonical:UbuntuServer:18.04-LTS:latest` (string)
* `location` - (Optional) Azure region to create the virtual machine. Default `westus` (string)
* `managed_disks` - (Optional) Configures VM and availability set for managed disks. Just for Rancher v2.3.x and above. Default `false` (bool)
* `no_public_ip` - (Optional) Do not create a public IP address for the machine. Default `false` (bool)
* `open_port` - (Optional) Make the specified port number accessible from the Internet. (list)
* `private_ip_address` - (Optional) Specify a static private IP address for the machine. (string)
* `resource_group` - (Optional) Azure Resource Group name (will be created if missing). Default `docker-machine` (string)
* `size` - (Optional) Size for Azure Virtual Machine. Default `Standard_A2` (string)
* `ssh_user` - (Optional) Username for SSH login (string)
* `static_public_ip` - (Optional) Assign a static public IP address to the machine. Default `false` (bool)
* `storage_type` - (Optional) Type of Storage Account to host the OS Disk for the machine. Default `Standard_LRS` (string)
* `subnet` - (Optional) Azure Subnet Name to be used within the Virtual Network. Default `docker-machine` (string)
* `subnet_prefix` - (Optional) Private CIDR block to be used for the new subnet, should comply RFC 1918. Default `192.168.0.0/16` (string)
* `update_domain_count` - (Optional) Update domain count to use for availability set. Default `5` (string)
* `use_private_ip` - (Optional) Use private IP address of the machine to connect. Default `false` (bool)
* `vnet` - (Optional) Azure Virtual Network name to connect the virtual machine (in [resourcegroup:]name format). Default `docker-machine-vnet` (string)

### `digitalocean_config`

#### Arguments

* `access_token` - (Optional/Sensitive) Digital Ocean access token. Mandatory on Rancher v2.0.x and v2.1.x. Use `rancher2_cloud_credential` from Rancher v2.2.x (string)
* `backups` - (Optional) Enable backups for droplet. Default `false` (bool)
* `image` - (Optional) Digital Ocean Image. Default `ubuntu-16-04-x64` (string)
* `ipv6` - (Optional) Enable ipv6 for droplet. Default `false` (bool)
* `monitoring` - (Optional) Enable monitoring for droplet. Default `false` (bool)
* `private_networking` - (Optional) Enable private networking for droplet. Default `false` (bool)
* `region` - (Optional) Digital Ocean region. Default `nyc3` (string)
* `size` - (Optional) Digital Ocean size. Default `s-1vcpu-1gb` (string)
* `ssh_key_fingerprint` - (Optional/Sensitive) SSH key fingerprint (string)
* `ssh_key_path` - (Optional) SSH private key path (string)
* `ssh_port` - (Optional) SSH port. Default `22` (string)
* `ssh_user` - (Optional) SSH username. Default `root` (string)
* `tags` - (Optional) Comma-separated list of tags to apply to the Droplet (string)
* `userdata` - (Optional) Path to file with cloud-init user-data (string)

### `openstack_config`

#### Arguments

* `auth_url` - (Required) OpenStack authentication URL (string)
* `availability_zone` - (Required) OpenStack availability zone (string)
* `region` - (Required) OpenStack region name (string)
* `username` - (Required) OpenStack username (string)
* `active_timeout`- (Optional) OpenStack active timeout Default `200` (string)
* `cacert` - (Optional) CA certificate bundle to verify against (string)
* `config_drive` - (Optional) Enables the OpenStack config drive for the instance. Default `false` (bool)
* `domain_id` - (Required*) OpenStack domain ID. Identity v3 only. Conflicts with `domain_name` (string)
* `domain_name` - (Required*) OpenStack domain name. Identity v3 only. Conflicts with `domain_id` (string)
* `endpoint_type` - (Optional) OpenStack endpoint type. adminURL, internalURL or publicURL (string)
* `flavor_id` - (Required*) OpenStack flavor id to use for the instance. Conflicts with `flavor_name` (string)
* `flavor_name` - (Required*) OpenStack flavor name to use for the instance. Conflicts with `flavor_id` (string)
* `floating_ip_pool` - (Optional) OpenStack floating IP pool to get an IP from to assign to the instance (string)
* `image_id` - (Required*) OpenStack image id to use for the instance. Conflicts with `image_name` (string)
* `image_name` - (Required*) OpenStack image name to use for the instance. Conflicts with `image_id` (string)
* `insecure` - (Optional) Disable TLS credential checking. Default `false` (bool)
* `ip_version` - (Optional) OpenStack version of IP address assigned for the machine Default `4` (string)
* `keypair_name` - (Optional) OpenStack keypair to use to SSH to the instance (string)
* `net_id` - (Required*) OpenStack network id the machine will be connected on. Conflicts with `net_name` (string)
* `net_name` - (Required*) OpenStack network name the machine will be connected on. Conflicts with `net_id` (string)
* `nova_network` - (Optional) Use the nova networking services instead of neutron (string)
* `password` - (Optional/Sensitive) OpenStack password. Mandatory on Rancher v2.0.x and v2.1.x. Use `rancher2_cloud_credential` from Rancher v2.2.x (string)
* `private_key_file` - (Optional) Private keyfile absolute path to use for SSH (string)
* `sec_groups` - (Optional) OpenStack comma separated security groups for the machine (string)
* `ssh_port` - (Optional) OpenStack SSH port * Default `22` (string)
* `ssh_user` - (Optional) OpenStack SSH user * Default: `root` (string)
* `tenant_id` - (Required*) OpenStack tenant id. Conflicts with `tenant_name` (string)
* `tenant_name` - (Required*) OpenStack tenant name. Conflicts with `tenant_id` (string)
* `user_data_file` - (Optional) File containing an openstack userdata script (string)

> **Note**: `Required*` denotes that either the _name or _id is required but you cannot use both.

### `vsphere_config`

#### Arguments

* `boot2docker_url` - (Optional) vSphere URL for boot2docker iso image. Default `https://releases.rancher.com/os/latest/rancheros-vmware.iso` (string)
* `cfgparam` - (Optional) vSphere vm configuration parameters (used for guestinfo) (list)
* `clone_from` - (Optional) If you choose creation type clone a name of what you want to clone is required. From Rancher v2.3.3 (string)
* `cloud_config` - (Optional) Filepath to a cloud-config yaml file to put into the ISO user-data. From Rancher v2.3.3 (string)
* `cloudinit` - (Optional) vSphere cloud-init file or url to set in the guestinfo (string)
* `content_library` - (Optional) If you choose to clone from a content library template specify the name of the library. From Rancher v2.3.3 (string)
* `cpu_count` - (Optional) vSphere CPU number for docker VM. Default `2` (string)
* `creation_type` - (Optional) Creation type when creating a new virtual machine. Supported values: vm, template, library, legacy. Default `legacy`. From Rancher v2.3.3 (string)
* `custom_attributes` - (Optional) vSphere custom attributes, format key/value e.g. `200=my custom value`. From Rancher v2.3.3 (List)
* `datacenter` - (Optional) vSphere datacenter for docker VM (string)
* `datastore` - (Optional) vSphere datastore for docker VM (string)
* `datastore_cluster` - (Optional) vSphere datastore cluster for virtual machine. From Rancher v2.3.3 (string)
* `disk_size` - (Optional) vSphere size of disk for docker VM (in MB). Default `20480` (string)
* `folder` - (Optional) vSphere folder for the docker VM. This folder must already exist in the datacenter (string)
* `hostsystem` - (Optional) vSphere compute resource where the docker VM will be instantiated. This can be omitted if using a cluster with DRS (string)
* `memory_size` - (Optional) vSphere size of memory for docker VM (in MB). Default `2048` (string)
* `network` - (Optional) vSphere network where the docker VM will be attached (list)
* `password` - (Optional/Sensitive) vSphere password. Mandatory on Rancher v2.0.x and v2.1.x. Use `rancher2_cloud_credential` from Rancher v2.2.x (string)
* `pool` - (Optional) vSphere resource pool for docker VM (string)
* `ssh_password` - (Optional) If using a non-B2D image you can specify the ssh password. Default `tcuser`. From Rancher v2.3.3 (string)
* `ssh_port` - (Optional) If using a non-B2D image you can specify the ssh port. Default `22`. From Rancher v2.3.3 (string)
* `ssh_user` - (Optional) If using a non-B2D image you can specify the ssh user. Default `docker`. From Rancher v2.3.3 (string)
* `ssh_user_group` - (Optional) If using a non-B2D image the uploaded keys will need chown'ed. Default `staff`. From Rancher v2.3.3 (string)
* `tags` - (Optional) vSphere tags id e.g. `urn:xxx`. From Rancher v2.3.3 (list)
* `username` - (Optional/Sensitive) vSphere username. Mandatory on Rancher v2.0.x and v2.1.x. Use `rancher2_cloud_credential` from Rancher v2.2.x (string)
* `vapp_ip_allocation_policy` - (Optional) vSphere vApp IP allocation policy. Supported values are: `dhcp`, `fixed`, `transient` and `fixedAllocated` (string)
* `vapp_ip_protocol` - (Optional) vSphere vApp IP protocol for this deployment. Supported values are: `IPv4` and `IPv6` (string)
* `vapp_property` - (Optional) vSphere vApp properties (list)
* `vapp_transport` - (Optional) vSphere OVF environment transports to use for properties. Supported values are: `iso` and `com.vmware.guestInfo` (string)
* `vcenter` - (Optional/Sensitive) vSphere IP/hostname for vCenter. Mandatory on Rancher v2.0.x and v2.1.x. Use `rancher2_cloud_credential` from Rancher v2.2.x (string)
* `vcenter_port` - (Optional/Sensitive) vSphere Port for vCenter. Mandatory on Rancher v2.0.x and v2.1.x. Use `rancher2_cloud_credential` from Rancher v2.2.x. Default `443` (string)

## Timeouts

`rancher2_node_template` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating node templates.
- `update` - (Default `10 minutes`) Used for node template modifications.
- `delete` - (Default `10 minutes`) Used for deleting node templates.

## Import

Node Template can be imported using the Rancher Node Template ID

```
$ terraform import rancher2_node_template.foo <node_template_id>
```

## Upgrading to Rancher v2.3.3

Due to [this feature](https://github.com/rancher/rancher/pull/23718) included on Rancher v2.3.3, `rancher2_node_template` are now global scope objects with RBAC around them, instead of user scope objects as they were. This means that existing node templates `id` field is changing on upgrade. Because the Terraform provider can not find the old `id`, it will try to recreate them.

As a workaround, if you are upgrading Rancher from previous releases to v2.3.3, you need to get node templates new id from Rancher API, refresh tfstate and import `rancher2_node_template` resources with new id.

```
$ curl -sk -X GET -H "Authorization: Bearer ${RANCHER_TOKEN_KEY}" ${RANCHER_URL}/v3/nodeTemplates | jq .data
$ terraform refresh
$ terraform import rancher2_node_template.<name> <new_id>
$ terraform apply
```
