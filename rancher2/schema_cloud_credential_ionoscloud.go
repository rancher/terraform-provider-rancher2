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
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			Description:  "Ionos Cloud API token. Mutually exclusive with username/password",
			AtLeastOneOf: []string{"ionoscloud_credential_config.0.token", "ionoscloud_credential_config.0.username"},
			ConflictsWith: []string{
				"ionoscloud_credential_config.0.username",
				"ionoscloud_credential_config.0.password",
			},
		},
		"username": {
			Type:          schema.TypeString,
			Optional:      true,
			Sensitive:     true,
			Description:   "Ionos Cloud username. Requires password; mutually exclusive with token",
			RequiredWith:  []string{"ionoscloud_credential_config.0.password"},
			ConflictsWith: []string{"ionoscloud_credential_config.0.token"},
		},
		"password": {
			Type:          schema.TypeString,
			Optional:      true,
			Sensitive:     true,
			Description:   "Ionos Cloud user password. Requires username; mutually exclusive with token",
			RequiredWith:  []string{"ionoscloud_credential_config.0.username"},
			ConflictsWith: []string{"ionoscloud_credential_config.0.token"},
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
