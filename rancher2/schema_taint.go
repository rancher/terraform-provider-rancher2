package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	taintEffectNoExecute        = "NoExecute"
	taintEffectNoSchedule       = "NoSchedule"
	taintEffectPreferNoSchedule = "PreferNoSchedule"
)

var (
	taintEffectTypes = []string{taintEffectNoExecute, taintEffectNoSchedule, taintEffectPreferNoSchedule}
)

//Schemas

func taintFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"key": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"value": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"effect": &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			Default:      taintEffectNoSchedule,
			ValidateFunc: validation.StringInSlice(taintEffectTypes, true),
		},
		"time_added": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
	}

	return s
}
