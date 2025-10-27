package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Types

func clusterV2RKEConfigUpgradeStrategyDrainOptionsFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Drain options enabled?",
		},
		"force": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Drain options force",
		},
		"ignore_daemon_sets": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Drain options ignore daemon sets",
		},
		"ignore_errors": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Drain options ignore errors",
		},
		"delete_empty_dir_data": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Drain options delete empty dir data",
		},
		"disable_eviction": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Drain options disable eviction",
		},
		"grace_period": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Drain options grace period",
		},
		"timeout": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Drain options timeout",
		},
		"skip_wait_for_delete_timeout_seconds": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Drain options skip wait for delete timeout seconds",
		},
	}

	return s
}

func clusterV2RKEConfigUpgradeStrategyFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"control_plane_concurrency": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "How many controlplane nodes should be upgrade at time, 0 is infinite. Percentages are also accepted",
		},
		"control_plane_drain_options": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Controlplane nodes drain options",
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigUpgradeStrategyDrainOptionsFields(),
			},
		},
		"worker_concurrency": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "How many worker nodes should be upgrade at time",
		},
		"worker_drain_options": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Worker nodes drain options",
			Elem: &schema.Resource{
				Schema: clusterV2RKEConfigUpgradeStrategyDrainOptionsFields(),
			},
		},
	}

	return s
}
