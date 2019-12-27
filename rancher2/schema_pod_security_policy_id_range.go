package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func podSecurityPolicyIDRangeFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"min": {
			Type:        schema.TypeInt,
			Description: "min is the start of the range, inclusive.",
			Required:    true,
		},
		"max": {
			Type:        schema.TypeInt,
			Description: "max is the end of the range, inclusive.",
			Required:    true,
		},
	}

	return s
}
