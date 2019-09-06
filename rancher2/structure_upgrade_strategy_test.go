package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testRollingUpdateConf        *managementClient.RollingUpdate
	testRollingUpdateInterface   []interface{}
	testUpgradeStrategyConf      *managementClient.UpgradeStrategy
	testUpgradeStrategyInterface []interface{}
)

func init() {
	testRollingUpdateConf = &managementClient.RollingUpdate{
		BatchSize: 10,
		Interval:  10,
	}
	testRollingUpdateInterface = []interface{}{
		map[string]interface{}{
			"batch_size": 10,
			"interval":   10,
		},
	}
	testUpgradeStrategyConf = &managementClient.UpgradeStrategy{
		RollingUpdate: testRollingUpdateConf,
	}
	testUpgradeStrategyInterface = []interface{}{
		map[string]interface{}{
			"rolling_update": testRollingUpdateInterface,
		},
	}
}

func TestFlattenRollingUpdate(t *testing.T) {

	cases := []struct {
		Input          *managementClient.RollingUpdate
		ExpectedOutput []interface{}
	}{
		{
			testRollingUpdateConf,
			testRollingUpdateInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRollingUpdate(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenUpgradeStrategy(t *testing.T) {

	cases := []struct {
		Input          *managementClient.UpgradeStrategy
		ExpectedOutput []interface{}
	}{
		{
			testUpgradeStrategyConf,
			testUpgradeStrategyInterface,
		},
	}

	for _, tc := range cases {
		output := flattenUpgradeStrategy(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRollingUpdate(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.RollingUpdate
	}{
		{
			testRollingUpdateInterface,
			testRollingUpdateConf,
		},
	}

	for _, tc := range cases {
		output := expandRollingUpdate(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandUpgradeStrategy(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.UpgradeStrategy
	}{
		{
			testUpgradeStrategyInterface,
			testUpgradeStrategyConf,
		},
	}

	for _, tc := range cases {
		output := expandUpgradeStrategy(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
