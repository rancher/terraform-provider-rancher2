package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const AuthConfigOpenLdapName = "openldap"

//Schemas

func authConfigOpenLdapFields() map[string]*schema.Schema {
	return authConfigLdapFields()
}
