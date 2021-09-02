package rancher2

import (
	"reflect"
	"testing"

	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
)

var (
	testClusterV2RKEConfigUpgradeStrategyDrainOptionsConf      rkev1.DrainOptions
	testClusterV2RKEConfigUpgradeStrategyDrainOptionsInterface []interface{}
	testClusterV2RKEConfigUpgradeStrategyConf                  rkev1.ClusterUpgradeStrategy
	testClusterV2RKEConfigUpgradeStrategyInterface             []interface{}
)

func init() {
	testClusterV2RKEConfigUpgradeStrategyDrainOptionsConf = rkev1.DrainOptions{
		Enabled:                         false,
		Force:                           true,
		IgnoreDaemonSets:                newTrue(),
		IgnoreErrors:                    true,
		DeleteEmptyDirData:              true,
		DisableEviction:                 true,
		GracePeriod:                     30,
		Timeout:                         20,
		SkipWaitForDeleteTimeoutSeconds: 10,
	}

	testClusterV2RKEConfigUpgradeStrategyDrainOptionsInterface = []interface{}{
		map[string]interface{}{
			"enabled":                              false,
			"force":                                true,
			"ignore_daemon_sets":                   true,
			"ignore_errors":                        true,
			"delete_empty_dir_data":                true,
			"disable_eviction":                     true,
			"grace_period":                         30,
			"timeout":                              20,
			"skip_wait_for_delete_timeout_seconds": 10,
		},
	}
	testClusterV2RKEConfigUpgradeStrategyConf = rkev1.ClusterUpgradeStrategy{
		ControlPlaneConcurrency:  "control_plane_concurrency",
		ControlPlaneDrainOptions: testClusterV2RKEConfigUpgradeStrategyDrainOptionsConf,
		WorkerConcurrency:        "worker_concurrency",
		WorkerDrainOptions:       testClusterV2RKEConfigUpgradeStrategyDrainOptionsConf,
	}

	testClusterV2RKEConfigUpgradeStrategyInterface = []interface{}{
		map[string]interface{}{
			"control_plane_concurrency":   "control_plane_concurrency",
			"control_plane_drain_options": testClusterV2RKEConfigUpgradeStrategyDrainOptionsInterface,
			"worker_concurrency":          "worker_concurrency",
			"worker_drain_options":        testClusterV2RKEConfigUpgradeStrategyDrainOptionsInterface,
		},
	}
}

func TestFlattenClusterV2RKEConfigUpgradeStrategyDrainOptions(t *testing.T) {

	cases := []struct {
		Input          rkev1.DrainOptions
		ExpectedOutput []interface{}
	}{
		{
			testClusterV2RKEConfigUpgradeStrategyDrainOptionsConf,
			testClusterV2RKEConfigUpgradeStrategyDrainOptionsInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterV2RKEConfigUpgradeStrategyDrainOptions(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterV2RKEConfigUpgradeStrategy(t *testing.T) {

	cases := []struct {
		Input          rkev1.ClusterUpgradeStrategy
		ExpectedOutput []interface{}
	}{
		{
			testClusterV2RKEConfigUpgradeStrategyConf,
			testClusterV2RKEConfigUpgradeStrategyInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterV2RKEConfigUpgradeStrategy(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterV2RKEConfigUpgradeStrategyDrainOptions(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput rkev1.DrainOptions
	}{
		{
			testClusterV2RKEConfigUpgradeStrategyDrainOptionsInterface,
			testClusterV2RKEConfigUpgradeStrategyDrainOptionsConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterV2RKEConfigUpgradeStrategyDrainOptions(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterV2RKEConfigUpgradeStrategy(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput rkev1.ClusterUpgradeStrategy
	}{
		{
			testClusterV2RKEConfigUpgradeStrategyInterface,
			testClusterV2RKEConfigUpgradeStrategyConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterV2RKEConfigUpgradeStrategy(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
