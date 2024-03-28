package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
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
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
