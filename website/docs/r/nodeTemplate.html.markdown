---
layout: "rancher2"
page_title: "Rancher2: rancher2_node_template"
sidebar_current: "docs-rancher2-resource-node_template"
description: |-
  Provides a Rancher v2 Node Template resource. This can be used to create Node template for rancher v2 rke clusters and retrieve their information.
---

# rancher2\_node\_template

Provides a Rancher v2 Node Template resource. This can be used to create Node Template for rancher v2 and retrieve their information. 

Only amazonec2, azure and digitalocean drivers are supported for node templates.

## Example Usage

```hcl
# Create a new rancher2 Node Template
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

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Node Template.
* `amazonec2_config` - (Optional) AWS config for the Node Template.
* `auth_certificate_authority` - (Optional/Sensitive) Auth certificate authority for the Node Template.
* `auth_key` - (Optional/Sensitive) Auth key for the Node Template.
* `azure_config` - (Optional) Azure config for the Node Template.
* `description` - (Optional) Description for the Node Template.
* `digitalocean_config` - (Optional) Digitalocean config for the Node Template.
* `docker_version` - (Optional) Docker version for the node template.
* `engine_env` - (Optional) Engine environment for the node template.
* `engine_insecure_registry` - (Optional) Insecure registry for the node template.
* `engine_install_url` - (Optional) Engine install URL for the node template.
* `engine_label` - (Optional) Engine label for the node template.
* `engine_opt` - (Optional) Engine options for the node template.
* `engine_registry_mirror` - (Optional) Engine registry mirror for the node template.
* `engine_storage_driver` - (Optional) Engine storage driver for the node template.
* `use_internal_ip_address` - (Optional) Engine storage driver for the node template (bool)
* `annotations` - (Optional) Annotations for Node Template object.
* `labels` - (Optional/Computed) Labels for Node Template object.


### Rancher `amazonec2_config`

The following arguments are supported:

* `access_key` - (Required/Sensitive) AWS access key (string)
* `secret_key` - (Required/Sensitive) AWS secret key (string)
* `ami` - (Required) AWS machine image (string)
* `region` - (Required) AWS region. (string)
* `security_group` - (Required) AWS VPC security group. (list)
* `subnet_id` - (Required) AWS VPC subnet id (string)
* `vpc_id` - (Required) AWS VPC id. (string)
* `zone` - (Required) AWS zone for instance (i.e. a,b,c,d,e) (string)
* `block_duration_minutes` - (Optional) AWS spot instance duration in minutes (60, 120, 180, 240, 300, or 360). Default `0`
* `device_name` - (Optional) AWS root device name. Default `/dev/sda1`
* `endpoint` - (Optional) Optional endpoint URL (hostname only or fully qualified URI) (string)
* `iam_instance_profile` - (Optional) AWS IAM Instance Profile (string)
* `insecure_transport` - (Optional) Disable SSL when sending requests (bool)
* `instance_type` - (Optional) AWS instance type. Default `t2.micro`
* `keypair_name` - (Optional) AWS keypair to use; requires --amazonec2-ssh-keypath (string)
* `monitoring` - (Optional) Set this flag to enable CloudWatch monitoring. Deafult `false`
* `open_port` - (Optional) Make the specified port number accessible from the Internet. (list)
* `private_address_only` - (Optional) Only use a private IP address. Default `false`
* `request_spot_instance` - (Optional) Set this flag to request spot instance. Default `false`
* `retries` - (Optional) Set retry count for recoverable failures (use -1 to disable). Default `5`
* `root_size` - (Optional) AWS root disk size (in GB). Default `16`
* `security_group_readonly` - (Optional) Skip adding default rules to security groups (bool)
* `session_token` - (Optional/Sensitive) AWS Session Token (string)
* `spot_price` - (Optional) AWS spot instance bid price (in dollar). Default `0.50`
* `ssh_keypath` - (Optional) SSH Key for Instance (string)
* `ssh_user` - (Optional) Set the name of the ssh user (string)
* `tags` - (Optional) AWS Tags (e.g. key1,value1,key2,value2) (string)
* `use_ebs_optimized_instance` - (Optional) Create an EBS optimized instance. Default `false`
* `use_private_address` - (Optional) Force the usage of private IP address. Default `false`
* `userdata` - (Optional) Path to file with cloud-init user data (string)
* `volume_type` - (Optional) Amazon EBS volume type. Default `gp2`


### Rancher `azure_config`

The following arguments are supported:

* `client_id` - (Required/Sensitive) Azure Service Principal Account ID
* `client_secret` - (Required/Sensitive) Azure Service Principal Account password
* `subscription_id` - (Required) Azure Subscription ID (string)
* `availability_set` - (Optional) Azure Availability Set to place the virtual machine into. Default `docker-machine`
* `custom_data` - (Optional) Path to file with custom-data (string)
* `dns` - (Optional) A unique DNS label for the public IP adddress (string)
* `docker_port` - (Optional) Port number for Docker engine. Default `2376`
* `environment` - (Optional) Azure environment (e.g. AzurePublicCloud, AzureChinaCloud). Default `AzurePublicCloud`
* `image` - (Optional) Azure virtual machine OS image. Default `canonical:UbuntuServer:16.04.0-LTS:latest`
* `location` - (Optional) Azure region to create the virtual machine. Default `westus`
* `no_public_ip` - (Optional) Do not create a public IP address for the machine. Default `false`
* `open_port` - (Optional) Make the specified port number accessible from the Internet. (list)
* `private_ip_address` - (Optional) Specify a static private IP address for the machine. (string)
* `resource_group` - (Optional) Azure Resource Group name (will be created if missing). Default `docker-machine`
* `size` - (Optional) Size for Azure Virtual Machine. Default `Standard_A2`
* `ssh_user` - (Optional) Username for SSH login (string)
* `static_public_ip` - (Optional) Assign a static public IP address to the machine. Default `false`
* `storage_type` - (Optional) Type of Storage Account to host the OS Disk for the machine. Default `Standard_LRS`
* `subnet` - (Optional) Azure Subnet Name to be used within the Virtual Network. Default `docker-machine`
* `subnet_prefix` - (Optional) Private CIDR block to be used for the new subnet, should comply RFC 1918. Default `192.168.0.0/16`
* `use_private_ip` - (Optional) Use private IP address of the machine to connect. Default `false`
* `vnet` - (Optional) Azure Virtual Network name to connect the virtual machine (in [resourcegroup:]name format). Default `docker-machine-vnet`


### Rancher `digitalocean_config`

The following arguments are supported:

* `access_token` - (Required/Sensitive) Digital Ocean access token
* `backups` - (Optional) Enable backups for droplet. Default `false`
* `image` - (Optional) Digital Ocean Image. Default `ubuntu-16-04-x64`
* `ipv6` - (Optional) Enable ipv6 for droplet. Default `false`
* `monitoring` - (Optional) Enable monitoring for droplet. Default `false`
* `private_networking` - (Optional) Enable private networking for droplet. Default `false`
* `region` - (Optional) Digital Ocean region. Default `nyc3`
* `size` - (Optional) Digital Ocean size. Default `s-1vcpu-1gb`
* `ssh_key_fingerprint` - (Optional/Sensitive) SSH key fingerprint
* `ssh_key_path` - (Optional) SSH private key path
* `ssh_port` - (Optional) SSH port. Default `22`
* `ssh_user` - (Optional) SSH username. Default `root`
* `tags` - (Optional) Comma-separated list of tags to apply to the Droplet
* `userdata` - (Optional) Path to file with cloud-init user-data

### Timeouts

`rancher2_node_template` provides the following
[Timeouts](/docs/configuration/resources.html#timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating node templates.
- `update` - (Default `10 minutes`) Used for node template modifications.
- `delete` - (Default `10 minutes`) Used for deleting node templates.

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource.
* `driver` - (Computed) The driver of the node template.

## Import

Node Pool can be imported using the rancher Node Template ID

```
$ terraform import rancher2_node_template.foo <node_template_id>
```

