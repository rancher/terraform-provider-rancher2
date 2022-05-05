package rancher2

import (
	"reflect"
	"testing"

	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
)

var (
	testClusterV2RKEConfigETCDSnapshotRestoreConf      *rkev1.ETCDSnapshotRestore
	testClusterV2RKEConfigETCDSnapshotRestoreInterface []interface{}
)

func init() {
	testClusterV2RKEConfigETCDSnapshotRestoreConf = &rkev1.ETCDSnapshotRestore{
		Name:             "SnapshotTestName",
		Generation:       1,
		RestoreRKEConfig: "all",
	}

	testClusterV2RKEConfigETCDSnapshotRestoreInterface = []interface{}{
		map[string]interface{}{
			"name":               "SnapshotTestName",
			"generation":         1,
			"restore_rke_config": "all",
		},
	}
}

func TestFlattenClusterV2RKEConfigETCDSnapshotRestore(t *testing.T) {
	cases := []struct {
		Input          *rkev1.ETCDSnapshotRestore
		ExpectedOutput []interface{}
	}{
		{
			testClusterV2RKEConfigETCDSnapshotRestoreConf,
			testClusterV2RKEConfigETCDSnapshotRestoreInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterV2RKEConfigETCDSnapshotRestore(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterV2RKEConfigETCDSnapshotRestore(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *rkev1.ETCDSnapshotRestore
	}{
		{
			testClusterV2RKEConfigETCDSnapshotRestoreInterface,
			testClusterV2RKEConfigETCDSnapshotRestoreConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterV2RKEConfigETCDSnapshotRestore(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
