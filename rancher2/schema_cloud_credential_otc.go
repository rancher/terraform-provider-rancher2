package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Types

type openTelekomCloudCredentialConfig struct {
	UserName string `json:"userName,omitempty" yaml:"userName,omitempty"`
	Password string `json:"password,omitempty" yaml:"password,omitempty"`

	AccessKey string `json:"accessKey,omitempty" yaml:"accessKey,omitempty"`
	SecretKey string `json:"secretKey,omitempty" yaml:"secretKey,omitempty"`

	Token string `json:"token,omitempty" yaml:"token,omitempty"`
}

//Schemas

func cloudCredentialOpenTelekomCloudFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"user_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Open Telekom Cloud username",
		},
		"password": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Open Telekom Cloud password",
		},
		"access_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Open Telekom Cloud access key",
		},
		"secret_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Open Telekom Cloud secret key",
		},
		"token": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			Description: "Open Telekom Cloud token",
		},
	}

	return s
}
