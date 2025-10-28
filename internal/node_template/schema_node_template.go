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
	HarvesterConfig     *harvesterConfig     `json:"harvesterConfig,omitempty" yaml:"harvesterConfig,omitempty"`
	HetznerConfig       *hetznerConfig       `json:"hetznerConfig,omitempty" yaml:"hetznerConfig,omitempty"`
	LinodeConfig        *linodeConfig        `json:"linodeConfig,omitempty" yaml:"linodeConfig,omitempty"`
	OpennebulaConfig    *opennebulaConfig    `json:"opennebulaConfig,omitempty" yaml:"opennebulaConfig,omitempty"`
	OpenstackConfig     *openstackConfig     `json:"openstackConfig,omitempty" yaml:"openstackConfig,omitempty"`
	VmwarevsphereConfig *vmwarevsphereConfig `json:"vmwarevsphereConfig,omitempty" yaml:"vmwarevsphereConfig,omitempty"`
	OutscaleConfig      *outscaleConfig      `json:"outscaleConfig,omitempty" yaml:"outscaleConfig,omitempty"`
}

//Schemas

var allNodeTemplateDriverConfigFields = []string{
	"amazonec2_config",
	"azure_config",
	"digitalocean_config",
	"harvester_config",
	"hetzner_config",
	"linode_config",
	"opennebula_config",
	"openstack_config",
	"vsphere_config",
	"outscale_config"}

func getConflicts(fieldNames []string, fieldName string) []string {
	conflicts := make([]string, 0, len(fieldNames)-1)
	for _, name := range fieldNames {
		if name != fieldName {
			conflicts = append(conflicts, name)
		}
	}
	return conflicts
}

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
			ConflictsWith: getConflicts(allNodeTemplateDriverConfigFields, "amazonec2_config"),
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
			ConflictsWith: getConflicts(allNodeTemplateDriverConfigFields, "azure_config"),
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
			ConflictsWith: getConflicts(allNodeTemplateDriverConfigFields, "digitalocean_config"),
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
		"harvester_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: getConflicts(allNodeTemplateDriverConfigFields, "harvester_config"),
			Elem: &schema.Resource{
				Schema: harvesterConfigFields(),
			},
		},
		"hetzner_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: getConflicts(allNodeTemplateDriverConfigFields, "hetzner_config"),
			Elem: &schema.Resource{
				Schema: hetznerConfigFields(),
			},
		},
		"linode_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: getConflicts(allNodeTemplateDriverConfigFields, "linode_config"),
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
		"opennebula_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: getConflicts(allNodeTemplateDriverConfigFields, "opennebula_config"),
			Elem: &schema.Resource{
				Schema: opennebulaConfigFields(),
			},
		},
		"openstack_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: getConflicts(allNodeTemplateDriverConfigFields, "openstack_config"),
			Elem: &schema.Resource{
				Schema: openstackConfigFields(),
			},
		},
		"outscale_config": {
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: getConflicts(allNodeTemplateDriverConfigFields, "outscale_config"),
			Elem: &schema.Resource{
				Schema: outscaleConfigFields(),
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
			ConflictsWith: getConflicts(allNodeTemplateDriverConfigFields, "vsphere_config"),
			Elem: &schema.Resource{
				Schema: vsphereConfigFields(),
			},
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
