package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	clusterK3SKind   = "k3s"
	clusterDriverK3S = "k3s"
)

//Schemas

func clusterK3SUpgradeStrategyConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"drain_server_nodes": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Drain server nodes",
		},
		"drain_worker_nodes": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Drain worker nodes",
		},
		"server_concurrency": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     1,
			Description: "Server concurrency",
		},
		"worker_concurrency": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     1,
			Description: "Worker concurrency",
		},
	}

	return s
}

func clusterK3SConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"upgrade_strategy": {
			Type:        schema.TypeList,
			MaxItems:    1,
			Optional:    true,
			Computed:    true,
			Description: "The K3S upgrade strategy",
			Elem: &schema.Resource{
				Schema: clusterK3SUpgradeStrategyConfigFields(),
			},
		},
		"version": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The K3S kubernetes version",
		},
	}

	return s
}
