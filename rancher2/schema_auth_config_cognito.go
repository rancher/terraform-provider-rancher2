package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const AuthConfigCognitoName = "cognito"

//Schemas

func authConfigCognitoFields() map[string]*schema.Schema {
	fields := oidcSchemaFields()

	fields["name_claim"].Default = "cognito:username"

	// Terraform doesn't allow defaults for Computed fields.
	//
	// It's not clear why this field is Computed in the Schema, but we'll set it
	// to false so we can provide a default.
	fields["groups_field"].Computed = false
	fields["groups_field"].Default = "cognito:groups"

	return fields
}
