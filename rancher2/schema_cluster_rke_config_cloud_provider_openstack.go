package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	openstackLBMonitorDelay      = "60s"
	openstackLBMonitorMaxRetries = 5
	openstackLBMonitorTimeout    = "30s"
)

//Schemas

func clusterRKEConfigCloudProviderOpenstackBlockStorageFields() map[string]*schema.Schema {
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

func clusterRKEConfigCloudProviderOpenstackGlobalFields() map[string]*schema.Schema {
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
		"tenant_id": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
			Computed:  true,
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

func clusterRKEConfigCloudProviderOpenstackLoadBalancerFields() map[string]*schema.Schema {
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
			Type:     schema.TypeString,
			Optional: true,
			Default:  openstackLBMonitorDelay,
		},
		"monitor_max_retries": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  openstackLBMonitorMaxRetries,
		},
		"monitor_timeout": {
			Type:     schema.TypeString,
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

func clusterRKEConfigCloudProviderOpenstackMetadataFields() map[string]*schema.Schema {
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

func clusterRKEConfigCloudProviderOpenstackRouteFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"router_id": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

func clusterRKEConfigCloudProviderOpenstackFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"global": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigCloudProviderOpenstackGlobalFields(),
			},
		},
		"block_storage": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigCloudProviderOpenstackBlockStorageFields(),
			},
		},
		"load_balancer": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigCloudProviderOpenstackLoadBalancerFields(),
			},
		},
		"metadata": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigCloudProviderOpenstackMetadataFields(),
			},
		},
		"route": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigCloudProviderOpenstackRouteFields(),
			},
		},
	}
	return s
}
