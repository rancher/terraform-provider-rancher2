package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// unknownValue mirrors hcl2shim.UnknownVariableValue, the in-band marker the SDK
// uses for a value that is not known until apply. That package is internal to
// the SDK, so the constant is repeated here.
const unknownValue = "74D93920-ED26-11E3-AC10-0800200C9A66"

const machineConfigNameKey = "rke_config.0.machine_pools.0.machine_config.0.name"

// clusterV2DiffConfig builds a minimal rancher2_cluster_v2 config with a single
// machine pool whose machine_config name is set to machineConfigName.
func clusterV2DiffConfig(machineConfigName string) *terraform.ResourceConfig {
	return clusterV2DiffConfigKind("VmwarevsphereConfig", machineConfigName)
}

// clusterV2DiffConfigKind is clusterV2DiffConfig with machine_config kind also
// parameterised, so the kind attribute can be driven unknown independently.
func clusterV2DiffConfigKind(machineConfigKind, machineConfigName string) *terraform.ResourceConfig {
	return terraform.NewResourceConfigRaw(map[string]interface{}{
		"name":               "test-cluster",
		"fleet_namespace":    "fleet-default",
		"kubernetes_version": "v1.31.5+rke2r1",
		"rke_config": []interface{}{
			map[string]interface{}{
				"machine_pools": []interface{}{
					map[string]interface{}{
						"name":                         "pool1",
						"cloud_credential_secret_name": "cattle-global-data:cc-test",
						"control_plane_role":           true,
						"etcd_role":                    true,
						"worker_role":                  false,
						"quantity":                     1,
						"machine_config": []interface{}{
							map[string]interface{}{
								"kind": machineConfigKind,
								"name": machineConfigName,
							},
						},
					},
				},
			},
		},
	})
}

// TestResourceRancher2ClusterV2DiffUnknownMachineConfigName covers
// rancher/terraform-provider-rancher2#1501: creating a rancher2_cluster_v2 in
// the same apply as the rancher2_machine_config_v2 it references.
//
// rancher2_machine_config_v2 takes generate_name and exposes .name as Computed
// only, so the name genuinely cannot be known at plan time. CustomizeDiff must
// leave that unknown in place. Normalizing the block with
// SetNew(flatten(expand(GetChange()))) plans a known "" instead, because
// GetChange collapses unknowns to the zero value, and Terraform then aborts the
// apply:
//
//	Provider produced inconsistent final plan ... produced an invalid new value
//	for .rke_config[0].machine_pools[0].machine_config[0].name: was
//	cty.StringVal(""), but now cty.StringVal("nc-mc-pool1-p5mqh")
func TestResourceRancher2ClusterV2DiffUnknownMachineConfigName(t *testing.T) {
	diff, err := resourceRancher2ClusterV2().Diff(nil, clusterV2DiffConfig(unknownValue), nil)
	if err != nil {
		t.Fatalf("Diff failed: %v", err)
	}
	if diff == nil {
		t.Fatal("expected a diff for a resource being created, got nil")
	}

	attr, ok := diff.Attributes[machineConfigNameKey]
	if !ok {
		// This is how the bug presents: SetNew planned "", which matched the
		// empty prior value, so the SDK dropped the attribute and the shim
		// reports a known "" to Terraform.
		t.Fatalf("%s was dropped from the planned diff; it must be planned as unknown", machineConfigNameKey)
	}
	if !attr.NewComputed {
		t.Errorf("%s was planned as the known value %q; it must stay unknown so apply can supply the generated name",
			machineConfigNameKey, attr.New)
	}
}

// TestResourceRancher2ClusterV2DiffKnownMachineConfigName is one half of the
// guard against the #1501 fix over-applying: a fully known machine_config name
// must be planned as that known value.
func TestResourceRancher2ClusterV2DiffKnownMachineConfigName(t *testing.T) {
	diff, err := resourceRancher2ClusterV2().Diff(nil, clusterV2DiffConfig("nc-mc-pool1-p5mqh"), nil)
	if err != nil {
		t.Fatalf("Diff failed: %v", err)
	}
	if diff == nil {
		t.Fatal("expected a diff for a resource being created, got nil")
	}

	attr, ok := diff.Attributes[machineConfigNameKey]
	if !ok {
		t.Fatalf("%s missing from the planned diff", machineConfigNameKey)
	}
	if attr.NewComputed {
		t.Errorf("%s was planned as unknown even though the config supplies a known name", machineConfigNameKey)
	}
	if attr.New != "nc-mc-pool1-p5mqh" {
		t.Errorf("%s planned as %q, want %q", machineConfigNameKey, attr.New, "nc-mc-pool1-p5mqh")
	}
}

// clusterV2PriorState is an existing cluster whose machine pool points at an
// already-created machine config named oldName.
func clusterV2PriorState(oldName string) *terraform.InstanceState {
	return &terraform.InstanceState{
		ID: "fleet-default/test-cluster",
		Attributes: map[string]string{
			"id":                 "fleet-default/test-cluster",
			"name":               "test-cluster",
			"fleet_namespace":    "fleet-default",
			"kubernetes_version": "v1.31.5+rke2r1",

			"rke_config.#":                                              "1",
			"rke_config.0.machine_pools.#":                              "1",
			"rke_config.0.machine_pools.0.name":                         "pool1",
			"rke_config.0.machine_pools.0.cloud_credential_secret_name": "cattle-global-data:cc-test",
			"rke_config.0.machine_pools.0.control_plane_role":           "true",
			"rke_config.0.machine_pools.0.etcd_role":                    "true",
			"rke_config.0.machine_pools.0.worker_role":                  "false",
			"rke_config.0.machine_pools.0.quantity":                     "1",
			"rke_config.0.machine_pools.0.machine_config.#":             "1",
			"rke_config.0.machine_pools.0.machine_config.0.kind":        "VmwarevsphereConfig",
			"rke_config.0.machine_pools.0.machine_config.0.name":        oldName,
		},
	}
}

// TestResourceRancher2ClusterV2DiffReplacedMachineConfigName covers the same defect
// reached by a different route: an existing cluster whose machine config is
// replaced, so generate_name yields a new name that is again unknown at plan
// time. Here the normalization plans the *old* name as the new value, which is
// worse than the create case -- "" at least looks obviously wrong.
func TestResourceRancher2ClusterV2DiffReplacedMachineConfigName(t *testing.T) {
	const oldName = "nc-mc-pool1-oldxx"

	diff, err := resourceRancher2ClusterV2().Diff(
		clusterV2PriorState(oldName), clusterV2DiffConfig(unknownValue), nil)
	if err != nil {
		t.Fatalf("Diff failed: %v", err)
	}
	if diff == nil {
		t.Fatal("expected a diff when the machine config name changes, got nil")
	}

	attr, ok := diff.Attributes[machineConfigNameKey]
	if !ok {
		t.Fatalf("%s was dropped from the planned diff; the cluster would keep pointing at the destroyed machine config %q", machineConfigNameKey, oldName)
	}
	if !attr.NewComputed {
		t.Errorf("%s was planned as the known value %q (old value was %q); it must stay unknown so apply can supply the regenerated name",
			machineConfigNameKey, attr.New, oldName)
	}
}

const machineConfigKindKey = "rke_config.0.machine_pools.0.machine_config.0.kind"

// TestResourceRancher2ClusterV2DiffUnknownMachineConfigKind covers the same
// defect on the sibling attribute. rancher2_machine_config_v2 exposes both kind
// and name as Computed only, so a config that references either of them --
// kind = rancher2_machine_config_v2.pool1.kind -- has an unknown here at plan
// time, and CustomizeDiff would collapse it to a known "" just the same.
//
// Fixing only name would leave this path broken; at least one reporter on #1501
// worked around it by hardcoding the kind literal.
func TestResourceRancher2ClusterV2DiffUnknownMachineConfigKind(t *testing.T) {
	diff, err := resourceRancher2ClusterV2().Diff(
		nil, clusterV2DiffConfigKind(unknownValue, "nc-mc-pool1-p5mqh"), nil)
	if err != nil {
		t.Fatalf("Diff failed: %v", err)
	}
	if diff == nil {
		t.Fatal("expected a diff for a resource being created, got nil")
	}

	attr, ok := diff.Attributes[machineConfigKindKey]
	if !ok {
		t.Fatalf("%s was dropped from the planned diff; it must be planned as unknown", machineConfigKindKey)
	}
	if !attr.NewComputed {
		t.Errorf("%s was planned as the known value %q; it must stay unknown so apply can supply the kind",
			machineConfigKindKey, attr.New)
	}
}
