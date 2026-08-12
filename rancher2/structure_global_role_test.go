package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testGlobalRolePolicyRulesConf                       []managementClient.PolicyRule
	testGlobalRolePolicyRulesInterface                  []any
	testGlobalRoleConf                                  *managementClient.GlobalRole
	testGlobalRoleInterface                             map[string]any
	testGlobalRoleWithInheritedClusterRolesConf         *managementClient.GlobalRole
	testGlobalRoleWithInheritedClusterRolesInterface    map[string]any
	testGlobalRoleWithInheritedNamespacedRulesConf      *managementClient.GlobalRole
	testGlobalRoleWithInheritedNamespacedRulesInterface map[string]any
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
	testGlobalRolePolicyRulesInterface = []any{
		map[string]any{
			"api_groups": []any{
				"api_group1",
				"api_group2",
			},
			"non_resource_urls": []any{
				"non_resource_urls1",
				"non_resource_urls2",
			},
			"resource_names": []any{
				"resource_names1",
				"resource_names2",
			},
			"resources": []any{
				"resources1",
				"resources2",
			},
			"verbs": []any{
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
	testGlobalRoleInterface = map[string]any{
		"new_user_default": true,
		"description":      "description",
		"name":             "name",
		"rules":            testGlobalRolePolicyRulesInterface,
		"annotations": map[string]any{
			"node_one": "one",
			"node_two": "two",
		},
		"labels": map[string]any{
			"option1": "value1",
			"option2": "value2",
		},
	}

	testGlobalRoleWithInheritedClusterRolesConf = &managementClient.GlobalRole{
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
		InheritedClusterRoles: []string{
			"cluster-owner",
		},
	}
	testGlobalRoleWithInheritedClusterRolesInterface = map[string]any{
		"new_user_default": true,
		"description":      "description",
		"name":             "name",
		"rules":            testGlobalRolePolicyRulesInterface,
		"annotations": map[string]any{
			"node_one": "one",
			"node_two": "two",
		},
		"labels": map[string]any{
			"option1": "value1",
			"option2": "value2",
		},
		"inherited_cluster_roles": []any{
			"cluster-owner",
		},
	}

	testGlobalRoleWithInheritedNamespacedRulesConf = &managementClient.GlobalRole{
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
		InheritedNamespacedRules: map[string][]managementClient.PolicyRule{
			"namespace-one": testGlobalRolePolicyRulesConf,
			"namespace-two": {},
		},
	}
	testGlobalRoleWithInheritedNamespacedRulesInterface = map[string]any{
		"new_user_default": true,
		"description":      "description",
		"name":             "name",
		"rules":            testGlobalRolePolicyRulesInterface,
		"annotations": map[string]any{
			"node_one": "one",
			"node_two": "two",
		},
		"labels": map[string]any{
			"option1": "value1",
			"option2": "value2",
		},
		"inherited_namespaced_rules": []any{
			map[string]any{
				"namespace": "namespace-one",
				"rules":     testGlobalRolePolicyRulesInterface,
			},
			map[string]any{
				"namespace": "namespace-two",
				"rules":     []any{},
			},
		},
	}
}

func TestFlattenGlobalRole(t *testing.T) {
	cases := []struct {
		Input          *managementClient.GlobalRole
		ExpectedOutput map[string]any
	}{
		{
			testGlobalRoleConf,
			testGlobalRoleInterface,
		},
		{
			testGlobalRoleWithInheritedClusterRolesConf,
			testGlobalRoleWithInheritedClusterRolesInterface,
		},
		{
			testGlobalRoleWithInheritedNamespacedRulesConf,
			testGlobalRoleWithInheritedNamespacedRulesInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, globalRoleFields(), tc.ExpectedOutput)
		err := flattenGlobalRole(output, tc.Input)
		if err != nil {
			assert.FailNow(t, "[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]any{}
		expectedValues := map[string]any{}
		for k := range tc.ExpectedOutput {
			if k == "inherited_namespaced_rules" {
				assert.ElementsMatch(t, tc.ExpectedOutput[k].([]any), output.Get(k).(*schema.Set).List())
				continue
			}
			expectedOutput[k] = output.Get(k)
			expectedValues[k] = tc.ExpectedOutput[k]
		}
		assert.Equal(t, expectedValues, expectedOutput, "Unexpected output from flattener.")
	}
}

func TestExpandGlobalRole(t *testing.T) {
	cases := []struct {
		Input          map[string]any
		ExpectedOutput *managementClient.GlobalRole
	}{
		{
			testGlobalRoleInterface,
			testGlobalRoleConf,
		},
		{
			testGlobalRoleWithInheritedClusterRolesInterface,
			testGlobalRoleWithInheritedClusterRolesConf,
		},
		{
			testGlobalRoleWithInheritedNamespacedRulesInterface,
			testGlobalRoleWithInheritedNamespacedRulesConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, globalRoleFields(), tc.Input)
		output := expandGlobalRole(inputResourceData)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
