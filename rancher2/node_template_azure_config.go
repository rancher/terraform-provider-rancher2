package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	azureConfigDriver = "azure"
)

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
			Required:    true,
			Sensitive:   true,
			Description: "Azure Service Principal Account ID (optional, browser auth is used if not specified)",
		},
		"client_secret": {
			Type:        schema.TypeString,
			Required:    true,
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
			Required:    true,
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

// Flatteners

func flattenAzureConfig(in *azureConfig) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	if len(in.AvailabilitySet) > 0 {
		obj["availability_set"] = in.AvailabilitySet
	}

	if len(in.ClientID) > 0 {
		obj["client_id"] = in.ClientID
	}

	if len(in.ClientSecret) > 0 {
		obj["client_secret"] = in.ClientSecret
	}

	if len(in.CustomData) > 0 {
		obj["custom_data"] = in.CustomData
	}

	if len(in.DNS) > 0 {
		obj["dns"] = in.DNS
	}

	if len(in.Environment) > 0 {
		obj["environment"] = in.Environment
	}

	if len(in.Image) > 0 {
		obj["image"] = in.Image
	}

	if len(in.Location) > 0 {
		obj["location"] = in.Location
	}

	obj["no_public_ip"] = in.NoPublicIP

	if len(in.OpenPort) > 0 {
		obj["open_port"] = toArrayInterface(in.OpenPort)
	}

	obj["private_address_only"] = in.PrivateAddressOnly

	if len(in.PrivateIPAddress) > 0 {
		obj["private_ip_address"] = in.PrivateIPAddress
	}

	if len(in.ResourceGroup) > 0 {
		obj["resource_group"] = in.ResourceGroup
	}

	if len(in.Size) > 0 {
		obj["size"] = in.Size
	}

	if len(in.SSHUser) > 0 {
		obj["ssh_user"] = in.SSHUser
	}

	obj["static_public_ip"] = in.StaticPublicIP

	if len(in.StorageType) > 0 {
		obj["storage_type"] = in.StorageType
	}

	if len(in.Subnet) > 0 {
		obj["subnet"] = in.Subnet
	}

	if len(in.SubnetPrefix) > 0 {
		obj["subnet_prefix"] = in.SubnetPrefix
	}

	if len(in.SubscriptionID) > 0 {
		obj["subscription_id"] = in.SubscriptionID
	}

	obj["use_private_ip"] = in.UsePrivateIP

	if len(in.Vnet) > 0 {
		obj["vnet"] = in.Vnet
	}

	return []interface{}{obj}
}

// Expanders

func expandAzureConfig(p []interface{}) *azureConfig {
	obj := &azureConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["availability_set"].(string); ok && len(v) > 0 {
		obj.AvailabilitySet = v
	}

	if v, ok := in["client_id"].(string); ok && len(v) > 0 {
		obj.ClientID = v
	}

	if v, ok := in["client_secret"].(string); ok && len(v) > 0 {
		obj.ClientSecret = v
	}

	if v, ok := in["custom_data"].(string); ok && len(v) > 0 {
		obj.CustomData = v
	}

	if v, ok := in["dns"].(string); ok && len(v) > 0 {
		obj.DNS = v
	}

	if v, ok := in["environment"].(string); ok && len(v) > 0 {
		obj.Environment = v
	}

	if v, ok := in["image"].(string); ok && len(v) > 0 {
		obj.Image = v
	}

	if v, ok := in["location"].(string); ok && len(v) > 0 {
		obj.Location = v
	}

	if v, ok := in["no_public_ip"].(bool); ok {
		obj.NoPublicIP = v
	}

	if v, ok := in["open_port"].([]interface{}); ok && len(v) > 0 {
		obj.OpenPort = toArrayString(v)
	}

	if v, ok := in["private_address_only"].(bool); ok {
		obj.PrivateAddressOnly = v
	}

	if v, ok := in["private_ip_address"].(string); ok && len(v) > 0 {
		obj.PrivateIPAddress = v
	}

	if v, ok := in["resource_group"].(string); ok && len(v) > 0 {
		obj.ResourceGroup = v
	}

	if v, ok := in["size"].(string); ok && len(v) > 0 {
		obj.Size = v
	}

	if v, ok := in["ssh_user"].(string); ok && len(v) > 0 {
		obj.SSHUser = v
	}

	if v, ok := in["static_public_ip"].(bool); ok {
		obj.StaticPublicIP = v
	}

	if v, ok := in["storage_type"].(string); ok && len(v) > 0 {
		obj.StorageType = v
	}

	if v, ok := in["subnet"].(string); ok && len(v) > 0 {
		obj.Subnet = v
	}

	if v, ok := in["subnet_prefix"].(string); ok && len(v) > 0 {
		obj.SubnetPrefix = v
	}

	if v, ok := in["subscription_id"].(string); ok && len(v) > 0 {
		obj.SubscriptionID = v
	}

	if v, ok := in["use_private_ip"].(bool); ok {
		obj.UsePrivateIP = v
	}

	if v, ok := in["vnet"].(string); ok && len(v) > 0 {
		obj.Vnet = v
	}

	return obj
}
