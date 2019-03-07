package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

//Schemas

func bastionHostFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"address": {
			Type:     schema.TypeString,
			Required: true,
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
			Type:     schema.TypeString,
			Required: true,
		},
	}
	return s
}

// Flatteners

func flattenBastionHost(in *managementClient.BastionHost) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.Address) > 0 {
		obj["address"] = in.Address
	}

	if len(in.Port) > 0 {
		obj["port"] = in.Port
	}

	obj["ssh_agent_auth"] = in.SSHAgentAuth

	if len(in.SSHKey) > 0 {
		obj["ssh_key"] = in.SSHKey
	}

	if len(in.SSHKeyPath) > 0 {
		obj["ssh_key_path"] = in.SSHKeyPath
	}

	if len(in.User) > 0 {
		obj["user"] = in.User
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandBastionHost(p []interface{}) (*managementClient.BastionHost, error) {
	obj := &managementClient.BastionHost{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["address"].(string); ok && len(v) > 0 {
		obj.Address = v
	}

	if v, ok := in["port"].(string); ok && len(v) > 0 {
		obj.Port = v
	}

	if v, ok := in["ssh_agent_auth"].(bool); ok {
		obj.SSHAgentAuth = v
	}

	if v, ok := in["ssh_key"].(string); ok && len(v) > 0 {
		obj.SSHKey = v
	}

	if v, ok := in["ssh_key_path"].(string); ok && len(v) > 0 {
		obj.SSHKeyPath = v
	}

	if v, ok := in["user"].(string); ok && len(v) > 0 {
		obj.User = v
	}

	return obj, nil
}
