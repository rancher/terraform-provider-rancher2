package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestValidateInheritedNamespacedRules(t *testing.T) {
	cases := []struct {
		name     string
		ruleSets []any
		wantErr  string
	}{
		{
			name: "duplicate namespaces",
			ruleSets: []any{
				map[string]any{
					"namespace": "default",
					"rules": []any{
						map[string]any{
							"resources": []any{"configmaps"},
							"verbs":     []any{"get"},
						},
					},
				},
				map[string]any{
					"namespace": "default",
					"rules": []any{
						map[string]any{
							"resources": []any{"secrets"},
							"verbs":     []any{"list"},
						},
					},
				},
			},
			wantErr: `inherited_namespaced_rules cannot contain duplicate entries for namespace "default"`,
		},
		{
			name: "unique namespaces",
			ruleSets: []any{
				map[string]any{
					"namespace": "default",
					"rules":     []any{},
				},
				map[string]any{
					"namespace": "kube-system",
					"rules":     []any{},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data := schema.TestResourceDataRaw(t, globalRoleFields(), map[string]any{
				"name":                       "test-role",
				"inherited_namespaced_rules": tc.ruleSets,
			})

			err := validateInheritedNamespacedRules(data.Get("inherited_namespaced_rules"))
			if tc.wantErr == "" {
				assert.NoError(t, err)
				return
			}

			assert.EqualError(t, err, tc.wantErr)
		})
	}
}

func TestValidateInheritedNamespacedRulesAllowsUnknownValues(t *testing.T) {
	assert.NoError(t, validateInheritedNamespacedRules(nil))

	ruleSets := schema.NewSet(func(any) int { return 0 }, []any{
		map[string]any{
			"namespace": nil,
		},
	})
	assert.NoError(t, validateInheritedNamespacedRules(ruleSets))
}
