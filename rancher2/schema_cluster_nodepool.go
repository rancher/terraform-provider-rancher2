package rancher2

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

type BaseNodePool struct {
	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Taints []K8sTaint        `json:"taints,omitempty" yaml:"taints,omitempty"`
	Name   string            `json:"name,omitempty" yaml:"name,omitempty"`

	// deprecated fields
	AddDefaultLabel  bool              `json:"addDefaultLabel,omitempty" yaml:"addDefaultLabel,omitempty"`
	AddDefaultTaint  bool              `json:"addDefaultTaint,omitempty" yaml:"addDefaultTaint,omitempty"`
	AdditionalLabels map[string]string `json:"additionalLabels,omitempty" yaml:"additionalLabels,omitempty"`
	AdditionalTaints []K8sTaint        `json:"additionalTaints,omitempty" yaml:"additionalTaints,omitempty"`
}

type K8sTaint struct {
	Effect string `json:"effect,omitempty" yaml:"effect,omitempty"`
	Key    string `json:"key,omitempty" yaml:"key,omitempty"`
	Value  string `json:"value,omitempty" yaml:"value,omitempty"`
}

func newNodePoolSchema(input map[string]*schema.Schema) map[string]*schema.Schema {
	input["labels"] = &schema.Schema{
		Type:        schema.TypeMap,
		Optional:    true,
		Description: "List of k8s labels",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}

	input["taints"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "List of k8s taints",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"effect": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "",
				},
				"key": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "",
				},
				"value": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "",
				},
			},
		},
	}

	input["name"] = &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Name of the worker pool",
	}

	// deprecated fields
	input["add_default_label"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Deprecated:  "Support for default label has been deprecated",
		Description: "Adds default label for EKS worker nodes",
	}

	input["add_default_taint"] = &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Deprecated:  "Support for default taint has been deprecated",
		Description: "Adds default taint for EKS worker nodes",
	}

	input["additional_labels"] = &schema.Schema{
		Type:        schema.TypeMap,
		Optional:    true,
		Deprecated:  "Use 'labels' field instead",
		Description: "List of additional labels",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}

	input["additional_taints"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Deprecated:  "Use 'taints' field instead",
		Description: "List of additional taints",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"effect": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "",
				},
				"key": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "",
				},
				"value": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "",
				},
			},
		},
	}

	return input
}
