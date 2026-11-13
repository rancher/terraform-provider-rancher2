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
		"wait_catalogs": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Wait until all catalogs are downloaded and active",
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

// clusterSyncFieldsV0 represents the schema used before the removal of the
// RKE1-only node_pool_ids and node_pool_id, node_template_id and ssh_user fields.
func clusterSyncFieldsV0() map[string]*schema.Schema {
	s := clusterSyncFields()

	s["node_pool_ids"] = &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		Description: "Cluster node pool ids",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	}

	s["nodes"] = &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: clusterNodeFieldsV0(),
		},
	}

	return s
}

// clusterNodeFieldsV0 represents the schema used before the removal of the
// RKE1-only node_pool_id, node_template_id and ssh_user fields
func clusterNodeFieldsV0() map[string]*schema.Schema {
	s := clusterNodeFields()

	s["node_pool_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	s["node_template_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	s["ssh_user"] = &schema.Schema{
		Type:      schema.TypeString,
		Computed:  true,
		Sensitive: true,
	}

	return s
}

func clusterNodeFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{

		"capacity": {
			Type:     schema.TypeMap,
			Computed: true,
		},
		"cluster_id": {
			Type:     schema.TypeString,
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
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"ip_address": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"provider_id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"requested_hostname": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"roles": {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"system_info": {
			Type:     schema.TypeMap,
			Computed: true,
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
