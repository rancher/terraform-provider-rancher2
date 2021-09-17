package rancher2

import (
	"reflect"
	"testing"

	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
)

var (
	testClusterV2RKEConfigLocalAuthEndpointConf      rkev1.LocalClusterAuthEndpoint
	testClusterV2RKEConfigLocalAuthEndpointInterface []interface{}
)

func init() {
	testClusterV2RKEConfigLocalAuthEndpointConf = rkev1.LocalClusterAuthEndpoint{
		CACerts: "ca_certs",
		Enabled: true,
		FQDN:    "fqdn",
	}

	testClusterV2RKEConfigLocalAuthEndpointInterface = []interface{}{
		map[string]interface{}{
			"ca_certs": "ca_certs",
			"enabled":  true,
			"fqdn":     "fqdn",
		},
	}
}

func TestFlattenClusterV2RKEConfigLocalAuthEndpoint(t *testing.T) {

	cases := []struct {
		Input          rkev1.LocalClusterAuthEndpoint
		ExpectedOutput []interface{}
	}{
		{
			testClusterV2RKEConfigLocalAuthEndpointConf,
			testClusterV2RKEConfigLocalAuthEndpointInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterV2RKEConfigLocalAuthEndpoint(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterV2RKEConfigLocalAuthEndpoint(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput rkev1.LocalClusterAuthEndpoint
	}{
		{
			testClusterV2RKEConfigLocalAuthEndpointInterface,
			testClusterV2RKEConfigLocalAuthEndpointConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterV2RKEConfigLocalAuthEndpoint(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
