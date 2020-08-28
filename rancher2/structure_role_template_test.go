package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testRoleTemplatePolicyRulesConf      []managementClient.PolicyRule
	testRoleTemplatePolicyRulesInterface []interface{}
	testRoleTemplateClusterConf          *managementClient.RoleTemplate
	testRoleTemplateClusterInterface     map[string]interface{}
	testRoleTemplateProjectConf          *managementClient.RoleTemplate
	testRoleTemplateProjectInterface     map[string]interface{}
)

func init() {
	testRoleTemplatePolicyRulesConf = []managementClient.PolicyRule{
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
	testRoleTemplatePolicyRulesInterface = []interface{}{
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
	testRoleTemplateClusterConf = &managementClient.RoleTemplate{
		Administrative:        true,
		Context:               "cluster",
		ClusterCreatorDefault: true,
		Description:           "description",
		External:              true,
		Hidden:                true,
		Locked:                true,
		Name:                  "name",
		RoleTemplateIDs: []string{
			"role_template1",
			"role_template2",
		},
		Rules: testRoleTemplatePolicyRulesConf,
		Annotations: map[string]string{
			"node_one": "one",
			"node_two": "two",
		},
		Labels: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
	}
	testRoleTemplateClusterInterface = map[string]interface{}{
		"administrative": true,
		"builtin":        false,
		"context":        "cluster",
		"default_role":   true,
		"description":    "description",
		"external":       true,
		"hidden":         true,
		"locked":         true,
		"name":           "name",
		"role_template_ids": []interface{}{
			"role_template1",
			"role_template2",
		},
		"rules": testRoleTemplatePolicyRulesInterface,
		"annotations": map[string]interface{}{
			"node_one": "one",
			"node_two": "two",
		},
		"labels": map[string]interface{}{
			"option1": "value1",
			"option2": "value2",
		},
	}
	testRoleTemplateProjectConf = &managementClient.RoleTemplate{
		Administrative:        true,
		Context:               "project",
		Description:           "description",
		External:              true,
		Hidden:                true,
		Locked:                true,
		Name:                  "name",
		ProjectCreatorDefault: true,
		RoleTemplateIDs: []string{
			"role_template1",
			"role_template2",
		},
		Rules: testRoleTemplatePolicyRulesConf,
		Annotations: map[string]string{
			"node_one": "one",
			"node_two": "two",
		},
		Labels: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
	}
	testRoleTemplateProjectInterface = map[string]interface{}{
		"administrative": true,
		"builtin":        false,
		"context":        "project",
		"default_role":   true,
		"description":    "description",
		"external":       true,
		"hidden":         true,
		"locked":         true,
		"name":           "name",
		"role_template_ids": []interface{}{
			"role_template1",
			"role_template2",
		},
		"rules": testRoleTemplatePolicyRulesInterface,
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

func TestFlattenRoleTemplate(t *testing.T) {

	cases := []struct {
		Input          *managementClient.RoleTemplate
		ExpectedOutput map[string]interface{}
	}{
		{
			testRoleTemplateClusterConf,
			testRoleTemplateClusterInterface,
		},
		{
			testRoleTemplateProjectConf,
			testRoleTemplateProjectInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, roleTemplateFields(), tc.ExpectedOutput)
		err := flattenRoleTemplate(output, tc.Input)
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

func TestExpandRoleTemplate(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.RoleTemplate
	}{
		{
			testRoleTemplateClusterInterface,
			testRoleTemplateClusterConf,
		},
		{
			testRoleTemplateProjectInterface,
			testRoleTemplateProjectConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, roleTemplateFields(), tc.Input)
		output := expandRoleTemplate(inputResourceData)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
