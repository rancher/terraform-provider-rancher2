package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func clusterRKEConfigServicesKubeControllerFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_cidr": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"extra_args": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"extra_binds": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"extra_env": {
			Type:     schema.TypeList,
			Optional: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"image": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"service_cluster_ip_range": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
	return s
}
