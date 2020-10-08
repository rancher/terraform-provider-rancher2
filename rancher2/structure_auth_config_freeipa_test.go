package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testAuthConfigFreeIpaConf      *managementClient.LdapConfig
	testAuthConfigFreeIpaInterface map[string]interface{}
)

func init() {
	testAuthConfigFreeIpaConf = &managementClient.LdapConfig{
		Name:                            AuthConfigFreeIpaName,
		Type:                            managementClient.FreeIpaConfigType,
		AccessMode:                      "access",
		AllowedPrincipalIDs:             []string{"allowed1", "allowed2"},
		Enabled:                         true,
		Servers:                         []string{"server1", "server2"},
		ServiceAccountDistinguishedName: "service_account_distinguished_name",
		UserSearchBase:                  "user_search_base",
		Certificate:                     "certificate",
		ConnectionTimeout:               10,
		GroupDNAttribute:                "group_dn_attribute",
		GroupMemberMappingAttribute:     "group_member_mapping_attribute",
		GroupMemberUserAttribute:        "group_member_user_attribute",
		GroupNameAttribute:              "group_name_attribute",
		GroupObjectClass:                "group_object_class",
		GroupSearchAttribute:            "group_search_attribute",
		GroupSearchBase:                 "group_search_base",
		NestedGroupMembershipEnabled:    true,
		Port:                            389,
		TLS:                             true,
		UserDisabledBitMask:             0,
		UserLoginAttribute:              "user_login_attribute",
		UserMemberAttribute:             "user_member_attribute",
		UserNameAttribute:               "user_name_attribute",
		UserObjectClass:                 "user_object_class",
		UserSearchAttribute:             "user_search_attribute",
	}
	testAuthConfigFreeIpaInterface = map[string]interface{}{
		"name":                               AuthConfigFreeIpaName,
		"type":                               managementClient.FreeIpaConfigType,
		"access_mode":                        "access",
		"allowed_principal_ids":              []interface{}{"allowed1", "allowed2"},
		"enabled":                            true,
		"servers":                            []interface{}{"server1", "server2"},
		"service_account_distinguished_name": "service_account_distinguished_name",
		"user_search_base":                   "user_search_base",
		"certificate":                        Base64Encode("certificate"),
		"connection_timeout":                 10,
		"group_dn_attribute":                 "group_dn_attribute",
		"group_member_mapping_attribute":     "group_member_mapping_attribute",
		"group_member_user_attribute":        "group_member_user_attribute",
		"group_name_attribute":               "group_name_attribute",
		"group_object_class":                 "group_object_class",
		"group_search_attribute":             "group_search_attribute",
		"group_search_base":                  "group_search_base",
		"nested_group_membership_enabled":    true,
		"port":                               389,
		"tls":                                true,
		"user_disabled_bit_mask":             0,
		"user_login_attribute":               "user_login_attribute",
		"user_member_attribute":              "user_member_attribute",
		"user_name_attribute":                "user_name_attribute",
		"user_object_class":                  "user_object_class",
		"user_search_attribute":              "user_search_attribute",
	}
}

func TestFlattenAuthConfigFreeIpa(t *testing.T) {

	cases := []struct {
		Input          *managementClient.LdapConfig
		ExpectedOutput map[string]interface{}
	}{
		{
			testAuthConfigFreeIpaConf,
			testAuthConfigFreeIpaInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, authConfigFreeIpaFields(), map[string]interface{}{})
		err := flattenAuthConfigFreeIpa(output, tc.Input)
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

func TestExpandAuthConfigFreeIpa(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.LdapConfig
	}{
		{
			testAuthConfigFreeIpaInterface,
			testAuthConfigFreeIpaConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, authConfigFreeIpaFields(), tc.Input)
		output, err := expandAuthConfigFreeIpa(inputResourceData)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
