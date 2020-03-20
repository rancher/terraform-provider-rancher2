package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	RunAsGroupStrategyMustRunAs        = "MustRunAs"
	RunAsGroupStrategyMustRunAsNonRoot = "MustRunAsNonRoot"
	RunAsGroupStrategyRunAsAny         = "RunAsAny"
)

var (
	runAsGroupStrategies = []string{
		RunAsGroupStrategyMustRunAs,
		RunAsGroupStrategyMustRunAsNonRoot,
		RunAsGroupStrategyRunAsAny,
	}
)

//Schemas

func podSecurityPolicyRunAsGroupFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"range": {
			Type:        schema.TypeList,
			Description: "ranges are the allowed ranges of gids that may be used. If you would like to force a single gid then supply a single range with the same start and end. Required for MustRunAs.",
			Optional:    true,
			Elem: &schema.Resource{
				Schema: podSecurityPolicyIDRangeFields(),
			},
		},
		"rule": {
			Type:         schema.TypeString,
			Description:  "rule is the strategy that will dictate the allowable RunAsGroup values that may be set.",
			Required:     true,
			ValidateFunc: validation.StringInSlice(runAsGroupStrategies, true),
		},
	}

	return s
}
