package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	clusterClient "github.com/rancher/types/client/cluster/v3"
)

//Schemas

func clusterNamespaceResourceQuotaFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"limit": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &schema.Resource{
				Schema: clusterResourceQuotaLimitFields(),
			},
		},
	}

	return s
}

// Flatteners

func flattenClusterNamespaceResourceQuota(in *clusterClient.NamespaceResourceQuota) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if in.Limit != nil {
		limit, err := flattenClusterResourceQuotaLimit(in.Limit)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["limit"] = limit
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandClusterNamespaceResourceQuota(p []interface{}) (*clusterClient.NamespaceResourceQuota, error) {
	obj := &clusterClient.NamespaceResourceQuota{}

	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["limit"].([]interface{}); ok && len(v) > 0 {
		limit, err := expandClusterResourceQuotaLimit(v)
		if err != nil {
			return obj, err
		}
		obj.Limit = limit
	}

	return obj, nil
}
