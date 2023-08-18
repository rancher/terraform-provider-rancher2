package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
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
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")

	}
}
