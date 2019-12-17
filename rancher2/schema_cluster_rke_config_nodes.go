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
		"labels": &schema.Schema{
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
