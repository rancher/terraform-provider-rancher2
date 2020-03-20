package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	SELinuxStrategyMustRunAs = "MustRunAs"
	SELinuxStrategyRunAsAny  = "RunAsAny"
)

var (
	seLinuxStrategies = []string{
		SELinuxStrategyMustRunAs,
		SELinuxStrategyRunAsAny,
	}
)

//Schemas

func podSecurityPolicySELinuxFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"se_linux_option": {
			Type:        schema.TypeList,
			Description: "seLinuxOptions required to run as; required for MustRunAs. More info: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/",
			Optional:    true,
			MaxItems:    1,
			Elem: &schema.Resource{
				Schema: podSecurityPolicySELinuxOptionsFields(),
			},
		},
		"rule": {
			Type:         schema.TypeString,
			Description:  "rule is the strategy that will dictate the allowable labels that may be set.",
			Required:     true,
			ValidateFunc: validation.StringInSlice(seLinuxStrategies, true),
		},
	}

	return s
}
