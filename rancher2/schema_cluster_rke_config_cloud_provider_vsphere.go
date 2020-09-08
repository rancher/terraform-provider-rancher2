package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func clusterRKEConfigCloudProviderVsphereDiskFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"scsi_controller_type": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

func clusterRKEConfigCloudProviderVsphereGlobalFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"datacenters": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"insecure_flag": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"password": {
			Type:      schema.TypeString,
			Optional:  true,
			Computed:  true,
			Sensitive: true,
		},
		"user": {
			Type:      schema.TypeString,
			Optional:  true,
			Computed:  true,
			Sensitive: true,
		},
		"port": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"soap_roundtrip_count": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

func clusterRKEConfigCloudProviderVsphereNetworkFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"public_network": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

func clusterRKEConfigCloudProviderVsphereVirtualCenterFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"datacenters": {
			Type:     schema.TypeString,
			Required: true,
		},
		"name": {
			Type:     schema.TypeString,
			Required: true,
		},
		"password": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"user": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"port": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"soap_roundtrip_count": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

func clusterRKEConfigCloudProviderVsphereWorkspaceFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"datacenter": {
			Type:     schema.TypeString,
			Required: true,
		},
		"folder": {
			Type:     schema.TypeString,
			Required: true,
		},
		"server": {
			Type:     schema.TypeString,
			Required: true,
		},
		"default_datastore": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"resourcepool_path": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

func clusterRKEConfigCloudProviderVsphereFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"virtual_center": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigCloudProviderVsphereVirtualCenterFields(),
			},
		},
		"workspace": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigCloudProviderVsphereWorkspaceFields(),
			},
		},
		"disk": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigCloudProviderVsphereDiskFields(),
			},
		},
		"global": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigCloudProviderVsphereGlobalFields(),
			},
		},
		"network": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigCloudProviderVsphereNetworkFields(),
			},
		},
	}
	return s
}
