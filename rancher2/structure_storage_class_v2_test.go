package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"k8s.io/api/core/v1"
	storV1 "k8s.io/api/storage/v1"
)

var (
	testStorageClassV2Conf      *StorageClassV2
	testStorageClassV2Interface map[string]interface{}
)

func init() {
	testStorageClassV2Conf = &StorageClassV2{}

	testStorageClassV2Conf.TypeMeta.Kind = storageClassV2Kind
	testStorageClassV2Conf.TypeMeta.APIVersion = storageClassV2APIVersion

	testStorageClassV2Conf.ObjectMeta.Name = "name"
	testStorageClassV2Conf.ObjectMeta.Annotations = map[string]string{
		"value1": "one",
		"value2": "two",
	}
	testStorageClassV2Conf.ObjectMeta.Labels = map[string]string{
		"label1": "one",
		"label2": "two",
	}
	testStorageClassV2Conf.Provisioner = "provisioner"
	testStorageClassV2Conf.AllowVolumeExpansion = newTrue()
	testStorageClassV2Conf.MountOptions = []string{"mount1", "mount2"}
	testStorageClassV2Conf.Parameters = map[string]string{
		"param1": "one",
		"param2": "two",
	}
	reclaim := v1.PersistentVolumeReclaimPolicy("reclaim_policy")
	testStorageClassV2Conf.ReclaimPolicy = &reclaim
	binding := storV1.VolumeBindingMode("volume_binding_mode")
	testStorageClassV2Conf.VolumeBindingMode = &binding

	testStorageClassV2Interface = map[string]interface{}{
		"name":                   "name",
		"k8s_provisioner":        "provisioner",
		"allow_volume_expansion": true,
		"mount_options":          []interface{}{"mount1", "mount2"},
		"parameters": map[string]interface{}{
			"param1": "one",
			"param2": "two",
		},
		"reclaim_policy":      "reclaim_policy",
		"volume_binding_mode": "volume_binding_mode",
		"annotations": map[string]interface{}{
			"value1": "one",
			"value2": "two",
		},
		"labels": map[string]interface{}{
			"label1": "one",
			"label2": "two",
		},
	}
}

func TestFlattenStorageClassV2(t *testing.T) {

	cases := []struct {
		Input          *StorageClassV2
		ExpectedOutput map[string]interface{}
	}{
		{
			testStorageClassV2Conf,
			testStorageClassV2Interface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, storageClassV2Fields(), tc.ExpectedOutput)
		err := flattenStorageClassV2(output, tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		if !reflect.DeepEqual(expectedOutput, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, expectedOutput)
		}
	}
}

func TestExpandStorageClassV2(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *StorageClassV2
	}{
		{
			testStorageClassV2Interface,
			testStorageClassV2Conf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, storageClassV2Fields(), tc.Input)
		output := expandStorageClassV2(inputResourceData)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
