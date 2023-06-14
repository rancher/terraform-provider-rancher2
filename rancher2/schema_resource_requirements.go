package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func resourceRequirementFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cpu_limit": {
			Type:        schema.TypeString,
			Description: "The maximum CPU limit for agent",
			Optional:    true,
		},
		"cpu_request": {
			Type:        schema.TypeString,
			Description: "The minimum CPU required for agent",
			Optional:    true,
		},
		"memory_limit": {
			Type:        schema.TypeString,
			Description: "The maximum memory limit for agent",
			Optional:    true,
		},
		"memory_request": {
			Type:        schema.TypeString,
			Description: "The minimum memory required for agent",
			Optional:    true,
		},
	}
	return s
}
