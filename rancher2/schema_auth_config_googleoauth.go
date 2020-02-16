package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const AuthConfigGoogleOauthName = "googleoauth"

//Schemas

func authConfigGoogleOauthFields() map[string]*schema.Schema {
	r := authConfigFields()
	s := map[string]*schema.Schema{
		"admin_email": {
			Type:     schema.TypeString,
			Required: true,
		},
		"hostname": {
			Type:     schema.TypeString,
			Required: true,
		},
		"oauth_credential": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"service_account_credential": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"nested_group_membership_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
	}

	for k, v := range r {
		s[k] = v
	}

	return s
}
