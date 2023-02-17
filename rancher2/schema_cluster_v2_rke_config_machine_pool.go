package rancher2

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

//Types

func clusterV2RKEConfigMachinePoolMachineConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"kind": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Machine config kind",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Machine config name",
		},
	}

	return s
}

func clusterV2RKEConfigMachinePoolRollingUpdateFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"max_unavailable": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Rolling update max unavailable",
		},
		"max_surge": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Rolling update max surge",
		},
	}

	return s
}

func clusterV2RKEConfigMachinePoolFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Machine pool name",
		},
		"cloud_credential_secret_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Machine pool cloud credential secret name",
		},
		"machine_config": {
			Type:        schema.TypeList,
			Required:    true,
			MaxItems:    1,
			Description: "Machine config data",
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigMachinePoolMachineConfigFields(),
			},
		},
		"control_plane_role": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Machine pool control plane role",
		},
		"etcd_role": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Machine pool etcd role",
		},
		"drain_before_delete": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Machine pool drain before delete",
		},
		"node_drain_timeout": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     0,
			Description: "seconds to wait for machine pool drain to complete before machine deletion",
		},
		"paused": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Machine pool paused",
		},
		"quantity": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      1,
			ValidateFunc: validation.IntAtLeast(1),
			Description:  "Machine pool quantity",
		},
		"rolling_update": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Machine pool rolling update",
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigMachinePoolRollingUpdateFields(),
			},
		},
		"taints": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Machine pool taints",
			Elem: &schema.Resource{
				Schema: taintV2Fields(),
			},
		},
		"worker_role": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Machine pool worker role",
		},
		"node_startup_timeout_seconds": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "seconds a new node has to become active before it is replaced",
		},
		"unhealthy_node_timeout_seconds": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "seconds an unhealthy node has to become active before it is replaced",
		},
		"max_unhealthy": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "max unhealthy nodes for automated replacement to be allowed",
		},
		"unhealthy_range": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "range of unhealthy nodes for automated replacement to be allowed",
		},
		"machine_labels": {
			Type:        schema.TypeMap,
			Optional:    true,
			Computed:    true,
			Description: "Labels of the machine",
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				// Supressing diff for labels containing cattle.io/
				if (strings.Contains(k, commonAnnotationLabelCattle) || strings.Contains(k, commonAnnotationLabelRancher)) && new == "" {
					return true
				}
				return false
			},
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
