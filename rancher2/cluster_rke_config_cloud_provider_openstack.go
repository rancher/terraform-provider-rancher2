package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	openstackLBMonitorDelay      = 60
	openstackLBMonitorMaxRetries = 5
	openstackLBMonitorTimeout    = 30
	cloudProviderOpenstackName   = "openstack"
)

//Schemas

func openstackBlockStorageCloudProviderFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"bs_version": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"ignore_volume_az": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"trust_device_path": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
	}

	return s
}

func openstackGlobalCloudProviderFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"auth_url": {
			Type:     schema.TypeString,
			Required: true,
		},
		"password": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"tenant_id": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"user_id": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"username": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"ca_file": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"domain_id": {
			Type:      schema.TypeString,
			Optional:  true,
			Computed:  true,
			Sensitive: true,
		},
		"domain_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"region": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"tenant_name": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"trust_id": {
			Type:      schema.TypeString,
			Optional:  true,
			Computed:  true,
			Sensitive: true,
		},
	}
	return s
}

func openstackLoadBalancerCloudProviderFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"create_monitor": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"floating_network_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"lb_method": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"lb_provider": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"lb_version": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"manage_security_groups": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"monitor_delay": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  openstackLBMonitorDelay,
		},
		"monitor_max_retries": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  openstackLBMonitorMaxRetries,
		},
		"monitor_timeout": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  openstackLBMonitorTimeout,
		},
		"subnet_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"use_octavia": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

func openstackMetadataCloudProviderFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"request_timeout": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"search_order": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

func openstackRouteCloudProviderFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"router_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

func openstackCloudProviderFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"block_storage": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: openstackBlockStorageCloudProviderFields(),
			},
		},
		"global": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &schema.Resource{
				Schema: openstackGlobalCloudProviderFields(),
			},
		},
		"load_balancer": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: openstackLoadBalancerCloudProviderFields(),
			},
		},
		"metadata": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: openstackMetadataCloudProviderFields(),
			},
		},
		"route": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: openstackRouteCloudProviderFields(),
			},
		},
	}
	return s
}

// Flatteners

func flattenOpenstackBlockStorageCloudProvider(in *managementClient.BlockStorageOpenstackOpts) ([]interface{}, error) {
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

func flattenOpenstackGlobalCloudProvider(in *managementClient.GlobalOpenstackOpts) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.AuthURL) > 0 {
		obj["auth_url"] = in.AuthURL
	}

	if len(in.Password) > 0 {
		obj["password"] = in.Password
	}

	if len(in.TenantID) > 0 {
		obj["tenant_id"] = in.TenantID
	}

	if len(in.UserID) > 0 {
		obj["user_id"] = in.UserID
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

	if len(in.TenantName) > 0 {
		obj["tenant_name"] = in.TenantName
	}

	if len(in.TrustID) > 0 {
		obj["trust_id"] = in.TrustID
	}

	return []interface{}{obj}, nil
}

func flattenOpenstackLoadBalancerCloudProvider(in *managementClient.LoadBalancerOpenstackOpts) ([]interface{}, error) {
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

	if in.MonitorDelay > 0 {
		obj["monitor_delay"] = int(in.MonitorDelay)
	}

	if in.MonitorMaxRetries > 0 {
		obj["monitor_max_retries"] = int(in.MonitorMaxRetries)
	}

	if in.MonitorTimeout > 0 {
		obj["monitor_timeout"] = int(in.MonitorTimeout)
	}

	if len(in.SubnetID) > 0 {
		obj["subnet_id"] = in.SubnetID
	}

	obj["use_octavia"] = in.UseOctavia

	return []interface{}{obj}, nil
}

func flattenOpenstackMetadataCloudProvider(in *managementClient.MetadataOpenstackOpts) ([]interface{}, error) {
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

func flattenOpenstackRouteCloudProvider(in *managementClient.RouteOpenstackOpts) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.RouterID) > 0 {
		obj["router_id"] = in.RouterID
	}

	return []interface{}{obj}, nil
}

func flattenOpenstackCloudProvider(in *managementClient.OpenstackCloudProvider) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if in.BlockStorage != nil {
		blockStorage, err := flattenOpenstackBlockStorageCloudProvider(in.BlockStorage)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["block_storage"] = blockStorage
	}

	if in.Global != nil {
		global, err := flattenOpenstackGlobalCloudProvider(in.Global)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["global"] = global
	}

	if in.LoadBalancer != nil {
		loadBalancer, err := flattenOpenstackLoadBalancerCloudProvider(in.LoadBalancer)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["load_balancer"] = loadBalancer
	}

	if in.Metadata != nil {
		metadata, err := flattenOpenstackMetadataCloudProvider(in.Metadata)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["metadata"] = metadata
	}

	if in.Route != nil {
		route, err := flattenOpenstackRouteCloudProvider(in.Route)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["route"] = route
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandOpenstackBlockStorageCloudProvider(p []interface{}) (*managementClient.BlockStorageOpenstackOpts, error) {
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

func expandOpenstackGlobalCloudProvider(p []interface{}) (*managementClient.GlobalOpenstackOpts, error) {
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

	if v, ok := in["tenant_id"].(string); ok && len(v) > 0 {
		obj.TenantID = v
	}

	if v, ok := in["user_id"].(string); ok && len(v) > 0 {
		obj.UserID = v
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

	if v, ok := in["tenant_name"].(string); ok && len(v) > 0 {
		obj.TenantName = v
	}

	if v, ok := in["trust_id"].(string); ok && len(v) > 0 {
		obj.TrustID = v
	}

	return obj, nil
}

func expandOpenstackLoadBalancerCloudProvider(p []interface{}) (*managementClient.LoadBalancerOpenstackOpts, error) {
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

	if v, ok := in["monitor_delay"].(int); ok && v > 0 {
		obj.MonitorDelay = int64(v)
	}

	if v, ok := in["monitor_max_retries"].(int); ok && v > 0 {
		obj.MonitorMaxRetries = int64(v)
	}

	if v, ok := in["monitor_timeout"].(int); ok && v > 0 {
		obj.MonitorTimeout = int64(v)
	}

	if v, ok := in["subnet_id"].(string); ok && len(v) > 0 {
		obj.SubnetID = v
	}

	if v, ok := in["use_octavia"].(bool); ok {
		obj.UseOctavia = v
	}

	return obj, nil
}

func expandOpenstackMetadataCloudProvider(p []interface{}) (*managementClient.MetadataOpenstackOpts, error) {
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

func expandOpenstackRouteCloudProvider(p []interface{}) (*managementClient.RouteOpenstackOpts, error) {
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

func expandOpenstackCloudProvider(p []interface{}) (*managementClient.OpenstackCloudProvider, error) {
	obj := &managementClient.OpenstackCloudProvider{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["block_storage"].([]interface{}); ok && len(v) > 0 {
		blockStorage, err := expandOpenstackBlockStorageCloudProvider(v)
		if err != nil {
			return obj, err
		}
		obj.BlockStorage = blockStorage
	}

	if v, ok := in["global"].([]interface{}); ok && len(v) > 0 {
		global, err := expandOpenstackGlobalCloudProvider(v)
		if err != nil {
			return obj, err
		}
		obj.Global = global
	}

	if v, ok := in["load_balancer"].([]interface{}); ok && len(v) > 0 {
		loadBalancer, err := expandOpenstackLoadBalancerCloudProvider(v)
		if err != nil {
			return obj, err
		}
		obj.LoadBalancer = loadBalancer
	}

	if v, ok := in["metadata"].([]interface{}); ok && len(v) > 0 {
		metadata, err := expandOpenstackMetadataCloudProvider(v)
		if err != nil {
			return obj, err
		}
		obj.Metadata = metadata
	}

	if v, ok := in["route"].([]interface{}); ok && len(v) > 0 {
		route, err := expandOpenstackRouteCloudProvider(v)
		if err != nil {
			return obj, err
		}
		obj.Route = route
	}

	return obj, nil
}
