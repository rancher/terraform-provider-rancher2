package rancher2

import (
	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
)

// Flatteners

func flattenClusterV2LocalAuthEndpoint(in rkev1.LocalClusterAuthEndpoint) []interface{} {
	empty := rkev1.LocalClusterAuthEndpoint{}
	if in == empty {
		return nil
	}

	obj := make(map[string]interface{})

	if len(in.CACerts) > 0 {
		obj["ca_certs"] = in.CACerts
	}
	obj["enabled"] = in.Enabled
	if len(in.FQDN) > 0 {
		obj["fqdn"] = in.FQDN
	}

	return []interface{}{obj}
}

// Expanders

func expandClusterV2LocalAuthEndpoint(p []interface{}) rkev1.LocalClusterAuthEndpoint {
	if p == nil || len(p) == 0 || p[0] == nil {
		return rkev1.LocalClusterAuthEndpoint{}
	}

	obj := rkev1.LocalClusterAuthEndpoint{}

	in := p[0].(map[string]interface{})

	if v, ok := in["ca_certs"].(string); ok && len(v) > 0 {
		obj.CACerts = v
	}
	if v, ok := in["enabled"].(bool); ok {
		obj.Enabled = v
	}
	if v, ok := in["fqdn"].(string); ok && len(v) > 0 {
		obj.FQDN = v
	}

	return obj
}
