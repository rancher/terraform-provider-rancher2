package rancher2

import (
	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
	"reflect"
	"testing"
)

var (
	testClusterV2RKEConfigRotateEncryptionKeysConf      *rkev1.RotateEncryptionKeys
	testClusterV2RKEConfigRotateEncryptionKeysInterface []interface{}
)

func init() {
	testClusterV2RKEConfigRotateEncryptionKeysConf = &rkev1.RotateEncryptionKeys{
		Generation: 2,
	}

	testClusterV2RKEConfigRotateEncryptionKeysInterface = []interface{}{
		map[string]interface{}{
			"generation": 2,
		},
	}
}

func Test_flattenClusterV2RKEConfigRotateEncryptionKeys(t *testing.T) {
	cases := []struct {
		Input          *rkev1.RotateEncryptionKeys
		ExpectedOutput []interface{}
	}{
		{
			testClusterV2RKEConfigRotateEncryptionKeysConf,
			testClusterV2RKEConfigRotateEncryptionKeysInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterV2RKEConfigRotateEncryptionKeys(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterV2RKEConfigRotateEncryptionKeys(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *rkev1.RotateEncryptionKeys
	}{
		{
			testClusterV2RKEConfigRotateEncryptionKeysInterface,
			testClusterV2RKEConfigRotateEncryptionKeysConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterV2RKEConfigRotateEncryptionKeys(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
