package rancher2

// Flatteners

func flattenLinodeConfig(in *linodeConfig) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	if len(in.AuthorizedUsers) > 0 {
		obj["authorized_users"] = in.AuthorizedUsers
	}

	obj["create_private_ip"] = in.CreatePrivateIP

	if len(in.DockerPort) > 0 {
		obj["docker_port"] = in.DockerPort
	}

	if len(in.Image) > 0 {
		obj["image"] = in.Image
	}

	if len(in.InstanceType) > 0 {
		obj["instance_type"] = in.InstanceType
	}

	if len(in.Label) > 0 {
		obj["label"] = in.Label
	}

	if len(in.Region) > 0 {
		obj["region"] = in.Region
	}

	if len(in.RootPass) > 0 {
		obj["root_pass"] = in.RootPass
	}

	if len(in.SSHPort) > 0 {
		obj["ssh_port"] = in.SSHPort
	}

	if len(in.SSHUser) > 0 {
		obj["ssh_user"] = in.SSHUser
	}

	if len(in.StackScript) > 0 {
		obj["stackscript"] = in.StackScript
	}

	if len(in.StackscriptData) > 0 {
		obj["stackscript_data"] = in.StackscriptData
	}

	if len(in.SwapSize) > 0 {
		obj["swap_size"] = in.SwapSize
	}

	if len(in.Tags) > 0 {
		obj["tags"] = in.Tags
	}

	if len(in.Token) > 0 {
		obj["token"] = in.Token
	}

	if len(in.UAPrefix) > 0 {
		obj["ua_prefix"] = in.UAPrefix
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

	if v, ok := in["authorized_users"].(string); ok && len(v) > 0 {
		obj.AuthorizedUsers = v
	}

	if v, ok := in["create_private_ip"].(bool); ok {
		obj.CreatePrivateIP = v
	}

	if v, ok := in["docker_port"].(string); ok && len(v) > 0 {
		obj.DockerPort = v
	}

	if v, ok := in["image"].(string); ok && len(v) > 0 {
		obj.Image = v
	}

	if v, ok := in["instance_type"].(string); ok && len(v) > 0 {
		obj.InstanceType = v
	}

	if v, ok := in["label"].(string); ok && len(v) > 0 {
		obj.Label = v
	}

	if v, ok := in["region"].(string); ok && len(v) > 0 {
		obj.Region = v
	}

	if v, ok := in["root_pass"].(string); ok && len(v) > 0 {
		obj.RootPass = v
	}

	if v, ok := in["ssh_port"].(string); ok && len(v) > 0 {
		obj.SSHPort = v
	}

	if v, ok := in["ssh_user"].(string); ok && len(v) > 0 {
		obj.SSHUser = v
	}

	if v, ok := in["stackscript"].(string); ok && len(v) > 0 {
		obj.StackScript = v
	}

	if v, ok := in["stackscript_data"].(string); ok && len(v) > 0 {
		obj.StackscriptData = v
	}

	if v, ok := in["swap_size"].(string); ok && len(v) > 0 {
		obj.SwapSize = v
	}

	if v, ok := in["tags"].(string); ok && len(v) > 0 {
		obj.Tags = v
	}

	if v, ok := in["token"].(string); ok && len(v) > 0 {
		obj.Token = v
	}

	if v, ok := in["ua_prefix"].(string); ok && len(v) > 0 {
		obj.UAPrefix = v
	}

	return obj
}
