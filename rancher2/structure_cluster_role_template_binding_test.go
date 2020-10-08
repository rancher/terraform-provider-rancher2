package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testClusterRoleTemplateBindingConf      *managementClient.ClusterRoleTemplateBinding
	testClusterRoleTemplateBindingInterface map[string]interface{}
)

func init() {
	testClusterRoleTemplateBindingConf = &managementClient.ClusterRoleTemplateBinding{
		ClusterID:        "cluster-test",
		RoleTemplateID:   "role-test",
		Name:             "test",
		GroupID:          "group-test",
		GroupPrincipalID: "group-principal-test",
		UserID:           "user-test",
		UserPrincipalID:  "user-principal-test",
	}
	testClusterRoleTemplateBindingInterface = map[string]interface{}{
		"cluster_id":         "cluster-test",
		"role_template_id":   "role-test",
		"name":               "test",
		"group_id":           "group-test",
		"group_principal_id": "group-principal-test",
		"user_id":            "user-test",
		"user_principal_id":  "user-principal-test",
	}
}

func TestFlattenClusterRoleTemplateBinding(t *testing.T) {

	cases := []struct {
		Input          *managementClient.ClusterRoleTemplateBinding
		ExpectedOutput map[string]interface{}
	}{
		{
			testClusterRoleTemplateBindingConf,
			testClusterRoleTemplateBindingInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, clusterRoleTemplateBindingFields(), map[string]interface{}{})
		err := flattenClusterRoleTemplateBinding(output, tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		if !reflect.DeepEqual(expectedOutput, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				expectedOutput, output)
		}
	}
}

func TestExpandClusterRoleTemplateBinding(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.ClusterRoleTemplateBinding
	}{
		{
			testClusterRoleTemplateBindingInterface,
			testClusterRoleTemplateBindingConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, clusterRoleTemplateBindingFields(), tc.Input)
		output := expandClusterRoleTemplateBinding(inputResourceData)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
