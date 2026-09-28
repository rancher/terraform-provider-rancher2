package rancher2

import (
	"bytes"
	"log"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/rancher/norman/clientbase"
	"github.com/stretchr/testify/assert"
)

// testUnknownVariableValue is the same sentinel Terraform's SDK uses
// internally (terraform-plugin-sdk/internal/configs/hcl2shim.UnknownVariableValue)
// to represent a value that is not yet known at plan time, for example an
// interpolation referencing a not-yet-created resource's computed attribute.
// It is unexported in the SDK, so it is duplicated here only for simulating
// that condition in tests.
const testUnknownVariableValue = "74D93920-ED26-11E3-AC10-0800200C9A66"

func testClusterV2LocalAuthEndpointResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"local_auth_endpoint": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Optional: true,
			Elem: &schema.Resource{
				Schema: clusterV2LocalAuthEndpointFields(),
			},
		},
	}
}

func TestClusterV2LocalAuthEndpointUseInternalCACerts(t *testing.T) {
	cases := []struct {
		name string
		raw  map[string]interface{}
		want bool
	}{
		{
			name: "no local_auth_endpoint block",
			raw:  map[string]interface{}{},
			want: false,
		},
		{
			name: "use_internal_ca_certs not set defaults to false",
			raw: map[string]interface{}{
				"local_auth_endpoint": []interface{}{
					map[string]interface{}{
						"enabled": true,
					},
				},
			},
			want: false,
		},
		{
			// Regression test for #2438: a legacy manually-supplied ca_certs that
			// happens to equal the cluster's own CA must never be reported as
			// use_internal_ca_certs = true. The helper only ever reads the stored
			// boolean, it never inspects ca_certs' content.
			name: "manually-supplied ca_certs does not imply use_internal_ca_certs",
			raw: map[string]interface{}{
				"local_auth_endpoint": []interface{}{
					map[string]interface{}{
						"enabled":  true,
						"ca_certs": "-----BEGIN CERTIFICATE-----\nFAKE\n-----END CERTIFICATE-----\n",
					},
				},
			},
			want: false,
		},
		{
			name: "use_internal_ca_certs true is read back",
			raw: map[string]interface{}{
				"local_auth_endpoint": []interface{}{
					map[string]interface{}{
						"enabled":               true,
						"fqdn":                  "ace.example.com",
						"use_internal_ca_certs": true,
					},
				},
			},
			want: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, testClusterV2LocalAuthEndpointResourceSchema(), tc.raw)
			assert.Equal(t, tc.want, clusterV2LocalAuthEndpointUseInternalCACerts(d))
		})
	}
}

func TestClusterV2LocalAuthEndpointShouldUseInternalCACerts(t *testing.T) {
	cases := []struct {
		name string
		raw  map[string]interface{}
		want bool
	}{
		{
			name: "enabled with internal CA",
			raw: map[string]interface{}{
				"local_auth_endpoint": []interface{}{map[string]interface{}{
					"enabled":               true,
					"use_internal_ca_certs": true,
				}},
			},
			want: true,
		},
		{
			name: "disabled with internal CA",
			raw: map[string]interface{}{
				"local_auth_endpoint": []interface{}{map[string]interface{}{
					"enabled":               false,
					"use_internal_ca_certs": true,
				}},
			},
			want: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, testClusterV2LocalAuthEndpointResourceSchema(), tc.raw)
			assert.Equal(t, tc.want, clusterV2LocalAuthEndpointShouldUseInternalCACerts(d))
		})
	}
}

func TestClusterV2LocalAuthEndpointSchemaLocations(t *testing.T) {
	t.Run("top-level schema exposes use_internal_ca_certs", func(t *testing.T) {
		_, ok := clusterV2LocalAuthEndpointFields()["use_internal_ca_certs"]
		assert.True(t, ok)
	})

	t.Run("deprecated current nested schema excludes use_internal_ca_certs", func(t *testing.T) {
		endpoint := clusterV2RKEConfigFields()["local_auth_endpoint"].Elem.(*schema.Resource)
		_, ok := endpoint.Schema["use_internal_ca_certs"]
		assert.False(t, ok)
	})

	t.Run("deprecated V0 nested schema excludes use_internal_ca_certs", func(t *testing.T) {
		endpoint := clusterV2RKEConfigFieldsV0()["local_auth_endpoint"].Elem.(*schema.Resource)
		_, ok := endpoint.Schema["use_internal_ca_certs"]
		assert.False(t, ok)
	})

	t.Run("cluster exposes internal CA synchronization state", func(t *testing.T) {
		syncRequired := clusterV2Fields()["local_auth_endpoint_ca_sync_required"]
		assert.NotNil(t, syncRequired)
		assert.True(t, syncRequired.Computed)
		assert.False(t, syncRequired.Optional)
	})
}

func TestClusterV2LocalAuthEndpointCASyncRequired(t *testing.T) {
	cases := []struct {
		name        string
		useInternal bool
		enabled     bool
		endpointCA  string
		internalCA  string
		want        bool
	}{
		{name: "matching CA", useInternal: true, enabled: true, endpointCA: "ca-a", internalCA: "ca-a", want: false},
		{name: "missing endpoint CA", useInternal: true, enabled: true, internalCA: "ca-a", want: true},
		{name: "rotated internal CA", useInternal: true, enabled: true, endpointCA: "ca-a", internalCA: "ca-b", want: true},
		{name: "internal CA unavailable", useInternal: true, enabled: true, endpointCA: "ca-a", want: false},
		{name: "endpoint disabled", useInternal: true, endpointCA: "ca-a", internalCA: "ca-b", want: false},
		{name: "manual CA mode", enabled: true, endpointCA: "ca-a", internalCA: "ca-b", want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cluster := &ClusterV2{}
			cluster.Spec.LocalClusterAuthEndpoint.Enabled = tc.enabled
			cluster.Spec.LocalClusterAuthEndpoint.CACerts = tc.endpointCA
			assert.Equal(t, tc.want, clusterV2LocalAuthEndpointCASyncRequired(tc.useInternal, cluster, tc.internalCA))
		})
	}
}

func TestGetClusterV2LocalAuthEndpointCACertSkipsEmptyClusterID(t *testing.T) {
	var logs bytes.Buffer
	oldOutput := log.Writer()
	log.SetOutput(&logs)
	t.Cleanup(func() { log.SetOutput(oldOutput) })

	called := false
	caCert, err := getClusterV2LocalAuthEndpointCACert("", func(string) (string, error) {
		called = true
		return "internal-ca", nil
	})

	assert.NoError(t, err)
	assert.Empty(t, caCert)
	assert.False(t, called)
	assert.Contains(t, logs.String(), "Skipping local_auth_endpoint CA certificate fetch because cluster_v1_id is not available")
}

func TestGetClusterV2LocalAuthEndpointCACertLogsClusterID(t *testing.T) {
	var logs bytes.Buffer
	oldOutput := log.Writer()
	log.SetOutput(&logs)
	t.Cleanup(func() { log.SetOutput(oldOutput) })

	_, err := getClusterV2LocalAuthEndpointCACert("c-m-test", func(string) (string, error) {
		return "internal-ca", nil
	})

	assert.NoError(t, err)
	assert.Contains(t, logs.String(), "Fetching cluster c-m-test CA certificate for local_auth_endpoint")
}

func TestSetClusterV2LocalAuthEndpointUseInternalCACerts(t *testing.T) {
	t.Run("reinjects the preserved value into an existing block", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, testClusterV2LocalAuthEndpointResourceSchema(), map[string]interface{}{
			"local_auth_endpoint": []interface{}{
				map[string]interface{}{
					"enabled": true,
					"fqdn":    "ace.example.com",
				},
			},
		})

		err := setClusterV2LocalAuthEndpointUseInternalCACerts(d, true)
		assert.NoError(t, err)
		assert.True(t, clusterV2LocalAuthEndpointUseInternalCACerts(d))
	})

	t.Run("restores true when there is no local_auth_endpoint block", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, testClusterV2LocalAuthEndpointResourceSchema(), map[string]interface{}{})

		err := setClusterV2LocalAuthEndpointUseInternalCACerts(d, true)
		assert.NoError(t, err)
		assert.True(t, clusterV2LocalAuthEndpointUseInternalCACerts(d))
	})

	t.Run("blanks ca_certs when use_internal_ca_certs is true", func(t *testing.T) {
		// ca_certs is Optional but not Computed. If the real, auto-populated
		// CA value were left in state while config never sets ca_certs,
		// Terraform's diff engine would treat "in state, not in config" as a
		// removal on every subsequent plan, and the plan would never
		// converge. See TestClusterV2LocalAuthEndpointUseInternalCACertsDoesNotProduceAPhantomDiff.
		d := schema.TestResourceDataRaw(t, testClusterV2LocalAuthEndpointResourceSchema(), map[string]interface{}{
			"local_auth_endpoint": []interface{}{
				map[string]interface{}{
					"enabled":  true,
					"fqdn":     "ace.example.com",
					"ca_certs": "-----BEGIN CERTIFICATE-----\nREAL\n-----END CERTIFICATE-----\n",
				},
			},
		})

		err := setClusterV2LocalAuthEndpointUseInternalCACerts(d, true)
		assert.NoError(t, err)

		v, _ := d.Get("local_auth_endpoint").([]interface{})
		m, ok := clusterV2LocalAuthEndpointRawMap(v)
		assert.True(t, ok)
		assert.Equal(t, "", m["ca_certs"])
	})

	t.Run("preserves ca_certs when use_internal_ca_certs is false", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, testClusterV2LocalAuthEndpointResourceSchema(), map[string]interface{}{
			"local_auth_endpoint": []interface{}{
				map[string]interface{}{
					"enabled":  true,
					"ca_certs": "manual-ca",
				},
			},
		})

		err := setClusterV2LocalAuthEndpointUseInternalCACerts(d, false)
		assert.NoError(t, err)

		v, _ := d.Get("local_auth_endpoint").([]interface{})
		m, ok := clusterV2LocalAuthEndpointRawMap(v)
		assert.True(t, ok)
		assert.Equal(t, "manual-ca", m["ca_certs"])
	})
}

func TestClusterV2LocalAuthEndpointUseInternalCACertsDoesNotProduceAPhantomDiff(t *testing.T) {
	// Regression test: once Read has captured a real, internally-populated
	// CA and use_internal_ca_certs stays true in config, the *next* plan
	// must not show a diff wanting to clear ca_certs back to empty.
	d := schema.TestResourceDataRaw(t, testClusterV2LocalAuthEndpointResourceSchema(), map[string]interface{}{
		"local_auth_endpoint": []interface{}{
			map[string]interface{}{
				"enabled":               true,
				"fqdn":                  "ace.example.com",
				"use_internal_ca_certs": true,
			},
		},
	})

	// Simulate flattenClusterV2LocalAuthEndpoint writing the real API value
	// into ca_certs during Read, before the preserve step runs.
	raw, _ := d.Get("local_auth_endpoint").([]interface{})
	m, ok := clusterV2LocalAuthEndpointRawMap(raw)
	assert.True(t, ok)
	m["ca_certs"] = "-----BEGIN CERTIFICATE-----\nREAL\n-----END CERTIFICATE-----\n"
	assert.NoError(t, d.Set("local_auth_endpoint", []interface{}{m}))

	assert.NoError(t, setClusterV2LocalAuthEndpointUseInternalCACerts(d, true))

	oldInterface, _ := d.Get("local_auth_endpoint").([]interface{})
	// What the next plan's new value looks like when nothing in config changed:
	// ca_certs is absent from config, so its planned value is the zero value.
	newInterface := []interface{}{
		map[string]interface{}{
			"enabled":               true,
			"fqdn":                  "ace.example.com",
			"ca_certs":              "",
			"use_internal_ca_certs": true,
		},
	}

	assert.True(t, clusterV2LocalAuthEndpointDiffEqual(oldInterface, newInterface))
}

func TestClusterV2LocalAuthEndpointDiffEqual(t *testing.T) {
	enabledTrue := map[string]interface{}{
		"enabled":               true,
		"fqdn":                  "ace.example.com",
		"ca_certs":              "",
		"use_internal_ca_certs": true,
	}
	enabledFalse := map[string]interface{}{
		"enabled":               true,
		"fqdn":                  "ace.example.com",
		"ca_certs":              "",
		"use_internal_ca_certs": false,
	}

	cases := []struct {
		name string
		old  []interface{}
		new  []interface{}
		want bool
	}{
		{
			name: "identical blocks including the flag are equal",
			old:  []interface{}{enabledTrue},
			new:  []interface{}{enabledTrue},
			want: true,
		},
		{
			// Regression test for #2436: toggling use_internal_ca_certs with
			// everything else equal must never be treated as diff noise, or the
			// CustomizeDiff would clear the change and the plan would never
			// converge.
			name: "flag toggled true to false with everything else equal is not equal",
			old:  []interface{}{enabledTrue},
			new:  []interface{}{enabledFalse},
			want: false,
		},
		{
			name: "flag toggled false to true with everything else equal is not equal",
			old:  []interface{}{enabledFalse},
			new:  []interface{}{enabledTrue},
			want: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, clusterV2LocalAuthEndpointDiffEqual(tc.old, tc.new))
		})
	}
}

func TestClusterV2CustomizeDiffAllowsUnknownFQDN(t *testing.T) {
	// Regression test: on the very first plan of a fresh configuration,
	// fqdn commonly interpolates a not-yet-created resource's attribute
	// (e.g. an aws_lb's dns_name). The SDK materializes such a value as the
	// zero value ("") via GetChange, which endpoint diff validation
	// would otherwise misinterpret as "fqdn not configured" and reject with
	// a false positive, even though fqdn IS configured (just not yet known).
	r := resourceRancher2ClusterV2()
	c := terraform.NewResourceConfigRaw(map[string]any{
		"name":               "test",
		"kubernetes_version": "v1.30.0",
		"local_auth_endpoint": []any{
			map[string]any{
				"enabled":               true,
				"fqdn":                  testUnknownVariableValue,
				"use_internal_ca_certs": true,
			},
		},
	})

	_, err := r.Diff(nil, c, nil)
	assert.NoError(t, err)
}

func TestClusterV2CustomizeDiffRejectsKnownCACertsWithUnknownFQDN(t *testing.T) {
	r := resourceRancher2ClusterV2()
	c := terraform.NewResourceConfigRaw(map[string]any{
		"name":               "test",
		"kubernetes_version": "v1.30.0",
		"local_auth_endpoint": []any{
			map[string]any{
				"enabled":               true,
				"fqdn":                  testUnknownVariableValue,
				"ca_certs":              "manual-ca",
				"use_internal_ca_certs": true,
			},
		},
	})

	_, err := r.Diff(nil, c, nil)
	assert.ErrorContains(t, err, `only one of "ca_certs" or "use_internal_ca_certs" can be set`)
}

func TestClusterV2CustomizeDiffAllowsUnknownCACerts(t *testing.T) {
	r := resourceRancher2ClusterV2()
	c := terraform.NewResourceConfigRaw(map[string]any{
		"name":               "test",
		"kubernetes_version": "v1.30.0",
		"local_auth_endpoint": []any{
			map[string]any{
				"enabled":               true,
				"fqdn":                  "ace.example.com",
				"ca_certs":              testUnknownVariableValue,
				"use_internal_ca_certs": true,
			},
		},
	})

	_, err := r.Diff(nil, c, nil)
	assert.NoError(t, err)
}

func TestClusterV2CustomizeDiffStillRejectsKnownEmptyFQDN(t *testing.T) {
	// Sanity check for the fix above: a genuinely empty (known, not
	// unknown) fqdn alongside use_internal_ca_certs=true must still be
	// rejected.
	r := resourceRancher2ClusterV2()
	c := terraform.NewResourceConfigRaw(map[string]any{
		"name":               "test",
		"kubernetes_version": "v1.30.0",
		"local_auth_endpoint": []any{
			map[string]any{
				"enabled":               true,
				"use_internal_ca_certs": true,
			},
		},
	})

	_, err := r.Diff(nil, c, nil)
	assert.ErrorContains(t, err, `"fqdn" is required in "local_auth_endpoint" when "use_internal_ca_certs" is true`)
}

func TestClusterV2CustomizeDiffSchedulesInternalCASync(t *testing.T) {
	cases := []struct {
		name           string
		syncRequired   bool
		enabled        bool
		useInternal    bool
		wantMarkerDiff bool
	}{
		{name: "synchronization required", syncRequired: true, enabled: true, useInternal: true, wantMarkerDiff: true},
		{name: "CA matches", enabled: true, useInternal: true},
		{name: "endpoint disabled", syncRequired: true, useInternal: true},
		{name: "manual CA mode", syncRequired: true, enabled: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := &schema.Resource{
				Schema: map[string]*schema.Schema{
					"local_auth_endpoint": {
						Type:     schema.TypeList,
						MaxItems: 1,
						Optional: true,
						Elem:     &schema.Resource{Schema: clusterV2LocalAuthEndpointFields()},
					},
					"local_auth_endpoint_ca_sync_required": {
						Type:     schema.TypeBool,
						Computed: true,
					},
				},
				CustomizeDiff: customizeClusterV2LocalAuthEndpointCASync,
			}
			state := &terraform.InstanceState{
				ID: "test",
				Attributes: map[string]string{
					"local_auth_endpoint.#":                       "1",
					"local_auth_endpoint.0.enabled":               strconv.FormatBool(tc.enabled),
					"local_auth_endpoint.0.fqdn":                  "ace.example.com",
					"local_auth_endpoint.0.ca_certs":              "",
					"local_auth_endpoint.0.use_internal_ca_certs": strconv.FormatBool(tc.useInternal),
					"local_auth_endpoint_ca_sync_required":        strconv.FormatBool(tc.syncRequired),
				},
			}
			config := terraform.NewResourceConfigRaw(map[string]interface{}{
				"local_auth_endpoint": []interface{}{map[string]interface{}{
					"enabled":               tc.enabled,
					"fqdn":                  "ace.example.com",
					"use_internal_ca_certs": tc.useInternal,
				}},
			})

			diff, err := r.Diff(state, config, nil)
			assert.NoError(t, err)
			if tc.wantMarkerDiff {
				if assert.NotNil(t, diff) {
					marker := diff.Attributes["local_auth_endpoint_ca_sync_required"]
					if assert.NotNil(t, marker) {
						assert.True(t, marker.NewComputed)
					}
				}
			} else if diff != nil {
				assert.Nil(t, diff.Attributes["local_auth_endpoint_ca_sync_required"])
			}
		})
	}
}

func TestUpdateClusterV2LocalAuthEndpointCACerts(t *testing.T) {
	latest := &ClusterV2{}
	latest.Spec.KubernetesVersion = "v1.31.0"
	fetches := 0
	updates := 0

	got, err := updateClusterV2LocalAuthEndpointCACerts(func() (*ClusterV2, error) {
		fetches++
		return latest, nil
	}, func(cluster *ClusterV2) (*ClusterV2, error) {
		updates++
		if updates == 1 {
			latest = &ClusterV2{}
			latest.Spec.KubernetesVersion = "v1.32.0"
			return nil, &clientbase.APIError{StatusCode: http.StatusConflict}
		}
		return cluster, nil
	}, "internal-ca", 0, time.Second)

	assert.NoError(t, err)
	assert.Equal(t, 2, fetches)
	assert.Equal(t, 2, updates)
	assert.Equal(t, "v1.32.0", got.Spec.KubernetesVersion)
	assert.Equal(t, "internal-ca", got.Spec.LocalClusterAuthEndpoint.CACerts)
}
