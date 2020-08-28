package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testPolicyRulesConf      []managementClient.PolicyRule
	testPolicyRulesInterface []interface{}
)

func init() {
	testPolicyRulesConf = []managementClient.PolicyRule{
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
	testPolicyRulesInterface = []interface{}{
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
}

func TestFlattenPolicyRules(t *testing.T) {

	cases := []struct {
		Input          []managementClient.PolicyRule
		ExpectedOutput []interface{}
	}{
		{
			testPolicyRulesConf,
			testPolicyRulesInterface,
		},
	}

	for _, tc := range cases {
		output := flattenPolicyRules(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandPolicyRules(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []managementClient.PolicyRule
	}{
		{
			testPolicyRulesInterface,
			testPolicyRulesConf,
		},
	}

	for _, tc := range cases {
		output := expandPolicyRules(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
