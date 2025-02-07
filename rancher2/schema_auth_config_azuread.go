package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const AuthConfigAzureADName = "azuread"

//Schemas

func authConfigAzureADFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"application_id": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"application_secret": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"auth_endpoint": {
			Type:     schema.TypeString,
			Required: true,
		},
		"endpoint": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "https://login.microsoftonline.com/",
		},
		"graph_endpoint": {
			Type:     schema.TypeString,
			Required: true,
		},
		"rancher_url": {
			Type:     schema.TypeString,
			Required: true,
		},
		"tenant_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"token_endpoint": {
			Type:     schema.TypeString,
			Required: true,
		},
	}

	for k, v := range authConfigFields() {
		s[k] = v
	}

	return s
}
