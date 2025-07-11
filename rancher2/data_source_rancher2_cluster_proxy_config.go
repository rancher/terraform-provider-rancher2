package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRancher2ClusterProxyConfig() *schema.Resource {
	return &schema.Resource{
		Read: resourceRancher2ClusterProxyConfigV2Read,
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Computed: true,
				Type:     schema.TypeString,
			},
			"enabled": {
				Computed: true,
				Type:     schema.TypeBool,
			},
		},
	}
}
