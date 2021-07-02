package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testAuthConfigActiveDirectoryConf      *managementClient.ActiveDirectoryConfig
	testAuthConfigActiveDirectoryInterface map[string]interface{}
)

func init() {
	testAuthConfigActiveDirectoryConf = &managementClient.ActiveDirectoryConfig{
		Name:                         AuthConfigActiveDirectoryName,
		Type:                         managementClient.ActiveDirectoryConfigType,
		AccessMode:                   "access",
		AllowedPrincipalIDs:          []string{"allowed1", "allowed2"},
		Enabled:                      true,
		Servers:                      []string{"server1", "server2"},
		ServiceAccountUsername:       "service_account_username",
		UserSearchBase:               "user_search_base",
		Certificate:                  "certificate",
		ConnectionTimeout:            10,
		DefaultLoginDomain:           "default_login_domain",
		GroupDNAttribute:             "group_dn_attribute",
		GroupMemberMappingAttribute:  "group_member_mapping_attribute",
		GroupMemberUserAttribute:     "group_member_user_attribute",
		GroupNameAttribute:           "group_name_attribute",
		GroupObjectClass:             "group_object_class",
		GroupSearchAttribute:         "group_search_attribute",
		GroupSearchBase:              "group_search_base",
		GroupSearchFilter:            "group_search_filter",
		NestedGroupMembershipEnabled: newTrue(),
		Port:                         389,
		StartTLS:                     true,
		TLS:                          true,
		UserDisabledBitMask:          0,
		UserEnabledAttribute:         "user_enabled_attribute",
		UserLoginAttribute:           "user_login_attribute",
		UserNameAttribute:            "user_name_attribute",
		UserObjectClass:              "user_object_class",
		UserSearchAttribute:          "user_search_attribute",
		UserSearchFilter:             "user_search_filter",
	}
	testAuthConfigActiveDirectoryInterface = map[string]interface{}{
		"name":                            AuthConfigActiveDirectoryName,
		"type":                            managementClient.ActiveDirectoryConfigType,
		"access_mode":                     "access",
		"allowed_principal_ids":           []interface{}{"allowed1", "allowed2"},
		"enabled":                         true,
		"servers":                         []interface{}{"server1", "server2"},
		"service_account_username":        "service_account_username",
		"user_search_base":                "user_search_base",
		"certificate":                     "certificate",
		"connection_timeout":              10,
		"default_login_domain":            "default_login_domain",
		"group_dn_attribute":              "group_dn_attribute",
		"group_member_mapping_attribute":  "group_member_mapping_attribute",
		"group_member_user_attribute":     "group_member_user_attribute",
		"group_name_attribute":            "group_name_attribute",
		"group_object_class":              "group_object_class",
		"group_search_attribute":          "group_search_attribute",
		"group_search_base":               "group_search_base",
		"group_search_filter":             "group_search_filter",
		"nested_group_membership_enabled": true,
		"port":                            389,
		"start_tls":                       true,
		"tls":                             true,
		"user_disabled_bit_mask":          0,
		"user_enabled_attribute":          "user_enabled_attribute",
		"user_login_attribute":            "user_login_attribute",
		"user_name_attribute":             "user_name_attribute",
		"user_object_class":               "user_object_class",
		"user_search_attribute":           "user_search_attribute",
		"user_search_filter":              "user_search_filter",
	}
}

func TestFlattenAuthConfigActiveDirectory(t *testing.T) {

	cases := []struct {
		Input          *managementClient.ActiveDirectoryConfig
		ExpectedOutput map[string]interface{}
	}{
		{
			testAuthConfigActiveDirectoryConf,
			testAuthConfigActiveDirectoryInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, authConfigActiveDirectoryFields(), map[string]interface{}{})
		err := flattenAuthConfigActiveDirectory(output, tc.Input)
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

func TestExpandAuthConfigActiveDirectory(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.ActiveDirectoryConfig
	}{
		{
			testAuthConfigActiveDirectoryInterface,
			testAuthConfigActiveDirectoryConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, authConfigActiveDirectoryFields(), tc.Input)
		output, err := expandAuthConfigActiveDirectory(inputResourceData)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
