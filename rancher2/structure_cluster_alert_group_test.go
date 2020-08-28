package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testClusterAlertGroupRecipientsConf      []managementClient.Recipient
	testClusterAlertGroupRecipientsInterface []interface{}
	testClusterAlertGroupConf                *managementClient.ClusterAlertGroup
	testClusterAlertGroupInterface           map[string]interface{}
)

func init() {
	testClusterAlertGroupRecipientsConf = []managementClient.Recipient{
		{
			NotifierID:   "notifier_id",
			NotifierType: "webhook",
			Recipient:    "recipient",
		},
	}
	testClusterAlertGroupRecipientsInterface = []interface{}{
		map[string]interface{}{
			"notifier_id":       "notifier_id",
			"notifier_type":     "webhook",
			"recipient":         "recipient",
			"default_recipient": false,
		},
	}
	testClusterAlertGroupConf = &managementClient.ClusterAlertGroup{
		Name:                  "name",
		ClusterID:             "cluster_id",
		Description:           "description",
		GroupIntervalSeconds:  300,
		GroupWaitSeconds:      300,
		Recipients:            testClusterAlertGroupRecipientsConf,
		RepeatIntervalSeconds: 6000,
	}
	testClusterAlertGroupInterface = map[string]interface{}{
		"name":                    "name",
		"cluster_id":              "cluster_id",
		"description":             "description",
		"group_interval_seconds":  300,
		"group_wait_seconds":      300,
		"recipients":              testClusterAlertGroupRecipientsInterface,
		"repeat_interval_seconds": 6000,
	}
}

func TestFlattenClusterAlertGroup(t *testing.T) {

	cases := []struct {
		Input          *managementClient.ClusterAlertGroup
		ExpectedOutput map[string]interface{}
	}{
		{
			testClusterAlertGroupConf,
			testClusterAlertGroupInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, clusterAlertGroupFields(), map[string]interface{}{})
		err := flattenClusterAlertGroup(output, tc.Input)
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

func TestExpandClusterAlertGroup(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.ClusterAlertGroup
	}{
		{
			testClusterAlertGroupInterface,
			testClusterAlertGroupConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, clusterAlertGroupFields(), tc.Input)
		output := expandClusterAlertGroup(inputResourceData)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
