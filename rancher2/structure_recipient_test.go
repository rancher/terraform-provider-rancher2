package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testRecipientsConf      []managementClient.Recipient
	testRecipientsInterface []interface{}
)

func init() {
	testRecipientsConf = []managementClient.Recipient{
		{
			NotifierID:   "notifier_id",
			NotifierType: "webhook",
			Recipient:    "recipient",
		},
	}
	testRecipientsInterface = []interface{}{
		map[string]interface{}{
			"notifier_id":   "notifier_id",
			"notifier_type": "webhook",
			"recipient":     "recipient",
		},
	}
}

func TestFlattenRecipients(t *testing.T) {

	cases := []struct {
		Input          []managementClient.Recipient
		ExpectedOutput []interface{}
	}{
		{
			testRecipientsConf,
			testRecipientsInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRecipients(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandRecipients(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []managementClient.Recipient
	}{
		{
			testRecipientsInterface,
			testRecipientsConf,
		},
	}

	for _, tc := range cases {
		output := expandRecipients(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
