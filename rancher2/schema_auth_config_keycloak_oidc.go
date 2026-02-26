package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const AuthConfigKeyCloakOIDCName = "keycloakoidc"

//Schemas

func authConfigKeyCloakOIDCFields() map[string]*schema.Schema {
	return oidcSchemaFields()
}
