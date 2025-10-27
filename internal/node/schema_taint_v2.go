package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	taintV2EffectNoExecute        = "NoExecute"
	taintV2EffectNoSchedule       = "NoSchedule"
	taintV2EffectPreferNoSchedule = "PreferNoSchedule"
)

var (
	taintV2EffectTypes = []string{taintV2EffectNoExecute, taintV2EffectNoSchedule, taintV2EffectPreferNoSchedule}
)

//Schemas

func taintV2Fields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"key": {
			Type:     schema.TypeString,
			Required: true,
		},
		"value": {
			Type:     schema.TypeString,
			Required: true,
		},
		"effect": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      taintEffectNoSchedule,
			ValidateFunc: validation.StringInSlice(taintV2EffectTypes, true),
		},
	}

	return s
}
