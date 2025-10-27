package rancher2

// Flatteners

func flattenHarvesterConfig(in *harvesterConfig) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	if len(in.VMNamespace) > 0 {
		obj["vm_namespace"] = in.VMNamespace
	}

	if len(in.VMAffinity) > 0 {
		obj["vm_affinity"] = in.VMAffinity
	}

	if len(in.CPUCount) > 0 {
		obj["cpu_count"] = in.CPUCount
	}

	if len(in.MemorySize) > 0 {
		obj["memory_size"] = in.MemorySize
	}

	if len(in.DiskSize) > 0 {
		obj["disk_size"] = in.DiskSize
	}

	if len(in.DiskBus) > 0 {
		obj["disk_bus"] = in.DiskBus
	}

	if len(in.ImageName) > 0 {
		obj["image_name"] = in.ImageName
	}

	if len(in.DiskInfo) > 0 {
		obj["disk_info"] = in.DiskInfo
	}

	if len(in.SSHUser) > 0 {
		obj["ssh_user"] = in.SSHUser
	}

	if len(in.SSHPassword) > 0 {
		obj["ssh_password"] = in.SSHPassword
	}

	if len(in.NetworkName) > 0 {
		obj["network_name"] = in.NetworkName
	}

	if len(in.NetworkModel) > 0 {
		obj["network_model"] = in.NetworkModel
	}

	if len(in.NetworkInfo) > 0 {
		obj["network_info"] = in.NetworkInfo
	}

	if len(in.UserData) > 0 {
		obj["user_data"] = in.UserData
	}

	if len(in.NetworkData) > 0 {
		obj["network_data"] = in.NetworkData
	}

	return []interface{}{obj}
}

// Expanders

func expandHarvestercloudConfig(p []interface{}) *harvesterConfig {
	obj := &harvesterConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["vm_namespace"].(string); ok && len(v) > 0 {
		obj.VMNamespace = v
	}

	if v, ok := in["vm_affinity"].(string); ok && len(v) > 0 {
		obj.VMAffinity = v
	}

	if v, ok := in["cpu_count"].(string); ok && len(v) > 0 {
		obj.CPUCount = v
	}

	if v, ok := in["memory_size"].(string); ok && len(v) > 0 {
		obj.MemorySize = v
	}

	if v, ok := in["disk_size"].(string); ok && len(v) > 0 {
		obj.DiskSize = v
	}

	if v, ok := in["disk_bus"].(string); ok && len(v) > 0 {
		obj.DiskBus = v
	}

	if v, ok := in["image_name"].(string); ok && len(v) > 0 {
		obj.ImageName = v
	}

	if v, ok := in["disk_info"].(string); ok && len(v) > 0 {
		obj.DiskInfo = v
	}

	if v, ok := in["ssh_user"].(string); ok && len(v) > 0 {
		obj.SSHUser = v
	}

	if v, ok := in["ssh_password"].(string); ok && len(v) > 0 {
		obj.SSHPassword = v
	}

	if v, ok := in["network_name"].(string); ok && len(v) > 0 {
		obj.NetworkName = v
	}

	if v, ok := in["network_model"].(string); ok && len(v) > 0 {
		obj.NetworkModel = v
	}

	if v, ok := in["network_info"].(string); ok && len(v) > 0 {
		obj.NetworkInfo = v
	}

	if v, ok := in["user_data"].(string); ok && len(v) > 0 {
		obj.UserData = v
	}

	if v, ok := in["network_data"].(string); ok && len(v) > 0 {
		obj.NetworkData = v
	}

	return obj
}
