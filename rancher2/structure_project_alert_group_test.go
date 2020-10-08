package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testProjectAlertGroupRecipientsConf      []managementClient.Recipient
	testProjectAlertGroupRecipientsInterface []interface{}
	testProjectAlertGroupConf                *managementClient.ProjectAlertGroup
	testProjectAlertGroupInterface           map[string]interface{}
)

func init() {
	testProjectAlertGroupRecipientsConf = []managementClient.Recipient{
		{
			NotifierID:   "notifier_id",
			NotifierType: "webhook",
			Recipient:    "recipient",
		},
	}
	testProjectAlertGroupRecipientsInterface = []interface{}{
		map[string]interface{}{
			"notifier_id":       "notifier_id",
			"notifier_type":     "webhook",
			"recipient":         "recipient",
			"default_recipient": false,
		},
	}
	testProjectAlertGroupConf = &managementClient.ProjectAlertGroup{
		Name:                  "name",
		ProjectID:             "project_id",
		Description:           "description",
		GroupIntervalSeconds:  300,
		GroupWaitSeconds:      300,
		Recipients:            testProjectAlertGroupRecipientsConf,
		RepeatIntervalSeconds: 6000,
	}
	testProjectAlertGroupInterface = map[string]interface{}{
		"name":                    "name",
		"project_id":              "project_id",
		"description":             "description",
		"group_interval_seconds":  300,
		"group_wait_seconds":      300,
		"recipients":              testProjectAlertGroupRecipientsInterface,
		"repeat_interval_seconds": 6000,
	}
}

func TestFlattenProjectAlertGroup(t *testing.T) {

	cases := []struct {
		Input          *managementClient.ProjectAlertGroup
		ExpectedOutput map[string]interface{}
	}{
		{
			testProjectAlertGroupConf,
			testProjectAlertGroupInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, projectAlertGroupFields(), map[string]interface{}{})
		err := flattenProjectAlertGroup(output, tc.Input)
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

func TestExpandProjectAlertGroup(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.ProjectAlertGroup
	}{
		{
			testProjectAlertGroupInterface,
			testProjectAlertGroupConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, projectAlertGroupFields(), tc.Input)
		output := expandProjectAlertGroup(inputResourceData)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
