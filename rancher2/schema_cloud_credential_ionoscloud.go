package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Types

type ionoscloudCredentialConfig struct {
	Token    string `json:"token,omitempty" yaml:"token,omitempty"`
	Username string `json:"username,omitempty" yaml:"username,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`
	Endpoint string `json:"endpoint,omitempty" yaml:"endpoint,omitempty"`
}

//Schemas

func cloudCredentialIonoscloudFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"token": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Ionos Cloud API token",
		},
		"username": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Ionos Cloud username",
		},
		"password": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Ionos Cloud user password",
		},
		"endpoint": {
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "https://api.ionos.com/cloudapi/v6",
			Description: "IonosCloud API endpoint",
		},
	}

	return s
}
