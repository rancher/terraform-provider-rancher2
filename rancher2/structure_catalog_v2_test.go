package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/apis/catalog.cattle.io/v1"
)

var (
	testCatalogV2Conf      *ClusterRepo
	testCatalogV2Interface map[string]interface{}
)

func init() {
	testCatalogV2Conf = &ClusterRepo{}

	testCatalogV2Conf.TypeMeta.Kind = catalogV2Kind
	testCatalogV2Conf.TypeMeta.APIVersion = catalogV2APIGroup + "/" + catalogV2APIVersion

	testCatalogV2Conf.ObjectMeta.Name = "name"
	testCatalogV2Conf.ObjectMeta.Annotations = map[string]string{
		"value1": "one",
		"value2": "two",
	}
	testCatalogV2Conf.ObjectMeta.Labels = map[string]string{
		"label1": "one",
		"label2": "two",
	}
	testCatalogV2Conf.Spec.CABundle = []byte("ca_bundle")
	testCatalogV2Conf.Spec.Enabled = newTrue()
	testCatalogV2Conf.Spec.GitBranch = "git_branch"
	testCatalogV2Conf.Spec.GitRepo = "git_repo"
	testCatalogV2Conf.Spec.InsecureSkipTLSverify = false
	testCatalogV2Conf.Spec.ClientSecret = &managementClient.SecretReference{
		Name:      "secret_name",
		Namespace: "secret_namespace",
	}
	testCatalogV2Conf.Spec.ServiceAccount = "service_account"
	testCatalogV2Conf.Spec.ServiceAccountNamespace = "service_account_namespace"
	testCatalogV2Conf.Spec.URL = "url"

	testCatalogV2Interface = map[string]interface{}{
		"name":                      "name",
		"ca_bundle":                 "ca_bundle",
		"enabled":                   true,
		"git_branch":                "git_branch",
		"git_repo":                  "git_repo",
		"insecure":                  false,
		"secret_name":               "secret_name",
		"secret_namespace":          "secret_namespace",
		"service_account":           "service_account",
		"service_account_namespace": "service_account_namespace",
		"url":                       "url",
		"annotations": map[string]interface{}{
			"value1": "one",
			"value2": "two",
		},
		"labels": map[string]interface{}{
			"label1": "one",
			"label2": "two",
		},
	}
}

func TestFlattenCatalogV2(t *testing.T) {

	cases := []struct {
		Input          *ClusterRepo
		ExpectedOutput map[string]interface{}
	}{
		{
			testCatalogV2Conf,
			testCatalogV2Interface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, catalogV2Fields(), tc.ExpectedOutput)
		err := flattenCatalogV2(output, tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		if !reflect.DeepEqual(expectedOutput, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, expectedOutput)
		}
	}
}

func TestExpandCatalogV2(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *ClusterRepo
	}{
		{
			testCatalogV2Interface,
			testCatalogV2Conf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, catalogV2Fields(), tc.Input)
		output := expandCatalogV2(inputResourceData)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
