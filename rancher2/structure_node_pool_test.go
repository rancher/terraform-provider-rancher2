package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testNodePoolNodeTaintsConf      []managementClient.Taint
	testNodePoolNodeTaintsInterface []interface{}
	testNodePoolConf                *managementClient.NodePool
	testNodePoolInterface           map[string]interface{}
)

func init() {
	testNodePoolNodeTaintsConf = []managementClient.Taint{
		{
			Key:       "key",
			Value:     "value",
			Effect:    "recipient",
			TimeAdded: "time_added",
		},
	}
	testNodePoolNodeTaintsInterface = []interface{}{
		map[string]interface{}{
			"key":        "key",
			"value":      "value",
			"effect":     "recipient",
			"time_added": "time_added",
		},
	}
	testNodePoolConf = &managementClient.NodePool{
		ClusterID:               "cluster-test",
		Name:                    "test",
		DeleteNotReadyAfterSecs: 0,
		DrainBeforeDelete:       true,
		HostnamePrefix:          "terraform-test",
		NodeTaints:              testNodePoolNodeTaintsConf,
		NodeTemplateID:          "node-test",
		Quantity:                3,
		ControlPlane:            true,
		Etcd:                    true,
		Worker:                  true,
		Annotations: map[string]string{
			"one": "one",
			"two": "two",
		},
		Labels: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
	}
	testNodePoolInterface = map[string]interface{}{
		"cluster_id":                  "cluster-test",
		"name":                        "test",
		"delete_not_ready_after_secs": 0,
		"drain_before_delete":         true,
		"hostname_prefix":             "terraform-test",
		"node_taints":                 testNodePoolNodeTaintsInterface,
		"node_template_id":            "node-test",
		"quantity":                    3,
		"control_plane":               true,
		"etcd":                        true,
		"worker":                      true,
		"annotations": map[string]interface{}{
			"one": "one",
			"two": "two",
		},
		"labels": map[string]interface{}{
			"option1": "value1",
			"option2": "value2",
		},
	}
}

func TestFlattenNodePool(t *testing.T) {

	cases := []struct {
		Input          *managementClient.NodePool
		ExpectedOutput map[string]interface{}
	}{
		{
			testNodePoolConf,
			testNodePoolInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, nodePoolFields(), map[string]interface{}{})
		err := flattenNodePool(output, tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		if !reflect.DeepEqual(expectedOutput, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				expectedOutput, tc.ExpectedOutput)
		}
	}
}

func TestExpandNodePool(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *managementClient.NodePool
	}{
		{
			testNodePoolInterface,
			testNodePoolConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, nodePoolFields(), tc.Input)
		output := expandNodePool(inputResourceData)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
