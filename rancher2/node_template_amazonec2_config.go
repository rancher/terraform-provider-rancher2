package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
)

const (
	amazonec2ConfigDriver = "amazonec2"
)

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

//Schemas

func amazonec2ConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"access_key": {
			Type:        schema.TypeString,
			Required:    true,
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
			Required:    true,
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

// Flatteners

func flattenAmazonec2Config(in *amazonec2Config) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	if len(in.AccessKey) > 0 {
		obj["access_key"] = in.AccessKey
	}

	if len(in.Ami) > 0 {
		obj["ami"] = in.Ami
	}

	if len(in.BlockDurationMinutes) > 0 {
		obj["block_duration_minutes"] = in.BlockDurationMinutes
	}

	if len(in.DeviceName) > 0 {
		obj["device_name"] = in.DeviceName
	}

	if len(in.Endpoint) > 0 {
		obj["endpoint"] = in.Endpoint
	}

	if len(in.IamInstanceProfile) > 0 {
		obj["iam_instance_profile"] = in.IamInstanceProfile
	}

	obj["insecure_transport"] = in.InsecureTransport

	if len(in.InstanceType) > 0 {
		obj["instance_type"] = in.InstanceType
	}

	if len(in.KeypairName) > 0 {
		obj["keypair_name"] = in.KeypairName
	}

	obj["monitoring"] = in.Monitoring

	if len(in.OpenPort) > 0 {
		obj["open_port"] = toArrayInterface(in.OpenPort)
	}

	obj["private_address_only"] = in.PrivateAddressOnly

	if len(in.Region) > 0 {
		obj["region"] = in.Region
	}

	obj["request_spot_instance"] = in.RequestSpotInstance

	if len(in.Retries) > 0 {
		obj["retries"] = in.Retries
	}

	if len(in.RootSize) > 0 {
		obj["root_size"] = in.RootSize
	}

	if len(in.SecretKey) > 0 {
		obj["secret_key"] = in.SecretKey
	}

	if len(in.SecurityGroup) > 0 {
		obj["security_group"] = toArrayInterface(in.SecurityGroup)
	}

	obj["security_group_readonly"] = in.SecurityGroupReadonly

	if len(in.SessionToken) > 0 {
		obj["session_token"] = in.SessionToken
	}

	if len(in.SpotPrice) > 0 {
		obj["spot_price"] = in.SpotPrice
	}

	if len(in.SSHKeypath) > 0 {
		obj["ssh_keypath"] = in.SSHKeypath
	}

	if len(in.SSHUser) > 0 {
		obj["ssh_user"] = in.SSHUser
	}

	if len(in.SubnetID) > 0 {
		obj["subnet_id"] = in.SubnetID
	}

	if len(in.Tags) > 0 {
		obj["tags"] = in.Tags
	}

	obj["use_ebs_optimized_instance"] = in.UseEbsOptimizedInstance

	obj["use_private_address"] = in.UsePrivateAddress

	if len(in.Userdata) > 0 {
		obj["userdata"] = in.Userdata
	}

	if len(in.VolumeType) > 0 {
		obj["volume_type"] = in.VolumeType
	}

	if len(in.VpcID) > 0 {
		obj["vpc_id"] = in.VpcID
	}

	if len(in.Zone) > 0 {
		obj["zone"] = in.Zone
	}

	return []interface{}{obj}
}

// Expanders

func expandAmazonec2Config(p []interface{}) *amazonec2Config {
	obj := &amazonec2Config{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["access_key"].(string); ok && len(v) > 0 {
		obj.AccessKey = v
	}

	if v, ok := in["ami"].(string); ok && len(v) > 0 {
		obj.Ami = v
	}

	if v, ok := in["block_duration_minutes"].(string); ok && len(v) > 0 {
		obj.BlockDurationMinutes = v
	}

	if v, ok := in["device_name"].(string); ok && len(v) > 0 {
		obj.DeviceName = v
	}

	if v, ok := in["endpoint"].(string); ok && len(v) > 0 {
		obj.Endpoint = v
	}

	if v, ok := in["iam_instance_profile"].(string); ok && len(v) > 0 {
		obj.IamInstanceProfile = v
	}

	if v, ok := in["insecure_transport"].(bool); ok {
		obj.InsecureTransport = v
	}

	if v, ok := in["instance_type"].(string); ok && len(v) > 0 {
		obj.InstanceType = v
	}

	if v, ok := in["keypair_name"].(string); ok && len(v) > 0 {
		obj.KeypairName = v
	}

	if v, ok := in["monitoring"].(bool); ok {
		obj.Monitoring = v
	}

	if v, ok := in["open_port"].([]interface{}); ok && len(v) > 0 {
		obj.OpenPort = toArrayString(v)
	}

	if v, ok := in["private_address_only"].(bool); ok {
		obj.PrivateAddressOnly = v
	}

	if v, ok := in["region"].(string); ok && len(v) > 0 {
		obj.Region = v
	}

	if v, ok := in["request_spot_instance"].(bool); ok {
		obj.RequestSpotInstance = v
	}

	if v, ok := in["retries"].(string); ok && len(v) > 0 {
		obj.Retries = v
	}

	if v, ok := in["root_size"].(string); ok && len(v) > 0 {
		obj.RootSize = v
	}

	if v, ok := in["secret_key"].(string); ok && len(v) > 0 {
		obj.SecretKey = v
	}

	if v, ok := in["security_group"].([]interface{}); ok && len(v) > 0 {
		obj.SecurityGroup = toArrayString(v)
	}

	if v, ok := in["security_group_readonly"].(bool); ok {
		obj.SecurityGroupReadonly = v
	}

	if v, ok := in["session_token"].(string); ok && len(v) > 0 {
		obj.SessionToken = v
	}

	if v, ok := in["spot_price"].(string); ok && len(v) > 0 {
		obj.SpotPrice = v
	}

	if v, ok := in["ssh_keypath"].(string); ok && len(v) > 0 {
		obj.SSHKeypath = v
	}

	if v, ok := in["ssh_user"].(string); ok && len(v) > 0 {
		obj.SSHUser = v
	}

	if v, ok := in["subnet_id"].(string); ok && len(v) > 0 {
		obj.SubnetID = v
	}

	if v, ok := in["tags"].(string); ok && len(v) > 0 {
		obj.Tags = v
	}

	if v, ok := in["use_ebs_optimized_instance"].(bool); ok {
		obj.UseEbsOptimizedInstance = v
	}

	if v, ok := in["use_private_address"].(bool); ok {
		obj.UsePrivateAddress = v
	}

	if v, ok := in["userdata"].(string); ok && len(v) > 0 {
		obj.Userdata = v
	}

	if v, ok := in["volume_type"].(string); ok && len(v) > 0 {
		obj.VolumeType = v
	}

	if v, ok := in["vpc_id"].(string); ok && len(v) > 0 {
		obj.VpcID = v
	}

	if v, ok := in["zone"].(string); ok && len(v) > 0 {
		obj.Zone = v
	}

	return obj
}
