package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Types

type amazonec2CredentialConfig struct {
	AccessKey     string `json:"accessKey,omitempty" yaml:"accessKey,omitempty"`
	SecretKey     string `json:"secretKey,omitempty" yaml:"secretKey,omitempty"`
	DefaultRegion string `json:"defaultRegion,omitempty" yaml:"defaultRegion,omitempty"`
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
		"default_region": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "AWS default region",
		},
	}

	return s
}
