package rancher2

import "testing"

func TestResourceRancher2ClusterStateUpgradeV2_RemovesLegacyRKE1Fields(t *testing.T) {
	initial := map[string]interface{}{
		"name":                         "test-cluster",
		"driver":                       "eks",
		"eks_config_v2":                []interface{}{map[string]interface{}{"region": "us-east-1"}},
		"rke_config":                   []interface{}{map[string]interface{}{"kubernetes_version": "v1.25.0-rancher1-1"}},
		"cluster_template_answers":     []interface{}{map[string]interface{}{"values": map[string]interface{}{"foo": "bar"}}},
		"cluster_template_id":          "ct-abcde",
		"cluster_template_questions":   []interface{}{map[string]interface{}{"variable": "a", "default": "b"}},
		"cluster_template_revision_id": "ctr-abcde",
	}

	upgraded, err := resourceRancher2ClusterStateUpgradeV2(initial, nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	legacyKeys := []string{
		"rke_config",
		"cluster_template_answers",
		"cluster_template_id",
		"cluster_template_questions",
		"cluster_template_revision_id",
	}
	for _, key := range legacyKeys {
		if _, exists := upgraded[key]; exists {
			t.Fatalf("expected key %q to be removed", key)
		}
	}

	if got, exists := upgraded["eks_config_v2"]; !exists || got == nil {
		t.Fatalf("expected eks_config_v2 to be preserved")
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
