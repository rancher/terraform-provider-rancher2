package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	outscaleConfigDriver = "outscale"
)

//Types

type outscaleConfig struct {
	AccessKey          string   `json:"accessKey,omitempty" yaml:"accessKey,omitempty"`
	SourceOmi          string   `json:"sourceOmi,omitempty" yaml:"sourceOmi,omitempty"`
	Region             string   `json:"region,omitempty" yaml:"region,omitempty"`
	SecretKey          string   `json:"secretKey,omitempty" yaml:"secretKey,omitempty"`
	SecurityGroupIds   []string `json:"securityGroupIds,omitempty" yaml:"securityGroupIds,omitempty"`
	ExtraTagsAll       []string `json:"extraTagsAll,omitempty" yaml:"extraTagsAll,omitempty"`
	ExtraTagsInstances []string `json:"extraTagsInstances,omitempty" yaml:"extraTagsInstances,omitempty"`
	InstanceType       string   `json:"instanceType,omitempty" yaml:"instanceType,omitempty"`
	RootDiskType       string   `json:"rootDiskType,omitempty" yaml:"rootDiskType,omitempty"`
	RootDiskSize       int      `json:"rootDiskSize,string,omitempty" yaml:"rootDiskSize,omitempty"`
	RootDiskIops       int      `json:"rootDiskIops,string,omitempty" yaml:"rootDiskIops,omitempty"`
}

//Schemas

func outscaleConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"source_omi": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "ami-2cf1fa3e",
			Description: "Outscale Machine Image to use as bootstrap for the VM",
		},
		"region": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "eu-west-2",
			Description: "Outscale Region",
		},
		"security_group_ids": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Ids of user defined Security Groups to add to the machine",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"access_key": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Outscale Access Key",
		},
		"instance_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "tinav2.c1r2p3",
			Description: "Outscale VM type",
		},
		"secret_key": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Outscale Secret Key",
		},
		"extra_tags_all": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Extra tags for all created resources (e.g. key1=value1,key2=value2)",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"extra_tags_instances": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Extra tags only for instances (e.g. key1=value1,key2=value2)",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"root_disk_type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Type of the Root Disk. Possible values are :'standard', 'gp2' or 'io1'.",
		},
		"root_disk_size": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Size of the Root Disk (in GB). From 1 to 14901.",
		},
		"root_disk_iops": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Iops for io1 Root Disk. From 1 to 13000.",
		},
	}

	return s
}
