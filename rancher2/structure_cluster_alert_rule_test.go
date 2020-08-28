package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testClusterAlertEventRuleConf              *managementClient.EventRule
	testClusterAlertEventRuleInterface         []interface{}
	testClusterAlertMetricRuleConf             *managementClient.MetricRule
	testClusterAlertMetricRuleInterface        []interface{}
	testClusterAlertNodeRuleConf               *managementClient.NodeRule
	testClusterAlertNodeRuleInterface          []interface{}
	testClusterAlertSystemServiceRuleConf      *managementClient.SystemServiceRule
	testClusterAlertSystemServiceRuleInterface []interface{}
	testClusterAlertRuleConf                   *managementClient.ClusterAlertRule
	testClusterAlertRuleInterface              map[string]interface{}
)

func init() {
	testClusterAlertEventRuleConf = &managementClient.EventRule{
		EventType:    eventRuleTypeWarning,
		ResourceKind: eventRuleResourceKindNode,
	}
	testClusterAlertEventRuleInterface = []interface{}{
		map[string]interface{}{
			"event_type":    eventRuleTypeWarning,
			"resource_kind": eventRuleResourceKindNode,
		},
	}
	testClusterAlertMetricRuleConf = &managementClient.MetricRule{
		Comparison:     metricRuleComparisonEqual,
		Duration:       "30",
		Expression:     "expression",
		ThresholdValue: float64(3.5),
		Description:    "description",
	}
	testClusterAlertMetricRuleInterface = []interface{}{
		map[string]interface{}{
			"comparison":      metricRuleComparisonEqual,
			"duration":        "30",
			"expression":      "expression",
			"threshold_value": float64(3.5),
			"description":     "description",
		},
	}
	testClusterAlertNodeRuleConf = &managementClient.NodeRule{
		CPUThreshold: int64(70),
		Condition:    nodeRuleConditionNotReady,
		MemThreshold: int64(70),
		NodeID:       "node_id",
		Selector: map[string]string{
			"selector1": "selector1",
			"selector2": "selector2",
		},
	}
	testClusterAlertNodeRuleInterface = []interface{}{
		map[string]interface{}{
			"cpu_threshold": 70,
			"condition":     nodeRuleConditionNotReady,
			"mem_threshold": 70,
			"node_id":       "node_id",
			"selector": map[string]interface{}{
				"selector1": "selector1",
				"selector2": "selector2",
			},
		},
	}
	testClusterAlertSystemServiceRuleConf = &managementClient.SystemServiceRule{
		Condition: systemServiceRuleConditionScheduler,
	}
	testClusterAlertSystemServiceRuleInterface = []interface{}{
		map[string]interface{}{
			"condition": systemServiceRuleConditionScheduler,
		},
	}
	testClusterAlertRuleConf = &managementClient.ClusterAlertRule{
		Name:                  "name",
		ClusterID:             "cluster_id",
		EventRule:             testClusterAlertEventRuleConf,
		GroupID:               "group_id",
		GroupIntervalSeconds:  300,
		GroupWaitSeconds:      300,
		Inherited:             newTrue(),
		MetricRule:            testClusterAlertMetricRuleConf,
		NodeRule:              testClusterAlertNodeRuleConf,
		RepeatIntervalSeconds: 6000,
		Severity:              alertRuleSeverityCritical,
		SystemServiceRule:     testClusterAlertSystemServiceRuleConf,
	}
	testClusterAlertRuleInterface = map[string]interface{}{
		"name":                    "name",
		"cluster_id":              "cluster_id",
		"event_rule":              testClusterAlertEventRuleInterface,
		"group_id":                "group_id",
		"group_interval_seconds":  300,
		"group_wait_seconds":      300,
		"inherited":               true,
		"metric_rule":             testClusterAlertMetricRuleInterface,
		"node_rule":               testClusterAlertNodeRuleInterface,
		"repeat_interval_seconds": 6000,
		"severity":                alertRuleSeverityCritical,
		"system_service_rule":     testClusterAlertSystemServiceRuleInterface,
	}
}

func TestFlattenClusterAlertRule(t *testing.T) {

	cases := []struct {
		Input          *managementClient.ClusterAlertRule
		ExpectedOutput map[string]interface{}
	}{
		{
			testClusterAlertRuleConf,
			testClusterAlertRuleInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, clusterAlertRuleFields(), map[string]interface{}{})
		err := flattenClusterAlertRule(output, tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		if !reflect.DeepEqual(expectedOutput, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				expectedOutput, tc.ExpectedOutput)
		}
	}
}

func TestExpandClusterAlertRule(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.ClusterAlertRule
	}{
		{
			testClusterAlertRuleInterface,
			testClusterAlertRuleConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, clusterAlertRuleFields(), tc.Input)
		output := expandClusterAlertRule(inputResourceData)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
