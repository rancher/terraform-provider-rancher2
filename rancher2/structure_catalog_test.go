package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testCatalogGlobalConf       *managementClient.Catalog
	testCatalogGlobalInterface  map[string]interface{}
	testCatalogClusterConf      *managementClient.ClusterCatalog
	testCatalogClusterInterface map[string]interface{}
	testCatalogProjectConf      *managementClient.ProjectCatalog
	testCatalogProjectInterface map[string]interface{}
)

func init() {
	testCatalogGlobalConf = &managementClient.Catalog{
		Name:        "catalog-test",
		URL:         "url",
		Description: "description",
		Kind:        "kind",
		Branch:      "branch",
		HelmVersion: "helm_v3",
	}
	testCatalogGlobalInterface = map[string]interface{}{
		"name":        "catalog-test",
		"url":         "url",
		"description": "description",
		"kind":        "kind",
		"branch":      "branch",
		"scope":       "global",
		"version":     "helm_v3",
	}
	testCatalogClusterConf = &managementClient.ClusterCatalog{
		Name:        "catalog-test",
		URL:         "url",
		Description: "description",
		Kind:        "kind",
		Branch:      "branch",
		ClusterID:   "cluster_id",
		HelmVersion: "helm_v3",
	}
	testCatalogClusterInterface = map[string]interface{}{
		"name":        "catalog-test",
		"url":         "url",
		"description": "description",
		"kind":        "kind",
		"branch":      "branch",
		"scope":       "cluster",
		"cluster_id":  "cluster_id",
		"version":     "helm_v3",
	}
	testCatalogProjectConf = &managementClient.ProjectCatalog{
		Name:        "catalog-test",
		URL:         "url",
		Description: "description",
		Kind:        "kind",
		Branch:      "branch",
		ProjectID:   "project_id",
		HelmVersion: "helm_v3",
	}
	testCatalogProjectInterface = map[string]interface{}{
		"name":        "catalog-test",
		"url":         "url",
		"description": "description",
		"kind":        "kind",
		"branch":      "branch",
		"scope":       "project",
		"project_id":  "project_id",
		"version":     "helm_v3",
	}
}

func TestFlattenCatalog(t *testing.T) {

	cases := []struct {
		Input          interface{}
		ExpectedOutput map[string]interface{}
	}{
		{
			testCatalogGlobalConf,
			testCatalogGlobalInterface,
		},
		{
			testCatalogClusterConf,
			testCatalogClusterInterface,
		},
		{
			testCatalogProjectConf,
			testCatalogProjectInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, catalogFields(), tc.ExpectedOutput)
		err := flattenCatalog(output, tc.Input)
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

func TestExpandCatalog(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput interface{}
	}{
		{
			testCatalogGlobalInterface,
			testCatalogGlobalConf,
		},
		{
			testCatalogClusterInterface,
			testCatalogClusterConf,
		},
		{
			testCatalogProjectInterface,
			testCatalogProjectConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, catalogFields(), tc.Input)
		output := expandCatalog(inputResourceData)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
