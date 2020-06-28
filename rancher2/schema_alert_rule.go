package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

const (
	alertRuleSeverityCritical                   = "critical"
	alertRuleSeverityInfo                       = "info"
	alertRuleSeverityWarning                    = "warning"
	eventRuleTypeWarning                        = "Warning"
	eventRuleTypeNormal                         = "Normal"
	eventRuleResourceKindDaemondSet             = "DaemonSet"
	eventRuleResourceKindDeployment             = "Deployment"
	eventRuleResourceKindNode                   = "Node"
	eventRuleResourceKindPod                    = "Pod"
	eventRuleResourceKindStatefulSet            = "StatefulSet"
	metricRuleComparisonEqual                   = "equal"
	metricRuleComparisonGreaterOrEqual          = "greater-or-equal"
	metricRuleComparisonGreaterThan             = "greater-than"
	metricRuleComparisonLessOrEqual             = "less-or-equal"
	metricRuleComparisonLessThan                = "less-than"
	metricRuleComparisonNotEqual                = "not-equal"
	metricRuleComparisonNotNull                 = "has-value"
	nodeRuleConditionCPU                        = "cpu"
	nodeRuleConditionMem                        = "mem"
	nodeRuleConditionNotReady                   = "notready"
	podRuleConditionNotRunning                  = "notrunning"
	podRuleConditionNotScheduled                = "notscheduled"
	podRuleConditionRestarts                    = "restarts"
	systemServiceRuleConditionControllerManager = "controller-manager"
	systemServiceRuleConditionEtcd              = "etcd"
	systemServiceRuleConditionScheduler         = "scheduler"
)

var (
	alertRuleSeverityTypes = []string{alertRuleSeverityCritical, alertRuleSeverityInfo, alertRuleSeverityWarning}
	eventRuleTypes         = []string{eventRuleTypeNormal, eventRuleTypeWarning}
	eventRuleResourceKinds = []string{
		eventRuleResourceKindDaemondSet,
		eventRuleResourceKindDeployment,
		eventRuleResourceKindNode,
		eventRuleResourceKindPod,
		eventRuleResourceKindStatefulSet,
	}
	metricRuleComparisons = []string{
		metricRuleComparisonEqual,
		metricRuleComparisonGreaterOrEqual,
		metricRuleComparisonGreaterThan,
		metricRuleComparisonLessOrEqual,
		metricRuleComparisonLessThan,
		metricRuleComparisonNotEqual,
		metricRuleComparisonNotNull,
	}
	nodeRuleConditions = []string{
		nodeRuleConditionCPU,
		nodeRuleConditionMem,
		nodeRuleConditionNotReady,
	}
	podRuleConditions = []string{
		podRuleConditionNotRunning,
		podRuleConditionNotScheduled,
		podRuleConditionRestarts,
	}
	systemServiceRuleConditions = []string{
		systemServiceRuleConditionControllerManager,
		systemServiceRuleConditionEtcd,
		systemServiceRuleConditionScheduler,
	}
)

//Schemas

func alertRuleFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"group_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Alert rule group ID",
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Alert rule name",
		},
		"group_interval_seconds": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     180,
			Description: "Alert rule interval seconds",
		},
		"group_wait_seconds": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     180,
			Description: "Alert rule wait seconds",
		},
		"inherited": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
			Description: "Alert rule inherited",
		},
		"repeat_interval_seconds": {
			Type:        schema.TypeInt,
			Optional:    true,
			Default:     3600,
			Description: "Alert rule repeat interval seconds",
		},
		"severity": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      alertRuleSeverityCritical,
			Description:  "Alert rule severity",
			ValidateFunc: validation.StringInSlice(alertRuleSeverityTypes, true),
		},
	}

	for k, v := range commonAnnotationLabelFields() {
		s[k] = v
	}

	return s
}

func eventRuleFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"event_type": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      eventRuleTypeWarning,
			Description:  "Event type",
			ValidateFunc: validation.StringInSlice(eventRuleTypes, true),
		},
		"resource_kind": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "Resource kind",
			ValidateFunc: validation.StringInSlice(eventRuleResourceKinds, true),
		},
	}
	return s
}

func metricRuleFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"comparison": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      metricRuleComparisonEqual,
			Description:  "Metric rule comparison",
			ValidateFunc: validation.StringInSlice(metricRuleComparisons, true),
		},
		"duration": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Metric rule duration",
		},
		"expression": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Metric rule expression",
		},
		"threshold_value": {
			Type:        schema.TypeFloat,
			Required:    true,
			Description: "Metric rule threshold value",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Metric rule description",
		},
	}
	return s
}

func nodeRuleFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cpu_threshold": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      70,
			Description:  "Node rule cpu threshold",
			ValidateFunc: validation.IntBetween(1, 100),
		},
		"condition": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      nodeRuleConditionNotReady,
			Description:  "Node rule condition",
			ValidateFunc: validation.StringInSlice(nodeRuleConditions, true),
		},
		"mem_threshold": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      70,
			Description:  "Node rule mem threshold",
			ValidateFunc: validation.IntBetween(1, 100),
		},
		"node_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Node ID",
		},
		"selector": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Node rule selector",
		},
	}
	return s
}

func podRuleFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"pod_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Pod ID",
		},
		"condition": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      podRuleConditionNotRunning,
			Description:  "Pod rule condition",
			ValidateFunc: validation.StringInSlice(podRuleConditions, true),
		},
		"restart_interval_seconds": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      300,
			Description:  "Pod rule restart interval seconds",
			ValidateFunc: validation.IntAtLeast(1),
		},
		"restart_times": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      3,
			Description:  "Pod rule restart times",
			ValidateFunc: validation.IntAtLeast(1),
		},
	}
	return s
}

func systemServiceRuleFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"condition": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      systemServiceRuleConditionScheduler,
			Description:  "System service rule condition",
			ValidateFunc: validation.StringInSlice(systemServiceRuleConditions, true),
		},
	}
	return s
}

func workloadRuleFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"available_percentage": {
			Type:         schema.TypeInt,
			Optional:     true,
			Default:      70,
			Description:  "Workload rule available percentage",
			ValidateFunc: validation.IntBetween(1, 100),
		},
		"selector": {
			Type:        schema.TypeMap,
			Optional:    true,
			Description: "Workload rule selector",
		},
		"workload_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Workload ID",
		},
	}
	return s
}
