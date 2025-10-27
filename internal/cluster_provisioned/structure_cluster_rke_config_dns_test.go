package rancher2

import (
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var (
	testClusterRKEConfigDNSTolerationsConf                      []managementClient.Toleration
	testClusterRKEConfigDNSTolerationsInterface                 []interface{}
	testClusterRKEConfigDNSNodelocalConf                        *managementClient.Nodelocal
	testClusterRKEConfigDNSLinearAutoscalerParamsConf           *managementClient.LinearAutoscalerParams
	testClusterRKEConfigDNSLinearAutoscalerParamsConfFalse      *managementClient.LinearAutoscalerParams
	testClusterRKEConfigDNSNodelocalInterface                   []interface{}
	testClusterRKEConfigDNSLinearAutoscalerParamsInterface      []interface{}
	testClusterRKEConfigDNSLinearAutoscalerParamsInterfaceFalse []interface{}
	testClusterRKEConfigDNSConf                                 *managementClient.DNSConfig
	testClusterRKEConfigDNSInterface                            []interface{}
)

func init() {
	seconds := int64(10)
	testClusterRKEConfigDNSTolerationsConf = []managementClient.Toleration{
		{
			Key:               "key",
			Value:             "value",
			Effect:            "recipient",
			Operator:          "operator",
			TolerationSeconds: &seconds,
		},
	}
	testClusterRKEConfigDNSTolerationsInterface = []interface{}{
		map[string]interface{}{
			"key":      "key",
			"value":    "value",
			"effect":   "recipient",
			"operator": "operator",
			"seconds":  10,
		},
	}
	testClusterRKEConfigDNSNodelocalConf = &managementClient.Nodelocal{
		NodeSelector: map[string]string{
			"sel1": "value1",
			"sel2": "value2",
		},
		IPAddress: "ip_address",
	}
	testClusterRKEConfigDNSLinearAutoscalerParamsConf = &managementClient.LinearAutoscalerParams{
		CoresPerReplica:           float64(128),
		Max:                       int64(0),
		Min:                       int64(1),
		NodesPerReplica:           float64(4),
		PreventSinglePointFailure: true,
	}
	testClusterRKEConfigDNSLinearAutoscalerParamsConfFalse = &managementClient.LinearAutoscalerParams{
		CoresPerReplica:           float64(64),
		Max:                       int64(10),
		Min:                       int64(1),
		NodesPerReplica:           float64(8),
		PreventSinglePointFailure: false,
	}
	testClusterRKEConfigDNSNodelocalInterface = []interface{}{
		map[string]interface{}{
			"node_selector": map[string]interface{}{
				"sel1": "value1",
				"sel2": "value2",
			},
			"ip_address": "ip_address",
		},
	}
	testClusterRKEConfigDNSLinearAutoscalerParamsInterface = []interface{}{
		map[string]interface{}{
			"cores_per_replica":            float64(128),
			"max":                          0,
			"min":                          1,
			"nodes_per_replica":            float64(4),
			"prevent_single_point_failure": true,
		},
	}
	testClusterRKEConfigDNSLinearAutoscalerParamsInterfaceFalse = []interface{}{
		map[string]interface{}{
			"cores_per_replica":            float64(64),
			"max":                          10,
			"min":                          1,
			"nodes_per_replica":            float64(8),
			"prevent_single_point_failure": false,
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
	testClusterRKEConfigDNSConf = &managementClient.DNSConfig{
		NodeSelector: map[string]string{
			"sel1": "value1",
			"sel2": "value2",
		},
		Nodelocal:              testClusterRKEConfigDNSNodelocalConf,
		LinearAutoscalerParams: testClusterRKEConfigDNSLinearAutoscalerParamsConf,
		Options: map[string]string{
			"opt1": "value1",
			"opt2": "value2",
		},
		Provider:            "kube-dns",
		ReverseCIDRs:        []string{"rev1", "rev2"},
		Tolerations:         testClusterRKEConfigDNSTolerationsConf,
		UpstreamNameservers: []string{"up1", "up2"},
		UpdateStrategy:      testDeploymentStrategyConf,
	}
	testClusterRKEConfigDNSInterface = []interface{}{
		map[string]interface{}{
			"node_selector": map[string]interface{}{
				"sel1": "value1",
				"sel2": "value2",
			},
			"nodelocal":                testClusterRKEConfigDNSNodelocalInterface,
			"linear_autoscaler_params": testClusterRKEConfigDNSLinearAutoscalerParamsInterface,
			"options": map[string]interface{}{
				"opt1": "value1",
				"opt2": "value2",
			},
			"provider":             "kube-dns",
			"reverse_cidrs":        []interface{}{"rev1", "rev2"},
			"tolerations":          testClusterRKEConfigDNSTolerationsInterface,
			"upstream_nameservers": []interface{}{"up1", "up2"},
			"update_strategy":      testDeploymentStrategyInterface,
		},
	}
}

func TestFlattenClusterRKEConfigDNSNodelocal(t *testing.T) {

	cases := []struct {
		Input          *managementClient.Nodelocal
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigDNSNodelocalConf,
			testClusterRKEConfigDNSNodelocalInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterRKEConfigDNSNodelocal(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestFlattenClusterRKEConfigDNSLinearAutoscalerParams(t *testing.T) {

	cases := []struct {
		Input          *managementClient.LinearAutoscalerParams
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigDNSLinearAutoscalerParamsConf,
			testClusterRKEConfigDNSLinearAutoscalerParamsInterface,
		},
		{
			testClusterRKEConfigDNSLinearAutoscalerParamsConfFalse,
			testClusterRKEConfigDNSLinearAutoscalerParamsInterfaceFalse,
		},
	}

	for _, tc := range cases {
		output := flattenClusterRKEConfigDNSLinearAutoscalerParams(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestFlattenClusterRKEConfigDNS(t *testing.T) {

	cases := []struct {
		Input          *managementClient.DNSConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigDNSConf,
			testClusterRKEConfigDNSInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigDNS(tc.Input)
		if err != nil {
			assert.FailNow(t, "[ERROR] on flattener: %#v", err)
		}
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandClusterRKEConfigDNSNodelocal(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.Nodelocal
	}{
		{
			testClusterRKEConfigDNSNodelocalInterface,
			testClusterRKEConfigDNSNodelocalConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterRKEConfigDNSNodelocal(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}

func TestExpandClusterRKEConfigDNSLinearAutoscalerParams(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.LinearAutoscalerParams
	}{
		{
			testClusterRKEConfigDNSLinearAutoscalerParamsInterface,
			testClusterRKEConfigDNSLinearAutoscalerParamsConf,
		},
		{
			testClusterRKEConfigDNSLinearAutoscalerParamsInterfaceFalse,
			testClusterRKEConfigDNSLinearAutoscalerParamsConfFalse,
		},
	}

	for _, tc := range cases {
		output := expandClusterRKEConfigDNSLinearAutoscalerParams(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}

func TestExpandClusterRKEConfigDNS(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.DNSConfig
	}{
		{
			testClusterRKEConfigDNSInterface,
			testClusterRKEConfigDNSConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigDNS(tc.Input)
		if err != nil {
			assert.FailNow(t, "[ERROR] on expander: %#v", err)
		}
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
