package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func machineConfigV2Amazonec2Fields() map[string]*schema.Schema {
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
			Default:     "t3a.medium",
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
		"ssh_key_contents": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "SSH Key file contents for sshKeyContents",
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
