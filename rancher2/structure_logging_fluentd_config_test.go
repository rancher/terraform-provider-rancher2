package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testLoggingFluentdConfigFluentServerConf      []managementClient.FluentServer
	testLoggingFluentdConfigFluentServerInterface []interface{}
	testLoggingFluentdConf                        *managementClient.FluentForwarderConfig
	testLoggingFluentdInterface                   []interface{}
)

func init() {
	testLoggingFluentdConfigFluentServerConf = []managementClient.FluentServer{
		managementClient.FluentServer{
			Endpoint:  "host.terraform.test",
			Hostname:  "hostname",
			Password:  "YYYYYYYY",
			SharedKey: "ZZZZZZZZ",
			Standby:   false,
			Username:  "user",
			Weight:    5,
		},
	}
	testLoggingFluentdConfigFluentServerInterface = []interface{}{
		map[string]interface{}{
			"endpoint":   "host.terraform.test",
			"hostname":   "hostname",
			"password":   "YYYYYYYY",
			"shared_key": "ZZZZZZZZ",
			"standby":    false,
			"username":   "user",
			"weight":     5,
		},
	}
	testLoggingFluentdConf = &managementClient.FluentForwarderConfig{
		FluentServers: testLoggingFluentdConfigFluentServerConf,
		Certificate:   "XXXXXXXX",
		Compress:      true,
		EnableTLS:     true,
	}
	testLoggingFluentdInterface = []interface{}{
		map[string]interface{}{
			"fluent_servers": testLoggingFluentdConfigFluentServerInterface,
			"certificate":    "XXXXXXXX",
			"compress":       true,
			"enable_tls":     true,
		},
	}
}

func TestFlattenLoggingFluentdConfig(t *testing.T) {

	cases := []struct {
		Input          *managementClient.FluentForwarderConfig
		ExpectedOutput []interface{}
	}{
		{
			testLoggingFluentdConf,
			testLoggingFluentdInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenLoggingFluentdConfig(tc.Input, []interface{}{})
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandLoggingFluentdConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.FluentForwarderConfig
	}{
		{
			testLoggingFluentdInterface,
			testLoggingFluentdConf,
		},
	}

	for _, tc := range cases {
		output, err := expandLoggingFluentdConfig(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
