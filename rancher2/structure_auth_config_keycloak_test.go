package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testAuthConfigKeyCloakConf      *managementClient.KeyCloakConfig
	testAuthConfigKeyCloakInterface map[string]interface{}
)

func init() {
	testAuthConfigKeyCloakConf = &managementClient.KeyCloakConfig{
		Name:                AuthConfigKeyCloakName,
		Type:                managementClient.KeyCloakConfigType,
		AccessMode:          "access",
		AllowedPrincipalIDs: []string{"allowed1", "allowed2"},
		Enabled:             true,
		DisplayNameField:    "display_name_field",
		GroupsField:         "groups_field",
		IDPMetadataContent:  "idp",
		RancherAPIHost:      "rancher_api_host",
		SpCert:              "sp_cert",
		UIDField:            "uid_field",
		UserNameField:       "user_name_field",
		EntityID:            "entity_id",
	}
	testAuthConfigKeyCloakInterface = map[string]interface{}{
		"name":                  AuthConfigKeyCloakName,
		"type":                  managementClient.KeyCloakConfigType,
		"access_mode":           "access",
		"allowed_principal_ids": []interface{}{"allowed1", "allowed2"},
		"enabled":               true,
		"display_name_field":    "display_name_field",
		"groups_field":          "groups_field",
		"idp_metadata_content":  "idp",
		"rancher_api_host":      "rancher_api_host",
		"sp_cert":               "sp_cert",
		"uid_field":             "uid_field",
		"user_name_field":       "user_name_field",
		"entity_id":             "entity_id",
	}
}

func TestFlattenAuthConfigKeyCloak(t *testing.T) {

	cases := []struct {
		Input          *managementClient.KeyCloakConfig
		ExpectedOutput map[string]interface{}
	}{
		{
			testAuthConfigKeyCloakConf,
			testAuthConfigKeyCloakInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, authConfigKeyCloakFields(), map[string]interface{}{})
		err := flattenAuthConfigKeyCloak(output, tc.Input)
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

func TestExpandAuthConfigKeyCloak(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.KeyCloakConfig
	}{
		{
			testAuthConfigKeyCloakInterface,
			testAuthConfigKeyCloakConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, authConfigKeyCloakFields(), tc.Input)
		output, err := expandAuthConfigKeyCloak(inputResourceData)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
