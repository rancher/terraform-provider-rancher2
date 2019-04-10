package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
)

//Types

type genericCredentialConfig struct {
	driverName string
	driverID   string
	config     map[string]interface{}
}

//Schemas

func cloudCredentialGenericFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"driver": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "NodeDriver Name for generic node driver",
		},
		"config": {
			Type:        schema.TypeMap,
			Required:    true,
			Description: "Cloud credential config for generic node driver",
		},
	}

	return s
}
