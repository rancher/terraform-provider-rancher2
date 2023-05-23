---
page_title: "rancher2_machine_config_v2 Resource"
---

# rancher2\_machine\_config\_v2 Resource

Provides a Rancher v2 Machine config v2 resource. This can be used to create Machine Config v2 for Rancher v2 and retrieve their information. This resource is available from Rancher v2.6.0 and above.

`amazonec2`, `azure`, `digitalocean`, `harvester`, `linode`, `openstack`, and `vsphere` cloud providers are supported for machine config V2

**Note** This resource is used by 

## Example Usage

```hcl
# Create amazonec2 machine config v2
resource "rancher2_machine_config_v2" "foo" {
  generate_name = "test-foo"
  amazonec2_config {
    ami =  "<AMI_ID>"
    region = "<REGION>"
    security_group = [<AWS_SG>]
    subnet_id = "<SUBNET_ID>"
    vpc_id = "<VPC_ID>"
    zone = "<ZONE>"
  }
}
```
### Using the Harvester Node Driver

```hcl
# Get imported harvester cluster info
data "rancher2_cluster_v2" "foo-harvester" {
  name = "foo-harvester"
}

# Create a new Cloud Credential for an imported Harvester cluster
resource "rancher2_cloud_credential" "foo-harvester" {
  name = "foo-harvester"
  harvester_credential_config {
    cluster_id = data.rancher2_cluster_v2.foo-harvester.cluster_v1_id
    cluster_type = "imported"
    kubeconfig_content = data.rancher2_cluster_v2.foo-harvester.kube_config
  }
}

# Create a new rancher2 machine config v2 using harvester node_driver
resource "rancher2_machine_config_v2" "foo-harvester-v2" {
  generate_name = "foo-harvester-v2"
  harvester_config {
    vm_namespace = "default"
    cpu_count = "2"
    memory_size = "4"
    disk_info = <<EOF
    {
        "disks": [{
            "imageName": "harvester-public/image-57hzg",
            "size": 40,
            "bootOrder": 1
        }]
    }
    EOF
    network_info = <<EOF
    {
        "interfaces": [{
            "networkName": "harvester-public/vlan1"
        }]
    }
    EOF
    ssh_user = "ubuntu"
    user_data = <<EOF
    package_update: true
    packages:
      - qemu-guest-agent
      - iptables
    runcmd:
      - - systemctl
        - enable
        - '--now'
        - qemu-guest-agent.service
    EOF
  }
}
```

## Argument Reference

The following arguments are supported:

* `generate_name` - (Required/ForceNew) Cluster V2 generate name. The pattern to generate machine config name. e.g  generate_name=\"prod-pool1\" will generate \"nc-prod-pool1-?????\" name computed at `name` attribute (string)
* `fleet_namespace` - (Optional/ForceNew) Cluster V2 fleet namespace
* `amazonec2_config` - (Optional) AWS config for the Machine Config V2. Conflicts with `azure_config`, `digitalocean_config`, `harvester_config`, `linode_config`, `openstack_config` and `vsphere_config` (list maxitems:1)
* `azure_config` - (Optional) Azure config for the Machine Config V2. Conflicts with `amazonec2_config`, `digitalocean_config`, `harvester_config`, `linode_config`, `openstack_config` and `vsphere_config` (list maxitems:1)
* `digitalocean_config` - (Optional) Digitalocean config for the Machine Config V2. Conflicts with `amazonec2_config`, `azure_config`, `harvester_config`, `linode_config`, `openstack_config` and `vsphere_config` (list maxitems:1)
* `harvester_config` - (Optional) Harvester config for the Machine Config V2. Conflicts with `amazonec2_config`, `azure_config`, `digitalocean_config`, `linode_config`, `openstack_config` and `vsphere_config` (list maxitems:1)
* `linode_config` - (Optional) Linode config for the Machine Config V2. Conflicts with `amazonec2_config`, `azure_config`, `digitalocean_config`, `harvester_config`, `openstack_config` and `vsphere_config` (list maxitems:1)
* `openstack_config` - (Optional) Openstack config for the Machine Config V2. Conflicts with `amazonec2_config`, `azure_config`, `digitalocean_config`, `harvester_config`, `linode_config` and `vsphere_config` (list maxitems:1)
* `vsphere_config` - (Optional) vSphere config for the Machine Config V2. Conflicts with `amazonec2_config`, `azure_config`, `digitalocean_config`, `harvester_config`, `linode_config` and `openstack_config` (list maxitems:1)
* `annotations` - (Optional) Annotations for Machine Config V2 object (map)
* `labels` - (Optional/Computed) Labels for Machine Config V2 object (map)

**Note** `labels` and `node_taints` will be applied to nodes deployed using the Machine Config V2

## Attributes Reference

The following attributes are exported:

* `id` - (Computed) The ID of the resource (string)
* `kind` - (Computed) The machine config kind (string)
* `name` - (Computed) The machine config name (string)
* `resource_version` - (Computed) The machine config k8s resource version (string)

## Nested blocks

### `amazonec2_config`

#### Arguments

* `ami` - (Required) AWS machine image (string)
* `region` - (Required) AWS region. (string)
* `security_group` - (Required) AWS VPC security group. (list)
* `subnet_id` - (Required) AWS VPC subnet id (string)
* `vpc_id` - (Required) AWS VPC id. (string)
* `zone` - (Required) AWS zone for instance (i.e. a,b,c,d,e) (string)
* `block_duration_minutes` - (Optional) AWS spot instance duration in minutes (60, 120, 180, 240, 300, or 360). Default `0` (string)
* `device_name` - (Optional) AWS root device name. Default `/dev/sda1` (string)
* `encrypt_ebs_volume` - (Optional) Encrypt EBS volume. Default `false` (bool)
* `endpoint` - (Optional) Optional endpoint URL (hostname only or fully qualified URI) (string)
* `http_endpoint` - (Optional) Enables or disables the HTTP metadata endpoint on your instances (string)
* `http_tokens` - (Optional) The state of token usage for your instance metadata requests (string)
* `iam_instance_profile` - (Optional) AWS IAM Instance Profile (string)
* `insecure_transport` - (Optional) Disable SSL when sending requests (bool)
* `instance_type` - (Optional) AWS instance type. Default `t3a.medium` (string)
* `keypair_name` - (Optional) AWS keypair to use; requires --amazonec2-ssh-keypath (string)
* `kms_key` - (Optional) Custom KMS key ID using the AWS Managed CMK (string)
* `monitoring` - (Optional) Set this flag to enable CloudWatch monitoring. Deafult `false` (bool)
* `open_port` - (Optional) Make the specified port number accessible from the Internet. (list)
* `private_address_only` - (Optional) Only use a private IP address. Default `false` (bool)
* `request_spot_instance` - (Optional) Set this flag to request spot instance. Default `false` (bool)
* `retries` - (Optional) Set retry count for recoverable failures (use -1 to disable). Default `5` (string)
* `root_size` - (Optional) AWS root disk size (in GB). Default `16` (string)
* `security_group_readonly` - (Optional) Skip adding default rules to security groups (bool)
* `session_token` - (Optional/Sensitive) AWS Session Token (string)
* `spot_price` - (Optional) AWS spot instance bid price (in dollar). Default `0.50` (string)
* `ssh_key_contents` - (Optional/Sensitive) SSH Key for Instance (string)
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
* `fault_domain_count` - (Optional) Fault domain count to use for availability set. Default `3` (string)
* `image` - (Optional) Azure virtual machine OS image. Default `canonical:UbuntuServer:18.04-LTS:latest` (string)
* `location` - (Optional) Azure region to create the virtual machine. Default `westus` (string)
* `managed_disks` - (Optional) Configures VM and availability set for managed disks. Just for Rancher v2.3.x and above. Default `false` (bool)
* `no_public_ip` - (Optional) Do not create a public IP address for the machine. Default `false` (bool)
* `nsg` - (Optional) Azure Network Security Group to assign this node to (accepts either a name or resource ID, default is to create a new NSG for each machine). Default `docker-machine-nsg` (string)
* `open_port` - (Optional) Make the specified port number accessible from the Internet. (list)
* `private_address_only` - (Optional) Only use a private IP address. Default `false` (bool)
* `private_ip_address` - (Optional) Specify a static private IP address for the machine. (string)
* `resource_group` - (Optional) Azure Resource Group name (will be created if missing). Default `docker-machine` (string)
* `size` - (Optional) Size for Azure Virtual Machine. Default `Standard_A2` (string)
* `ssh_user` - (Optional) Username for SSH login (string)
* `static_public_ip` - (Optional) Assign a static public IP address to the machine. Default `false` (bool)
* `storage_type` - (Optional) Type of Storage Account to host the OS Disk for the machine. Default `Standard_LRS` (string)
* `subnet` - (Optional) Azure Subnet Name to be used within the Virtual Network. Default `docker-machine` (string)
* `subnet_prefix` - (Optional) Private CIDR block to be used for the new subnet, should comply RFC 1918. Default `192.168.0.0/16` (string)
* `subscription_id` - (Optional) Azure Subscription ID (string)
* `tenant_id` - (Optional) Azure Tenant ID (string)
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
* `ssh_key_contents` - (Optional/Sensitive) SSH private key contents (string)
* `ssh_key_fingerprint` - (Optional/Sensitive) SSH key fingerprint (string)
* `ssh_port` - (Optional) SSH port. Default `22` (string)
* `ssh_user` - (Optional) SSH username. Default `root` (string)
* `tags` - (Optional) Comma-separated list of tags to apply to the Droplet (string)
* `userdata` - (Optional) Path to file with cloud-init user-data (string)

### `harvester_config`

#### Arguments

* `vm_namespace` - (Required) Virtual machine namespace e.g. `default` (string)
* `cpu_count` - (Optional) CPU count, Default `2` (string)
* `memory_size` - (Optional) Memory size (in GiB), Default `4` (string)
* `disk_size` - (Deprecated) Use `disk_info` instead
* `disk_bus` - (Deprecated) Use `disk_info` instead
* `image_name` - (Deprecated) Use `disk_info` instead
* `disk_info` - (Required) A JSON string specifying info for the disks e.g. `{\"disks\":[{\"imageName\":\"harvester-public/image-57hzg\",\"bootOrder\":1,\"size\":40},{\"storageClassName\":\"node-driver-test\",\"bootOrder\":2,\"size\":1}]}` (string)
* `ssh_user` - (Required) SSH username e.g. `ubuntu` (string)
* `ssh_password` - (Optional/Sensitive) SSH password (string)
* `network_name` - (Deprecated) Use `network_info` instead
* `network_model` - (Deprecated) Use `network_info` instead
* `network_info` - (Required) A JSON string specifying info for the networks e.g. `{\"interfaces\":[{\"networkName\":\"harvester-public/vlan1\"},{\"networkName\":\"harvester-public/vlan2\"}]}` (string)
* `user_data` - (Optional) UserData content of cloud-init, base64 is supported. If the image does not contain the qemu-guest-agent package, you must install and start qemu-guest-agent using userdata (string)
* `network_data` - (Optional) NetworkData content of cloud-init, base64 is supported (string)
* `vm_affinity` - (Optional) Virtual machine affinity, only base64 format is supported. For Rancher v2.6.7 or above (string)

### `linode_config`

#### Arguments

* `authorized_users` - (Optional) Linode user accounts (seperated by commas) whose Linode SSH keys will be permitted root access to the created node. (string)
* `create_private_ip` - (Optional) Create private IP for the instance. Default `false` (bool)
* `docker_port` - (Optional) Docker Port. Default `2376` (string)
* `image` - (Optional) Specifies the Linode Instance image which determines the OS distribution and base files. Default `linode/ubuntu18.04` (string)
* `instance_type` - (Optional) Specifies the Linode Instance type which determines CPU, memory, disk size, etc. Default `g6-standard-4` (string)
* `label` - (Optional) Linode Instance Label. (string)
* `region` - (Optional) Specifies the region (location) of the Linode instance. Default `us-east` (string)
* `root_pass` - (Optional/Sensitive) Root Password (string)
* `ssh_port` - (Optional) SSH port. Default `22` (string)
* `ssh_user` - (Optional) SSH username. Default `root` (string)
* `stackscript` - (Optional) Specifies the Linode StackScript to use to create the instance. (string)
* `stackscript_data` - (Optional) A JSON string specifying data for the selected StackScript. (string)
* `swap_size` - (Optional) Linode Instance Swap Size (MB). Default `512` (string)
* `tags` - (Optional) A comma separated list of tags to apply to the the Linode resource (string)
* `token` - (Optional/Sensitive) Linode API token. Mandatory on Rancher v2.0.x and v2.1.x. Use `rancher2_cloud_credential` from Rancher v2.2.x (string)
* `ua_prefix` - (Optional) Prefix the User-Agent in Linode API calls with some 'product/version' (string)

### `openstack_config`

#### Arguments

* `auth_url` - (Required) OpenStack authentication URL (string)
* `availability_zone` - (Required) OpenStack availability zone (string)
* `region` - (Required) OpenStack region name (string)
* `username` - (Required++) OpenStack username (string)
* `active_timeout`- (Optional) OpenStack active timeout Default `200` (string)
* `cacert` - (Optional) CA certificate bundle to verify against (string)
* `config_drive` - (Optional) Enables the OpenStack config drive for the instance. Default `false` (bool)
* `domain_id` - (Required++) OpenStack domain ID. Identity v3 only. Conflicts with `domain_name` (string)
* `domain_name` - (Required++) OpenStack domain name. Identity v3 only. Conflicts with `domain_id` (string)
* `endpoint_type` - (Optional) OpenStack endpoint type. adminURL, internalURL or publicURL (string)
* `flavor_id` - (Required+) OpenStack flavor id to use for the instance. Conflicts with `flavor_name` (string)
* `flavor_name` - (Required+) OpenStack flavor name to use for the instance. Conflicts with `flavor_id` (string)
* `floating_ip_pool` - (Optional) OpenStack floating IP pool to get an IP from to assign to the instance (string)
* `image_id` - (Required+) OpenStack image id to use for the instance. Conflicts with `image_name` (string)
* `image_name` - (Required+) OpenStack image name to use for the instance. Conflicts with `image_id` (string)
* `insecure` - (Optional) Disable TLS credential checking. Default `false` (bool)
* `ip_version` - (Optional) OpenStack version of IP address assigned for the machine Default `4` (string)
* `keypair_name` - (Optional) OpenStack keypair to use to SSH to the instance (string)
* `net_id` - (Required+) OpenStack network id the machine will be connected on. Conflicts with `net_name` (string)
* `net_name` - (Required+) OpenStack network name the machine will be connected on. Conflicts with `net_id` (string)
* `nova_network` - (Optional) Use the nova networking services instead of neutron (string)
* `password` - (Optional/Sensitive) OpenStack password. Mandatory on Rancher v2.0.x and v2.1.x. Use `rancher2_cloud_credential` from Rancher v2.2.x (string)
* `private_key_file` - (Optional/Sensitive) Private key content to use for SSH (string)
* `sec_groups` - (Optional) OpenStack comma separated security groups for the machine (string)
* `ssh_port` - (Optional) OpenStack SSH port * Default `22` (string)
* `ssh_user` - (Optional) OpenStack SSH user * Default: `root` (string)
* `tenant_id` - (Required++) OpenStack tenant id. Conflicts with `tenant_name` (string)
* `tenant_name` - (Required++) OpenStack tenant name. Conflicts with `tenant_id` (string)
* `tenant_domain_id` - (Required++) OpenStack tenant domain id. Conflicts with `tenant_domain_name` (string)
* `tenant_domain_name` - (Required++) OpenStack tenant domain name. Conflicts with `tenant_domain_id` (string)
* `user_data_file` - (Optional) File containing an openstack userdata script (string)
* `user_domain_id` - (Required++) OpenStack user domain id. Conflicts with `user_domain_name` (string)
* `user_domain_name` - (Required++) OpenStack user domain name. Conflicts with `user_domain_id` (string)
* `application_credential_id` - (Optional) OpenStack application credential id. Conflicts with `application_credential_name` (string)
* `application_credential_name` - (Optional) OpenStack application credential name. Conflicts with `application_credential_id` (string)
* `application_credential_secret` - (Optional) OpenStack application credential secret (string)
* `boot_from_volume` - (Optional) Enable booting from volume. Default is `false` (bool)
* `volume_size` - (Optional) OpenStack volume size (GiB). Required when `boot_from_volume` is `true` (string)
* `volume_type` - (Optional) OpenStack volume type. Required when `boot_from_volume` is `true` and openstack cloud does not have a default volume type (string)
* `volume_id` - (Optional) OpenStack volume id of existing volume. Applicable only when `boot_from_volume` is `true` (string)
* `volume_name` - (Optional) OpenStack volume name of existing volume. Applicable only when `boot_from_volume` is `true` (string)
* `volume_device_path` - (Optional) OpenStack volume device path (attaching). Applicable only when `boot_from_volume` is `true`. Omit for auto `/dev/vdb`. (string)
> **Note**: `Required+` denotes that either the _name or _id is required but you cannot use both.
> **Note**: `Required++` denotes that either the _name or _id is required unless `application_credential_id` is defined.

### `vsphere_config`

#### Arguments

* `boot2docker_url` - (Optional) vSphere URL for boot2docker iso image. Default `https://releases.rancher.com/os/latest/rancheros-vmware.iso` (string)
* `cfgparam` - (Optional) vSphere vm configuration parameters (used for guestinfo) (list)
* `clone_from` - (Optional) If you choose creation type vm (clone vm) a name of what vm you want to clone is required (string)
* `cloud_config` - (Optional) Filepath to a cloud-config yaml file to put into the ISO user-data (string)
* `cloudinit` - (Optional) vSphere cloud-init file or url to set in the guestinfo (string)
* `content_library` - (Optional) If you choose to clone from a content library template specify the name of the library (string)
* `cpu_count` - (Optional) vSphere CPU number for docker VM. Default `2` (string)
* `creation_type` - (Optional) Creation type when creating a new virtual machine. Supported values: vm, template, library, legacy. Default `legacy` (string)
* `custom_attributes` - (Optional) vSphere custom attributes, format key/value e.g. `200=my custom value` (List)
* `datacenter` - (Optional) vSphere datacenter for docker VM (string)
* `datastore` - (Optional) vSphere datastore for docker VM (string)
* `datastore_cluster` - (Optional) vSphere datastore cluster for virtual machine (string)
* `disk_size` - (Optional) vSphere size of disk for docker VM (in MB). Default `20480` (string)
* `folder` - (Optional) vSphere folder for the docker VM. This folder must already exist in the datacenter (string)
* `hostsystem` - (Optional) vSphere compute resource where the docker VM will be instantiated. This can be omitted if using a cluster with DRS (string)
* `memory_size` - (Optional) vSphere size of memory for docker VM (in MB). Default `2048` (string)
* `network` - (Optional) vSphere network where the docker VM will be attached (list)
* `pool` - (Optional) vSphere resource pool for docker VM (string)
* `ssh_password` - (Optional) If using a non-B2D image you can specify the ssh password. Default `tcuser` (string)
* `ssh_port` - (Optional) If using a non-B2D image you can specify the ssh port. Default `22` (string)
* `ssh_user` - (Optional) If using a non-B2D image you can specify the ssh user. Default `docker`. (string)
* `ssh_user_group` - (Optional) If using a non-B2D image the uploaded keys will need chown'ed. Default `staff` (string)
* `tags` - (Optional) vSphere tags id e.g. `urn:xxx` (list)
* `vapp_ip_allocation_policy` - (Optional) vSphere vApp IP allocation policy. Supported values are: `dhcp`, `fixed`, `transient` and `fixedAllocated` (string)
* `vapp_ip_protocol` - (Optional) vSphere vApp IP protocol for this deployment. Supported values are: `IPv4` and `IPv6` (string)
* `vapp_property` - (Optional) vSphere vApp properties (list)
* `vapp_transport` - (Optional) vSphere OVF environment transports to use for properties. Supported values are: `iso` and `com.vmware.guestInfo` (string)
* `vcenter` - (Optional/Sensitive) vSphere IP/hostname for vCenter (string)
* `vcenter_port` - (Optional/Sensitive) vSphere Port for vCenter Default `443` (string)

## Timeouts

`rancher2_machine_config_v2` provides the following
[Timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) configuration options:

- `create` - (Default `10 minutes`) Used for creating machine configs.
- `update` - (Default `10 minutes`) Used for machine config modifications.
- `delete` - (Default `10 minutes`) Used for deleting machine configs.
