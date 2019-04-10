package rancher2

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
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
	genericConfig       *genericNodeTemplateConfig
}

func isTypedNodeDriver(driverName string) bool {
	return driverName == amazonec2ConfigDriver ||
		driverName == azureConfigDriver ||
		driverName == digitaloceanConfigDriver ||
		driverName == openstackConfigDriver ||
		driverName == vmwarevsphereConfigDriver
}

func (n *NodeTemplate) UnmarshalJSON(data []byte) error {
	type Alias NodeTemplate
	var dest Alias
	if err := json.Unmarshal(data, &dest); err != nil {
		return err
	}

	var rawValues map[string]interface{}
	if err := json.Unmarshal(data, &rawValues); err != nil {
		return err
	}

	var driverName string
	if rawDriver, ok := rawValues["driver"]; ok {
		driverName = rawDriver.(string)
	}

	if driverName != "" && !isTypedNodeDriver(driverName) {
		driverConfigName := fmt.Sprintf("%sConfig", driverName)
		if v, ok := rawValues[driverConfigName]; ok {
			if cv, ok := v.(map[string]interface{}); ok {
				dest.genericConfig = &genericNodeTemplateConfig{
					driverName: driverName,
					config:     cv,
				}
			}
		}
	}

	*n = NodeTemplate(dest)
	return nil
}

func (n *NodeTemplate) MarshalJSON() ([]byte, error) {
	type Alias NodeTemplate
	data, err := json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(n),
	})
	if err != nil {
		return nil, err
	}

	var results map[string]interface{}
	if err := json.Unmarshal(data, &results); err != nil {
		return nil, err
	}

	if n.genericConfig != nil {
		driverConfigName := fmt.Sprintf("%sConfig", n.Driver)
		results[driverConfigName] = n.genericConfig.config
	}

	return json.Marshal(results)
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
			ConflictsWith: []string{"azure_config", "digitalocean_config", "generic_config", "openstack_config", "vsphere_config"},
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
			ConflictsWith: []string{"amazonec2_config", "digitalocean_config", "generic_config", "openstack_config", "vsphere_config"},
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
			ConflictsWith: []string{"amazonec2_config", "azure_config", "generic_config", "openstack_config", "vsphere_config"},
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
		"generic_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"amazonec2_config", "azure_config", "digitalocean_config", "vsphere_config", "openstack_config"},
			Elem: &schema.Resource{
				Schema: genericNodeConfigFields(),
			},
		},
		"openstack_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"amazonec2_config", "azure_config", "digitalocean_config", "generic_config", "vsphere_config"},
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
			ConflictsWith: []string{"amazonec2_config", "azure_config", "digitalocean_config", "generic_config", "openstack_config"},
			Elem: &schema.Resource{
				Schema: vsphereConfigFields(),
			},
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
