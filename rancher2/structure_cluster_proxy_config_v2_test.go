package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	norman "github.com/rancher/norman/types"
	managementv3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	testClusterProxyConfigV2Conf      *ClusterProxyConfigV2
	testClusterProxyConfigV2Interface map[string]interface{}
)

func init() {
	testClusterProxyConfigV2Conf = &ClusterProxyConfigV2{
		Resource: norman.Resource{
			ID: "c-m-xxxxx",
		},
		ClusterProxyConfig: managementv3.ClusterProxyConfig{
			TypeMeta: metav1.TypeMeta{
				Kind:       clusterProxyConfigV2Kind,
				APIVersion: clusterProxyConfigV2APIVersion,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      clusterProxyConfigV2Name,
				Namespace: "c-m-xxxxx",
			},
			Enabled: true,
		},
	}

	testClusterProxyConfigV2Interface = map[string]interface{}{
		"cluster_id": "c-m-xxxxx",
		"enabled":    true,
	}
}

func TestFlattenClusterProxyConfigV2(t *testing.T) {
	cases := []struct {
		Input          *ClusterProxyConfigV2
		ExpectedOutput map[string]interface{}
	}{
		{
			testClusterProxyConfigV2Conf,
			testClusterProxyConfigV2Interface,
		},
	}

	for _, tc := range cases {
		resourceData := schema.TestResourceDataRaw(t, clusterProxyConfigV2Fields(), map[string]interface{}{})
		err := flattenClusterProxyConfigV2(resourceData, tc.Input)
		if err != nil {
			t.Fatalf("Error flattening ClusterProxyConfigV2: %v", err)
		}

		for k, v := range tc.ExpectedOutput {
			actual := resourceData.Get(k)
			if !reflect.DeepEqual(actual, v) {
				t.Fatalf("flattenClusterProxyConfigV2: Expected %#v for key %s, got %#v", v, k, actual)
			}
		}
	}
}

func TestFlattenClusterProxyConfigV2Nil(t *testing.T) {
	resourceData := schema.TestResourceDataRaw(t, clusterProxyConfigV2Fields(), map[string]interface{}{})
	err := flattenClusterProxyConfigV2(resourceData, nil)
	if err == nil {
		t.Fatal("Expected error when flattening nil ClusterProxyConfigV2, got nil")
	}
}

func TestExpandClusterProxyConfigV2(t *testing.T) {
	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput *ClusterProxyConfigV2
	}{
		{
			testClusterProxyConfigV2Interface,
			testClusterProxyConfigV2Conf,
		},
	}

	for _, tc := range cases {
		resourceData := schema.TestResourceDataRaw(t, clusterProxyConfigV2Fields(), tc.Input)
		output, err := expandClusterProxyConfigV2(resourceData)
		if err != nil {
			t.Fatalf("Error expanding ClusterProxyConfigV2: %v", err)
		}

		// Check basic fields
		if output.TypeMeta.APIVersion != tc.ExpectedOutput.TypeMeta.APIVersion {
			t.Fatalf("expandClusterProxyConfigV2: Expected APIVersion %s, got %s", tc.ExpectedOutput.TypeMeta.APIVersion, output.TypeMeta.APIVersion)
		}

		if output.TypeMeta.Kind != tc.ExpectedOutput.TypeMeta.Kind {
			t.Fatalf("expandClusterProxyConfigV2: Expected Kind %s, got %s", tc.ExpectedOutput.TypeMeta.Kind, output.TypeMeta.Kind)
		}

		if output.ObjectMeta.Name != tc.ExpectedOutput.ObjectMeta.Name {
			t.Fatalf("expandClusterProxyConfigV2: Expected Name %s, got %s", tc.ExpectedOutput.ObjectMeta.Name, output.ObjectMeta.Name)
		}

		if output.ObjectMeta.Namespace != tc.ExpectedOutput.ObjectMeta.Namespace {
			t.Fatalf("expandClusterProxyConfigV2: Expected Namespace %s, got %s", tc.ExpectedOutput.ObjectMeta.Namespace, output.ObjectMeta.Namespace)
		}

		// Check enabled field
		if output.Enabled != tc.ExpectedOutput.Enabled {
			t.Fatalf("expandClusterProxyConfigV2: Expected enabled %v, got %v", tc.ExpectedOutput.Enabled, output.Enabled)
		}

	}
}

func TestExpandClusterProxyConfigV2Nil(t *testing.T) {
	_, err := expandClusterProxyConfigV2(nil)
	if err == nil {
		t.Fatal("Expected error when expanding nil ResourceData, got nil")
	}
}
