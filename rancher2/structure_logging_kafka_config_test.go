package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testLoggingKafkaConf      *managementClient.KafkaConfig
	testLoggingKafkaInterface []interface{}
)

func init() {
	testLoggingKafkaConf = &managementClient.KafkaConfig{
		BrokerEndpoints:   []string{"hostname.test"},
		Certificate:       "XXXXXXXX",
		ClientCert:        "YYYYYYYY",
		ClientKey:         "ZZZZZZZZ",
		Topic:             "test",
		ZookeeperEndpoint: "zookeeper",
	}
	testLoggingKafkaInterface = []interface{}{
		map[string]interface{}{
			"broker_endpoints":   []interface{}{"hostname.test"},
			"certificate":        "XXXXXXXX",
			"client_cert":        "YYYYYYYY",
			"client_key":         "ZZZZZZZZ",
			"topic":              "test",
			"zookeeper_endpoint": "zookeeper",
		},
	}
}

func TestFlattenLoggingKafkaConfig(t *testing.T) {

	cases := []struct {
		Input          *managementClient.KafkaConfig
		ExpectedOutput []interface{}
	}{
		{
			testLoggingKafkaConf,
			testLoggingKafkaInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenLoggingKafkaConfig(tc.Input, []interface{}{})
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandLoggingKafkaConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.KafkaConfig
	}{
		{
			testLoggingKafkaInterface,
			testLoggingKafkaConf,
		},
	}

	for _, tc := range cases {
		output, err := expandLoggingKafkaConfig(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
