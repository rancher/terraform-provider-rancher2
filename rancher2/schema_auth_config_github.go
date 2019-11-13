package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const AuthConfigGithubName = "github"

//Schemas

func authConfigGithubFields() map[string]*schema.Schema {
	r := authConfigFields()
	s := map[string]*schema.Schema{
		"client_id": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"client_secret": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"hostname": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "github.com",
		},
		"tls": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
	}

	for k, v := range r {
		s[k] = v
	}

	return s
}
