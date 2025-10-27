package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var (
	testRollingUpdateInterface           []interface{}
	testRollingUpdateDeploymentConf      *managementClient.RollingUpdateDeployment
	testRollingUpdateDeploymentInterface []interface{}
	testRollingUpdateDaemonSetConf       *managementClient.RollingUpdateDaemonSet
	testRollingUpdateDaemonSetInterface  []interface{}
	testUpgradeStrategyInterface         []interface{}
	testDeploymentStrategyConf           *managementClient.DeploymentStrategy
	testDeploymentStrategyInterface      []interface{}
	testDaemonSetStrategyConf            *managementClient.DaemonSetUpdateStrategy
	testDaemonSetStrategyInterface       []interface{}
)

func init() {
	testRollingUpdateInterface = []interface{}{
		map[string]interface{}{
			"batch_size": 10,
			"interval":   10,
		},
	}
	testRollingUpdateDeploymentConf = &managementClient.RollingUpdateDeployment{
		MaxSurge:       intstr.FromInt(10),
		MaxUnavailable: intstr.FromInt(10),
	}
	testRollingUpdateDeploymentInterface = []interface{}{
		map[string]interface{}{
			"max_surge":       10,
			"max_unavailable": 10,
		},
	}
	testRollingUpdateDaemonSetConf = &managementClient.RollingUpdateDaemonSet{
		MaxUnavailable: intstr.FromInt(10),
	}
	testRollingUpdateDaemonSetInterface = []interface{}{
		map[string]interface{}{
			"max_unavailable": 10,
		},
	}
	testUpgradeStrategyInterface = []interface{}{
		map[string]interface{}{
			"rolling_update": testRollingUpdateInterface,
		},
	}
	testDeploymentStrategyConf = &managementClient.DeploymentStrategy{
		RollingUpdate: testRollingUpdateDeploymentConf,
		Strategy:      "strategy",
	}
	testDeploymentStrategyInterface = []interface{}{
		map[string]interface{}{
			"rolling_update": testRollingUpdateDeploymentInterface,
			"strategy":       "strategy",
		},
	}
	testDaemonSetStrategyConf = &managementClient.DaemonSetUpdateStrategy{
		RollingUpdate: testRollingUpdateDaemonSetConf,
		Strategy:      "strategy",
	}
	testDaemonSetStrategyInterface = []interface{}{
		map[string]interface{}{
			"rolling_update": testRollingUpdateDaemonSetInterface,
			"strategy":       "strategy",
		},
	}
}

func TestFlattenRollingUpdateDaemonSet(t *testing.T) {

	cases := []struct {
		Input          *managementClient.RollingUpdateDaemonSet
		ExpectedOutput []interface{}
	}{
		{
			testRollingUpdateDaemonSetConf,
			testRollingUpdateDaemonSetInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRollingUpdateDaemonSet(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestFlattenRollingUpdateDeployment(t *testing.T) {

	cases := []struct {
		Input          *managementClient.RollingUpdateDeployment
		ExpectedOutput []interface{}
	}{
		{
			testRollingUpdateDeploymentConf,
			testRollingUpdateDeploymentInterface,
		},
	}

	for _, tc := range cases {
		output := flattenRollingUpdateDeployment(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestFlattenDaemonSetStrategy(t *testing.T) {

	cases := []struct {
		Input          *managementClient.DaemonSetUpdateStrategy
		ExpectedOutput []interface{}
	}{
		{
			testDaemonSetStrategyConf,
			testDaemonSetStrategyInterface,
		},
	}

	for _, tc := range cases {
		output := flattenDaemonSetStrategy(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestFlattenDeploymentStrategy(t *testing.T) {

	cases := []struct {
		Input          *managementClient.DeploymentStrategy
		ExpectedOutput []interface{}
	}{
		{
			testDeploymentStrategyConf,
			testDeploymentStrategyInterface,
		},
	}

	for _, tc := range cases {
		output := flattenDeploymentStrategy(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandRollingUpdateDaemonSet(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.RollingUpdateDaemonSet
	}{
		{
			testRollingUpdateDaemonSetInterface,
			testRollingUpdateDaemonSetConf,
		},
	}

	for _, tc := range cases {
		output := expandRollingUpdateDaemonSet(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}

func TestExpandRollingUpdateDeployment(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.RollingUpdateDeployment
	}{
		{
			testRollingUpdateDeploymentInterface,
			testRollingUpdateDeploymentConf,
		},
	}

	for _, tc := range cases {
		output := expandRollingUpdateDeployment(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}

func TestExpandDaemonSetStrategy(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.DaemonSetUpdateStrategy
	}{
		{
			testDaemonSetStrategyInterface,
			testDaemonSetStrategyConf,
		},
	}

	for _, tc := range cases {
		output := expandDaemonSetStrategy(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}

func TestExpandDeploymentStrategy(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.DeploymentStrategy
	}{
		{
			testDeploymentStrategyInterface,
			testDeploymentStrategyConf,
		},
	}

	for _, tc := range cases {
		output := expandDeploymentStrategy(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
