package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testAuthConfigGoogleOauthConf      *managementClient.GoogleOauthConfig
	testAuthConfigGoogleOauthInterface map[string]interface{}
)

func init() {
	testAuthConfigGoogleOauthConf = &managementClient.GoogleOauthConfig{
		Name:                         AuthConfigGoogleOauthName,
		Type:                         managementClient.GoogleOauthConfigType,
		AccessMode:                   "access",
		AllowedPrincipalIDs:          []string{"allowed1", "allowed2"},
		Enabled:                      true,
		AdminEmail:                   "admin_email",
		Hostname:                     "hostname",
		OauthCredential:              "oauth_credential",
		ServiceAccountCredential:     "service_account_credential",
		NestedGroupMembershipEnabled: true,
	}
	testAuthConfigGoogleOauthInterface = map[string]interface{}{
		"name":                            AuthConfigGoogleOauthName,
		"type":                            managementClient.GoogleOauthConfigType,
		"access_mode":                     "access",
		"allowed_principal_ids":           []interface{}{"allowed1", "allowed2"},
		"enabled":                         true,
		"admin_email":                     "admin_email",
		"hostname":                        "hostname",
		"oauth_credential":                "oauth_credential",
		"service_account_credential":      "service_account_credential",
		"nested_group_membership_enabled": true,
	}
}

func TestFlattenAuthConfigGoogleOauth(t *testing.T) {

	cases := []struct {
		Input          *managementClient.GoogleOauthConfig
		ExpectedOutput map[string]interface{}
	}{
		{
			testAuthConfigGoogleOauthConf,
			testAuthConfigGoogleOauthInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, authConfigGoogleOauthFields(), map[string]interface{}{})
		err := flattenAuthConfigGoogleOauth(output, tc.Input)
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

func TestExpandAuthConfigGoogleOauth(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.GoogleOauthConfig
	}{
		{
			testAuthConfigGoogleOauthInterface,
			testAuthConfigGoogleOauthConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, authConfigGoogleOauthFields(), tc.Input)
		output, err := expandAuthConfigGoogleOauth(inputResourceData)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
