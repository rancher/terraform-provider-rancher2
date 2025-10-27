package rancher2

import (
	"sort"

	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
)

// Flatteners

func flattenClusterV2RKEConfigRegistryConfigs(p map[string]rkev1.RegistryConfig) []interface{} {
	if p == nil {
		return nil
	}
	sorted := make([]string, len(p))
	i := 0
	for k := range p {
		sorted[i] = k
		i++
	}
	sort.Strings(sorted)
	out := make([]interface{}, len(p))
	for i, k := range sorted {
		in := p[k]
		obj := map[string]interface{}{}

		obj["hostname"] = k

		if len(in.AuthConfigSecretName) > 0 {
			obj["auth_config_secret_name"] = in.AuthConfigSecretName
		}
		if len(in.TLSSecretName) > 0 {
			obj["tls_secret_name"] = in.TLSSecretName
		}
		if len(in.CABundle) > 0 {
			obj["ca_bundle"] = string(in.CABundle)
		}
		obj["insecure"] = in.InsecureSkipVerify
		out[i] = obj
	}

	return out
}

func flattenClusterV2RKEConfigRegistryMirrors(p map[string]rkev1.Mirror) []interface{} {
	if p == nil {
		return nil
	}
	sorted := make([]string, len(p))
	i := 0
	for k := range p {
		sorted[i] = k
		i++
	}
	sort.Strings(sorted)
	out := make([]interface{}, len(p))
	for i, k := range sorted {
		in := p[k]
		obj := map[string]interface{}{}

		obj["hostname"] = k

		if len(in.Endpoints) > 0 {
			obj["endpoints"] = toArrayInterface(in.Endpoints)
		}
		if len(in.Rewrites) > 0 {
			obj["rewrites"] = toMapInterface(in.Rewrites)
		}
		out[i] = obj
	}

	return out
}

func flattenClusterV2RKEConfigRegistry(in *rkev1.Registry) []interface{} {
	if in == nil {
		return nil
	}

	obj := make(map[string]interface{})
	if len(in.Configs) > 0 {
		obj["configs"] = flattenClusterV2RKEConfigRegistryConfigs(in.Configs)
	}
	if len(in.Mirrors) > 0 {
		obj["mirrors"] = flattenClusterV2RKEConfigRegistryMirrors(in.Mirrors)
	}

	return []interface{}{obj}
}

// Expanders

func expandClusterV2RKEConfigRegistryConfigs(p []interface{}) map[string]rkev1.RegistryConfig {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}

	out := map[string]rkev1.RegistryConfig{}
	for i := range p {
		in := p[i].(map[string]interface{})
		obj := rkev1.RegistryConfig{}

		if v, ok := in["auth_config_secret_name"].(string); ok && len(v) > 0 {
			obj.AuthConfigSecretName = v
		}
		if v, ok := in["tls_secret_name"].(string); ok && len(v) > 0 {
			obj.TLSSecretName = v
		}
		if v, ok := in["ca_bundle"].(string); ok && len(v) > 0 {
			obj.CABundle = []byte(v)
		}
		if v, ok := in["insecure"].(bool); ok {
			obj.InsecureSkipVerify = v
		}
		out[in["hostname"].(string)] = obj
	}

	return out
}

func expandClusterV2RKEConfigRegistryMirrors(p []interface{}) map[string]rkev1.Mirror {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}

	out := map[string]rkev1.Mirror{}
	for i := range p {
		in := p[i].(map[string]interface{})
		obj := rkev1.Mirror{}

		if v, ok := in["endpoints"].([]interface{}); ok && len(v) > 0 {
			obj.Endpoints = toArrayString(v)
		}
		if v, ok := in["rewrites"].(map[string]interface{}); ok && len(v) > 0 {
			obj.Rewrites = toMapString(v)
		}
		out[in["hostname"].(string)] = obj
	}

	return out
}

func expandClusterV2RKEConfigRegistry(p []interface{}) *rkev1.Registry {
	if p == nil || len(p) == 0 || p[0] == nil {
		return nil
	}

	obj := &rkev1.Registry{}

	in := p[0].(map[string]interface{})

	if v, ok := in["configs"].([]interface{}); ok && len(v) > 0 {
		obj.Configs = expandClusterV2RKEConfigRegistryConfigs(v)
	}
	if v, ok := in["mirrors"].([]interface{}); ok && len(v) > 0 {
		obj.Mirrors = expandClusterV2RKEConfigRegistryMirrors(v)
	}

	return obj
}
