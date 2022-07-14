package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	amazonec2ConfigDriver = "amazonec2"
)

//Types

type amazonec2Config struct {
	AccessKey               string   `json:"accessKey,omitempty" yaml:"accessKey,omitempty"`
	Ami                     string   `json:"ami,omitempty" yaml:"ami,omitempty"`
	BlockDurationMinutes    string   `json:"blockDurationMinutes,omitempty" yaml:"blockDurationMinutes,omitempty"`
	DeviceName              string   `json:"deviceName,omitempty" yaml:"deviceName,omitempty"`
	EncryptEbsVolume        bool     `json:"encryptEbsVolume,omitempty" yaml:"encryptEbsVolume,omitempty"`
	Endpoint                string   `json:"endpoint,omitempty" yaml:"endpoint,omitempty"`
	HTTPEndpoint            string   `json:"httpEndpoint,omitempty" yaml:"httpEndpoint,omitempty"`
	HTTPTokens              string   `json:"httpTokens,omitempty" yaml:"httpTokens,omitempty"`
	IamInstanceProfile      string   `json:"iamInstanceProfile,omitempty" yaml:"iamInstanceProfile,omitempty"`
	InsecureTransport       bool     `json:"insecureTransport,omitempty" yaml:"insecureTransport,omitempty"`
	InstanceType            string   `json:"instanceType,omitempty" yaml:"instanceType,omitempty"`
	KeypairName             string   `json:"keypairName,omitempty" yaml:"keypairName,omitempty"`
	KmsKey                  string   `json:"kmsKey,omitempty" yaml:"kmsKey,omitempty"`
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

//Schemas

func amazonec2ConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"ami": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "AWS machine image",
		},
		"region": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "AWS Region",
		},
		"security_group": {
			Type:        schema.TypeList,
			Required:    true,
			Description: "AWS VPC security group",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"subnet_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "AWS VPC subnet id",
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
		"access_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "AWS Access Key",
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
		"encrypt_ebs_volume": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Encrypt EBS volume",
		},
		"endpoint": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Optional endpoint URL (hostname only or fully qualified URI)",
		},
		"http_endpoint": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Enables or disables the HTTP metadata endpoint on your instances",
		},
		"http_tokens": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The state of token usage for your instance metadata requests",
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
		"kms_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Custom KMS key ID using the AWS Managed CMK",
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
	}

	return s
}
