package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testMembersConf      []managementClient.Member
	testMembersInterface []interface{}
)

func init() {
	testMembersConf = []managementClient.Member{
		{
			AccessType:       "access_type",
			GroupPrincipalID: "group_principal_id",
			UserPrincipalID:  "user_principal_id",
		},
	}
	testMembersInterface = []interface{}{
		map[string]interface{}{
			"access_type":        "access_type",
			"group_principal_id": "group_principal_id",
			"user_principal_id":  "user_principal_id",
		},
	}
}

func TestFlattenMembers(t *testing.T) {

	cases := []struct {
		Input          []managementClient.Member
		ExpectedOutput []interface{}
	}{
		{
			testMembersConf,
			testMembersInterface,
		},
	}

	for _, tc := range cases {
		output := flattenMembers(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandMembers(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []managementClient.Member
	}{
		{
			testMembersInterface,
			testMembersConf,
		},
	}

	for _, tc := range cases {
		output := expandMembers(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
