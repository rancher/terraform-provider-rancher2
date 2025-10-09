package rancher2

import (
	"maps"

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
			ForceNew:    true,
			Description: "Indicates whether downstream proxy requests for service account tokens is enabled",
		},
	}

	maps.Copy(s, commonAnnotationLabelFields())

	return s
}
