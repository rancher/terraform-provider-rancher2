package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	memberAccessTypeMember = "member"
	memberAccessTypeOwner  = "owner"
	memberAccessTypeRO     = "read-only"
)

var (
	memberAccessTypeKinds = []string{memberAccessTypeMember, memberAccessTypeOwner, memberAccessTypeRO}
)

//Schemas

func memberFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"access_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Description:  "Member access type: " + memberAccessTypeMember + ", " + memberAccessTypeOwner + ", " + memberAccessTypeRO,
			ValidateFunc: validation.StringInSlice(memberAccessTypeKinds, true),
		},
		"group_principal_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Member group principal id",
		},
		"user_principal_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Member user principal id",
		},
	}

	return s
}
