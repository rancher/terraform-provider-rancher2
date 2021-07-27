package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var (
	testClusterRKEConfigIngressTolerationsConf            []managementClient.Toleration
	testClusterRKEConfigIngressTolerationsInterface       []interface{}
	testClusterRKEConfigIngressUpdateDaemonSetConf        *managementClient.RollingUpdateDaemonSet
	testClusterRKEConfigIngressUpdateDaemonSetInterface   []interface{}
	testClusterRKEConfigIngresstDaemonSetStrategyConf     *managementClient.DaemonSetUpdateStrategy
	testClusterRKEConfigIngressDaemonSetStrategyInterface []interface{}
	testClusterRKEConfigIngressConf                       *managementClient.IngressConfig
	testClusterRKEConfigIngressInterface                  []interface{}
)

func init() {
	seconds := int64(10)
	testClusterRKEConfigIngressTolerationsConf = []managementClient.Toleration{
		{
			Key:               "key",
			Value:             "value",
			Effect:            "recipient",
			Operator:          "operator",
			TolerationSeconds: &seconds,
		},
	}
	testClusterRKEConfigIngressTolerationsInterface = []interface{}{
		map[string]interface{}{
			"key":      "key",
			"value":    "value",
			"effect":   "recipient",
			"operator": "operator",
			"seconds":  10,
		},
	}
	testClusterRKEConfigIngressUpdateDaemonSetConf = &managementClient.RollingUpdateDaemonSet{
		MaxUnavailable: intstr.FromInt(10),
	}
	testClusterRKEConfigIngressUpdateDaemonSetInterface = []interface{}{
		map[string]interface{}{
			"max_unavailable": 10,
		},
	}
	testClusterRKEConfigIngresstDaemonSetStrategyConf = &managementClient.DaemonSetUpdateStrategy{
		RollingUpdate: testClusterRKEConfigIngressUpdateDaemonSetConf,
		Strategy:      "strategy",
	}
	testClusterRKEConfigIngressDaemonSetStrategyInterface = []interface{}{
		map[string]interface{}{
			"rolling_update": testClusterRKEConfigIngressUpdateDaemonSetInterface,
			"strategy":       "strategy",
		},
	}
	testClusterRKEConfigIngressConf = &managementClient.IngressConfig{
		DefaultBackend: newFalse(),
		DNSPolicy:      "test",
		ExtraArgs: map[string]string{
			"arg_one": "one",
			"arg_two": "two",
		},
		HTTPPort:    80,
		HTTPSPort:   443,
		NetworkMode: "test",
		NodeSelector: map[string]string{
			"node_one": "one",
			"node_two": "two",
		},
		Options: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
		Provider:       "test",
		Tolerations:    testClusterRKEConfigIngressTolerationsConf,
		UpdateStrategy: testClusterRKEConfigIngresstDaemonSetStrategyConf,
	}
	testClusterRKEConfigIngressInterface = []interface{}{
		map[string]interface{}{
			"default_backend": false,
			"dns_policy":      "test",
			"extra_args": map[string]interface{}{
				"arg_one": "one",
				"arg_two": "two",
			},
			"http_port":    80,
			"https_port":   443,
			"network_mode": "test",
			"node_selector": map[string]interface{}{
				"node_one": "one",
				"node_two": "two",
			},
			"options": map[string]interface{}{
				"option1": "value1",
				"option2": "value2",
			},
			"provider":        "test",
			"tolerations":     testClusterRKEConfigIngressTolerationsInterface,
			"update_strategy": testClusterRKEConfigIngressDaemonSetStrategyInterface,
		},
	}
}

func TestFlattenClusterRKEConfigIngress(t *testing.T) {

	cases := []struct {
		Input          *managementClient.IngressConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigIngressConf,
			testClusterRKEConfigIngressInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigIngress(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigIngress(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.IngressConfig
	}{
		{
			testClusterRKEConfigIngressInterface,
			testClusterRKEConfigIngressConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigIngress(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
