package rancher2

// Flatteners

func flattenCloudCredentialVsphere(in *vmwarevsphereCredentialConfig, p []interface{}) []interface{} {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}
	}

	if len(in.Password) > 0 {
		obj["password"] = in.Password
	}

	if len(in.Username) > 0 {
		obj["username"] = in.Username
	}

	if len(in.Vcenter) > 0 {
		obj["vcenter"] = in.Vcenter
	}

	if len(in.VcenterPort) > 0 {
		obj["vcenter_port"] = in.VcenterPort
	}

	return []interface{}{obj}
}

// Expanders

func expandCloudCredentialVsphere(p []interface{}) *vmwarevsphereCredentialConfig {
	obj := &vmwarevsphereCredentialConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["password"].(string); ok && len(v) > 0 {
		obj.Password = v
	}

	if v, ok := in["username"].(string); ok && len(v) > 0 {
		obj.Username = v
	}

	if v, ok := in["vcenter"].(string); ok && len(v) > 0 {
		obj.Vcenter = v
	}

	if v, ok := in["vcenter_port"].(string); ok && len(v) > 0 {
		obj.VcenterPort = v
	}

	return obj
}
