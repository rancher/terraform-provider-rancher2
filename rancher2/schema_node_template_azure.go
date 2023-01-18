package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	azureConfigDriver = "azure"
)

//Types

type azureConfig struct {
	AvailabilitySet        string   `json:"availabilitySet,omitempty" yaml:"availabilitySet,omitempty"`
	ClientID               string   `json:"clientId,omitempty" yaml:"clientId,omitempty"`
	ClientSecret           string   `json:"clientSecret,omitempty" yaml:"clientSecret,omitempty"`
	CustomData             string   `json:"customData,omitempty" yaml:"customData,omitempty"`
	DiskSize               string   `json:"diskSize,omitempty" yaml:"diskSize,omitempty"`
	DNS                    string   `json:"dns,omitempty" yaml:"dns,omitempty"`
	Environment            string   `json:"environment,omitempty" yaml:"environment,omitempty"`
	FaultDomainCount       string   `json:"faultDomainCount,omitempty" yaml:"faultDomainCount,omitempty"`
	Image                  string   `json:"image,omitempty" yaml:"image,omitempty"`
	Location               string   `json:"location,omitempty" yaml:"location,omitempty"`
	ManagedDisks           bool     `json:"managedDisks,omitempty" yaml:"managedDisks,omitempty"`
	NoPublicIP             bool     `json:"noPublicIp,omitempty" yaml:"noPublicIp,omitempty"`
	NSG                    string   `json:"nsg,omitempty" yaml:"nsg,omitempty"`
	Plan                   string   `json:"plan,omitempty" yaml:"plan,omitempty"`
	OpenPort               []string `json:"openPort,omitempty" yaml:"openPort,omitempty"`
	PrivateAddressOnly     bool     `json:"privateAddressOnly,omitempty" yaml:"privateAddressOnly,omitempty"`
	PrivateIPAddress       string   `json:"privateIpAddress,omitempty" yaml:"privateIpAddress,omitempty"`
	ResourceGroup          string   `json:"resourceGroup,omitempty" yaml:"resourceGroup,omitempty"`
	Size                   string   `json:"size,omitempty" yaml:"size,omitempty"`
	SSHUser                string   `json:"sshUser,omitempty" yaml:"sshUser,omitempty"`
	StaticPublicIP         bool     `json:"staticPublicIp,omitempty" yaml:"staticPublicIp,omitempty"`
	StorageType            string   `json:"storageType,omitempty" yaml:"storageType,omitempty"`
	Subnet                 string   `json:"subnet,omitempty" yaml:"subnet,omitempty"`
	SubnetPrefix           string   `json:"subnetPrefix,omitempty" yaml:"subnetPrefix,omitempty"`
	SubscriptionID         string   `json:"subscriptionId,omitempty" yaml:"subscriptionId,omitempty"`
	UpdateDomainCount      string   `json:"updateDomainCount,omitempty" yaml:"updateDomainCount,omitempty"`
	UsePrivateIP           bool     `json:"usePrivateIp,omitempty" yaml:"usePrivateIp,omitempty"`
	Vnet                   string   `json:"vnet,omitempty" yaml:"vnet,omitempty"`
	Tags                   string   `json:"tags,omitempty" yaml:"tags,omitempty"`
	AcceleratedNetworking  bool     `json:"acceleratedNetworking,omitempty" yaml:"acceleratedNetworking,omitempty"`
	AvailabilityZone       string   `json:"availabilityZone,omitempty" yaml:"availabilityZone,omitempty"`
	UsePublicIPStandardSKU bool     `json:"enablePublicIpStandardSku,omitempty" yaml:"enablePublicIpStandardSku,omitempty"`
}

//Schemas

func azureConfigFields() map[string]*schema.Schema {
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
		"plan": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Purchase plan for Azure Virtual Machine (in <publisher>:<product>:<plan> format)",
		},
		"open_port": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Make the specified port number accessible from the Internet",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
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
			Default:     "Standard_A2",
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
		"tags": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "Tags to be applied to the Azure VM instance (e.g. key1,value1,key2,value2)",
		},
		"accelerated_networking": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     "",
			Description: "Enable Accelerated Networking when creating an Azure Network Interface",
		},
		"availability_zone": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "",
			Description: "The Azure Availability Zone the VM should be created in",
		},
		"use_public_ip_standard_sku": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Use the Standard SKU when creating a public IP for an Azure VM",
		},
	}

	return s
}
