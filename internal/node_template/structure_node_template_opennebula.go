package rancher2

// Flatteners

func flattenOpenNebulaConfig(in *opennebulaConfig) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	obj["b2d_size"] = in.B2dSize
	obj["dev_prefix"] = in.DevPrefix
	obj["disk_resize"] = in.DiskResize
	obj["image_name"] = in.ImageName
	obj["memory"] = in.Memory
	obj["network_name"] = in.NetworkName
	obj["password"] = in.Password
	obj["template_id"] = in.TemplateID
	obj["user"] = in.User
	obj["xml_rpc_url"] = in.XMLRPCURL
	obj["cpu"] = in.CPU
	obj["disable_vnc"] = in.DisableVnc
	obj["image_id"] = in.ImageID
	obj["image_owner"] = in.ImageOwner
	obj["network_id"] = in.NetworkID
	obj["network_owner"] = in.NetworkOwner
	obj["ssh_user"] = in.SSHUser
	obj["template_name"] = in.TemplateName
	obj["vcpu"] = in.Vcpu

	return []interface{}{obj}
}

// Expanders

func expandOpennebulaConfig(p []interface{}) *opennebulaConfig {
	obj := &opennebulaConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["b2d_size"].(string); ok && len(v) > 0 {
		obj.B2dSize = v
	}
	if v, ok := in["dev_prefix"].(string); ok && len(v) > 0 {
		obj.DevPrefix = v
	}
	if v, ok := in["disk_resize"].(string); ok && len(v) > 0 {
		obj.DiskResize = v
	}
	if v, ok := in["image_name"].(string); ok && len(v) > 0 {
		obj.ImageName = v
	}
	if v, ok := in["memory"].(string); ok && len(v) > 0 {
		obj.Memory = v
	}
	if v, ok := in["network_name"].(string); ok && len(v) > 0 {
		obj.NetworkName = v
	}
	if v, ok := in["password"].(string); ok && len(v) > 0 {
		obj.Password = v
	}
	if v, ok := in["template_id"].(string); ok && len(v) > 0 {
		obj.TemplateID = v
	}
	if v, ok := in["user"].(string); ok && len(v) > 0 {
		obj.User = v
	}
	if v, ok := in["xml_rpc_url"].(string); ok && len(v) > 0 {
		obj.XMLRPCURL = v
	}
	if v, ok := in["cpu"].(string); ok && len(v) > 0 {
		obj.CPU = v
	}
	if v, ok := in["disable_vnc"].(bool); ok {
		obj.DisableVnc = v
	}
	if v, ok := in["image_id"].(string); ok && len(v) > 0 {
		obj.ImageID = v
	}
	if v, ok := in["image_owner"].(string); ok && len(v) > 0 {
		obj.ImageOwner = v
	}
	if v, ok := in["network_id"].(string); ok && len(v) > 0 {
		obj.NetworkID = v
	}
	if v, ok := in["network_owner"].(string); ok && len(v) > 0 {
		obj.NetworkOwner = v
	}
	if v, ok := in["ssh_user"].(string); ok && len(v) > 0 {
		obj.SSHUser = v
	}
	if v, ok := in["template_name"].(string); ok && len(v) > 0 {
		obj.TemplateName = v
	}
	if v, ok := in["vcpu"].(string); ok && len(v) > 0 {
		obj.Vcpu = v
	}

	return obj
}
