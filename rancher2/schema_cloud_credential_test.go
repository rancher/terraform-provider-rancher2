package rancher2

import (
	"encoding/json"
	"reflect"
	"testing"

	managementClient "github.com/rancher/types/client/management/v3"
)

func TestCloudCredentialMarshalJSON(t *testing.T) {
	cases := []struct {
		Input          *CloudCredential
		ExpectedOutput map[string]interface{}
	}{
		{
			Input: &CloudCredential{
				CloudCredential: managementClient.CloudCredential{
					Name: "azure",
				},
				AzureCredentialConfig: &azureCredentialConfig{
					ClientID:       "XXXXXXXXXXXXXXXXXXXX",
					ClientSecret:   "XXXXXXXXXXXXXXXXXXXX",
					SubscriptionID: "XXXXXXXXXXXXXXXXXXXX",
				},
			},
			ExpectedOutput: map[string]interface{}{
				"name":    "azure",
				"actions": nil,
				"links":   nil,
				"azurecredentialConfig": map[string]interface{}{
					"clientId":       "XXXXXXXXXXXXXXXXXXXX",
					"clientSecret":   "XXXXXXXXXXXXXXXXXXXX",
					"subscriptionId": "XXXXXXXXXXXXXXXXXXXX",
				},
			},
		},
		{
			Input: &CloudCredential{
				CloudCredential: managementClient.CloudCredential{
					Name: "test",
				},
				genericCredentialConfig: &genericCredentialConfig{
					driverName: "test",
					config: map[string]interface{}{
						"foo": "1",
						"bar": 2,
					},
				},
			},
			ExpectedOutput: map[string]interface{}{
				"name":    "test",
				"actions": nil,
				"links":   nil,
				"testcredentialConfig": map[string]interface{}{
					"foo": "1",
					"bar": float64(2),
				},
			},
		},
	}

	for _, tc := range cases {
		data, err := json.Marshal(tc.Input)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		var marshaled map[string]interface{}
		err = json.Unmarshal(data, &marshaled)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
		if !reflect.DeepEqual(marshaled, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from marshaler.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, marshaled)
		}
	}
}

func TestCloudCredentialUnmarshalJSON(t *testing.T) {
	cases := []struct {
		Input          string
		ExpectedOutput CloudCredential
	}{
		{
			Input: `{"name":"test","testcredentialConfig":{"bar":2,"foo":"1"}}`,
			ExpectedOutput: CloudCredential{
				CloudCredential: managementClient.CloudCredential{
					Name: "test",
				},
				genericCredentialConfig: &genericCredentialConfig{
					driverName: "test",
					config: map[string]interface{}{
						"foo": "1",
						"bar": 2,
					},
				},
			},
		},
		{
			Input: `{
				"name":"azure",
				"azurecredentialConfig": {
					"clientId":"XXXXXXXXXXXXXXXXXXXX",
					"clientSecret":"XXXXXXXXXXXXXXXXXXXX",
					"subscriptionId":"XXXXXXXXXXXXXXXXXXXX"
				}
			}`,
			ExpectedOutput: CloudCredential{
				CloudCredential: managementClient.CloudCredential{
					Name: "azure",
				},
				AzureCredentialConfig: &azureCredentialConfig{
					ClientID:       "XXXXXXXXXXXXXXXXXXXX",
					ClientSecret:   "XXXXXXXXXXXXXXXXXXXX",
					SubscriptionID: "XXXXXXXXXXXXXXXXXXXX",
				},
			},
		},
	}

	for _, expect := range cases {
		var cc CloudCredential
		err := json.Unmarshal([]byte(expect.Input), &cc)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}
		if !reflect.DeepEqual(cc.CloudCredential, expect.ExpectedOutput.CloudCredential) {
			t.Fatalf("Unexpected output from unmarshaler .\nExpected: %#v\nGiven:    %#v",
				expect.ExpectedOutput, cc)
		}
	}
}
