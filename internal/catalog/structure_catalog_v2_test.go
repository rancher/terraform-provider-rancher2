package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/apis/catalog.cattle.io/v1"
	"github.com/stretchr/testify/assert"
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
	testCatalogV2Conf.Spec.ExponentialBackOffValues = &managementClient.ExponentialBackOffValues{
		MinWait:    2,
		MaxWait:    10,
		MaxRetries: 5,
	}
	testCatalogV2Conf.Spec.CABundle = []byte("test DER data")
	testCatalogV2Conf.Spec.Enabled = newTrue()
	testCatalogV2Conf.Spec.GitBranch = "git_branch"
	testCatalogV2Conf.Spec.GitRepo = "git_repo"
	testCatalogV2Conf.Spec.InsecurePlainHTTP = false
	testCatalogV2Conf.Spec.InsecureSkipTLSverify = false
	testCatalogV2Conf.Spec.ClientSecret = &managementClient.SecretReference{
		Name:      "secret_name",
		Namespace: "secret_namespace",
	}
	testCatalogV2Conf.Spec.ServiceAccount = "service_account"
	testCatalogV2Conf.Spec.ServiceAccountNamespace = "service_account_namespace"
	testCatalogV2Conf.Spec.URL = "url"

	testCatalogV2Interface = map[string]interface{}{
		"name":                            "name",
		"ca_bundle":                       "dGVzdCBERVIgZGF0YQ==",
		"enabled":                         true,
		"exponential_backoff_min_wait":    2,
		"exponential_backoff_max_wait":    10,
		"exponential_backoff_max_retries": 5,
		"git_branch":                      "git_branch",
		"git_repo":                        "git_repo",
		"insecure":                        false,
		"insecure_plain_http":             false,
		"secret_name":                     "secret_name",
		"secret_namespace":                "secret_namespace",
		"service_account":                 "service_account",
		"service_account_namespace":       "service_account_namespace",
		"url":                             "url",
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
			assert.FailNow(t, "[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		assert.Equal(t, tc.ExpectedOutput, expectedOutput, "Unexpected output from flattener.")
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
		output, err := expandCatalogV2(inputResourceData)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
