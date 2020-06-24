package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func projectAlertGroupFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"project_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Alert group Project ID",
		},
	}

	for k, v := range alertGroupFields() {
		s[k] = v
	}

	return s
}
