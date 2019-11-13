package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func clusterSyncFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_id": &schema.Schema{
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Cluster id to sync",
		},
		"node_pool_ids": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Cluster node pool ids",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"synced": &schema.Schema{
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"default_project_id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"kube_config": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"system_project_id": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	return s
}
