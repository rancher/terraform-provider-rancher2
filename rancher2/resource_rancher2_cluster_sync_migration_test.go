package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestResourceRancher2ClusterSyncStateUpgraders(t *testing.T) {
	r := resourceRancher2ClusterSync()

	if r.SchemaVersion != 1 {
		t.Fatalf("expected schema version 1, got %d", r.SchemaVersion)
	}

	if len(r.StateUpgraders) != 1 {
		t.Fatalf("expected 1 state upgrader, got %d", len(r.StateUpgraders))
	}

	if r.StateUpgraders[0].Version != 0 {
		t.Fatalf("expected upgrader to have version 0, got %d", r.StateUpgraders[0].Version)
	}

	expectedType := resourceRancher2ClusterSyncResourceV0().CoreConfigSchema().ImpliedType()
	if !r.StateUpgraders[0].Type.Equals(expectedType) {
		t.Fatalf("expected upgrader to use the V0 schema type")
	}

	if err := r.InternalValidate(nil, true); err != nil {
		t.Fatalf("resource failed internal validation: %v", err)
	}
}

func TestResourceRancher2ClusterSyncStateUpgradeV0(t *testing.T) {
	rawState := map[string]interface{}{
		"cluster_id":    "test-id",
		"node_pool_ids": []interface{}{"np1", "np2"},
		"nodes": []interface{}{
			map[string]interface{}{
				"id":               "node1",
				"node_pool_id":     "np1",
				"node_template_id": "nt1",
				"ssh_user":         "root",
			},
		},
	}

	upgraded, err := resourceRancher2ClusterSyncStateUpgradeV0(rawState, nil)
	if err != nil {
		t.Fatalf("failed to upgrade V0 state: %v", err)
	}

	if _, exists := upgraded["node_pool_ids"]; exists {
		t.Fatalf("expected node_pool_ids to be removed from upgraded state")
	}

	node := upgraded["nodes"].([]interface{})[0].(map[string]interface{})
	for _, key := range []string{"node_pool_id", "node_template_id", "ssh_user"} {
		if _, exists := node[key]; exists {
			t.Fatalf("expected key %q to be removed from upgraded node", key)
		}
	}

	if got := upgraded["cluster_id"]; got != "test-id" {
		t.Fatalf("expected cluster_id to be preserved, got %#v", got)
	}
	if got := node["id"]; got != "node1" {
		t.Fatalf("expected node id to be preserved, got %#v", got)
	}
}

func TestResourceRancher2ClusterSyncStateUpgrade(t *testing.T) {
	r := resourceRancher2ClusterSync()

	// Prevent the actual resource read from requiring a Rancher client.
	r.Read = func(*schema.ResourceData, interface{}) error {
		return nil
	}

	state := &terraform.InstanceState{
		ID: "test-id",
		Attributes: map[string]string{
			"id":              "test-id",
			"cluster_id":      "test-id",
			"node_pool_ids.#": "1",
			"node_pool_ids.0": "np1",
			"nodes.#":         "1",
			"nodes.0.id":      "node1",
		},
		Meta: map[string]interface{}{
			"schema_version": "0",
		},
	}

	upgraded, err := r.Refresh(state, nil)
	if err != nil {
		t.Fatalf("failed to upgrade state: %v", err)
	}

	if upgraded == nil {
		t.Fatal("expected upgraded state")
	}

	if got := upgraded.Meta["schema_version"]; got != "1" && got != 1 {
		t.Fatalf("expected schema version 1, got %#v", got)
	}

	if _, exists := upgraded.Attributes["node_pool_ids.#"]; exists {
		t.Fatalf("expected node_pool_ids to be removed from upgraded state")
	}

	if got := upgraded.Attributes["cluster_id"]; got != "test-id" {
		t.Fatalf("expected cluster_id to be preserved, got %q", got)
	}
}
