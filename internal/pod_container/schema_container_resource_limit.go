package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func containerResourceLimitFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"limits_cpu": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"limits_memory": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"requests_cpu": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"requests_memory": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	return s
}
