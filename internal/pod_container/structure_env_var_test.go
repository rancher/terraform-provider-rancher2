package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testEnvVarConf      []managementClient.EnvVar
	testEnvVarInterface []interface{}
)

func init() {
	testEnvVarConf = []managementClient.EnvVar{
		{
			Name:  "name1",
			Value: "value1",
		},
		{
			Name:  "name2",
			Value: "value2",
		},
	}
	testEnvVarInterface = []interface{}{
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

func TestFlattenEnvVars(t *testing.T) {

	cases := []struct {
		Input          []managementClient.EnvVar
		ExpectedOutput []interface{}
	}{
		{
			testEnvVarConf,
			testEnvVarInterface,
		},
	}
	for _, tc := range cases {
		output := flattenEnvVars(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandEnvVars(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []managementClient.EnvVar
	}{
		{
			testEnvVarInterface,
			testEnvVarConf,
		},
	}
	for _, tc := range cases {
		output := expandEnvVars(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
