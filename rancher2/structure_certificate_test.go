package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	projectClient "github.com/rancher/types/client/project/v3"
)

var (
	testProjectCertificateConf         *projectClient.Certificate
	testProjectCertificateInterface    map[string]interface{}
	testNamespacedCertificateConf      *projectClient.NamespacedCertificate
	testNamespacedCertificateInterface map[string]interface{}
)

func init() {
	testProjectCertificateConf = &projectClient.Certificate{
		Certs:       "certificates",
		Key:         "key",
		ProjectID:   "project:test",
		Name:        "name",
		Description: "description",
		Annotations: map[string]string{
			"node_one": "one",
			"node_two": "two",
		},
		Labels: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
	}
	testProjectCertificateInterface = map[string]interface{}{
		"certs":       Base64Encode("certificates"),
		"key":         Base64Encode("key"),
		"project_id":  "project:test",
		"name":        "name",
		"description": "description",
		"annotations": map[string]interface{}{
			"node_one": "one",
			"node_two": "two",
		},
		"labels": map[string]interface{}{
			"option1": "value1",
			"option2": "value2",
		},
	}
	testNamespacedCertificateConf = &projectClient.NamespacedCertificate{
		Certs:       "certificates",
		Key:         "key",
		ProjectID:   "project:test",
		Name:        "name",
		Description: "description",
		NamespaceId: "namespace_id",
		Annotations: map[string]string{
			"node_one": "one",
			"node_two": "two",
		},
		Labels: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
	}
	testNamespacedCertificateInterface = map[string]interface{}{
		"certs":        Base64Encode("certificates"),
		"key":          Base64Encode("key"),
		"project_id":   "project:test",
		"name":         "name",
		"description":  "description",
		"namespace_id": "namespace_id",
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

func TestFlattenCertificate(t *testing.T) {

	cases := []struct {
		Input          interface{}
		ExpectedOutput map[string]interface{}
	}{
		{
			testProjectCertificateConf,
			testProjectCertificateInterface,
		},
		{
			testNamespacedCertificateConf,
			testNamespacedCertificateInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, certificateFields(), tc.ExpectedOutput)
		err := flattenCertificate(output, tc.Input)
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

func TestExpandCertificate(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput interface{}
	}{
		{
			testProjectCertificateInterface,
			testProjectCertificateConf,
		},
		{
			testNamespacedCertificateInterface,
			testNamespacedCertificateConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, certificateFields(), tc.Input)
		output, err := expandCertificate(inputResourceData)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
