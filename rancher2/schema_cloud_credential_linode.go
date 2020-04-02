package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Types

type linodeCredentialConfig struct {
	Token string `json:"token,omitempty" yaml:"token,omitempty"`
}

// Schemas

func cloudCredentialLinodeFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"token": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Linode API token",
		},
	}

	return s
}
