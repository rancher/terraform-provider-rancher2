package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testClusterRKEConfigServicesConf      *managementClient.RKEConfigServices
	testClusterRKEConfigServicesInterface []interface{}
)

func init() {
	testClusterRKEConfigServicesConf = &managementClient.RKEConfigServices{
		Etcd:           testClusterRKEConfigServicesETCDConf,
		KubeAPI:        testClusterRKEConfigServicesKubeAPIConf,
		KubeController: testClusterRKEConfigServicesKubeControllerConf,
		Kubelet:        testClusterRKEConfigServicesKubeletConf,
		Kubeproxy:      testClusterRKEConfigServicesKubeproxyConf,
		Scheduler:      testClusterRKEConfigServicesSchedulerConf,
	}
	testClusterRKEConfigServicesInterface = []interface{}{
		map[string]interface{}{
			"etcd":            testClusterRKEConfigServicesETCDInterface,
			"kube_api":        testClusterRKEConfigServicesKubeAPIInterface,
			"kube_controller": testClusterRKEConfigServicesKubeControllerInterface,
			"kubelet":         testClusterRKEConfigServicesKubeletInterface,
			"kubeproxy":       testClusterRKEConfigServicesKubeproxyInterface,
			"scheduler":       testClusterRKEConfigServicesSchedulerInterface,
		},
	}
}

func TestFlattenClusterRKEConfigServices(t *testing.T) {

	cases := []struct {
		Input          *managementClient.RKEConfigServices
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigServicesConf,
			testClusterRKEConfigServicesInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigServices(tc.Input, testClusterRKEConfigServicesInterface)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigServices(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.RKEConfigServices
	}{
		{
			testClusterRKEConfigServicesInterface,
			testClusterRKEConfigServicesConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigServices(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
