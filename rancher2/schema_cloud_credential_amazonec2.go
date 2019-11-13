package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Types

type amazonec2CredentialConfig struct {
	AccessKey string `json:"accessKey,omitempty" yaml:"accessKey,omitempty"`
	SecretKey string `json:"secretKey,omitempty" yaml:"secretKey,omitempty"`
}

//Schemas

func cloudCredentialAmazonec2Fields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"access_key": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "AWS Access Key",
		},
		"secret_key": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "AWS Secret Key",
		},
	}

	return s
}
