package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testLoggingSyslogConf      *managementClient.SyslogConfig
	testLoggingSyslogInterface []interface{}
)

func init() {
	testLoggingSyslogConf = &managementClient.SyslogConfig{
		Endpoint:    "hostname.test",
		Certificate: "XXXXXXXX",
		ClientCert:  "YYYYYYYY",
		ClientKey:   "ZZZZZZZZ",
		Program:     "program",
		Protocol:    "tcp",
		Severity:    "emergency",
		SSLVerify:   true,
		Token:       "XXXXXXXXXXXX",
	}
	testLoggingSyslogInterface = []interface{}{
		map[string]interface{}{
			"endpoint":    "hostname.test",
			"certificate": "XXXXXXXX",
			"client_cert": "YYYYYYYY",
			"client_key":  "ZZZZZZZZ",
			"program":     "program",
			"protocol":    "tcp",
			"severity":    "emergency",
			"ssl_verify":  true,
			"token":       "XXXXXXXXXXXX",
		},
	}
}

func TestFlattenLoggingSyslogConfig(t *testing.T) {

	cases := []struct {
		Input          *managementClient.SyslogConfig
		ExpectedOutput []interface{}
	}{
		{
			testLoggingSyslogConf,
			testLoggingSyslogInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenLoggingSyslogConfig(tc.Input, []interface{}{})
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandLoggingSyslogConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.SyslogConfig
	}{
		{
			testLoggingSyslogInterface,
			testLoggingSyslogConf,
		},
	}

	for _, tc := range cases {
		output, err := expandLoggingSyslogConfig(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
