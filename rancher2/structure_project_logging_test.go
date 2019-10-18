package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testProjectLoggingConfElasticSearch      *managementClient.ProjectLogging
	testProjectLoggingInterfaceElasticSearch map[string]interface{}
	testProjectLoggingConfFluentd            *managementClient.ProjectLogging
	testProjectLoggingInterfaceFluentd       map[string]interface{}
	testProjectLoggingConfKafka              *managementClient.ProjectLogging
	testProjectLoggingInterfaceKafka         map[string]interface{}
	testProjectLoggingConfSplunk             *managementClient.ProjectLogging
	testProjectLoggingInterfaceSplunk        map[string]interface{}
	testProjectLoggingConfSyslog             *managementClient.ProjectLogging
	testProjectLoggingInterfaceSyslog        map[string]interface{}
)

func init() {
	testProjectLoggingConfElasticSearch = &managementClient.ProjectLogging{
		ProjectID:           "project-test",
		Name:                "test",
		ElasticsearchConfig: testLoggingElasticsearchConf,
		NamespaceId:         "namespace-test",
		OutputFlushInterval: 10,
		OutputTags: map[string]string{
			"outputTag1": "value1",
			"outputTag2": "value2",
		},
	}
	testProjectLoggingInterfaceElasticSearch = map[string]interface{}{
		"project_id":            "project-test",
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
	testProjectLoggingConfFluentd = &managementClient.ProjectLogging{
		ProjectID:             "project-test",
		Name:                  "test",
		FluentForwarderConfig: testLoggingFluentdConf,
		NamespaceId:           "namespace-test",
		OutputFlushInterval:   10,
		OutputTags: map[string]string{
			"outputTag1": "value1",
			"outputTag2": "value2",
		},
	}
	testProjectLoggingInterfaceFluentd = map[string]interface{}{
		"project_id":            "project-test",
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
	testProjectLoggingConfKafka = &managementClient.ProjectLogging{
		ProjectID:           "project-test",
		Name:                "test",
		KafkaConfig:         testLoggingKafkaConf,
		NamespaceId:         "namespace-test",
		OutputFlushInterval: 10,
		OutputTags: map[string]string{
			"outputTag1": "value1",
			"outputTag2": "value2",
		},
	}
	testProjectLoggingInterfaceKafka = map[string]interface{}{
		"project_id":            "project-test",
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
	testProjectLoggingConfSplunk = &managementClient.ProjectLogging{
		ProjectID:           "project-test",
		Name:                "test",
		SplunkConfig:        testLoggingSplunkConf,
		NamespaceId:         "namespace-test",
		OutputFlushInterval: 10,
		OutputTags: map[string]string{
			"outputTag1": "value1",
			"outputTag2": "value2",
		},
	}
	testProjectLoggingInterfaceSplunk = map[string]interface{}{
		"project_id":            "project-test",
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
	testProjectLoggingConfSyslog = &managementClient.ProjectLogging{
		ProjectID:           "project-test",
		Name:                "test",
		SyslogConfig:        testLoggingSyslogConf,
		NamespaceId:         "namespace-test",
		OutputFlushInterval: 10,
		OutputTags: map[string]string{
			"outputTag1": "value1",
			"outputTag2": "value2",
		},
	}
	testProjectLoggingInterfaceSyslog = map[string]interface{}{
		"project_id":            "project-test",
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

func TestFlattenProjectLogging(t *testing.T) {

	cases := []struct {
		Input          *managementClient.ProjectLogging
		ExpectedOutput map[string]interface{}
	}{
		{
			testProjectLoggingConfElasticSearch,
			testProjectLoggingInterfaceElasticSearch,
		},
		{
			testProjectLoggingConfFluentd,
			testProjectLoggingInterfaceFluentd,
		},
		{
			testProjectLoggingConfKafka,
			testProjectLoggingInterfaceKafka,
		},
		{
			testProjectLoggingConfSplunk,
			testProjectLoggingInterfaceSplunk,
		},
		{
			testProjectLoggingConfSyslog,
			testProjectLoggingInterfaceSyslog,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, projectLoggingFields(), map[string]interface{}{})
		err := flattenProjectLogging(output, tc.Input)
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

func TestExpandProjectLogging(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.ProjectLogging
	}{
		{
			testProjectLoggingInterfaceElasticSearch,
			testProjectLoggingConfElasticSearch,
		},
		{
			testProjectLoggingInterfaceFluentd,
			testProjectLoggingConfFluentd,
		},
		{
			testProjectLoggingInterfaceKafka,
			testProjectLoggingConfKafka,
		},
		{
			testProjectLoggingInterfaceSplunk,
			testProjectLoggingConfSplunk,
		},
		{
			testProjectLoggingInterfaceSyslog,
			testProjectLoggingConfSyslog,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, projectLoggingFields(), tc.Input)
		output, err := expandProjectLogging(inputResourceData)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
