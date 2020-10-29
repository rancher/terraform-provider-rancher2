package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	projectClient "github.com/rancher/rancher/pkg/client/generated/project/v3"
)

var (
	testAppConfGlobal            *projectClient.App
	testAppInterfaceGlobal       map[string]interface{}
	testAppConfClusterLocal      *projectClient.App
	testAppInterfaceClusterLocal map[string]interface{}
	testAppConfCluster           *projectClient.App
	testAppInterfaceCluster      map[string]interface{}
	testAppConfProject           *projectClient.App
	testAppInterfaceProject      map[string]interface{}
)

func init() {
	testAppConfGlobal = &projectClient.App{
		ExternalID:      "catalog://?catalog=test&template=test&version=1.23.0",
		Name:            "name",
		ProjectID:       "project:test",
		TargetNamespace: "target_namespace",
		Answers: map[string]string{
			"answers1": "one",
			"answers2": "two",
		},
		Description:   "description",
		AppRevisionID: "revision_id",
		ValuesYaml:    "values_yaml",
		Annotations: map[string]string{
			"node_one": "one",
			"node_two": "two",
		},
		Labels: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
		Timeout: 1200,
	}
	testAppInterfaceGlobal = map[string]interface{}{
		"catalog_name": "test",
		//"external_id":      "catalog://?catalog=test&template=test&version=1.23.0",
		"name":             "name",
		"project_id":       "project:test",
		"target_namespace": "target_namespace",
		"template_name":    "test",
		"answers": map[string]interface{}{
			"answers1": "one",
			"answers2": "two",
		},
		"description":      "description",
		"revision_id":      "revision_id",
		"template_version": "1.23.0",
		"values_yaml":      Base64Encode("values_yaml"),
		"annotations": map[string]interface{}{
			"node_one": "one",
			"node_two": "two",
		},
		"labels": map[string]interface{}{
			"option1": "value1",
			"option2": "value2",
		},
	}
	testAppConfClusterLocal = &projectClient.App{
		ExternalID:      "catalog://?catalog=local/test&type=clusterCatalog&template=test&version=1.23.0",
		Name:            "name",
		ProjectID:       "project:test",
		TargetNamespace: "target_namespace",
		Answers: map[string]string{
			"answers1": "one",
			"answers2": "two",
		},
		Description:   "description",
		AppRevisionID: "revision_id",
		ValuesYaml:    "values_yaml",
		Annotations: map[string]string{
			"node_one": "one",
			"node_two": "two",
		},
		Labels: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
		Timeout: 1200,
	}
	testAppInterfaceClusterLocal = map[string]interface{}{
		"catalog_name":     "local:test",
		"name":             "name",
		"project_id":       "project:test",
		"target_namespace": "target_namespace",
		"template_name":    "test",
		"answers": map[string]interface{}{
			"answers1": "one",
			"answers2": "two",
		},
		"description":      "description",
		"revision_id":      "revision_id",
		"template_version": "1.23.0",
		"values_yaml":      Base64Encode("values_yaml"),
		"annotations": map[string]interface{}{
			"node_one": "one",
			"node_two": "two",
		},
		"labels": map[string]interface{}{
			"option1": "value1",
			"option2": "value2",
		},
	}
	testAppConfCluster = &projectClient.App{
		ExternalID:      "catalog://?catalog=c-XXXXX/test&type=clusterCatalog&template=test&version=1.23.0",
		Name:            "name",
		ProjectID:       "project:test",
		TargetNamespace: "target_namespace",
		Answers: map[string]string{
			"answers1": "one",
			"answers2": "two",
		},
		Description:   "description",
		AppRevisionID: "revision_id",
		ValuesYaml:    "values_yaml",
		Annotations: map[string]string{
			"node_one": "one",
			"node_two": "two",
		},
		Labels: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
		Timeout: 1200,
	}
	testAppInterfaceCluster = map[string]interface{}{
		"catalog_name":     "c-XXXXX:test",
		"name":             "name",
		"project_id":       "project:test",
		"target_namespace": "target_namespace",
		"template_name":    "test",
		"answers": map[string]interface{}{
			"answers1": "one",
			"answers2": "two",
		},
		"description":      "description",
		"revision_id":      "revision_id",
		"template_version": "1.23.0",
		"values_yaml":      Base64Encode("values_yaml"),
		"annotations": map[string]interface{}{
			"node_one": "one",
			"node_two": "two",
		},
		"labels": map[string]interface{}{
			"option1": "value1",
			"option2": "value2",
		},
	}
	testAppConfProject = &projectClient.App{
		ExternalID:      "catalog://?catalog=p-XXXXX/test&type=projectCatalog&template=test&version=1.23.0",
		Name:            "name",
		ProjectID:       "project:test",
		TargetNamespace: "target_namespace",
		Answers: map[string]string{
			"answers1": "one",
			"answers2": "two",
		},
		Description:   "description",
		AppRevisionID: "revision_id",
		ValuesYaml:    "values_yaml",
		Annotations: map[string]string{
			"node_one": "one",
			"node_two": "two",
		},
		Labels: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
		Timeout: 1200,
	}
	testAppInterfaceProject = map[string]interface{}{
		"catalog_name":     "p-XXXXX:test",
		"name":             "name",
		"project_id":       "project:test",
		"target_namespace": "target_namespace",
		"template_name":    "test",
		"answers": map[string]interface{}{
			"answers1": "one",
			"answers2": "two",
		},
		"description":      "description",
		"revision_id":      "revision_id",
		"template_version": "1.23.0",
		"values_yaml":      Base64Encode("values_yaml"),
		"annotations": map[string]interface{}{
			"node_one": "one",
			"node_two": "two",
		},
		"labels": map[string]interface{}{
			"option1": "value1",
			"option2": "value2",
		},
	}
}

func TestFlattenApp(t *testing.T) {

	cases := []struct {
		Input          *projectClient.App
		ExpectedOutput map[string]interface{}
	}{
		{
			testAppConfGlobal,
			testAppInterfaceGlobal,
		},
		{
			testAppConfClusterLocal,
			testAppInterfaceClusterLocal,
		},
		{
			testAppConfCluster,
			testAppInterfaceCluster,
		},
		{
			testAppConfProject,
			testAppInterfaceProject,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, appFields(), tc.ExpectedOutput)
		err := flattenApp(output, tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		if !reflect.DeepEqual(expectedOutput, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				expectedOutput, tc.ExpectedOutput)
		}
	}
}

func TestExpandApp(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput interface{}
	}{
		{
			testAppInterfaceGlobal,
			testAppConfGlobal,
		},
		{
			testAppInterfaceClusterLocal,
			testAppConfClusterLocal,
		},
		{
			testAppInterfaceCluster,
			testAppConfCluster,
		},
		{
			testAppInterfaceProject,
			testAppConfProject,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, appFields(), tc.Input)
		output, err := expandApp(inputResourceData)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
