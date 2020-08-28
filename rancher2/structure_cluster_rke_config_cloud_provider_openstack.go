package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenClusterRKEConfigCloudProviderOpenstackBlockStorage(in *managementClient.BlockStorageOpenstackOpts) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.BSVersion) > 0 {
		obj["bs_version"] = in.BSVersion
	}

	obj["ignore_volume_az"] = in.IgnoreVolumeAZ
	obj["trust_device_path"] = in.TrustDevicePath

	return []interface{}{obj}, nil
}

func flattenClusterRKEConfigCloudProviderOpenstackGlobal(in *managementClient.GlobalOpenstackOpts, p []interface{}) ([]interface{}, error) {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.AuthURL) > 0 {
		obj["auth_url"] = in.AuthURL
	}

	if len(in.Password) > 0 {
		obj["password"] = in.Password
	}

	if len(in.Username) > 0 {
		obj["username"] = in.Username
	}

	if len(in.CAFile) > 0 {
		obj["ca_file"] = in.CAFile
	}

	if len(in.DomainID) > 0 {
		obj["domain_id"] = in.DomainID
	}

	if len(in.DomainName) > 0 {
		obj["domain_name"] = in.DomainName
	}

	if len(in.Region) > 0 {
		obj["region"] = in.Region
	}

	if len(in.TenantID) > 0 {
		obj["tenant_id"] = in.TenantID
	}

	if len(in.TenantName) > 0 {
		obj["tenant_name"] = in.TenantName
	}

	if len(in.TrustID) > 0 {
		obj["trust_id"] = in.TrustID
	}

	return []interface{}{obj}, nil
}

func flattenClusterRKEConfigCloudProviderOpenstackLoadBalancer(in *managementClient.LoadBalancerOpenstackOpts) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	obj["create_monitor"] = in.CreateMonitor

	if len(in.FloatingNetworkID) > 0 {
		obj["floating_network_id"] = in.FloatingNetworkID
	}

	if len(in.LBMethod) > 0 {
		obj["lb_method"] = in.LBMethod
	}

	if len(in.LBProvider) > 0 {
		obj["lb_provider"] = in.LBProvider
	}

	if len(in.LBVersion) > 0 {
		obj["lb_version"] = in.LBVersion
	}

	obj["manage_security_groups"] = in.ManageSecurityGroups

	if len(in.MonitorDelay) > 0 {
		obj["monitor_delay"] = in.MonitorDelay
	}

	if in.MonitorMaxRetries > 0 {
		obj["monitor_max_retries"] = int(in.MonitorMaxRetries)
	}

	if len(in.MonitorTimeout) > 0 {
		obj["monitor_timeout"] = in.MonitorTimeout
	}

	if len(in.SubnetID) > 0 {
		obj["subnet_id"] = in.SubnetID
	}

	obj["use_octavia"] = in.UseOctavia

	return []interface{}{obj}, nil
}

func flattenClusterRKEConfigCloudProviderOpenstackMetadata(in *managementClient.MetadataOpenstackOpts) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if in.RequestTimeout > 0 {
		obj["request_timeout"] = int(in.RequestTimeout)
	}

	if len(in.SearchOrder) > 0 {
		obj["search_order"] = in.SearchOrder
	}

	return []interface{}{obj}, nil
}

func flattenClusterRKEConfigCloudProviderOpenstackRoute(in *managementClient.RouteOpenstackOpts) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.RouterID) > 0 {
		obj["router_id"] = in.RouterID
	}

	return []interface{}{obj}, nil
}

func flattenClusterRKEConfigCloudProviderOpenstack(in *managementClient.OpenstackCloudProvider, p []interface{}) ([]interface{}, error) {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}, nil
	}

	if in.BlockStorage != nil {
		blockStorage, err := flattenClusterRKEConfigCloudProviderOpenstackBlockStorage(in.BlockStorage)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["block_storage"] = blockStorage
	}

	if in.Global != nil {
		v, ok := obj["global"].([]interface{})
		if !ok {
			v = []interface{}{}
		}
		global, err := flattenClusterRKEConfigCloudProviderOpenstackGlobal(in.Global, v)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["global"] = global
	}

	if in.LoadBalancer != nil {
		loadBalancer, err := flattenClusterRKEConfigCloudProviderOpenstackLoadBalancer(in.LoadBalancer)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["load_balancer"] = loadBalancer
	}

	if in.Metadata != nil {
		metadata, err := flattenClusterRKEConfigCloudProviderOpenstackMetadata(in.Metadata)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["metadata"] = metadata
	}

	if in.Route != nil {
		route, err := flattenClusterRKEConfigCloudProviderOpenstackRoute(in.Route)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["route"] = route
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandClusterRKEConfigCloudProviderOpenstackBlockStorage(p []interface{}) (*managementClient.BlockStorageOpenstackOpts, error) {
	obj := &managementClient.BlockStorageOpenstackOpts{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["bs_version"].(string); ok && len(v) > 0 {
		obj.BSVersion = v
	}

	if v, ok := in["ignore_volume_az"].(bool); ok {
		obj.IgnoreVolumeAZ = v
	}

	if v, ok := in["trust_device_path"].(bool); ok {
		obj.TrustDevicePath = v
	}

	return obj, nil
}

func expandClusterRKEConfigCloudProviderOpenstackGlobal(p []interface{}) (*managementClient.GlobalOpenstackOpts, error) {
	obj := &managementClient.GlobalOpenstackOpts{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["auth_url"].(string); ok && len(v) > 0 {
		obj.AuthURL = v
	}

	if v, ok := in["password"].(string); ok && len(v) > 0 {
		obj.Password = v
	}

	if v, ok := in["username"].(string); ok && len(v) > 0 {
		obj.Username = v
	}

	if v, ok := in["ca_file"].(string); ok && len(v) > 0 {
		obj.CAFile = v
	}

	if v, ok := in["domain_id"].(string); ok && len(v) > 0 {
		obj.DomainID = v
	}

	if v, ok := in["domain_name"].(string); ok && len(v) > 0 {
		obj.DomainName = v
	}

	if v, ok := in["region"].(string); ok && len(v) > 0 {
		obj.Region = v
	}

	if v, ok := in["tenant_id"].(string); ok && len(v) > 0 {
		obj.TenantID = v
	}

	if v, ok := in["tenant_name"].(string); ok && len(v) > 0 {
		obj.TenantName = v
	}

	if v, ok := in["trust_id"].(string); ok && len(v) > 0 {
		obj.TrustID = v
	}

	return obj, nil
}

func expandClusterRKEConfigCloudProviderOpenstackLoadBalancer(p []interface{}) (*managementClient.LoadBalancerOpenstackOpts, error) {
	obj := &managementClient.LoadBalancerOpenstackOpts{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["create_monitor"].(bool); ok {
		obj.CreateMonitor = v
	}

	if v, ok := in["floating_network_id"].(string); ok && len(v) > 0 {
		obj.FloatingNetworkID = v
	}

	if v, ok := in["lb_method"].(string); ok && len(v) > 0 {
		obj.LBMethod = v
	}

	if v, ok := in["lb_provider"].(string); ok && len(v) > 0 {
		obj.LBProvider = v
	}

	if v, ok := in["lb_version"].(string); ok && len(v) > 0 {
		obj.LBVersion = v
	}

	if v, ok := in["manage_security_groups"].(bool); ok {
		obj.ManageSecurityGroups = v
	}

	if v, ok := in["monitor_delay"].(string); ok && len(v) > 0 {
		obj.MonitorDelay = v
	}

	if v, ok := in["monitor_max_retries"].(int); ok && v > 0 {
		obj.MonitorMaxRetries = int64(v)
	}

	if v, ok := in["monitor_timeout"].(string); ok && len(v) > 0 {
		obj.MonitorTimeout = v
	}

	if v, ok := in["subnet_id"].(string); ok && len(v) > 0 {
		obj.SubnetID = v
	}

	if v, ok := in["use_octavia"].(bool); ok {
		obj.UseOctavia = v
	}

	return obj, nil
}

func expandClusterRKEConfigCloudProviderOpenstackMetadata(p []interface{}) (*managementClient.MetadataOpenstackOpts, error) {
	obj := &managementClient.MetadataOpenstackOpts{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["request_timeout"].(int); ok && v > 0 {
		obj.RequestTimeout = int64(v)
	}

	if v, ok := in["search_order"].(string); ok && len(v) > 0 {
		obj.SearchOrder = v
	}

	return obj, nil
}

func expandClusterRKEConfigCloudProviderOpenstackRoute(p []interface{}) (*managementClient.RouteOpenstackOpts, error) {
	obj := &managementClient.RouteOpenstackOpts{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["router_id"].(string); ok && len(v) > 0 {
		obj.RouterID = v
	}

	return obj, nil
}

func expandClusterRKEConfigCloudProviderOpenstack(p []interface{}) (*managementClient.OpenstackCloudProvider, error) {
	obj := &managementClient.OpenstackCloudProvider{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["block_storage"].([]interface{}); ok && len(v) > 0 {
		blockStorage, err := expandClusterRKEConfigCloudProviderOpenstackBlockStorage(v)
		if err != nil {
			return obj, err
		}
		obj.BlockStorage = blockStorage
	}

	if v, ok := in["global"].([]interface{}); ok && len(v) > 0 {
		global, err := expandClusterRKEConfigCloudProviderOpenstackGlobal(v)
		if err != nil {
			return obj, err
		}
		obj.Global = global
	}

	if v, ok := in["load_balancer"].([]interface{}); ok && len(v) > 0 {
		loadBalancer, err := expandClusterRKEConfigCloudProviderOpenstackLoadBalancer(v)
		if err != nil {
			return obj, err
		}
		obj.LoadBalancer = loadBalancer
	}

	if v, ok := in["metadata"].([]interface{}); ok && len(v) > 0 {
		metadata, err := expandClusterRKEConfigCloudProviderOpenstackMetadata(v)
		if err != nil {
			return obj, err
		}
		obj.Metadata = metadata
	}

	if v, ok := in["route"].([]interface{}); ok && len(v) > 0 {
		route, err := expandClusterRKEConfigCloudProviderOpenstackRoute(v)
		if err != nil {
			return obj, err
		}
		obj.Route = route
	}

	return obj, nil
}
