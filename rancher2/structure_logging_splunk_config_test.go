package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testLoggingSplunkConf      *managementClient.SplunkConfig
	testLoggingSplunkInterface []interface{}
)

func init() {
	testLoggingSplunkConf = &managementClient.SplunkConfig{
		Endpoint:      "hostname.test",
		Certificate:   "XXXXXXXX",
		ClientCert:    "YYYYYYYY",
		ClientKey:     "ZZZZZZZZ",
		ClientKeyPass: "pass",
		Index:         "index",
		Source:        "source",
		SSLVerify:     true,
		Token:         "XXXXXXXXXXXX",
	}
	testLoggingSplunkInterface = []interface{}{
		map[string]interface{}{
			"endpoint":        "hostname.test",
			"certificate":     "XXXXXXXX",
			"client_cert":     "YYYYYYYY",
			"client_key":      "ZZZZZZZZ",
			"client_key_pass": "pass",
			"index":           "index",
			"source":          "source",
			"ssl_verify":      true,
			"token":           "XXXXXXXXXXXX",
		},
	}
}

func TestFlattenLoggingSplunkConfig(t *testing.T) {

	cases := []struct {
		Input          *managementClient.SplunkConfig
		ExpectedOutput []interface{}
	}{
		{
			testLoggingSplunkConf,
			testLoggingSplunkInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenLoggingSplunkConfig(tc.Input, []interface{}{})
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandLoggingSplunkConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.SplunkConfig
	}{
		{
			testLoggingSplunkInterface,
			testLoggingSplunkConf,
		},
	}

	for _, tc := range cases {
		output, err := expandLoggingSplunkConfig(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
