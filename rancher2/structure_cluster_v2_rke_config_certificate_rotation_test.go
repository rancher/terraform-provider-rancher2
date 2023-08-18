package rancher2

import (
	"testing"

	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
	"github.com/stretchr/testify/assert"
)

var (
	testClusterV2RKEConfigRotateCertificatesConf      *rkev1.RotateCertificates
	testClusterV2RKEConfigRotateCertificatesInterface []interface{}
)

func init() {
	testClusterV2RKEConfigRotateCertificatesConf = &rkev1.RotateCertificates{
		Generation: 2,
		Services: []string{
			"etcd",
			"kube-proxy",
		},
	}

	testClusterV2RKEConfigRotateCertificatesInterface = []interface{}{
		map[string]interface{}{
			"generation": 2,
			"services":   []interface{}{"etcd", "kube-proxy"},
		},
	}
}

func Test_flattenClusterV2RKEConfigRotateCertificates(t *testing.T) {
	cases := []struct {
		Input          *rkev1.RotateCertificates
		ExpectedOutput []interface{}
	}{
		{
			testClusterV2RKEConfigRotateCertificatesConf,
			testClusterV2RKEConfigRotateCertificatesInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterV2RKEConfigRotateCertificates(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandClusterV2RKEConfigRotateCertificates(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *rkev1.RotateCertificates
	}{
		{
			testClusterV2RKEConfigRotateCertificatesInterface,
			testClusterV2RKEConfigRotateCertificatesConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterV2RKEConfigRotateCertificates(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")

	}
}
