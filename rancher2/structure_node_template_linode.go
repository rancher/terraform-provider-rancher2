package rancher2

// Flatteners

func flattenLinodeConfig(in *linodeConfig) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	if len(in.AuthorizedUsers) > 0 {
		obj["authorizedUsers"] = in.AuthorizedUsers
	}

	obj["createPrivateIp"] = in.CreatePrivateIp

	if len(in.DockerPort) > 0 {
		obj["dockerPort"] = in.DockerPort
	}

	if len(in.Image) > 0 {
		obj["image"] = in.Image
	}

	if len(in.InstanceType) > 0 {
		obj["instanceType"] = in.InstanceType
	}

	if len(in.Label) > 0 {
		obj["label"] = in.Label
	}

	if len(in.Region) > 0 {
		obj["region"] = in.Region
	}

	if len(in.RootPass) > 0 {
		obj["rootPass"] = in.RootPass
	}

	if len(in.SSHPort) > 0 {
		obj["sshPort"] = in.SSHPort
	}

	if len(in.SSHUser) > 0 {
		obj["sshUser"] = in.SSHUser
	}

	if len(in.StackScript) > 0 {
		obj["stackscript"] = in.StackScript
	}

	if len(in.StackscriptData) > 0 {
		obj["stackscriptData"] = in.StackscriptData
	}

	if len(in.SwapSize) > 0 {
		obj["swapSize"] = in.SwapSize
	}

	if len(in.Tags) > 0 {
		obj["tags"] = in.Tags
	}

	if len(in.Token) > 0 {
		obj["token"] = in.Token
	}

	if len(in.UAPrefix) > 0 {
		obj["uaPrefix"] = in.UAPrefix
	}

	return []interface{}{obj}
}

// Expanders

func expandLinodeConfig(p []interface{}) *linodeConfig {
	obj := &linodeConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["authorizedUsers"].(string); ok && len(v) > 0 {
		obj.AuthorizedUsers = v
	}

	if v, ok := in["createPrivateIp"].(bool); ok {
		obj.CreatePrivateIp = v
	}

	if v, ok := in["dockerPort"].(string); ok && len(v) > 0 {
		obj.DockerPort = v
	}

	if v, ok := in["image"].(string); ok && len(v) > 0 {
		obj.Image = v
	}

	if v, ok := in["instanceType"].(string); ok && len(v) > 0 {
		obj.InstanceType = v
	}

	if v, ok := in["label"].(string); ok && len(v) > 0 {
		obj.Label = v
	}

	if v, ok := in["region"].(string); ok && len(v) > 0 {
		obj.Region = v
	}

	if v, ok := in["rootPass"].(string); ok && len(v) > 0 {
		obj.RootPass = v
	}

	if v, ok := in["sshPort"].(string); ok && len(v) > 0 {
		obj.SSHPort = v
	}

	if v, ok := in["sshUser"].(string); ok && len(v) > 0 {
		obj.SSHUser = v
	}

	if v, ok := in["stackscript"].(string); ok && len(v) > 0 {
		obj.StackScript = v
	}

	if v, ok := in["stackscriptData"].(string); ok && len(v) > 0 {
		obj.StackscriptData = v
	}

	if v, ok := in["swapSize"].(string); ok && len(v) > 0 {
		obj.SwapSize = v
	}

	if v, ok := in["tags"].(string); ok && len(v) > 0 {
		obj.Tags = v
	}

	if v, ok := in["token"].(string); ok && len(v) > 0 {
		obj.Token = v
	}

	if v, ok := in["uaPrefix"].(string); ok && len(v) > 0 {
		obj.UAPrefix = v
	}

	return obj
}
