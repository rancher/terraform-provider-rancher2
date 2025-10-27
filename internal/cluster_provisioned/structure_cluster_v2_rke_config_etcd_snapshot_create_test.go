package rancher2

import (
	"testing"

	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
	"github.com/stretchr/testify/assert"
)

var (
	testClusterV2RKEConfigETCDSnapshotCreateConf      *rkev1.ETCDSnapshotCreate
	testClusterV2RKEConfigETCDSnapshotCreateInterface []interface{}
)

func init() {
	testClusterV2RKEConfigETCDSnapshotCreateConf = &rkev1.ETCDSnapshotCreate{
		Generation: 1,
	}

	testClusterV2RKEConfigETCDSnapshotCreateInterface = []interface{}{
		map[string]interface{}{
			"generation": 1,
		},
	}
}

func TestFlattenClusterV2RKEConfigETCDSnapshotCreate(t *testing.T) {
	cases := []struct {
		Input          *rkev1.ETCDSnapshotCreate
		ExpectedOutput []interface{}
	}{
		{
			testClusterV2RKEConfigETCDSnapshotCreateConf,
			testClusterV2RKEConfigETCDSnapshotCreateInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterV2RKEConfigETCDSnapshotCreate(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandClusterV2RKEConfigETCDSnapshotCreate(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *rkev1.ETCDSnapshotCreate
	}{
		{
			testClusterV2RKEConfigETCDSnapshotCreateInterface,
			testClusterV2RKEConfigETCDSnapshotCreateConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterV2RKEConfigETCDSnapshotCreate(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
