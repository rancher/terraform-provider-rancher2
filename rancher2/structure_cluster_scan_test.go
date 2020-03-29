package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testClusterScanCisConfigConf      *managementClient.CisScanConfig
	testClusterScanCisConfigInterface []interface{}
	testClusterScanConfigConf         *managementClient.ClusterScanConfig
	testClusterScanConfigInterface    []interface{}
	testClusterScanConf               *managementClient.ClusterScan
	testClusterScanInterface          map[string]interface{}
)

func init() {
	testClusterScanCisConfigConf = &managementClient.CisScanConfig{
		DebugMaster:              true,
		DebugWorker:              true,
		OverrideBenchmarkVersion: "override_benchmark_version",
		OverrideSkip:             []string{"skip1", "skip2"},
		Profile:                  "profile",
	}
	testClusterScanCisConfigInterface = []interface{}{
		map[string]interface{}{
			"debug_master":               true,
			"debug_worker":               true,
			"override_benchmark_version": "override_benchmark_version",
			"override_skip":              []interface{}{"skip1", "skip2"},
			"profile":                    "profile",
		},
	}
	testClusterScanConfigConf = &managementClient.ClusterScanConfig{
		CisScanConfig: testClusterScanCisConfigConf,
	}
	testClusterScanConfigInterface = []interface{}{
		map[string]interface{}{
			"cis_scan_config": testClusterScanCisConfigInterface,
		},
	}
	testClusterScanConf = &managementClient.ClusterScan{
		ClusterID:  "cluster-test",
		Name:       "test",
		RunType:    "run_type",
		ScanConfig: testClusterScanConfigConf,
		ScanType:   "scan_type",
	}
	testClusterScanInterface = map[string]interface{}{
		"cluster_id":  "cluster-test",
		"name":        "test",
		"run_type":    "run_type",
		"scan_config": testClusterScanConfigInterface,
		"scan_type":   "scan_type",
	}
}

func TestFlattenClusterScanCisConfig(t *testing.T) {

	cases := []struct {
		Input          *managementClient.CisScanConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterScanCisConfigConf,
			testClusterScanCisConfigInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterScanCisConfig(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterScanConfig(t *testing.T) {

	cases := []struct {
		Input          *managementClient.ClusterScanConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterScanConfigConf,
			testClusterScanConfigInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterScanConfig(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterScan(t *testing.T) {

	cases := []struct {
		Input          *managementClient.ClusterScan
		ExpectedOutput map[string]interface{}
	}{
		{
			testClusterScanConf,
			testClusterScanInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, clusterScanFields(), tc.ExpectedOutput)
		err := flattenClusterScan(output, tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		if !reflect.DeepEqual(expectedOutput, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, expectedOutput)
		}
	}
}

func TestExpandClusterScanCisConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.CisScanConfig
	}{
		{
			testClusterScanCisConfigInterface,
			testClusterScanCisConfigConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterScanCisConfig(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterScanConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.ClusterScanConfig
	}{
		{
			testClusterScanConfigInterface,
			testClusterScanConfigConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterScanConfig(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterScan(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.ClusterScan
	}{
		{
			testClusterScanInterface,
			testClusterScanConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, clusterScanFields(), tc.Input)
		output := expandClusterScan(inputResourceData)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
