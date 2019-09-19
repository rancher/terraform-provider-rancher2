package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testLoggingElasticsearchConf      *managementClient.ElasticsearchConfig
	testLoggingElasticsearchInterface []interface{}
)

func init() {
	testLoggingElasticsearchConf = &managementClient.ElasticsearchConfig{
		Endpoint:      "hostname.test",
		DateFormat:    "YYYY-MM-DD",
		IndexPrefix:   "index",
		AuthPassword:  "password",
		AuthUserName:  "user",
		Certificate:   "XXXXXXXX",
		ClientCert:    "YYYYYYYY",
		ClientKey:     "ZZZZZZZZ",
		ClientKeyPass: "pass",
		SSLVerify:     true,
		SSLVersion:    "test",
	}
	testLoggingElasticsearchInterface = []interface{}{
		map[string]interface{}{
			"endpoint":        "hostname.test",
			"date_format":     "YYYY-MM-DD",
			"index_prefix":    "index",
			"auth_password":   "password",
			"auth_username":   "user",
			"certificate":     "XXXXXXXX",
			"client_cert":     "YYYYYYYY",
			"client_key":      "ZZZZZZZZ",
			"client_key_pass": "pass",
			"ssl_verify":      true,
			"ssl_version":     "test",
		},
	}
}

func TestFlattenLoggingElasticsearchConfig(t *testing.T) {

	cases := []struct {
		Input          *managementClient.ElasticsearchConfig
		ExpectedOutput []interface{}
	}{
		{
			testLoggingElasticsearchConf,
			testLoggingElasticsearchInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenLoggingElasticsearchConfig(tc.Input, []interface{}{})
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandLoggingElasticsearchConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.ElasticsearchConfig
	}{
		{
			testLoggingElasticsearchInterface,
			testLoggingElasticsearchConf,
		},
	}

	for _, tc := range cases {
		output, err := expandLoggingElasticsearchConfig(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
