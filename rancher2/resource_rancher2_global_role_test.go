package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestValidateInheritedNamespacedRules(t *testing.T) {
	cases := []struct {
		name     string
		ruleSets []interface{}
		wantErr  string
	}{
		{
			name: "duplicate namespaces",
			ruleSets: []interface{}{
				map[string]interface{}{
					"namespace": "default",
					"rules": []interface{}{
						map[string]interface{}{
							"resources": []interface{}{"configmaps"},
							"verbs":     []interface{}{"get"},
						},
					},
				},
				map[string]interface{}{
					"namespace": "default",
					"rules": []interface{}{
						map[string]interface{}{
							"resources": []interface{}{"secrets"},
							"verbs":     []interface{}{"list"},
						},
					},
				},
			},
			wantErr: `inherited_namespaced_rules cannot contain duplicate entries for namespace "default"`,
		},
		{
			name: "unique namespaces",
			ruleSets: []interface{}{
				map[string]interface{}{
					"namespace": "default",
					"rules":     []interface{}{},
				},
				map[string]interface{}{
					"namespace": "kube-system",
					"rules":     []interface{}{},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			data := schema.TestResourceDataRaw(t, globalRoleFields(), map[string]interface{}{
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

	ruleSets := schema.NewSet(func(interface{}) int { return 0 }, []interface{}{
		map[string]interface{}{
			"namespace": nil,
		},
	})
	assert.NoError(t, validateInheritedNamespacedRules(ruleSets))
}
