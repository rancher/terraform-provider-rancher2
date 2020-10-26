package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testNotifierMSTeamsConfigConf      *managementClient.MSTeamsConfig
	testNotifierMSTeamsConfigInterface []interface{}
)

func init() {
	testNotifierMSTeamsConfigConf = &managementClient.MSTeamsConfig{
		URL:      "url",
		ProxyURL: "proxy_url",
	}
	testNotifierMSTeamsConfigInterface = []interface{}{
		map[string]interface{}{
			"url":       "url",
			"proxy_url": "proxy_url",
		},
	}
}

func TestFlattenNotifierMSTeamsConfig(t *testing.T) {

	cases := []struct {
		Input          *managementClient.MSTeamsConfig
		ExpectedOutput []interface{}
	}{
		{
			testNotifierMSTeamsConfigConf,
			testNotifierMSTeamsConfigInterface,
		},
	}

	for _, tc := range cases {
		output := flattenNotifierMSTeamsConfig(tc.Input, tc.ExpectedOutput)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandNotifierMSTeamsConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.MSTeamsConfig
	}{
		{
			testNotifierMSTeamsConfigInterface,
			testNotifierMSTeamsConfigConf,
		},
	}

	for _, tc := range cases {
		output := expandNotifierMSTeamsConfig(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
