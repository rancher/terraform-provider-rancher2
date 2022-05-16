package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testUserConf                                                                   *managementClient.User
	testFlattenUserTokenConf, testExpandUserTokenConf                              *managementClient.Token
	testUserInterface, testFlattenUserTokenInterface, testExpandUserTokenInterface map[string]interface{}
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
	testFlattenUserTokenConf = &managementClient.Token{
		Enabled: newTrue(),
		Expired: false,
		Name:    "token-99999",
		Token:   "token-99999:very_secret_key",
	}
	testFlattenUserTokenInterface = map[string]interface{}{
		"token_enabled": true,
		"token_expired": false,
		"token_name":    "token-99999",
		"auth_token":    "token-99999:very_secret_key",
		"access_key":    "token-99999",
		"secret_key":    "very_secret_key",
	}
	testExpandUserTokenConf = &managementClient.Token{
		Description: "foo",
		TTLMillis:   60000,
		ClusterID:   "c-99999",
	}
	testExpandUserTokenInterface = map[string]interface{}{
		"token_config": []interface{}{
			map[string]interface{}{
				"ttl":         60,
				"description": "foo",
				"cluster_id":  "c-99999",
			},
		},
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

func TestFlattenUserToken(t *testing.T) {

	cases := []struct {
		Input          *managementClient.Token
		ExpectedOutput map[string]interface{}
	}{
		{
			testFlattenUserTokenConf,
			testFlattenUserTokenInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, userFields(), map[string]interface{}{})
		err := flattenUserToken(output, tc.Input, false)
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

func TestExpandUserToken(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.Token
	}{
		{
			testExpandUserTokenInterface,
			testExpandUserTokenConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, userFields(), tc.Input)
		output := expandUserToken(inputResourceData, false)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
