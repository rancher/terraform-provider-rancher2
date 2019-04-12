package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testClusterRKEConfigMonitoringConf      *managementClient.MonitoringConfig
	testClusterRKEConfigMonitoringInterface []interface{}
)

func init() {
	testClusterRKEConfigMonitoringConf = &managementClient.MonitoringConfig{
		Options: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
		Provider: "test",
	}
	testClusterRKEConfigMonitoringInterface = []interface{}{
		map[string]interface{}{
			"options": map[string]interface{}{
				"option1": "value1",
				"option2": "value2",
			},
			"provider": "test",
		},
	}
}

func TestFlattenClusterRKEConfigMonitoring(t *testing.T) {

	cases := []struct {
		Input          *managementClient.MonitoringConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigMonitoringConf,
			testClusterRKEConfigMonitoringInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigMonitoring(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigMonitoring(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.MonitoringConfig
	}{
		{
			testClusterRKEConfigMonitoringInterface,
			testClusterRKEConfigMonitoringConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigMonitoring(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
