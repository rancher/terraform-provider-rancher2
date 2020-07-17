package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Schemas
func GlobalDNSEntryFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"fqdn": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"provider_id": {
			Type:     schema.TypeString,
			Required: true,
			ForceNew: true,
		},
		"project_ids": {
			Type:          schema.TypeList,
			Optional:      true,
			ForceNew:      false,
			ConflictsWith: []string{"multi_cluster_app_id"},
			Elem:          &schema.Schema{Type: schema.TypeString},
		},
		"multi_cluster_app_id": {
			Type:          schema.TypeString,
			Optional:      true,
			ForceNew:      false,
			ConflictsWith: []string{"project_ids"},
		},
		"annotations": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"labels": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
	}

	return s
}
