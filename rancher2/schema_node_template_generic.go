package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
)

type genericNodeTemplateConfig struct {
	driverName string
	driverID   string
	config     map[string]interface{}
}

//Schemas

func genericNodeConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"driver": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Custom NodeDriver Name",
		},
		"config": {
			Type:        schema.TypeMap,
			Required:    true,
			Description: "Driver config for custom node driver",
		},
	}

	return s
}
