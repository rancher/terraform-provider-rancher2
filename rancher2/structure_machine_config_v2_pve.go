package rancher2

// Flatteners

func flattenMachineConfigV2Pve(in map[string]interface{}, p []interface{}) []interface{} {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if v, ok := in["pveUrl"].(string); ok && len(v) > 0 {
		obj["pve_url"] = v
	}
	if v, ok := in["pveTokenId"].(string); ok && len(v) > 0 {
		obj["pve_token_id"] = v
	}
	if v, ok := in["pveTokenSecret"].(string); ok && len(v) > 0 {
		obj["pve_token_secret"] = v
	}
	if v, ok := in["pveInsecureTls"].(bool); ok {
		obj["pve_insecure_tls"] = v
	}
	if v, ok := in["pveResourcePool"].(string); ok && len(v) > 0 {
		obj["pve_resource_pool"] = v
	}
	if v, ok := in["pveTemplateId"].(float64); ok {
		obj["pve_template_id"] = int(v)
	}
	if v, ok := in["pveIsoDevice"].(string); ok && len(v) > 0 {
		obj["pve_iso_device"] = v
	}
	if v, ok := in["pveNetworkInterface"].(string); ok && len(v) > 0 {
		obj["pve_network_interface"] = v
	}
	if v, ok := in["pveSshUser"].(string); ok && len(v) > 0 {
		obj["pve_ssh_user"] = v
	}
	if v, ok := in["pveSshPort"].(float64); ok {
		obj["pve_ssh_port"] = int(v)
	}
	if v, ok := in["pveProcessorSockets"].(string); ok && len(v) > 0 {
		obj["pve_processor_sockets"] = v
	}
	if v, ok := in["pveProcessorCores"].(string); ok && len(v) > 0 {
		obj["pve_processor_cores"] = v
	}
	if v, ok := in["pveMemory"].(string); ok && len(v) > 0 {
		obj["pve_memory"] = v
	}

	return []interface{}{obj}
}

// Expanders

func expandMachineConfigV2Pve(p []interface{}) map[string]interface{} {
	obj := make(map[string]interface{})
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["pve_url"].(string); ok && len(v) > 0 {
		obj["pveUrl"] = v
	}
	if v, ok := in["pve_token_id"].(string); ok && len(v) > 0 {
		obj["pveTokenId"] = v
	}
	if v, ok := in["pve_token_secret"].(string); ok && len(v) > 0 {
		obj["pveTokenSecret"] = v
	}
	if v, ok := in["pve_insecure_tls"].(bool); ok {
		obj["pveInsecureTls"] = v
	}
	if v, ok := in["pve_resource_pool"].(string); ok && len(v) > 0 {
		obj["pveResourcePool"] = v
	}
	if v, ok := in["pve_template_id"].(int); ok && v > 0 {
		obj["pveTemplateId"] = v
	}
	if v, ok := in["pve_iso_device"].(string); ok && len(v) > 0 {
		obj["pveIsoDevice"] = v
	}
	if v, ok := in["pve_network_interface"].(string); ok && len(v) > 0 {
		obj["pveNetworkInterface"] = v
	}
	if v, ok := in["pve_ssh_user"].(string); ok && len(v) > 0 {
		obj["pveSshUser"] = v
	}
	if v, ok := in["pve_ssh_port"].(int); ok && v > 0 {
		obj["pveSshPort"] = v
	}
	if v, ok := in["pve_processor_sockets"].(string); ok && len(v) > 0 {
		obj["pveProcessorSockets"] = v
	}
	if v, ok := in["pve_processor_cores"].(string); ok && len(v) > 0 {
		obj["pveProcessorCores"] = v
	}
	if v, ok := in["pve_memory"].(string); ok && len(v) > 0 {
		obj["pveMemory"] = v
	}

	return obj
}
