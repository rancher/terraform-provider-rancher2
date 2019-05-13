package rancher2

// Flatteners

func flattenVsphereConfig(in *vmwarevsphereConfig) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	if len(in.Boot2dockerURL) > 0 {
		obj["boot2docker_url"] = in.Boot2dockerURL
	}
	if len(in.Cfgparam) > 0 {
		obj["cfgparam"] = toArrayInterface(in.Cfgparam)
	}
	if len(in.Cloudinit) > 0 {
		obj["cloudinit"] = in.Cloudinit
	}
	if len(in.CPUCount) > 0 {
		obj["cpu_count"] = in.CPUCount
	}
	if len(in.Datacenter) > 0 {
		obj["datacenter"] = in.Datastore
	}
	if len(in.Datastore) > 0 {
		obj["datastore"] = in.Datastore
	}
	if len(in.DiskSize) > 0 {
		obj["disk_size"] = in.DiskSize
	}
	if len(in.Folder) > 0 {
		obj["folder"] = in.Folder
	}
	if len(in.Hostsystem) > 0 {
		obj["hostsystem"] = in.Hostsystem
	}
	if len(in.MemorySize) > 0 {
		obj["memory_size"] = in.MemorySize
	}
	if len(in.Network) > 0 {
		obj["network"] = toArrayInterface(in.Network)
	}
	if len(in.Password) > 0 {
		obj["password"] = in.Password
	}
	if len(in.Pool) > 0 {
		obj["pool"] = in.Pool
	}
	if len(in.Username) > 0 {
		obj["username"] = in.Username
	}
	if len(in.VappIpallocationpolicy) > 0 {
		obj["vapp_ip_allocation_policy"] = in.VappIpallocationpolicy
	}
	if len(in.VappIpprotocol) > 0 {
		obj["vapp_ip_protocol"] = in.VappIpprotocol
	}
	if len(in.VappProperty) > 0 {
		obj["vapp_property"] = toArrayInterface(in.VappProperty)
	}
	if len(in.VappTransport) > 0 {
		obj["vapp_transport"] = in.VappTransport
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

func expandVsphereConfig(p []interface{}) *vmwarevsphereConfig {
	obj := &vmwarevsphereConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["boot2docker_url"].(string); ok && len(v) > 0 {
		obj.Boot2dockerURL = v
	}
	if v, ok := in["cfgparam"].([]interface{}); ok && len(v) > 0 {
		obj.Cfgparam = toArrayString(v)
	}
	if v, ok := in["cloudinit"].(string); ok && len(v) > 0 {
		obj.Cloudinit = v
	}
	if v, ok := in["cpu_count"].(string); ok && len(v) > 0 {
		obj.CPUCount = v
	}
	if v, ok := in["datacenter"].(string); ok && len(v) > 0 {
		obj.Datacenter = v
	}
	if v, ok := in["datastore"].(string); ok && len(v) > 0 {
		obj.Datastore = v
	}
	if v, ok := in["disk_size"].(string); ok && len(v) > 0 {
		obj.DiskSize = v
	}
	if v, ok := in["folder"].(string); ok && len(v) > 0 {
		obj.Folder = v
	}
	if v, ok := in["hostsystem"].(string); ok && len(v) > 0 {
		obj.Hostsystem = v
	}
	if v, ok := in["memory_size"].(string); ok && len(v) > 0 {
		obj.MemorySize = v
	}
	if v, ok := in["network"].([]interface{}); ok && len(v) > 0 {
		obj.Network = toArrayString(v)
	}
	if v, ok := in["password"].(string); ok && len(v) > 0 {
		obj.Password = v
	}
	if v, ok := in["pool"].(string); ok && len(v) > 0 {
		obj.Pool = v
	}
	if v, ok := in["username"].(string); ok && len(v) > 0 {
		obj.Username = v
	}
	if v, ok := in["vapp_ip_allocation_policy"].(string); ok && len(v) > 0 {
		obj.VappIpallocationpolicy = v
	}
	if v, ok := in["vapp_ip_protocol"].(string); ok && len(v) > 0 {
		obj.VappIpprotocol = v
	}
	if v, ok := in["vapp_property"].([]interface{}); ok && len(v) > 0 {
		obj.VappProperty = toArrayString(v)
	}
	if v, ok := in["vapp_transport"].(string); ok && len(v) > 0 {
		obj.VappTransport = v
	}
	if v, ok := in["vcenter"].(string); ok && len(v) > 0 {
		obj.Vcenter = v
	}
	if v, ok := in["vcenter_port"].(string); ok && len(v) > 0 {
		obj.VcenterPort = v
	}

	return obj
}
