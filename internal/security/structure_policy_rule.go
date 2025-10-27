package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenPolicyRules(p []managementClient.PolicyRule) []interface{} {
	if len(p) == 0 {
		return []interface{}{}
	}

	out := make([]interface{}, len(p))
	for i, in := range p {
		obj := make(map[string]interface{})

		if len(in.APIGroups) > 0 {
			obj["api_groups"] = toArrayInterface(in.APIGroups)
		}

		if len(in.NonResourceURLs) > 0 {
			obj["non_resource_urls"] = toArrayInterface(in.NonResourceURLs)
		}

		if len(in.ResourceNames) > 0 {
			obj["resource_names"] = toArrayInterface(in.ResourceNames)
		}

		if len(in.Resources) > 0 {
			obj["resources"] = toArrayInterface(in.Resources)
		}

		if len(in.Verbs) > 0 {
			obj["verbs"] = toArrayInterface(in.Verbs)
		}

		out[i] = obj
	}

	return out
}

// Expanders

func expandPolicyRules(p []interface{}) []managementClient.PolicyRule {
	if len(p) == 0 || p[0] == nil {
		return []managementClient.PolicyRule{}
	}

	obj := make([]managementClient.PolicyRule, len(p))

	for i := range p {
		in := p[i].(map[string]interface{})

		if v, ok := in["api_groups"].([]interface{}); ok && len(v) > 0 {
			obj[i].APIGroups = toArrayString(v)
		}

		if v, ok := in["non_resource_urls"].([]interface{}); ok && len(v) > 0 {
			obj[i].NonResourceURLs = toArrayString(v)
		}

		if v, ok := in["resource_names"].([]interface{}); ok && len(v) > 0 {
			obj[i].ResourceNames = toArrayString(v)
		}

		if v, ok := in["resources"].([]interface{}); ok && len(v) > 0 {
			obj[i].Resources = toArrayString(v)
		}

		if v, ok := in["verbs"].([]interface{}); ok && len(v) > 0 {
			obj[i].Verbs = toArrayString(v)
		}
	}

	return obj
}
