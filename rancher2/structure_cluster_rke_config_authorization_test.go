package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testClusterRKEConfigAuthorizationConf      *managementClient.AuthzConfig
	testClusterRKEConfigAuthorizationInterface []interface{}
)

func init() {
	testClusterRKEConfigAuthorizationConf = &managementClient.AuthzConfig{
		Mode: "rbac",
		Options: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
	}
	testClusterRKEConfigAuthorizationInterface = []interface{}{
		map[string]interface{}{
			"mode": "rbac",
			"options": map[string]interface{}{
				"option1": "value1",
				"option2": "value2",
			},
		},
	}
}

func TestFlattenClusterRKEConfigAuthorization(t *testing.T) {

	cases := []struct {
		Input          *managementClient.AuthzConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigAuthorizationConf,
			testClusterRKEConfigAuthorizationInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigAuthorization(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigAuthorization(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.AuthzConfig
	}{
		{
			testClusterRKEConfigAuthorizationInterface,
			testClusterRKEConfigAuthorizationConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigAuthorization(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
