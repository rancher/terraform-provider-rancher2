package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	azureConfigDriver = "azure"
)

//Types

type azureConfig struct {
	AvailabilitySet    string   `json:"availabilitySet,omitempty" yaml:"availabilitySet,omitempty"`
	ClientID           string   `json:"clientId,omitempty" yaml:"clientId,omitempty"`
	ClientSecret       string   `json:"clientSecret,omitempty" yaml:"clientSecret,omitempty"`
	CustomData         string   `json:"customData,omitempty" yaml:"customData,omitempty"`
	DNS                string   `json:"dns,omitempty" yaml:"dns,omitempty"`
	Environment        string   `json:"environment,omitempty" yaml:"environment,omitempty"`
	Image              string   `json:"image,omitempty" yaml:"image,omitempty"`
	Location           string   `json:"location,omitempty" yaml:"location,omitempty"`
	NoPublicIP         bool     `json:"noPublicIp,omitempty" yaml:"noPublicIp,omitempty"`
	OpenPort           []string `json:"openPort,omitempty" yaml:"openPort,omitempty"`
	PrivateAddressOnly bool     `json:"privateAddressOnly,omitempty" yaml:"privateAddressOnly,omitempty"`
	PrivateIPAddress   string   `json:"privateIpAddress,omitempty" yaml:"privateIpAddress,omitempty"`
	ResourceGroup      string   `json:"resourceGroup,omitempty" yaml:"resourceGroup,omitempty"`
	Size               string   `json:"size,omitempty" yaml:"size,omitempty"`
	SSHUser            string   `json:"sshUser,omitempty" yaml:"sshUser,omitempty"`
	StaticPublicIP     bool     `json:"staticPublicIp,omitempty" yaml:"staticPublicIp,omitempty"`
	StorageType        string   `json:"storageType,omitempty" yaml:"storageType,omitempty"`
	Subnet             string   `json:"subnet,omitempty" yaml:"subnet,omitempty"`
	SubnetPrefix       string   `json:"subnetPrefix,omitempty" yaml:"subnetPrefix,omitempty"`
	SubscriptionID     string   `json:"subscriptionId,omitempty" yaml:"subscriptionId,omitempty"`
	UsePrivateIP       bool     `json:"usePrivateIp,omitempty" yaml:"usePrivateIp,omitempty"`
	Vnet               string   `json:"vnet,omitempty" yaml:"vnet,omitempty"`
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
		"image": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "canonical:UbuntuServer:16.04.0-LTS:latest",
			Description: "Azure virtual machine OS image",
		},
		"location": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "westus",
			Description: "Azure region to create the virtual machine",
		},
		"no_public_ip": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Do not create a public IP address for the machine",
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
		"use_private_ip": {
			Type:        schema.TypeString,
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
	}

	return s
}
