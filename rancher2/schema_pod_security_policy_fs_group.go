package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	fSGroupStrategyMayRunAs  = "MayRunAs"
	fSGroupStrategyMustRunAs = "MustRunAs"
	fSGroupStrategyRunAsAny  = "RunAsAny"
)

var (
	fsGroupStrategies = []string{
		fSGroupStrategyMayRunAs,
		fSGroupStrategyMustRunAs,
		fSGroupStrategyRunAsAny,
	}
)

//Schemas

func podSecurityPolicyAllowedFsGroupFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"range": {
			Type:        schema.TypeList,
			Description: "ranges are the allowed ranges of fs groups.  If you would like to force a single fs group then supply a single range with the same start and end. Required for MustRunAs.",
			Optional:    true,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: podSecurityPolicyIDRangeFields(),
			},
		},
		"rule": {
			Type:         schema.TypeString,
			Description:  "rule is the strategy that will dictate what FSGroup is used in the SecurityContext.",
			Optional:     true,
			ValidateFunc: validation.StringInSlice(fsGroupStrategies, true),
		},
	}

	return s
}
