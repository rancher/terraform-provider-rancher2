package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

var (
	testPodSecurityPolicyRuntimeClassStrategyConf           *managementClient.RuntimeClassStrategyOptions
	testPodSecurityPolicyRuntimeClassStrategyInterface      []interface{}
	testNilPodSecurityPolicyRuntimeClassStrategyConf        *managementClient.RuntimeClassStrategyOptions
	testEmptyPodSecurityPolicyRuntimeClassStrategyInterface []interface{}
)

func init() {
	testPodSecurityPolicyRuntimeClassStrategyConf = &managementClient.RuntimeClassStrategyOptions{
		AllowedRuntimeClassNames: []string{"foo", "bar"},
		DefaultRuntimeClassName:  "foo",
	}
	testPodSecurityPolicyRuntimeClassStrategyInterface = []interface{}{
		map[string]interface{}{
			"allowed_runtime_class_names": toArrayInterface([]string{"foo", "bar"}),
			"default_runtime_class_name":  "foo",
		},
	}
	testEmptyPodSecurityPolicyRuntimeClassStrategyInterface = []interface{}{}
}

func TestFlattenPodSecurityPolicyRuntimeClassStrategy(t *testing.T) {

	cases := []struct {
		Input          *managementClient.RuntimeClassStrategyOptions
		ExpectedOutput []interface{}
	}{
		{
			testPodSecurityPolicyRuntimeClassStrategyConf,
			testPodSecurityPolicyRuntimeClassStrategyInterface,
		},
		{
			testNilPodSecurityPolicyRuntimeClassStrategyConf,
			testEmptyPodSecurityPolicyRuntimeClassStrategyInterface,
		},
	}

	for _, tc := range cases {
		output := flattenPodSecurityPolicyRuntimeClassStrategy(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandPodSecurityPolicyRuntimeClassStrategy(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.RuntimeClassStrategyOptions
	}{
		{
			testPodSecurityPolicyRuntimeClassStrategyInterface,
			testPodSecurityPolicyRuntimeClassStrategyConf,
		},
	}

	for _, tc := range cases {
		output := expandPodSecurityPolicyRuntimeClassStrategy(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
