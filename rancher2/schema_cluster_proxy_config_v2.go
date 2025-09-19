package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func clusterProxyConfigV2Fields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cluster_id": {
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
			Description: "Cluster ID where the ClusterProxyConfig should be created",
		},
		"enabled": {
			Type:        schema.TypeBool,
			Required:    true,
			Description: "Indicates whether downstream proxy requests for service account tokens is enabled",
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}
