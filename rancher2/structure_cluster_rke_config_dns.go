package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenClusterRKEConfigDNSNodelocal(in *managementClient.Nodelocal) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return nil
	}

	if len(in.IPAddress) > 0 {
		obj["ip_address"] = in.IPAddress
	}

	if len(in.NodeSelector) > 0 {
		obj["node_selector"] = toMapInterface(in.NodeSelector)
	}

	return []interface{}{obj}
}

func flattenClusterRKEConfigDNS(in *managementClient.DNSConfig) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.NodeSelector) > 0 {
		obj["node_selector"] = toMapInterface(in.NodeSelector)
	}

	if in.Nodelocal != nil {
		obj["nodelocal"] = flattenClusterRKEConfigDNSNodelocal(in.Nodelocal)
	}

	if len(in.Provider) > 0 {
		obj["provider"] = in.Provider
	}

	if len(in.ReverseCIDRs) > 0 {
		obj["reverse_cidrs"] = toArrayInterface(in.ReverseCIDRs)
	}

	if len(in.UpstreamNameservers) > 0 {
		obj["upstream_nameservers"] = toArrayInterface(in.UpstreamNameservers)
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandClusterRKEConfigDNSNodelocal(p []interface{}) *managementClient.Nodelocal {
	obj := &managementClient.Nodelocal{}
	if len(p) == 0 || p[0] == nil {
		return nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["ip_address"].(string); ok && len(v) > 0 {
		obj.IPAddress = v
	}

	if v, ok := in["node_selector"].(map[string]interface{}); ok && len(v) > 0 {
		obj.NodeSelector = toMapString(v)
	}

	return obj
}

func expandClusterRKEConfigDNS(p []interface{}) (*managementClient.DNSConfig, error) {
	obj := &managementClient.DNSConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["node_selector"].(map[string]interface{}); ok && len(v) > 0 {
		obj.NodeSelector = toMapString(v)
	}

	if v, ok := in["nodelocal"].([]interface{}); ok && len(v) > 0 {
		obj.Nodelocal = expandClusterRKEConfigDNSNodelocal(v)
	}

	if v, ok := in["provider"].(string); ok && len(v) > 0 {
		obj.Provider = v
	}

	if v, ok := in["reverse_cidrs"].([]interface{}); ok && len(v) > 0 {
		obj.ReverseCIDRs = toArrayString(v)
	}

	if v, ok := in["upstream_nameservers"].([]interface{}); ok && len(v) > 0 {
		obj.UpstreamNameservers = toArrayString(v)
	}

	return obj, nil
}
