package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	networkPluginCalicoName  = "calico"
	networkPluginCanalName   = "canal"
	networkPluginFlannelName = "flannel"
	networkPluginWeaveName   = "weave"
)

var (
	networkPluginDefault = networkPluginCanalName
	networkPluginList    = []string{networkPluginCalicoName, networkPluginCanalName, networkPluginFlannelName, networkPluginWeaveName}
)

//Schemas

func clusterRKEConfigNetworkCalicoFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cloud_provider": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

func clusterRKEConfigNetworkCanalFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"iface": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

func clusterRKEConfigNetworkFlannelFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"iface": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}
	return s
}

func clusterRKEConfigNetworkWeaveFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"password": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
	return s
}

func clusterRKEConfigNetworkFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"calico_network_provider": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigNetworkCalicoFields(),
			},
		},
		"canal_network_provider": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigNetworkCanalFields(),
			},
		},
		"flannel_network_provider": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigNetworkFlannelFields(),
			},
		},
		"weave_network_provider": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigNetworkWeaveFields(),
			},
		},
		"options": {
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"plugin": {
			Type:         schema.TypeString,
			Optional:     true,
			Computed:     true,
			ValidateFunc: validation.StringInSlice(networkPluginList, true),
		},
	}
	return s
}
