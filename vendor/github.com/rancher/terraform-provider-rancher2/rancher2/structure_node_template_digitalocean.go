package rancher2

// Flatteners

func flattenDigitaloceanConfig(in *digitaloceanConfig) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	if len(in.AccessToken) > 0 {
		obj["access_token"] = in.AccessToken
	}

	obj["backups"] = in.Backups

	if len(in.Image) > 0 {
		obj["image"] = in.Image
	}

	obj["ipv6"] = in.IPV6
	obj["monitoring"] = in.Monitoring
	obj["private_networking"] = in.PrivateNetworking

	if len(in.Region) > 0 {
		obj["region"] = in.Region
	}

	if len(in.Size) > 0 {
		obj["size"] = in.Size
	}

	if len(in.SSHKeyFingerprint) > 0 {
		obj["ssh_key_fingerprint"] = in.SSHKeyFingerprint
	}

	if len(in.SSHKeyPath) > 0 {
		obj["ssh_key_path"] = in.SSHKeyPath
	}

	if len(in.SSHPort) > 0 {
		obj["ssh_port"] = in.SSHPort
	}

	if len(in.SSHUser) > 0 {
		obj["ssh_user"] = in.SSHUser
	}

	if len(in.Tags) > 0 {
		obj["tags"] = in.Tags
	}

	if len(in.Userdata) > 0 {
		obj["userdata"] = in.Userdata
	}

	return []interface{}{obj}
}

// Expanders

func expandDigitaloceanConfig(p []interface{}) *digitaloceanConfig {
	obj := &digitaloceanConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["access_token"].(string); ok && len(v) > 0 {
		obj.AccessToken = v
	}

	if v, ok := in["backups"].(bool); ok {
		obj.Backups = v
	}

	if v, ok := in["image"].(string); ok && len(v) > 0 {
		obj.Image = v
	}

	if v, ok := in["ipv6"].(bool); ok {
		obj.IPV6 = v
	}

	if v, ok := in["monitoring"].(bool); ok {
		obj.Monitoring = v
	}
	if v, ok := in["private_networking"].(bool); ok {
		obj.PrivateNetworking = v
	}

	if v, ok := in["region"].(string); ok && len(v) > 0 {
		obj.Region = v
	}

	if v, ok := in["size"].(string); ok && len(v) > 0 {
		obj.Size = v
	}

	if v, ok := in["ssh_key_fingerprint"].(string); ok && len(v) > 0 {
		obj.SSHKeyFingerprint = v
	}

	if v, ok := in["ssh_key_path"].(string); ok && len(v) > 0 {
		obj.SSHKeyPath = v
	}

	if v, ok := in["ssh_port"].(string); ok && len(v) > 0 {
		obj.SSHPort = v
	}

	if v, ok := in["ssh_user"].(string); ok && len(v) > 0 {
		obj.SSHUser = v
	}

	if v, ok := in["tags"].(string); ok && len(v) > 0 {
		obj.Tags = v
	}

	if v, ok := in["userdata"].(string); ok && len(v) > 0 {
		obj.Userdata = v
	}

	return obj
}
