package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func monitoringInputFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"answers": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Answers for monitor input",
		},
		"version": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Monitoring version",
		},
	}

	return s
}
