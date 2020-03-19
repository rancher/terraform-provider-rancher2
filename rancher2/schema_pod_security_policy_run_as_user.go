package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	RunAsUserStrategyMustRunAs        = "MustRunAs"
	RunAsUserStrategyMustRunAsNonRoot = "MustRunAsNonRoot"
	RunAsUserStrategyRunAsAny         = "RunAsAny"
)

var (
	runAsUserStrategies = []string{
		RunAsUserStrategyMustRunAs,
		RunAsUserStrategyMustRunAsNonRoot,
		RunAsUserStrategyRunAsAny,
	}
)

//Schemas

func podSecurityPolicyRunAsUserFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"range": {
			Type:        schema.TypeList,
			Description: "ranges are the allowed ranges of uids that may be used. If you would like to force a single uid then supply a single range with the same start and end. Required for MustRunAs.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: podSecurityPolicyIDRangeFields(),
			},
		},
		"rule": {
			Type:         schema.TypeString,
			Description:  "rule is the strategy that will dictate the allowable RunAsUser values that may be set.",
			Required:     true,
			ValidateFunc: validation.StringInSlice(runAsUserStrategies, true),
		},
	}

	return s
}
