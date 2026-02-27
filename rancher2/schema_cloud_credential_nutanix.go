package rancher2

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

// Types

type nutanixCredentialConfig struct {
	Endpoint string `json:"endpoint,omitempty" yaml:"endpoint,omitempty"`
	Username string `json:"username,omitempty" yaml:"username,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
	Port     string `json:"port,omitempty" yaml:"port,omitempty"`
}

// Schemas

func cloudCredentialNutanixFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"endpoint": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Nutanix management endpoint IP address/FQDN",
		},
		"username": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Nutanix management username",
		},
		"password": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Nutanix management password",
		},
		"port": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "9440",
			Description: "Nutanix management endpoint port",
		},
	}

	return s
}
