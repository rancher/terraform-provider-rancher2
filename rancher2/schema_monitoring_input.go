package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func monitoringInputFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"answers": &schema.Schema{
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Answers for monitor input",
		},
	}

	return s
}
