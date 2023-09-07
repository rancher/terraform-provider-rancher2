package rancher2

import (
	"testing"

	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
	"github.com/stretchr/testify/assert"
)

var (
	testEnvVarV2Conf      []rkev1.EnvVar
	testEnvVarV2Interface []interface{}
)

func init() {
	testEnvVarV2Conf = []rkev1.EnvVar{
		{
			Name:  "name1",
			Value: "value1",
		},
		{
			Name:  "name2",
			Value: "value2",
		},
	}
	testEnvVarV2Interface = []interface{}{
		map[string]interface{}{
			"name":  "name1",
			"value": "value1",
		},
		map[string]interface{}{
			"name":  "name2",
			"value": "value2",
		},
	}
}

func TestFlattenEnvVarsV2(t *testing.T) {

	cases := []struct {
		Input          []rkev1.EnvVar
		ExpectedOutput []interface{}
	}{
		{
			testEnvVarV2Conf,
			testEnvVarV2Interface,
		},
	}

	for _, tc := range cases {
		output := flattenEnvVarsV2(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandEnvVarsV2(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []rkev1.EnvVar
	}{
		{
			testEnvVarV2Interface,
			testEnvVarV2Conf,
		},
	}

	for _, tc := range cases {
		output := expandEnvVarsV2(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
