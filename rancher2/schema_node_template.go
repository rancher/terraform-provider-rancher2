package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	amazonec2ConfigDriver    = "amazonec2"
	azureConfigDriver        = "azure"
	digitaloceanConfigDriver = "digitalocean"
)

//Types

type amazonec2Config struct {
	AccessKey               string   `json:"accessKey,omitempty" yaml:"accessKey,omitempty"`
	Ami                     string   `json:"ami,omitempty" yaml:"ami,omitempty"`
	BlockDurationMinutes    string   `json:"blockDurationMinutes,omitempty" yaml:"blockDurationMinutes,omitempty"`
	DeviceName              string   `json:"deviceName,omitempty" yaml:"deviceName,omitempty"`
	Endpoint                string   `json:"endpoint,omitempty" yaml:"endpoint,omitempty"`
	IamInstanceProfile      string   `json:"iamInstanceProfile,omitempty" yaml:"iamInstanceProfile,omitempty"`
	InsecureTransport       bool     `json:"insecureTransport,omitempty" yaml:"insecureTransport,omitempty"`
	InstanceType            string   `json:"instanceType,omitempty" yaml:"instanceType,omitempty"`
	KeypairName             string   `json:"keypairName,omitempty" yaml:"keypairName,omitempty"`
	Monitoring              bool     `json:"monitoring,omitempty" yaml:"monitoring,omitempty"`
	OpenPort                []string `json:"openPort,omitempty" yaml:"openPort,omitempty"`
	PrivateAddressOnly      bool     `json:"privateAddressOnly,omitempty" yaml:"privateAddressOnly,omitempty"`
	Region                  string   `json:"region,omitempty" yaml:"region,omitempty"`
	RequestSpotInstance     bool     `json:"requestSpotInstance,omitempty" yaml:"requestSpotInstance,omitempty"`
	Retries                 string   `json:"retries,omitempty" yaml:"retries,omitempty"`
	RootSize                string   `json:"rootSize,omitempty" yaml:"rootSize,omitempty"`
	SecretKey               string   `json:"secretKey,omitempty" yaml:"secretKey,omitempty"`
	SecurityGroup           []string `json:"securityGroup,omitempty" yaml:"securityGroup,omitempty"`
	SecurityGroupReadonly   bool     `json:"securityGroupReadonly,omitempty" yaml:"securityGroupReadonly,omitempty"`
	SessionToken            string   `json:"sessionToken,omitempty" yaml:"sessionToken,omitempty"`
	SpotPrice               string   `json:"spotPrice,omitempty" yaml:"spotPrice,omitempty"`
	SSHKeypath              string   `json:"sshKeypath,omitempty" yaml:"sshKeypath,omitempty"`
	SSHUser                 string   `json:"sshUser,omitempty" yaml:"sshUser,omitempty"`
	SubnetID                string   `json:"subnetId,omitempty" yaml:"subnetId,omitempty"`
	Tags                    string   `json:"tags,omitempty" yaml:"tags,omitempty"`
	UseEbsOptimizedInstance bool     `json:"useEbsOptimizedInstance,omitempty" yaml:"useEbsOptimizedInstance,omitempty"`
	UsePrivateAddress       bool     `json:"usePrivateAddress,omitempty" yaml:"usePrivateAddress,omitempty"`
	Userdata                string   `json:"userdata,omitempty" yaml:"userdata,omitempty"`
	VolumeType              string   `json:"volumeType,omitempty" yaml:"volumeType,omitempty"`
	VpcID                   string   `json:"vpcId,omitempty" yaml:"vpcId,omitempty"`
	Zone                    string   `json:"zone,omitempty" yaml:"zone,omitempty"`
}

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

type digitaloceanConfig struct {
	AccessToken       string `json:"accessToken,omitempty" yaml:"accessToken,omitempty"`
	Backups           bool   `json:"backups,omitempty" yaml:"backups,omitempty"`
	Image             string `json:"image,omitempty" yaml:"image,omitempty"`
	IPV6              bool   `json:"ipv6,omitempty" yaml:"ipv6,omitempty"`
	Monitoring        bool   `json:"monitoring,omitempty" yaml:"monitoring,omitempty"`
	PrivateNetworking bool   `json:"privateNetworking,omitempty" yaml:"privateNetworking,omitempty"`
	Region            string `json:"region,omitempty" yaml:"region,omitempty"`
	Size              string `json:"size,omitempty" yaml:"size,omitempty"`
	SSHKeyFingerprint string `json:"sshKeyFingerprint,omitempty" yaml:"sshKeyFingerprint,omitempty"`
	SSHKeyPath        string `json:"sshKeyPath,omitempty" yaml:"sshKeyPath,omitempty"`
	SSHPort           string `json:"sshPort,omitempty" yaml:"sshPort,omitempty"`
	SSHUser           string `json:"sshUser,omitempty" yaml:"sshUser,omitempty"`
	Tags              string `json:"tags,omitempty" yaml:"tags,omitempty"`
	Userdata          string `json:"userdata,omitempty" yaml:"userdata,omitempty"`
}

type NodeTemplate struct {
	managementClient.NodeTemplate
	Amazonec2Config    *amazonec2Config    `json:"amazonec2Config,omitempty" yaml:"amazonec2Config,omitempty"`
	AzureConfig        *azureConfig        `json:"azureConfig,omitempty" yaml:"azureConfig,omitempty"`
	DigitaloceanConfig *digitaloceanConfig `json:"digitaloceanConfig,omitempty" yaml:"digitaloceanConfig,omitempty"`
}

//Schemas

func amazonec2ConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"access_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "AWS Access Key",
		},
		"ami": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "AWS machine image",
		},
		"block_duration_minutes": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "0",
			Description: "AWS spot instance duration in minutes (60, 120, 180, 240, 300, or 360)",
		},
		"device_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "/dev/sda1",
			Description: "AWS root device name",
		},
		"endpoint": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Optional endpoint URL (hostname only or fully qualified URI)",
		},
		"iam_instance_profile": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "AWS IAM Instance Profile",
		},
		"insecure_transport": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Disable SSL when sending requests",
		},
		"instance_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "t2.micro",
			Description: "AWS instance type",
		},
		"keypair_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "AWS keypair to use; requires --amazonec2-ssh-keypath",
		},
		"monitoring": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Set this flag to enable CloudWatch monitoring",
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
		"region": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "AWS Region",
		},
		"request_spot_instance": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Set this flag to request spot instance",
		},
		"retries": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "5",
			Description: "Set retry count for recoverable failures (use -1 to disable)",
		},
		"root_size": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "16",
			Description: "AWS root disk size (in GB)",
		},
		"secret_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "AWS Secret Key",
		},
		"security_group": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "AWS VPC security group",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"security_group_readonly": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Skip adding default rules to security groups",
		},
		"session_token": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "AWS Session Token",
		},
		"spot_price": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "0.50",
			Description: "AWS spot instance bid price (in dollar)",
		},
		"ssh_keypath": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "SSH Key for Instance",
		},
		"ssh_user": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "ubuntu",
			Description: "Set the name of the ssh user",
		},
		"subnet_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "AWS VPC subnet id",
		},
		"tags": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "AWS Tags (e.g. key1,value1,key2,value2)",
		},
		"use_ebs_optimized_instance": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Create an EBS optimized instance",
		},
		"use_private_address": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Force the usage of private IP address",
		},
		"userdata": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Path to file with cloud-init user data",
		},
		"volume_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "gp2",
			Description: "Amazon EBS volume type",
		},
		"vpc_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "AWS VPC id",
		},
		"zone": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "AWS zone for instance (i.e. a,b,c,d,e)",
		},
	}

	return s
}

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

func digitaloceanConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"access_token": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Digital Ocean access token",
		},
		"backups": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable backups for droplet",
		},
		"image": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "ubuntu-16-04-x64",
			Description: "Digital Ocean Image",
		},
		"ipv6": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable ipv6 for droplet",
		},
		"monitoring": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable monitoring for droplet",
		},
		"private_networking": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable private networking for droplet",
		},
		"region": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "nyc3",
			Description: "Digital Ocean region",
		},
		"size": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "s-1vcpu-1gb",
			Description: "Digital Ocean size",
		},
		"ssh_key_fingerprint": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "SSH key fingerprint",
		},
		"ssh_key_path": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "SSH private key path",
		},
		"ssh_port": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "22",
			Description: "SSH port",
		},
		"ssh_user": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "root",
			Description: "SSH username",
		},
		"tags": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Comma-separated list of tags to apply to the Droplet",
		},
		"userdata": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "docker-user",
			Description: "Path to file with cloud-init user-data",
		},
	}

	return s
}

func nodeTemplateFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"amazonec2_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"azure_config", "digitalocean_config"},
			Elem: &schema.Resource{
				Schema: amazonec2ConfigFields(),
			},
		},
		"auth_certificate_authority": &schema.Schema{
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"auth_key": &schema.Schema{
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"azure_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"amazonec2_config", "digitalocean_config"},
			Elem: &schema.Resource{
				Schema: azureConfigFields(),
			},
		},
		"cloud_credential_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"digitalocean_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"amazonec2_config", "azure_config"},
			Elem: &schema.Resource{
				Schema: digitaloceanConfigFields(),
			},
		},
		"docker_version": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"driver": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"engine_env": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
		},
		"engine_insecure_registry": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"engine_install_url": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"engine_label": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
		},
		"engine_opt": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
		},
		"engine_registry_mirror": &schema.Schema{
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"engine_storage_driver": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"use_internal_ip_address": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"annotations": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
		},
		"labels": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
	}

	return s
}
