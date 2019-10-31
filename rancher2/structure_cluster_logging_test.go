package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testClusterLoggingConfElasticSearch      *managementClient.ClusterLogging
	testClusterLoggingInterfaceElasticSearch map[string]interface{}
	testClusterLoggingConfFluentd            *managementClient.ClusterLogging
	testClusterLoggingInterfaceFluentd       map[string]interface{}
	testClusterLoggingConfKafka              *managementClient.ClusterLogging
	testClusterLoggingInterfaceKafka         map[string]interface{}
	testClusterLoggingConfSplunk             *managementClient.ClusterLogging
	testClusterLoggingInterfaceSplunk        map[string]interface{}
	testClusterLoggingConfSyslog             *managementClient.ClusterLogging
	testClusterLoggingInterfaceSyslog        map[string]interface{}
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
	testClusterLoggingConfElasticSearch = &managementClient.ClusterLogging{
		ClusterID:           "cluster-test",
		Name:                "test",
		ElasticsearchConfig: testLoggingElasticsearchConf,
		NamespaceId:         "namespace-test",
		OutputFlushInterval: 10,
		OutputTags: map[string]string{
			"outputTag1": "value1",
			"outputTag2": "value2",
		},
	}
	testClusterLoggingInterfaceElasticSearch = map[string]interface{}{
		"cluster_id":            "cluster-test",
		"name":                  "test",
		"kind":                  loggingElasticsearchKind,
		"elasticsearch_config":  testLoggingElasticsearchInterface,
		"namespace_id":          "namespace-test",
		"output_flush_interval": 10,
		"output_tags": map[string]interface{}{
			"outputTag1": "value1",
			"outputTag2": "value2",
		},
	}
	testClusterLoggingConfFluentd = &managementClient.ClusterLogging{
		ClusterID:             "cluster-test",
		Name:                  "test",
		FluentForwarderConfig: testLoggingFluentdConf,
		NamespaceId:           "namespace-test",
		OutputFlushInterval:   10,
		OutputTags: map[string]string{
			"outputTag1": "value1",
			"outputTag2": "value2",
		},
	}
	testClusterLoggingInterfaceFluentd = map[string]interface{}{
		"cluster_id":            "cluster-test",
		"name":                  "test",
		"kind":                  loggingFluentdKind,
		"fluentd_config":        testLoggingFluentdInterface,
		"namespace_id":          "namespace-test",
		"output_flush_interval": 10,
		"output_tags": map[string]interface{}{
			"outputTag1": "value1",
			"outputTag2": "value2",
		},
	}
	testClusterLoggingConfKafka = &managementClient.ClusterLogging{
		ClusterID:           "cluster-test",
		Name:                "test",
		KafkaConfig:         testLoggingKafkaConf,
		NamespaceId:         "namespace-test",
		OutputFlushInterval: 10,
		OutputTags: map[string]string{
			"outputTag1": "value1",
			"outputTag2": "value2",
		},
	}
	testClusterLoggingInterfaceKafka = map[string]interface{}{
		"cluster_id":            "cluster-test",
		"name":                  "test",
		"kind":                  loggingKafkaKind,
		"kafka_config":          testLoggingKafkaInterface,
		"namespace_id":          "namespace-test",
		"output_flush_interval": 10,
		"output_tags": map[string]interface{}{
			"outputTag1": "value1",
			"outputTag2": "value2",
		},
	}
	testClusterLoggingConfSplunk = &managementClient.ClusterLogging{
		ClusterID:           "cluster-test",
		Name:                "test",
		SplunkConfig:        testLoggingSplunkConf,
		NamespaceId:         "namespace-test",
		OutputFlushInterval: 10,
		OutputTags: map[string]string{
			"outputTag1": "value1",
			"outputTag2": "value2",
		},
	}
	testClusterLoggingInterfaceSplunk = map[string]interface{}{
		"cluster_id":            "cluster-test",
		"name":                  "test",
		"kind":                  loggingSplunkKind,
		"splunk_config":         testLoggingSplunkInterface,
		"namespace_id":          "namespace-test",
		"output_flush_interval": 10,
		"output_tags": map[string]interface{}{
			"outputTag1": "value1",
			"outputTag2": "value2",
		},
	}
	testClusterLoggingConfSyslog = &managementClient.ClusterLogging{
		ClusterID:           "cluster-test",
		Name:                "test",
		SyslogConfig:        testLoggingSyslogConf,
		NamespaceId:         "namespace-test",
		OutputFlushInterval: 10,
		OutputTags: map[string]string{
			"outputTag1": "value1",
			"outputTag2": "value2",
		},
	}
	testClusterLoggingInterfaceSyslog = map[string]interface{}{
		"cluster_id":            "cluster-test",
		"name":                  "test",
		"kind":                  loggingSyslogKind,
		"syslog_config":         testLoggingSyslogInterface,
		"namespace_id":          "namespace-test",
		"output_flush_interval": 10,
		"output_tags": map[string]interface{}{
			"outputTag1": "value1",
			"outputTag2": "value2",
		},
	}
}

func TestFlattenClusterLogging(t *testing.T) {

	cases := []struct {
		Input          *managementClient.ClusterLogging
		ExpectedOutput map[string]interface{}
	}{
		{
			testClusterLoggingConfElasticSearch,
			testClusterLoggingInterfaceElasticSearch,
		},
		{
			testClusterLoggingConfFluentd,
			testClusterLoggingInterfaceFluentd,
		},
		{
			testClusterLoggingConfKafka,
			testClusterLoggingInterfaceKafka,
		},
		{
			testClusterLoggingConfSplunk,
			testClusterLoggingInterfaceSplunk,
		},
		{
			testClusterLoggingConfSyslog,
			testClusterLoggingInterfaceSyslog,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, clusterLoggingFields(), map[string]interface{}{})
		err := flattenClusterLogging(output, tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		if !reflect.DeepEqual(expectedOutput, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				expectedOutput, output)
		}
	}
}

func TestExpandClusterLogging(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.ClusterLogging
	}{
		{
			testClusterLoggingInterfaceElasticSearch,
			testClusterLoggingConfElasticSearch,
		},
		{
			testClusterLoggingInterfaceFluentd,
			testClusterLoggingConfFluentd,
		},
		{
			testClusterLoggingInterfaceKafka,
			testClusterLoggingConfKafka,
		},
		{
			testClusterLoggingInterfaceSplunk,
			testClusterLoggingConfSplunk,
		},
		{
			testClusterLoggingInterfaceSyslog,
			testClusterLoggingConfSyslog,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, clusterLoggingFields(), tc.Input)
		output, err := expandClusterLogging(inputResourceData)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
