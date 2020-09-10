package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testClusterLoggingConfCustomTarget       *managementClient.ClusterLogging
	testClusterLoggingInterfaceCustomTarget  map[string]interface{}
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
		{
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
		Compress:      newTrue(),
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
		EnableTLS:   true,
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
			"enable_tls":  true,
			"token":       "XXXXXXXXXXXX",
		},
	}
	testClusterLoggingConfCustomTarget = &managementClient.ClusterLogging{
		ClusterID:           "cluster-test",
		Name:                "test",
		EnableJSONParsing:   true,
		CustomTargetConfig:  testLoggingCustomTargetConf,
		NamespaceId:         "namespace-test",
		OutputFlushInterval: 10,
		OutputTags: map[string]string{
			"outputTag1": "value1",
			"outputTag2": "value2",
		},
	}
	testClusterLoggingInterfaceCustomTarget = map[string]interface{}{
		"cluster_id":            "cluster-test",
		"name":                  "test",
		"kind":                  loggingCustomTargetKind,
		"custom_target_config":  testLoggingCustomTargetInterface,
		"enable_json_parsing":   true,
		"namespace_id":          "namespace-test",
		"output_flush_interval": 10,
		"output_tags": map[string]interface{}{
			"outputTag1": "value1",
			"outputTag2": "value2",
		},
	}
	testClusterLoggingConfElasticSearch = &managementClient.ClusterLogging{
		ClusterID:           "cluster-test",
		Name:                "test",
		EnableJSONParsing:   true,
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
		"enable_json_parsing":   true,
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
		EnableJSONParsing:     true,
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
		"enable_json_parsing":   true,
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
		EnableJSONParsing:   true,
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
		"enable_json_parsing":   true,
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
		EnableJSONParsing:   true,
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
		"enable_json_parsing":   true,
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
		EnableJSONParsing:   true,
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
		"enable_json_parsing":   true,
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
			testClusterLoggingConfCustomTarget,
			testClusterLoggingInterfaceCustomTarget,
		},
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
			testClusterLoggingInterfaceCustomTarget,
			testClusterLoggingConfCustomTarget,
		},
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
