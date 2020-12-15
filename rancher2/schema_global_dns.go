package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Schemas

func GlobalDNSFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"fqdn": {
			Type:     schema.TypeString,
			Required: true,
		},
		"provider_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
			ForceNew: true,
		},
		"multi_cluster_app_id": {
			Type:          schema.TypeString,
			Optional:      true,
			ConflictsWith: []string{"project_ids"},
		},
		"project_ids": {
			Type:          schema.TypeList,
			Optional:      true,
			ConflictsWith: []string{"multi_cluster_app_id"},
			Elem:          &schema.Schema{Type: schema.TypeString},
		},
		"ttl": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  300,
		},
	}
	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
