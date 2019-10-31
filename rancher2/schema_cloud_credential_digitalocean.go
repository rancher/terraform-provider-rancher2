package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Types

type digitaloceanCredentialConfig struct {
	AccessToken string `json:"accessToken,omitempty" yaml:"accessToken,omitempty"`
}

//Schemas

func cloudCredentialDigitaloceanFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"access_token": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Digital Ocean access token",
		},
	}

	return s
}
