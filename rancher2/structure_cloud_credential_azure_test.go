package rancher2

import (
	"reflect"
	"testing"
)

var (
	testCloudCredentialAzureConf      *azureCredentialConfig
	testCloudCredentialAzureInterface []interface{}
)

func init() {
	testCloudCredentialAzureConf = &azureCredentialConfig{
		ClientID:       "client_id",
		ClientSecret:   "client_secret",
		Environment:    "environment",
		SubscriptionID: "subscription_id",
		TenantID:       "tenant_id",
	}
	testCloudCredentialAzureInterface = []interface{}{
		map[string]interface{}{
			"client_id":       "client_id",
			"client_secret":   "client_secret",
			"environment":     "environment",
			"subscription_id": "subscription_id",
			"tenant_id":       "tenant_id",
		},
	}
}

func TestFlattenCloudCredentialAzure(t *testing.T) {

	cases := []struct {
		Input          *azureCredentialConfig
		ExpectedOutput []interface{}
	}{
		{
			testCloudCredentialAzureConf,
			testCloudCredentialAzureInterface,
		},
	}

	for _, tc := range cases {
		output := flattenCloudCredentialAzure(tc.Input, tc.ExpectedOutput)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandCloudCredentialAzure(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *azureCredentialConfig
	}{
		{
			testCloudCredentialAzureInterface,
			testCloudCredentialAzureConf,
		},
	}

	for _, tc := range cases {
		output := expandCloudCredentialAzure(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
