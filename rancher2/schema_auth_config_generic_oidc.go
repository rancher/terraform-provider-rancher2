package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const AuthConfigGenericOIDCName = "genericoidc"

//Schemas

func authConfigGenericOIDCFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"client_id": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "The OIDC Client ID.",
		},
		"client_secret": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "The OIDC Client Secret.",
		},
		"issuer": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The OIDC issuer URL.",
		},
		"rancher_url": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "The URL of the Rancher server. This is used as the redirect URI for the OIDC provider.",
		},
		"auth_endpoint": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The OIDC Auth Endpoint URL.",
		},
		"token_endpoint": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The OIDC Token Endpoint URL.",
		},
		"userinfo_endpoint": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The OIDC User Info Endpoint URL.",
		},
		"jwks_url": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The OIDC JWKS URL.",
		},
		"scopes": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The OIDC scopes to request. Defaults to `openid profile email`.",
		},
		"group_search_enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable group search.",
		},
		"groups_field": {
			Type:        schema.TypeString,
			Optional:    true,
			Computed:    true,
			Description: "The name of the OIDC claim to use for the user's group memberships.",
		},
		"certificate": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			StateFunc:   TrimSpace,
			Description: "A PEM-encoded CA certificate for the OIDC provider.",
		},
		"private_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
			StateFunc:   TrimSpace,
			Description: "A PEM-encoded private key for the OIDC provider.",
		},
	}

	for k, v := range authConfigFields() {
		s[k] = v
	}

	return s
}
