package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	tolerationEffectNoExecute        = "NoExecute"
	tolerationEffectNoSchedule       = "NoSchedule"
	tolerationEffectPreferNoSchedule = "PreferNoSchedule"
	tolerationOperatorEqual          = "Equal"
	tolerationOperatorExists         = "Exists"
)

var (
	tolerationEffectTypes   = []string{tolerationEffectNoExecute, tolerationEffectNoSchedule, tolerationEffectPreferNoSchedule}
	tolerationOperatorTypes = []string{tolerationOperatorEqual, tolerationOperatorExists}
)

//Schemas

func tolerationFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"key": {
			Type:     schema.TypeString,
			Required: true,
		},
		"effect": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      tolerationEffectNoSchedule,
			ValidateFunc: validation.StringInSlice(tolerationEffectTypes, true),
		},
		"operator": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      tolerationOperatorEqual,
			ValidateFunc: validation.StringInSlice(tolerationOperatorTypes, true),
		},
		"seconds": {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		"value": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	return s
}
