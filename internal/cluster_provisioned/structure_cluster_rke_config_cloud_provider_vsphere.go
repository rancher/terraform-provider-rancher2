package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenClusterRKEConfigCloudProviderVsphereDisk(in *managementClient.DiskVsphereOpts) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.SCSIControllerType) > 0 {
		obj["scsi_controller_type"] = in.SCSIControllerType
	}

	return []interface{}{obj}, nil
}

func flattenClusterRKEConfigCloudProviderVsphereGlobal(in *managementClient.GlobalVsphereOpts, p []interface{}) ([]interface{}, error) {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.Datacenters) > 0 {
		obj["datacenters"] = in.Datacenters
	}

	obj["insecure_flag"] = in.InsecureFlag

	if len(in.Password) > 0 {
		obj["password"] = in.Password
	}

	if len(in.VCenterPort) > 0 {
		obj["port"] = in.VCenterPort
	}

	if len(in.User) > 0 {
		obj["user"] = in.User
	}

	if in.RoundTripperCount > 0 {
		obj["soap_roundtrip_count"] = int(in.RoundTripperCount)
	}

	return []interface{}{obj}, nil
}

func flattenClusterRKEConfigCloudProviderVsphereNetwork(in *managementClient.NetworkVshpereOpts) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.PublicNetwork) > 0 {
		obj["public_network"] = in.PublicNetwork
	}

	return []interface{}{obj}, nil
}

func flattenClusterRKEConfigCloudProviderVsphereVirtualCenter(in map[string]managementClient.VirtualCenterConfig, p []interface{}) ([]interface{}, error) {
	if len(in) == 0 {
		return []interface{}{}, nil
	}

	out := make([]interface{}, len(in))
	lenP := len(p)
	i := 0
	for key := range in {
		var obj map[string]interface{}
		if lenP <= i {
			obj = make(map[string]interface{})
		} else {
			obj = p[i].(map[string]interface{})
		}

		obj["name"] = key
		if len(in[key].Datacenters) > 0 {
			obj["datacenters"] = in[key].Datacenters
		}

		if len(in[key].Password) > 0 {
			obj["password"] = in[key].Password
		}

		if len(in[key].VCenterPort) > 0 {
			obj["port"] = in[key].VCenterPort
		}

		if len(in[key].User) > 0 {
			obj["user"] = in[key].User
		}

		if in[key].RoundTripperCount > 0 {
			obj["soap_roundtrip_count"] = int(in[key].RoundTripperCount)
		}
		out[i] = obj
		i++
	}

	return out, nil
}

func flattenClusterRKEConfigCloudProviderVsphereWorkspace(in *managementClient.WorkspaceVsphereOpts) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.Datacenter) > 0 {
		obj["datacenter"] = in.Datacenter
	}

	if len(in.Folder) > 0 {
		obj["folder"] = in.Folder
	}

	if len(in.VCenterIP) > 0 {
		obj["server"] = in.VCenterIP
	}

	if len(in.DefaultDatastore) > 0 {
		obj["default_datastore"] = in.DefaultDatastore
	}

	if len(in.ResourcePoolPath) > 0 {
		obj["resourcepool_path"] = in.ResourcePoolPath
	}

	return []interface{}{obj}, nil
}

func flattenClusterRKEConfigCloudProviderVsphere(in *managementClient.VsphereCloudProvider, p []interface{}) ([]interface{}, error) {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}, nil
	}

	if in.Disk != nil {
		disk, err := flattenClusterRKEConfigCloudProviderVsphereDisk(in.Disk)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["disk"] = disk
	}

	if in.Global != nil {
		v, ok := obj["global"].([]interface{})
		if !ok {
			v = []interface{}{}
		}
		global, err := flattenClusterRKEConfigCloudProviderVsphereGlobal(in.Global, v)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["global"] = global
	}

	if in.Network != nil {
		network, err := flattenClusterRKEConfigCloudProviderVsphereNetwork(in.Network)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["network"] = network
	}

	if in.VirtualCenter != nil {
		v, ok := obj["virtual_center"].([]interface{})
		if !ok {
			v = []interface{}{}
		}
		vc, err := flattenClusterRKEConfigCloudProviderVsphereVirtualCenter(in.VirtualCenter, v)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["virtual_center"] = vc
	}

	if in.Workspace != nil {
		workspace, err := flattenClusterRKEConfigCloudProviderVsphereWorkspace(in.Workspace)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["workspace"] = workspace
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandClusterRKEConfigCloudProviderVsphereDisk(p []interface{}) (*managementClient.DiskVsphereOpts, error) {
	obj := &managementClient.DiskVsphereOpts{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["scsi_controller_type"].(string); ok && len(v) > 0 {
		obj.SCSIControllerType = v
	}

	return obj, nil
}

func expandClusterRKEConfigCloudProviderVsphereGlobal(p []interface{}) (*managementClient.GlobalVsphereOpts, error) {
	obj := &managementClient.GlobalVsphereOpts{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["datacenters"].(string); ok && len(v) > 0 {
		obj.Datacenters = v
	}

	if v, ok := in["insecure_flag"].(bool); ok {
		obj.InsecureFlag = v
	}

	if v, ok := in["password"].(string); ok && len(v) > 0 {
		obj.Password = v
	}

	if v, ok := in["port"].(string); ok && len(v) > 0 {
		obj.VCenterPort = v
	}

	if v, ok := in["user"].(string); ok && len(v) > 0 {
		obj.User = v
	}

	if v, ok := in["soap_roundtrip_count"].(int); ok && v > 0 {
		obj.RoundTripperCount = int64(v)
	}

	return obj, nil
}

func expandClusterRKEConfigCloudProviderVsphereNetwork(p []interface{}) (*managementClient.NetworkVshpereOpts, error) {
	obj := &managementClient.NetworkVshpereOpts{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["public_network"].(string); ok && len(v) > 0 {
		obj.PublicNetwork = v
	}

	return obj, nil
}

func expandClusterRKEConfigCloudProviderVsphereVirtualCenter(p []interface{}) (map[string]managementClient.VirtualCenterConfig, error) {
	if len(p) == 0 || p[0] == nil {
		return map[string]managementClient.VirtualCenterConfig{}, nil
	}

	obj := make(map[string]managementClient.VirtualCenterConfig)

	for i := range p {
		in := p[i].(map[string]interface{})
		aux := managementClient.VirtualCenterConfig{}
		key := in["name"].(string)

		if v, ok := in["datacenters"].(string); ok && len(v) > 0 {
			aux.Datacenters = v
		}

		if v, ok := in["password"].(string); ok && len(v) > 0 {
			aux.Password = v
		}

		if v, ok := in["port"].(string); ok && len(v) > 0 {
			aux.VCenterPort = v
		}

		if v, ok := in["user"].(string); ok && len(v) > 0 {
			aux.User = v
		}

		if v, ok := in["soap_roundtrip_count"].(int); ok && v > 0 {
			aux.RoundTripperCount = int64(v)
		}

		obj[key] = aux
	}

	return obj, nil
}

func expandClusterRKEConfigCloudProviderVsphereWorkspace(p []interface{}) (*managementClient.WorkspaceVsphereOpts, error) {
	obj := &managementClient.WorkspaceVsphereOpts{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["datacenter"].(string); ok && len(v) > 0 {
		obj.Datacenter = v
	}

	if v, ok := in["folder"].(string); ok && len(v) > 0 {
		obj.Folder = v
	}

	if v, ok := in["server"].(string); ok && len(v) > 0 {
		obj.VCenterIP = v
	}

	if v, ok := in["default_datastore"].(string); ok && len(v) > 0 {
		obj.DefaultDatastore = v
	}

	if v, ok := in["resourcepool_path"].(string); ok && len(v) > 0 {
		obj.ResourcePoolPath = v
	}

	return obj, nil
}

func expandClusterRKEConfigCloudProviderVsphere(p []interface{}) (*managementClient.VsphereCloudProvider, error) {
	obj := &managementClient.VsphereCloudProvider{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["disk"].([]interface{}); ok && len(v) > 0 {
		disk, err := expandClusterRKEConfigCloudProviderVsphereDisk(v)
		if err != nil {
			return obj, err
		}
		obj.Disk = disk
	}

	if v, ok := in["global"].([]interface{}); ok && len(v) > 0 {
		global, err := expandClusterRKEConfigCloudProviderVsphereGlobal(v)
		if err != nil {
			return obj, err
		}
		obj.Global = global
	}

	if v, ok := in["network"].([]interface{}); ok && len(v) > 0 {
		network, err := expandClusterRKEConfigCloudProviderVsphereNetwork(v)
		if err != nil {
			return obj, err
		}
		obj.Network = network
	}

	if v, ok := in["virtual_center"].([]interface{}); ok && len(v) > 0 {
		vc, err := expandClusterRKEConfigCloudProviderVsphereVirtualCenter(v)
		if err != nil {
			return obj, err
		}
		obj.VirtualCenter = vc
	}

	if v, ok := in["workspace"].([]interface{}); ok && len(v) > 0 {
		workspace, err := expandClusterRKEConfigCloudProviderVsphereWorkspace(v)
		if err != nil {
			return obj, err
		}
		obj.Workspace = workspace
	}

	return obj, nil
}
