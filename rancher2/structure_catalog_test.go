package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testCatalogConf      *managementClient.Catalog
	testCatalogInterface map[string]interface{}
)

func init() {
	testCatalogConf = &managementClient.Catalog{
		Name:        "catalog-test",
		URL:         "url",
		Description: "description",
		Kind:        "kind",
		Branch:      "branch",
	}
	testCatalogInterface = map[string]interface{}{
		"name":        "catalog-test",
		"url":         "url",
		"description": "description",
		"kind":        "kind",
		"branch":      "branch",
	}
}

func TestFlattenCatalog(t *testing.T) {

	cases := []struct {
		Input          *managementClient.Catalog
		ExpectedOutput map[string]interface{}
	}{
		{
			testCatalogConf,
			testCatalogInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, catalogFields(), map[string]interface{}{})
		err := flattenCatalog(output, tc.Input)
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

func TestExpandCatalog(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.Catalog
	}{
		{
			testCatalogInterface,
			testCatalogConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, catalogFields(), tc.Input)
		output := expandCatalog(inputResourceData)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
