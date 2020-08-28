package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testLoggingCustomTargetConf      *managementClient.CustomTargetConfig
	testLoggingCustomTargetInterface []interface{}
)

func init() {
	testLoggingCustomTargetConf = &managementClient.CustomTargetConfig{
		Certificate: "certificate",
		ClientCert:  "client_cert",
		ClientKey:   "client_key",
		Content:     "content",
	}
	testLoggingCustomTargetInterface = []interface{}{
		map[string]interface{}{
			"certificate": "certificate",
			"client_cert": "client_cert",
			"client_key":  "client_key",
			"content":     "content",
		},
	}
}

func TestFlattenLoggingCustomTargetConfig(t *testing.T) {

	cases := []struct {
		Input          *managementClient.CustomTargetConfig
		ExpectedOutput []interface{}
	}{
		{
			testLoggingCustomTargetConf,
			testLoggingCustomTargetInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenLoggingCustomTargetConfig(tc.Input, []interface{}{})
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandLoggingCustomTargetConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.CustomTargetConfig
	}{
		{
			testLoggingCustomTargetInterface,
			testLoggingCustomTargetConf,
		},
	}

	for _, tc := range cases {
		output, err := expandLoggingCustomTargetConfig(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
