package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testUserConf      *managementClient.User
	testUserInterface map[string]interface{}
)

func init() {
	testUserConf = &managementClient.User{
		Name:     "name",
		Username: "username",
		Enabled:  newTrue(),
	}
	testUserInterface = map[string]interface{}{
		"name":     "name",
		"username": "username",
		"enabled":  true,
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
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
