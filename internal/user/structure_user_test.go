package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testUserConf      *managementClient.User
	testUserInterface map[string]interface{}
)

func init() {
	testUserConf = &managementClient.User{
		Name:               "name",
		Username:           "username",
		Enabled:            newTrue(),
		MustChangePassword: *newTrue(),
	}
	testUserInterface = map[string]interface{}{
		"name":                 "name",
		"username":             "username",
		"enabled":              true,
		"must_change_password": true,
	}
}

func TestFlattenUser(t *testing.T) {

	cases := []struct {
		Input          *managementClient.User
		ExpectedOutput map[string]interface{}
	}{
		{
			testUserConf,
			testUserInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, userFields(), map[string]interface{}{})
		err := flattenUser(output, tc.Input)
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

func TestExpandUser(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.User
	}{
		{
			testUserInterface,
			testUserConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, userFields(), tc.Input)
		output := expandUser(inputResourceData)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
