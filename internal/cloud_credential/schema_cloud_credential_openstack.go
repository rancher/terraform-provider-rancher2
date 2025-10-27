package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Types

type openstackCredentialConfig struct {
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
}

//Schemas

func cloudCredentialOpenstackFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"password": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "OpenStack password",
		},
	}

	return s
}
