package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testAuthConfigAzureADConf      *managementClient.AzureADConfig
	testAuthConfigAzureADInterface map[string]interface{}
)

func init() {
	testAuthConfigAzureADConf = &managementClient.AzureADConfig{
		Name:                AuthConfigAzureADName,
		Type:                managementClient.AzureADConfigType,
		AccessMode:          "access",
		AllowedPrincipalIDs: []string{"allowed1", "allowed2"},
		Enabled:             true,
		ApplicationID:       "application_id",
		AuthEndpoint:        "auth_endpoint",
		Endpoint:            "endpoint",
		GraphEndpoint:       "graph_endpoint",
		RancherURL:          "rancher_url",
		TenantID:            "tenant_id",
		TokenEndpoint:       "token_endpoint",
	}
	testAuthConfigAzureADInterface = map[string]interface{}{
		"name":                  AuthConfigAzureADName,
		"type":                  managementClient.AzureADConfigType,
		"access_mode":           "access",
		"allowed_principal_ids": []interface{}{"allowed1", "allowed2"},
		"enabled":               true,
		"application_id":        "application_id",
		"auth_endpoint":         "auth_endpoint",
		"endpoint":              "endpoint",
		"graph_endpoint":        "graph_endpoint",
		"rancher_url":           "rancher_url",
		"tenant_id":             "tenant_id",
		"token_endpoint":        "token_endpoint",
	}
}

func TestFlattenAuthConfigAzureAD(t *testing.T) {

	cases := []struct {
		Input          *managementClient.AzureADConfig
		ExpectedOutput map[string]interface{}
	}{
		{
			testAuthConfigAzureADConf,
			testAuthConfigAzureADInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, authConfigAzureADFields(), map[string]interface{}{})
		err := flattenAuthConfigAzureAD(output, tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		if !reflect.DeepEqual(expectedOutput, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, expectedOutput)
		}
	}
}

func TestExpandAuthConfigAzureAD(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.AzureADConfig
	}{
		{
			testAuthConfigAzureADInterface,
			testAuthConfigAzureADConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, authConfigAzureADFields(), tc.Input)
		output, err := expandAuthConfigAzureAD(inputResourceData)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
