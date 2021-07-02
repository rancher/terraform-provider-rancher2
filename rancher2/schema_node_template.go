package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

//Types

type NodeTemplate struct {
	managementClient.NodeTemplate
	Amazonec2Config     *amazonec2Config     `json:"amazonec2Config,omitempty" yaml:"amazonec2Config,omitempty"`
	AzureConfig         *azureConfig         `json:"azureConfig,omitempty" yaml:"azureConfig,omitempty"`
	DigitaloceanConfig  *digitaloceanConfig  `json:"digitaloceanConfig,omitempty" yaml:"digitaloceanConfig,omitempty"`
	LinodeConfig        *linodeConfig        `json:"linodeConfig,omitempty" yaml:"linodeConfig,omitempty"`
	OpenstackConfig     *openstackConfig     `json:"openstackConfig,omitempty" yaml:"openstackConfig,omitempty"`
	VmwarevsphereConfig *vmwarevsphereConfig `json:"vmwarevsphereConfig,omitempty" yaml:"vmwarevsphereConfig,omitempty"`
	OpennebulaConfig    *opennebulaConfig    `json:"opennebulaConfig,omitempty" yaml:"opennebulaConfig,omitempty"`
	HetznerConfig       *hetznerConfig       `json:"hetznerConfig,omitempty" yaml:"hetznerConfig,omitempty"`
}

//Schemas

func nodeTemplateFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"amazonec2_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"azure_config", "digitalocean_config", "opennebula_config", "openstack_config", "hetzner_config", "vsphere_config"},
			Elem: &schema.Resource{
				Schema: amazonec2ConfigFields(),
			},
		},
		"auth_certificate_authority": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"auth_key": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"azure_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"amazonec2_config", "digitalocean_config", "opennebula_config", "openstack_config", "hetzner_config", "vsphere_config"},
			Elem: &schema.Resource{
				Schema: azureConfigFields(),
			},
		},
		"cloud_credential_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"description": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"digitalocean_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"amazonec2_config", "azure_config", "opennebula_config", "openstack_config", "hetzner_config", "vsphere_config"},
			Elem: &schema.Resource{
				Schema: digitaloceanConfigFields(),
			},
		},
		"driver": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"driver_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"engine_env": {
			Type:     schema.TypeMap,
			Optional: true,
		},
		"engine_insecure_registry": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"engine_install_url": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"engine_label": {
			Type:     schema.TypeMap,
			Optional: true,
		},
		"engine_opt": {
			Type:     schema.TypeMap,
			Optional: true,
		},
		"engine_registry_mirror": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"engine_storage_driver": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"linode_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"amazonec2_config", "azure_config", "digitalocean_config", "opennebula_config", "openstack_config", "hetzner_config", "vsphere_config"},
			Elem: &schema.Resource{
				Schema: linodeConfigFields(),
			},
		},
		"node_taints": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Resource{
				Schema: taintFields(),
			},
		},
		"openstack_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"amazonec2_config", "azure_config", "digitalocean_config", "opennebula_config", "hetzner_config", "vsphere_config"},
			Elem: &schema.Resource{
				Schema: openstackConfigFields(),
			},
		},
		"hetzner_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"amazonec2_config", "azure_config", "digitalocean_config", "opennebula_config", "openstack_config", "vsphere_config"},
			Elem: &schema.Resource{
				Schema: hetznerConfigFields(),
			},
		},
		"use_internal_ip_address": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"vsphere_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"amazonec2_config", "azure_config", "digitalocean_config", "opennebula_config", "hetzner_config", "openstack_config"},
			Elem: &schema.Resource{
				Schema: vsphereConfigFields(),
			},
		},
		"opennebula_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"amazonec2_config", "azure_config", "digitalocean_config", "openstack_config", "vsphere_config"},
			Elem: &schema.Resource{
				Schema: opennebulaConfigFields(),
			},
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
