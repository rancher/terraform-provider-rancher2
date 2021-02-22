package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

//Schemas

func clusterSyncFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Cluster id to sync",
		},
		"state_confirm": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      1,
			Description:  "Wait until active status is confirmed a number of times (wait interval of 5s)",
			ValidateFunc: validation.IntAtLeast(1),
		},
		"wait_monitoring": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Wait until monitoring is up and running",
		},
		"node_pool_ids": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Cluster node pool ids",
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"nodes": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterNodeFields(),
			},
		},
		"synced": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"default_project_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"kube_config": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},
		"system_project_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}

	return s
}

func clusterNodeFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{

		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"cluster_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"control_plane": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"etcd": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"external_ip_address": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"hostname": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ip_address": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"imported": {
			Type:     schema.TypeBool,
			Computed: true,
		},
		"info": {
			Type:     schema.TypeMap,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"node_name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"node_pool_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"node_template_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"pod_cidr": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"pod_cidrs": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"provider_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"requested_hostname": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ssh_user": {
			Type:      schema.TypeString,
			Computed:  true,
			Sensitive: true,
		},
		"worker": {
			Type:     schema.TypeBool,
			Computed: true,
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
