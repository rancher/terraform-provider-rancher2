package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testNotifierSMTPConfigConf      *managementClient.SMTPConfig
	testNotifierSMTPConfigInterface []interface{}
)

func init() {
	testNotifierSMTPConfigConf = &managementClient.SMTPConfig{
		DefaultRecipient: "default_recipient",
		Host:             "url",
		Port:             int64(25),
		Sender:           "sender",
		Password:         "password",
		TLS:              newTrue(),
		Username:         "username",
	}
	testNotifierSMTPConfigInterface = []interface{}{
		map[string]interface{}{
			"default_recipient": "default_recipient",
			"host":              "host",
			"port":              25,
			"sender":            "sender",
			"password":          "password",
			"tls":               newTrue(),
			"username":          "username",
		},
	}
}

func TestFlattenNotifierSMTPConfig(t *testing.T) {

	cases := []struct {
		Input          *managementClient.SMTPConfig
		ExpectedOutput []interface{}
	}{
		{
			testNotifierSMTPConfigConf,
			testNotifierSMTPConfigInterface,
		},
	}

	for _, tc := range cases {
		output := flattenNotifierSMTPConfig(tc.Input, testNotifierSMTPConfigInterface)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandNotifierSMTPConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.SMTPConfig
	}{
		{
			testNotifierSMTPConfigInterface,
			testNotifierSMTPConfigConf,
		},
	}

	for _, tc := range cases {
		output := expandNotifierSMTPConfig(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
