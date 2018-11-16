package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

//Schemas

func projectResourceQuotaFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"project_limit": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &schema.Resource{
				Schema: projectResourceQuotaLimitFields(),
			},
		},
		"namespace_default_limit": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &schema.Resource{
				Schema: projectResourceQuotaLimitFields(),
			},
		},
	}

	return s
}

// Flatteners

func flattenProjectResourceQuota(pQuota *managementClient.ProjectResourceQuota, nsQuota *managementClient.NamespaceResourceQuota) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if pQuota == nil || nsQuota == nil {
		return []interface{}{}, nil
	}

	if pQuota.Limit != nil {
		limit, err := flattenProjectResourceQuotaLimit(pQuota.Limit)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["project_limit"] = limit
	}

	if nsQuota.Limit != nil {
		limit, err := flattenProjectResourceQuotaLimit(nsQuota.Limit)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["namespace_default_limit"] = limit
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandProjectResourceQuota(p []interface{}) (*managementClient.ProjectResourceQuota, *managementClient.NamespaceResourceQuota, error) {
	pQuota := &managementClient.ProjectResourceQuota{}
	nsQuota := &managementClient.NamespaceResourceQuota{}

	if len(p) == 0 || p[0] == nil {
		return pQuota, nsQuota, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["project_limit"].([]interface{}); ok && len(v) > 0 {
		pLimit, err := expandProjectResourceQuotaLimit(v)
		if err != nil {
			return nil, nil, err
		}
		pQuota.Limit = pLimit
	}

	if v, ok := in["namespace_default_limit"].([]interface{}); ok && len(v) > 0 {
		nsLimit, err := expandProjectResourceQuotaLimit(v)
		if err != nil {
			return nil, nil, err
		}
		nsQuota.Limit = nsLimit
	}

	return pQuota, nsQuota, nil
}
