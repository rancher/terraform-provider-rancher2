package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func machineConfigV2AzureFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"availability_set": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "docker-machine",
			Description: "Azure Availability Set to place the virtual machine into",
		},
		"client_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Azure Service Principal Account ID (optional, browser auth is used if not specified)",
		},
		"client_secret": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Azure Service Principal Account password (optional, browser auth is used if not specified)",
		},
		"custom_data": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Path to file with custom-data",
		},
		"disk_size": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "30",
			Description: "Disk size if using managed disk",
		},
		"dns": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "A unique DNS label for the public IP adddress",
		},
		"docker_port": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "2376",
			Description: "Port number for Docker engine",
		},
		"environment": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "AzurePublicCloud",
			Description: "Azure environment (e.g. AzurePublicCloud, AzureChinaCloud)",
		},
		"fault_domain_count": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "3",
			Description: "Fault domain count to use for availability set",
		},
		"image": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "canonical:UbuntuServer:18.04-LTS:latest",
			Description: "Azure virtual machine OS image",
		},
		"location": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "westus",
			Description: "Azure region to create the virtual machine",
		},
		"managed_disks": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Configures VM and availability set for managed disks",
		},
		"no_public_ip": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Do not create a public IP address for the machine",
		},
		"nsg": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "docker-machine-nsg",
			Description: "Azure Network Security Group to assign this node to (accepts either a name or resource ID, default is to create a new NSG for each machine)",
		},
		"open_port": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Make the specified port number accessible from the Internet",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"private_address_only": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Only use a private IP address",
		},
		"private_ip_address": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Specify a static private IP address for the machine",
		},
		"resource_group": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "docker-machine",
			Description: "Azure Resource Group name (will be created if missing)",
		},
		"size": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "Standard_D2_v2",
			Description: "Size for Azure Virtual Machine",
		},
		"ssh_user": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "docker-user",
			Description: "Username for SSH login",
		},
		"static_public_ip": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Assign a static public IP address to the machine",
		},
		"storage_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "Standard_LRS",
			Description: "Type of Storage Account to host the OS Disk for the machine",
		},
		"subnet": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "docker-machine",
			Description: "Azure Subnet Name to be used within the Virtual Network",
		},
		"subnet_prefix": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "192.168.0.0/16",
			Description: "Private CIDR block to be used for the new subnet, should comply RFC 1918",
		},
		"subscription_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Azure Subscription ID",
		},
		"tenant_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Azure Tenant ID",
		},
		"update_domain_count": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "5",
			Description: "Update domain count to use for availability set",
		},
		"use_private_ip": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Use private IP address of the machine to connect",
		},
		"vnet": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "docker-machine-vnet",
			Description: "Azure Virtual Network name to connect the virtual machine (in [resourcegroup:]name format)",
		},
		"availability_zone": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "The Availability Zone that the Azure VM should be created in",
		},
		"tags": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "Tags to be applied to the Azure VM instance (e.g. key1,value1,key2,value2)",
		},
		"use_public_ip_standard_sku": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Use the standard SKU when creating a Public IP for the Azure VM instance",
		},
		"accelerated_networking": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Use Accelerated Networking when creating a network interface for the Azure VM",
		},
	}

	return s
}
