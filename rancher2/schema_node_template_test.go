package rancher2

import (
	"encoding/json"
	"reflect"
	"testing"

	managementClient "github.com/rancher/types/client/management/v3"
)

func TestNodeTemplateMarshalJSON(t *testing.T) {
	cases := []struct {
		Input          *NodeTemplate
		ExpectedOutput map[string]interface{}
	}{
		{
			Input: &NodeTemplate{
				NodeTemplate: managementClient.NodeTemplate{
					Driver: "test",
				},
				genericConfig: &genericNodeTemplateConfig{
					driverID:   "test",
					driverName: "test",
					config: map[string]interface{}{
						"foo": "1",
						"bar": 2,
					},
				},
			},
			ExpectedOutput: map[string]interface{}{
				"actions": nil,
				"links":   nil,
				"driver":  "test",
				"testConfig": map[string]interface{}{
					"foo": "1",
					"bar": float64(2),
				},
			},
		},
		{
			Input: &NodeTemplate{
				NodeTemplate: managementClient.NodeTemplate{
					Driver: "azure",
				},
				AzureConfig: &azureConfig{
					ClientID:       "XXXXXXXXXXXXXXXXXXXX",
					ClientSecret:   "XXXXXXXXXXXXXXXXXXXX",
					SubscriptionID: "XXXXXXXXXXXXXXXXXXXX",
				},
			},
			ExpectedOutput: map[string]interface{}{
				"actions": nil,
				"links":   nil,
				"driver":  "azure",
				"azureConfig": map[string]interface{}{
					"clientId":       "XXXXXXXXXXXXXXXXXXXX",
					"clientSecret":   "XXXXXXXXXXXXXXXXXXXX",
					"subscriptionId": "XXXXXXXXXXXXXXXXXXXX",
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

func TestNodeTemplateUnmarshalJSON(t *testing.T) {
	cases := []struct {
		Input          string
		ExpectedOutput NodeTemplate
	}{
		{
			Input: `{"driver":"test","testConfig":{"bar":2,"foo":"test"}}`,
			ExpectedOutput: NodeTemplate{
				NodeTemplate: managementClient.NodeTemplate{
					Driver: "test",
				},
				genericConfig: &genericNodeTemplateConfig{
					driverName: "test",
					config: map[string]interface{}{
						"foo": "test",
						"bar": float64(2),
					},
				},
			},
		},
		{
			Input: `{
				"driver": "azure",
				"azureConfig": {
					"clientId": "XXXXXXXXXXXXXXXXXXXX",
			    	"clientSecret": "XXXXXXXXXXXXXXXXXXXX",
    				"subscriptionId": "XXXXXXXXXXXXXXXXXXXX"
				}
			}`,
			ExpectedOutput: NodeTemplate{
				NodeTemplate: managementClient.NodeTemplate{
					Driver: "azure",
				},
				AzureConfig: &azureConfig{
					ClientID:       "XXXXXXXXXXXXXXXXXXXX",
					ClientSecret:   "XXXXXXXXXXXXXXXXXXXX",
					SubscriptionID: "XXXXXXXXXXXXXXXXXXXX",
				},
			},
		},
	}

	for _, tc := range cases {
		var nt NodeTemplate
		err := json.Unmarshal([]byte(tc.Input), &nt)
		if err != nil {
			t.Fatalf("Unexpected error: %s", err)
		}

		if !reflect.DeepEqual(nt, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from unmarshaler .\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, nt)
		}
	}
}
