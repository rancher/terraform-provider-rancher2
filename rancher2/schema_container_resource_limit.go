package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func containerResourceLimitFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"limits_cpu": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"limits_memory": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"requests_cpu": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"requests_memory": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	return s
}
