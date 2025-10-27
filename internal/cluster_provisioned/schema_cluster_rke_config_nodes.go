package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var (
	// RKEConfigNodesRoles available RKE roles for nodes
	RKEConfigNodesRoles = []string{"controlplane", "etcd", "worker"}
)

//Schemas

func clusterRKEConfigNodesFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"address": {
			Type:     schema.TypeString,
			Required: true,
		},
		"role": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringInSlice(RKEConfigNodesRoles, true),
			},
		},
		"docker_socket": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"hostname_override": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"internal_address": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"labels": {
			Type:     schema.TypeMap,
			Optional: true,
		},
		"node_id": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"port": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "22",
		},
		"ssh_agent_auth": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"ssh_key": {
			Type:      schema.TypeString,
			Optional:  true,
			Computed:  true,
			Sensitive: true,
		},
		"ssh_key_path": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"user": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
	}
	return s
}

func clusterRKEConfigNodeDrainInputFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"delete_local_data": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"force": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"grace_period": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  -1,
		},
		"ignore_daemon_sets": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"timeout": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      60,
			ValidateFunc: validation.IntBetween(1, 10800),
		},
	}
	return s
}

func clusterRKEConfigNodeUpgradeStrategyFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"drain": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"drain_input": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRKEConfigNodeDrainInputFields(),
			},
		},
		"max_unavailable_controlplane": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "1",
		},
		"max_unavailable_worker": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "10%",
		},
	}
	return s
}
