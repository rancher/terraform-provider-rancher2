package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testProjectRoleTemplateBindingConf      *managementClient.ProjectRoleTemplateBinding
	testProjectRoleTemplateBindingInterface map[string]interface{}
)

func init() {
	testProjectRoleTemplateBindingConf = &managementClient.ProjectRoleTemplateBinding{
		ProjectID:        "project-test",
		RoleTemplateID:   "role-test",
		Name:             "test",
		GroupID:          "group-test",
		GroupPrincipalID: "group-principal-test",
		UserID:           "user-test",
		UserPrincipalID:  "user-principal-test",
	}
	testProjectRoleTemplateBindingInterface = map[string]interface{}{
		"project_id":         "project-test",
		"role_template_id":   "role-test",
		"name":               "test",
		"group_id":           "group-test",
		"group_principal_id": "group-principal-test",
		"user_id":            "user-test",
		"user_principal_id":  "user-principal-test",
	}
}

func TestFlattenProjectRoleTemplateBinding(t *testing.T) {

	cases := []struct {
		Input          *managementClient.ProjectRoleTemplateBinding
		ExpectedOutput map[string]interface{}
	}{
		{
			testProjectRoleTemplateBindingConf,
			testProjectRoleTemplateBindingInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, projectRoleTemplateBindingFields(), map[string]interface{}{})
		err := flattenProjectRoleTemplateBinding(output, tc.Input)
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

func TestExpandProjectRoleTemplateBinding(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.ProjectRoleTemplateBinding
	}{
		{
			testProjectRoleTemplateBindingInterface,
			testProjectRoleTemplateBindingConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, projectRoleTemplateBindingFields(), tc.Input)
		output := expandProjectRoleTemplateBinding(inputResourceData)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
