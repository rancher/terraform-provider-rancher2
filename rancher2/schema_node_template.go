package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

//Types

type NodeTemplate struct {
	managementClient.NodeTemplate
	Amazonec2Config     *amazonec2Config     `json:"amazonec2Config,omitempty" yaml:"amazonec2Config,omitempty"`
	AzureConfig         *azureConfig         `json:"azureConfig,omitempty" yaml:"azureConfig,omitempty"`
	DigitaloceanConfig  *digitaloceanConfig  `json:"digitaloceanConfig,omitempty" yaml:"digitaloceanConfig,omitempty"`
	OpenstackConfig     *openstackConfig     `json:"openstackConfig,omitempty" yaml:"openstackConfig,omitempty"`
	VmwarevsphereConfig *vmwarevsphereConfig `json:"vmwarevsphereConfig,omitempty" yaml:"vmwarevsphereConfig,omitempty"`
}

//Schemas

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
			ConflictsWith: []string{"azure_config", "digitalocean_config", "openstack_config", "vsphere_config"},
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
			ConflictsWith: []string{"amazonec2_config", "digitalocean_config", "openstack_config", "vsphere_config"},
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
			ConflictsWith: []string{"amazonec2_config", "azure_config", "openstack_config", "vsphere_config"},
			Elem: &schema.Resource{
				Schema: digitaloceanConfigFields(),
			},
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
			Default:  "https://releases.rancher.com/install-docker/18.09.sh",
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
		"openstack_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"amazonec2_config", "azure_config", "digitalocean_config", "vsphere_config"},
			Elem: &schema.Resource{
				Schema: openstackConfigFields(),
			},
		},
		"use_internal_ip_address": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"vsphere_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"amazonec2_config", "azure_config", "digitalocean_config", "openstack_config"},
			Elem: &schema.Resource{
				Schema: vsphereConfigFields(),
			},
		},
		"annotations": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"labels": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
	}

	return s
}
