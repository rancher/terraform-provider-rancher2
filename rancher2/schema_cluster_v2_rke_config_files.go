package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Types

func clusterV2RKEConfigKeyToPathFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"key": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The key of the item(file) to retrieve",
		},
		"path": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The path to put the file in the target node",
		},
		"dynamic": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "If ture, the file is ignored when determining whether the node should be drained before updating the node plan (default: true).",
		},
		"permissions": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The numeric representation of the file permissions",
		},
		"hash": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The base64 encoded value of the SHA256 checksum of the file's content",
		},
	}

	return s
}

func clusterV2RKEConfigK8sObjectFileSourceFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The name of the K8s object",
		},
		"items": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Items(files) to retrieve from the K8s object",
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigKeyToPathFields(),
			},
		},
		"default_permissions": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "The default permissions to be applied when they are not set at the item level",
		},
	}

	return s
}
func clusterV2RKEConfigFileSourceFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"secret": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "The secret which is the source of files",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigK8sObjectFileSourceFields(),
			},
		},
		"configmap": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "The configmap which is the source of files",
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigK8sObjectFileSourceFields(),
			},
		},
	}

	return s
}

func clusterV2RKEConfigMachineSelectorFilesFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"machine_label_selector": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Description: "Machine label selector",
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigSystemConfigLabelSelectorFields(),
			},
		},
		"file_sources": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "File sources",
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigFileSourceFields(),
			},
		},
	}

	return s
}
