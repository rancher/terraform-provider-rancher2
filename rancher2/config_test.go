package rancher2

import (
	"testing"

	"github.com/rancher/norman/types"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

func TestFindUserByPrincipalID(t *testing.T) {
	users := []managementClient.User{
		{
			Resource:     types.Resource{ID: "user-local"},
			PrincipalIDs: []string{"local://local:user-local"},
		},
		{
			Resource:     types.Resource{ID: "user-remote"},
			PrincipalIDs: []string{"saml://remote-user", "local://user-remote"},
		},
	}

	tests := []struct {
		name        string
		principalID string
		expected    *managementClient.User
	}{
		{
			name:        "local user principal",
			principalID: "local://local:user-local",
			expected:    &users[0],
		},
		{
			name:        "external user principal",
			principalID: "saml://remote-user",
			expected:    &users[1],
		},
		{
			name:        "secondary local principal",
			principalID: "local://user-remote",
			expected:    &users[1],
		},
		{
			name:        "not found",
			principalID: "unknown://nobody",
			expected:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findUserByPrincipalID(users, tt.principalID)
			assert.Equal(t, tt.expected, got)
		})
	}
}
