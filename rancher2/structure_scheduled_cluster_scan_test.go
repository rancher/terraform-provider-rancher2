package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testScheduledClusterScanConfigConf      *managementClient.ScheduledClusterScanConfig
	testScheduledClusterScanConfigInterface []interface{}
	testScheduledClusterScanConf            *managementClient.ScheduledClusterScan
	testScheduledClusterScanInterface       []interface{}
)

func init() {
	testScheduledClusterScanConfigConf = &managementClient.ScheduledClusterScanConfig{
		CronSchedule: "cron_schedule",
		Retention:    5,
	}
	testScheduledClusterScanConfigInterface = []interface{}{
		map[string]interface{}{
			"cron_schedule": "cron_schedule",
			"retention":     5,
		},
	}
	testScheduledClusterScanConf = &managementClient.ScheduledClusterScan{
		Enabled:        true,
		ScanConfig:     testClusterScanConfigConf,
		ScheduleConfig: testScheduledClusterScanConfigConf,
	}
	testScheduledClusterScanInterface = []interface{}{
		map[string]interface{}{
			"enabled":         true,
			"scan_config":     testClusterScanConfigInterface,
			"schedule_config": testScheduledClusterScanConfigInterface,
		},
	}
}

func TestFlattenScheduledClusterScanConfig(t *testing.T) {

	cases := []struct {
		Input          *managementClient.ScheduledClusterScanConfig
		ExpectedOutput []interface{}
	}{
		{
			testScheduledClusterScanConfigConf,
			testScheduledClusterScanConfigInterface,
		},
	}

	for _, tc := range cases {
		output := flattenScheduledClusterScanConfig(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenScheduledClusterScan(t *testing.T) {

	cases := []struct {
		Input          *managementClient.ScheduledClusterScan
		ExpectedOutput []interface{}
	}{
		{
			testScheduledClusterScanConf,
			testScheduledClusterScanInterface,
		},
	}

	for _, tc := range cases {
		output := flattenScheduledClusterScan(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandScheduledClusterScanConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.ScheduledClusterScanConfig
	}{
		{
			testScheduledClusterScanConfigInterface,
			testScheduledClusterScanConfigConf,
		},
	}

	for _, tc := range cases {
		output := expandScheduledClusterScanConfig(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandScheduledClusterScan(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.ScheduledClusterScan
	}{
		{
			testScheduledClusterScanInterface,
			testScheduledClusterScanConf,
		},
	}

	for _, tc := range cases {
		output := expandScheduledClusterScan(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
