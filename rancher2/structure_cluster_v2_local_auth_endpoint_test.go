package rancher2

import (
	"testing"

	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
	"github.com/stretchr/testify/assert"
)

var (
	testClusterV2LocalAuthEndpointConf                 rkev1.LocalClusterAuthEndpoint
	testClusterV2LocalAuthEndpointInterface            []interface{}
	testClusterV2LocalAuthEndpointInterfaceUseInternal []interface{}
	testClusterV2LocalAuthEndpointInterfaceWithFlag    []interface{}
)

func init() {
	testClusterV2LocalAuthEndpointConf = rkev1.LocalClusterAuthEndpoint{
		CACerts: "ca_certs",
		Enabled: true,
		FQDN:    "fqdn",
	}

	testClusterV2LocalAuthEndpointInterface = []interface{}{
		map[string]interface{}{
			"ca_certs": "ca_certs",
			"enabled":  true,
			"fqdn":     "fqdn",
		},
	}
	testClusterV2LocalAuthEndpointInterfaceUseInternal = []interface{}{
		map[string]interface{}{
			"ca_certs":              "ca_certs",
			"enabled":               true,
			"fqdn":                  "fqdn",
			"use_internal_ca_certs": true,
		},
	}
	testClusterV2LocalAuthEndpointInterfaceWithFlag = []interface{}{
		map[string]interface{}{
			"ca_certs":              "ca_certs",
			"enabled":               true,
			"fqdn":                  "fqdn",
			"use_internal_ca_certs": false,
		},
	}
}

func TestFlattenClusterV2LocalAuthEndpoint(t *testing.T) {

	cases := []struct {
		Input          rkev1.LocalClusterAuthEndpoint
		ExpectedOutput []interface{}
	}{
		{
			testClusterV2LocalAuthEndpointConf,
			testClusterV2LocalAuthEndpointInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterV2LocalAuthEndpoint(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from flattener.")
	}
}

func TestExpandClusterV2LocalAuthEndpoint(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput rkev1.LocalClusterAuthEndpoint
	}{
		{
			testClusterV2LocalAuthEndpointInterfaceUseInternal,
			testClusterV2LocalAuthEndpointConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterV2LocalAuthEndpoint(tc.Input)
		assert.Equal(t, tc.ExpectedOutput, output, "Unexpected output from expander.")
	}
}
