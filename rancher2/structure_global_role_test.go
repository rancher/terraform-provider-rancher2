package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testGlobalRolePolicyRulesConf      []managementClient.PolicyRule
	testGlobalRolePolicyRulesInterface []interface{}
	testGlobalRoleConf                 *managementClient.GlobalRole
	testGlobalRoleInterface            map[string]interface{}
)

func init() {
	testGlobalRolePolicyRulesConf = []managementClient.PolicyRule{
		{
			APIGroups: []string{
				"api_group1",
				"api_group2",
			},
			NonResourceURLs: []string{
				"non_resource_urls1",
				"non_resource_urls2",
			},
			ResourceNames: []string{
				"resource_names1",
				"resource_names2",
			},
			Resources: []string{
				"resources1",
				"resources2",
			},
			Verbs: []string{
				"verbs1",
				"verbs2",
			},
		},
	}
	testGlobalRolePolicyRulesInterface = []interface{}{
		map[string]interface{}{
			"api_groups": []interface{}{
				"api_group1",
				"api_group2",
			},
			"non_resource_urls": []interface{}{
				"non_resource_urls1",
				"non_resource_urls2",
			},
			"resource_names": []interface{}{
				"resource_names1",
				"resource_names2",
			},
			"resources": []interface{}{
				"resources1",
				"resources2",
			},
			"verbs": []interface{}{
				"verbs1",
				"verbs2",
			},
		},
	}

	testGlobalRoleConf = &managementClient.GlobalRole{
		Description:    "description",
		Name:           "name",
		NewUserDefault: true,
		Rules:          testGlobalRolePolicyRulesConf,
		Annotations: map[string]string{
			"node_one": "one",
			"node_two": "two",
		},
		Labels: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
	}
	testGlobalRoleInterface = map[string]interface{}{
		"new_user_default": true,
		"description":      "description",
		"name":             "name",
		"rules":            testGlobalRolePolicyRulesInterface,
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

func TestFlattenGlobalRole(t *testing.T) {

	cases := []struct {
		Input          *managementClient.GlobalRole
		ExpectedOutput map[string]interface{}
	}{
		{
			testGlobalRoleConf,
			testGlobalRoleInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, globalRoleFields(), tc.ExpectedOutput)
		err := flattenGlobalRole(output, tc.Input)
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

func TestExpandGlobalRole(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.GlobalRole
	}{
		{
			testGlobalRoleInterface,
			testGlobalRoleConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, globalRoleFields(), tc.Input)
		output := expandGlobalRole(inputResourceData)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
