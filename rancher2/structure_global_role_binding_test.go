package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testGlobalRoleBindingConf      *managementClient.GlobalRoleBinding
	testGlobalRoleBindingInterface map[string]interface{}
)

func init() {
	testGlobalRoleBindingConf = &managementClient.GlobalRoleBinding{
		GlobalRoleID:     "global_role_id",
		GroupPrincipalID: "group_principal_id",
		UserID:           "user-test",
		Name:             "test",
	}
	testGlobalRoleBindingInterface = map[string]interface{}{
		"global_role_id":     "global_role_id",
		"group_principal_id": "group_principal_id",
		"user_id":            "user-test",
		"name":               "test",
	}
}

func TestFlattenGlobalRoleBinding(t *testing.T) {

	cases := []struct {
		Input          *managementClient.GlobalRoleBinding
		ExpectedOutput map[string]interface{}
	}{
		{
			testGlobalRoleBindingConf,
			testGlobalRoleBindingInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, globalRoleBindingFields(), map[string]interface{}{})
		err := flattenGlobalRoleBinding(output, tc.Input)
		if err != nil {
			assert.FailNow(t, "[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		if !reflect.DeepEqual(expectedOutput, tc.ExpectedOutput) {
			assert.FailNow(t, "Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				expectedOutput, output)
		}
	}
}

func TestExpandGlobalRoleBinding(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.GlobalRoleBinding
	}{
		{
			testGlobalRoleBindingInterface,
			testGlobalRoleBindingConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, globalRoleBindingFields(), tc.Input)
		output := expandGlobalRoleBinding(inputResourceData)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
