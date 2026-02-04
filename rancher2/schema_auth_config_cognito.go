package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const AuthConfigCognitoName = "cognito"

//Schemas

func authConfigCognitoFields() map[string]*schema.Schema {
	return oidcSchemaFields()
}
