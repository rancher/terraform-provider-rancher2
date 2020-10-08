package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testTokenConf      *managementClient.Token
	testTokenInterface map[string]interface{}
)

func init() {
	testTokenConf = &managementClient.Token{
		ClusterID:   "cluster_id",
		Description: "description",
		TTLMillis:   10000,
	}
	testTokenInterface = map[string]interface{}{
		"cluster_id":  "cluster_id",
		"description": "description",
		"ttl":         10,
	}
}

func TestFlattenToken(t *testing.T) {

	cases := []struct {
		Input          *managementClient.Token
		ExpectedOutput map[string]interface{}
	}{
		{
			testTokenConf,
			testTokenInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, tokenFields(), map[string]interface{}{})
		err := flattenToken(output, tc.Input, false)
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

func TestExpandToken(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.Token
	}{
		{
			testTokenInterface,
			testTokenConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, tokenFields(), tc.Input)
		output, err := expandToken(inputResourceData, false)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
