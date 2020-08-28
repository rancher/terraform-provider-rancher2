package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testAlertEventRuleConf              *managementClient.EventRule
	testAlertEventRuleInterface         []interface{}
	testAlertMetricRuleConf             *managementClient.MetricRule
	testAlertMetricRuleInterface        []interface{}
	testAlertNodeRuleConf               *managementClient.NodeRule
	testAlertNodeRuleInterface          []interface{}
	testAlertPodRuleConf                *managementClient.PodRule
	testAlertPodRuleInterface           []interface{}
	testAlertSystemServiceRuleConf      *managementClient.SystemServiceRule
	testAlertSystemServiceRuleInterface []interface{}
	testAlertWorkloadRuleConf           *managementClient.WorkloadRule
	testAlertWorkloadRuleInterface      []interface{}
)

func init() {
	testAlertEventRuleConf = &managementClient.EventRule{
		EventType:    eventRuleTypeWarning,
		ResourceKind: eventRuleResourceKindNode,
	}
	testAlertEventRuleInterface = []interface{}{
		map[string]interface{}{
			"event_type":    eventRuleTypeWarning,
			"resource_kind": eventRuleResourceKindNode,
		},
	}
	testAlertMetricRuleConf = &managementClient.MetricRule{
		Comparison:     metricRuleComparisonEqual,
		Duration:       "30",
		Expression:     "expression",
		ThresholdValue: float64(3.5),
		Description:    "description",
	}
	testAlertMetricRuleInterface = []interface{}{
		map[string]interface{}{
			"comparison":      metricRuleComparisonEqual,
			"duration":        "30",
			"expression":      "expression",
			"threshold_value": float64(3.5),
			"description":     "description",
		},
	}
	testAlertNodeRuleConf = &managementClient.NodeRule{
		CPUThreshold: int64(70),
		Condition:    nodeRuleConditionNotReady,
		MemThreshold: int64(70),
		NodeID:       "node_id",
		Selector: map[string]string{
			"selector1": "selector1",
			"selector2": "selector2",
		},
	}
	testAlertNodeRuleInterface = []interface{}{
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
	testAlertPodRuleConf = &managementClient.PodRule{
		Condition:              podRuleConditionNotRunning,
		PodID:                  "pod_id",
		RestartIntervalSeconds: int64(70),
		RestartTimes:           int64(5),
	}
	testAlertPodRuleInterface = []interface{}{
		map[string]interface{}{
			"condition":                podRuleConditionNotRunning,
			"pod_id":                   "pod_id",
			"restart_interval_seconds": 70,
			"restart_times":            5,
		},
	}
	testAlertSystemServiceRuleConf = &managementClient.SystemServiceRule{
		Condition: systemServiceRuleConditionScheduler,
	}
	testAlertSystemServiceRuleInterface = []interface{}{
		map[string]interface{}{
			"condition": systemServiceRuleConditionScheduler,
		},
	}
	testAlertWorkloadRuleConf = &managementClient.WorkloadRule{
		AvailablePercentage: int64(70),
		Selector: map[string]string{
			"selector1": "selector1",
			"selector2": "selector2",
		},
		WorkloadID: "workload_id",
	}
	testAlertWorkloadRuleInterface = []interface{}{
		map[string]interface{}{
			"available_percentage": 70,
			"selector": map[string]interface{}{
				"selector1": "selector1",
				"selector2": "selector2",
			},
			"workload_id": "workload_id",
		},
	}
}

func TestFlattenEventRule(t *testing.T) {

	cases := []struct {
		Input          *managementClient.EventRule
		ExpectedOutput []interface{}
	}{
		{
			testAlertEventRuleConf,
			testAlertEventRuleInterface,
		},
	}

	for _, tc := range cases {
		output := flattenEventRule(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenMetricRule(t *testing.T) {

	cases := []struct {
		Input          *managementClient.MetricRule
		ExpectedOutput []interface{}
	}{
		{
			testAlertMetricRuleConf,
			testAlertMetricRuleInterface,
		},
	}

	for _, tc := range cases {
		output := flattenMetricRule(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenNodeRule(t *testing.T) {

	cases := []struct {
		Input          *managementClient.NodeRule
		ExpectedOutput []interface{}
	}{
		{
			testAlertNodeRuleConf,
			testAlertNodeRuleInterface,
		},
	}

	for _, tc := range cases {
		output := flattenNodeRule(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenPodRule(t *testing.T) {

	cases := []struct {
		Input          *managementClient.PodRule
		ExpectedOutput []interface{}
	}{
		{
			testAlertPodRuleConf,
			testAlertPodRuleInterface,
		},
	}

	for _, tc := range cases {
		output := flattenPodRule(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenSystemServiceRule(t *testing.T) {

	cases := []struct {
		Input          *managementClient.SystemServiceRule
		ExpectedOutput []interface{}
	}{
		{
			testAlertSystemServiceRuleConf,
			testAlertSystemServiceRuleInterface,
		},
	}

	for _, tc := range cases {
		output := flattenSystemServiceRule(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenWorkloadRule(t *testing.T) {

	cases := []struct {
		Input          *managementClient.WorkloadRule
		ExpectedOutput []interface{}
	}{
		{
			testAlertWorkloadRuleConf,
			testAlertWorkloadRuleInterface,
		},
	}

	for _, tc := range cases {
		output := flattenWorkloadRule(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandEventRule(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.EventRule
	}{
		{
			testAlertEventRuleInterface,
			testAlertEventRuleConf,
		},
	}

	for _, tc := range cases {
		output := expandEventRule(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandMetricRule(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.MetricRule
	}{
		{
			testAlertMetricRuleInterface,
			testAlertMetricRuleConf,
		},
	}

	for _, tc := range cases {
		output := expandMetricRule(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandNodeRule(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.NodeRule
	}{
		{
			testAlertNodeRuleInterface,
			testAlertNodeRuleConf,
		},
	}

	for _, tc := range cases {
		output := expandNodeRule(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandPodRule(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.PodRule
	}{
		{
			testAlertPodRuleInterface,
			testAlertPodRuleConf,
		},
	}

	for _, tc := range cases {
		output := expandPodRule(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandSystemServiceRule(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.SystemServiceRule
	}{
		{
			testAlertSystemServiceRuleInterface,
			testAlertSystemServiceRuleConf,
		},
	}

	for _, tc := range cases {
		output := expandSystemServiceRule(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandWorkloadRule(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.WorkloadRule
	}{
		{
			testAlertWorkloadRuleInterface,
			testAlertWorkloadRuleConf,
		},
	}

	for _, tc := range cases {
		output := expandWorkloadRule(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
