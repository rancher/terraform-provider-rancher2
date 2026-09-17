package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/zclconf/go-cty/cty"
)

func TestResourceRancher2ClusterStateUpgraders(t *testing.T) {
	r := resourceRancher2Cluster()

	if r.SchemaVersion != 3 {
		t.Fatalf("expected schema version 3, got %d", r.SchemaVersion)
	}

	if len(r.StateUpgraders) != 3 {
		t.Fatalf("expected 3 state upgraders, got %d", len(r.StateUpgraders))
	}
	expectedVersions := []int{0, 1, 2}
	expectedTypes := []cty.Type{
		resourceRancher2ClusterResourceV3().CoreConfigSchema().ImpliedType(),
		resourceRancher2ClusterResourceV0().CoreConfigSchema().ImpliedType(),
		resourceRancher2ClusterResourceV2().CoreConfigSchema().ImpliedType(),
	}

	for i, expected := range expectedVersions {
		if r.StateUpgraders[i].Version != expected {
			t.Fatalf("expected upgrader %d to have version %d, got %d", i, expected, r.StateUpgraders[i].Version)
		}

		if !r.StateUpgraders[i].Type.Equals(expectedTypes[i]) {
			t.Fatalf("expected upgrader %d to use the schema type for version %d", i, expected)
		}
	}

	if err := r.InternalValidate(nil, true); err != nil {
		t.Fatalf("resource failed internal validation: %v", err)
	}
}

func TestResourceRancher2ClusterStateUpgradeV0(t *testing.T) {
	eventRateConfiguration := map[string]interface{}{
		"apiVersion": "eventratelimit.admission.k8s.io/v1alpha1",
		"kind":       "Configuration",
		"limits": []interface{}{
			map[string]interface{}{
				"type":  "Server",
				"qps":   10,
				"burst": 20,
			},
		},
	}

	secretsEncryptionConfiguration := map[string]interface{}{
		"apiVersion": "apiserver.config.k8s.io/v1",
		"kind":       "EncryptionConfiguration",
		"resources": []interface{}{
			map[string]interface{}{
				"resources": []interface{}{"secrets"},
				"providers": []interface{}{
					map[string]interface{}{
						"aescbc": map[string]interface{}{
							"keys": []interface{}{
								map[string]interface{}{
									"name":   "key1",
									"secret": "secret",
								},
							},
						},
					},
				},
			},
		},
	}

	admissionConfiguration := map[string]interface{}{
		"api_version": "apiserver.config.k8s.io/v1",
		"kind":        "AdmissionConfiguration",
		"plugins": []interface{}{
			map[string]interface{}{
				"name": "PodSecurity",
			},
		},
	}

	rawState := map[string]interface{}{
		"rke_config": []interface{}{
			map[string]interface{}{
				"services": []interface{}{
					map[string]interface{}{
						"kube_api": []interface{}{
							map[string]interface{}{
								"event_rate_limit": []interface{}{
									map[string]interface{}{
										"configuration": eventRateConfiguration,
									},
								},
								"secrets_encryption_config": []interface{}{
									map[string]interface{}{
										"custom_config": secretsEncryptionConfiguration,
									},
								},
								"admission_configuration": admissionConfiguration,
							},
						},
					},
				},
			},
		},
	}

	upgraded, err := resourceRancher2ClusterStateUpgradeV0(rawState, nil)
	if err != nil {
		t.Fatalf("failed to upgrade V0 state: %v", err)
	}

	rkeConfigs := upgraded["rke_config"].([]interface{})
	rkeConfig := rkeConfigs[0].(map[string]interface{})
	services := rkeConfig["services"].([]interface{})
	service := services[0].(map[string]interface{})
	kubeAPIs := service["kube_api"].([]interface{})
	kubeAPI := kubeAPIs[0].(map[string]interface{})

	eventRates := kubeAPI["event_rate_limit"].([]interface{})
	eventRate := eventRates[0].(map[string]interface{})
	eventConfiguration, ok := eventRate["configuration"].(string)
	if !ok {
		t.Fatalf("expected event_rate_limit.configuration to be string, got %T", eventRate["configuration"])
	}

	convertedEventConfiguration, err := ghodssyamlToMapInterface(eventConfiguration)
	if err != nil {
		t.Fatalf("failed to parse converted event_rate_limit.configuration: %v", err)
	}
	expectedEventConfiguration, err := interfaceToMap(eventRateConfiguration)
	if err != nil {
		t.Fatalf("failed to normalize expected event_rate_limit.configuration: %v", err)
	}
	if !reflect.DeepEqual(convertedEventConfiguration, expectedEventConfiguration) {
		t.Fatalf(
			"unexpected event_rate_limit.configuration:\nexpected %#v\ngot %#v",
			expectedEventConfiguration,
			convertedEventConfiguration,
		)
	}

	secretEncs := kubeAPI["secrets_encryption_config"].([]interface{})
	secretEnc := secretEncs[0].(map[string]interface{})
	customConfig, ok := secretEnc["custom_config"].(string)
	if !ok {
		t.Fatalf("expected secrets_encryption_config.custom_config to be string, got %T", secretEnc["custom_config"])
	}

	convertedCustomConfig, err := ghodssyamlToMapInterface(customConfig)
	if err != nil {
		t.Fatalf("failed to parse converted secrets_encryption_config.custom_config: %v", err)
	}
	if !reflect.DeepEqual(convertedCustomConfig, secretsEncryptionConfiguration) {
		t.Fatalf(
			"unexpected secrets_encryption_config.custom_config:\nexpected %#v\ngot %#v",
			secretsEncryptionConfiguration,
			convertedCustomConfig,
		)
	}

	admissionConfigurations, ok := kubeAPI["admission_configuration"].([]map[string]interface{})
	if !ok {
		t.Fatalf(
			"expected admission_configuration to be []map[string]interface{}, got %T",
			kubeAPI["admission_configuration"],
		)
	}

	if len(admissionConfigurations) != 1 {
		t.Fatalf("expected one admission configuration, got %d", len(admissionConfigurations))
	}

	if !reflect.DeepEqual(admissionConfigurations[0], admissionConfiguration) {
		t.Fatalf(
			"unexpected admission_configuration:\nexpected %#v\ngot %#v",
			admissionConfiguration,
			admissionConfigurations[0],
		)
	}
}

func TestResourceRancher2ClusterStateUpgrade(t *testing.T) {
	tests := []struct {
		name       string
		version    int
		attributes map[string]string
	}{
		{
			name:    "version 0",
			version: 0,
			attributes: map[string]string{
				"id":      "test-id",
				"name":    "test-cluster",
				"driver":  "gke",
				"ca_cert": "test-ca-cert",
			},
		},
		{
			name:    "version 1",
			version: 1,
			attributes: map[string]string{
				"id":                              "test-id",
				"name":                            "test-cluster",
				"driver":                          "rancherKubernetesEngine",
				"gke_config.#":                    "1",
				"gke_config.0.project_id":         "test-project",
				"rke_config.#":                    "1",
				"rke_config.0.kubernetes_version": "v1.24.10-rancher1-1",
				"cluster_template_id":             "ct-abcde",
			},
		},
		{
			name:    "version 2",
			version: 2,
			attributes: map[string]string{
				"id":                              "test-id",
				"name":                            "test-cluster",
				"driver":                          "rancherKubernetesEngine",
				"ca_cert":                         "test-ca-cert",
				"gke_config_v2.#":                 "1",
				"gke_config_v2.0.project_id":      "test-project",
				"rke_config.#":                    "1",
				"rke_config.0.kubernetes_version": "v1.24.10-rancher1-1",
				"cluster_template_id":             "ct-abcde",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := resourceRancher2Cluster()

			// Prevent the actual resource read from requiring a Rancher client.
			r.Read = func(*schema.ResourceData, interface{}) error {
				return nil
			}

			state := &terraform.InstanceState{
				ID:         "test-id",
				Attributes: tt.attributes,
				Meta: map[string]interface{}{
					"schema_version": tt.version,
				},
			}

			upgraded, err := r.Refresh(state, nil)
			if err != nil {
				t.Fatalf("failed to upgrade state: %v", err)
			}

			if upgraded == nil {
				t.Fatal("expected upgraded state")
			}

			if got := upgraded.Meta["schema_version"]; got != "3" && got != 3 {
				t.Fatalf("expected schema version 3, got %#v", got)
			}

			for _, key := range []string{
				"rke_config",
				"gke_config",
				"cluster_template_answers",
				"cluster_template_id",
				"cluster_template_questions",
				"cluster_template_revision_id",
			} {
				if _, exists := upgraded.Attributes[key]; exists {
					t.Fatalf("expected key %q to be removed from upgraded state", key)
				}
			}

			if got := upgraded.Attributes["name"]; got != "test-cluster" {
				t.Fatalf("expected name to be preserved, got %q", got)
			}
		})
	}
}

func TestResourceRancher2ClusterStateUpgradeV2_IdempotentWhenLegacyKeysMissing(t *testing.T) {
	initial := map[string]interface{}{
		"name":          "test-cluster",
		"driver":        "gke",
		"gke_config_v2": []interface{}{map[string]interface{}{"project_id": "my-project"}},
	}

	upgraded, err := resourceRancher2ClusterStateUpgradeV2(initial, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got, exists := upgraded["gke_config_v2"]; !exists || got == nil {
		t.Fatalf("expected gke_config_v2 to be preserved")
	}

	if got, exists := upgraded["name"]; !exists || got != "test-cluster" {
		t.Fatalf("expected name to be preserved")
	}
}
