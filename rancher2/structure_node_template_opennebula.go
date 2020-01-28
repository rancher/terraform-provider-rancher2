package rancher2

// Flatteners

func flattenOpenNebulaConfig(in *opennebulaConfig) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	obj["b2dSize"] = in.b2dsize
	obj["devPrefix"] = in.devprefix
	obj["diskResize"] = in.diskresize
	obj["imageName"] = in.imagename
	obj["memory"] = in.memory
	obj["networkName"] = in.networkname
	obj["password"] = in.password
	obj["templateId"] = in.templateid
	obj["user"] = in.user
	obj["xmlrpcurl"] = in.xmlrpcurl
	obj["cpu"] = in.cpu
	obj["disableVnc"] = in.disablevnc
	obj["imageId"] = in.imageid
	obj["imageOwner"] = in.imageowner
	obj["networkId"] = in.networkid
	obj["networkOwner"] = in.networkowner
	obj["sshUser"] = in.sshuser
	obj["templateName"] = in.templatename
	obj["vcpu"] = in.vcpu

	return []interface{}{obj}
}

// Expanders

func expandOpennebulaConfig(p []interface{}) *opennebulaConfig {
	obj := &opennebulaConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["b2dSize"].(string); ok && len(v) > 0 {
		obj.b2dsize = v
	}
	if v, ok := in["devPrefix"].(string); ok && len(v) > 0 {
		obj.devprefix = v
	}
	if v, ok := in["diskResize"].(string); ok && len(v) > 0 {
		obj.diskresize = v
	}
	if v, ok := in["imageName"].(string); ok && len(v) > 0 {
		obj.imagename = v
	}
	if v, ok := in["memory"].(string); ok && len(v) > 0 {
		obj.memory = v
	}
	if v, ok := in["networkName"].(string); ok && len(v) > 0 {
		obj.networkname = v
	}
	if v, ok := in["password"].(string); ok && len(v) > 0 {
		obj.password = v
	}
	if v, ok := in["templateId"].(string); ok && len(v) > 0 {
		obj.templateid = v
	}
	if v, ok := in["user"].(string); ok && len(v) > 0 {
		obj.user = v
	}
	if v, ok := in["xmlrpcurl"].(string); ok && len(v) > 0 {
		obj.xmlrpcurl = v
	}
	if v, ok := in["cpu"].(string); ok && len(v) > 0 {
		obj.cpu = v
	}
	if v, ok := in["disableVnc"].(bool); ok {
		obj.disablevnc = v
	}
	if v, ok := in["imageId"].(string); ok && len(v) > 0 {
		obj.imageid = v
	}
	if v, ok := in["imageOwner"].(string); ok && len(v) > 0 {
		obj.imageowner = v
	}
	if v, ok := in["networkId"].(string); ok && len(v) > 0 {
		obj.networkid = v
	}
	if v, ok := in["networkOwner"].(string); ok && len(v) > 0 {
		obj.networkowner = v
	}
	if v, ok := in["sshUser"].(string); ok && len(v) > 0 {
		obj.sshuser = v
	}
	if v, ok := in["templateName"].(string); ok && len(v) > 0 {
		obj.templatename = v
	}
	if v, ok := in["vcpu"].(string); ok && len(v) > 0 {
		obj.vcpu = v
	}

	return obj
}
