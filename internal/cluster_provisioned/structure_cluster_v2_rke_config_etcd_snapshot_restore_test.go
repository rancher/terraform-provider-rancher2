package rancher2

import (
	"testing"

	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
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
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
