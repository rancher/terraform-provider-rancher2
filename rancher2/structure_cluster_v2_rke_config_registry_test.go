package rancher2

import (
	"reflect"
	"testing"

	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
)

var (
	testClusterV2RKEConfigRegistryConfigsConf      map[string]rkev1.RegistryConfig
	testClusterV2RKEConfigRegistryConfigsInterface []interface{}
	testClusterV2RKEConfigRegistryMirrorsConf      map[string]rkev1.Mirror
	testClusterV2RKEConfigRegistryMirrorsInterface []interface{}
	testClusterV2RKEConfigRegistryConf             *rkev1.Registry
	testClusterV2RKEConfigRegistryInterface        []interface{}
)

func init() {
	testClusterV2RKEConfigRegistryConfigsConf = map[string]rkev1.RegistryConfig{
		"test": {
			AuthConfigSecretName: "auth_config_secret_name",
			TLSSecretName:        "tls_secret_name",
			CABundle:             []byte("ca_bundle"),
			InsecureSkipVerify:   true,
		},
	}

	testClusterV2RKEConfigRegistryConfigsInterface = []interface{}{
		map[string]interface{}{
			"hostname":                "test",
			"auth_config_secret_name": "auth_config_secret_name",
			"tls_secret_name":         "tls_secret_name",
			"ca_bundle":               "ca_bundle",
			"insecure":                true,
		},
	}
	testClusterV2RKEConfigRegistryMirrorsConf = map[string]rkev1.Mirror{
		"test": {
			Endpoints: []string{"value1", "value2"},
			Rewrites: map[string]string{
				"rw_one": "one",
				"rw_two": "two",
			},
		},
	}
	testClusterV2RKEConfigRegistryMirrorsInterface = []interface{}{
		map[string]interface{}{
			"hostname":  "test",
			"endpoints": []interface{}{"value1", "value2"},
			"rewrites": map[string]interface{}{
				"rw_one": "one",
				"rw_two": "two",
			},
		},
	}
	testClusterV2RKEConfigRegistryConf = &rkev1.Registry{
		Configs: testClusterV2RKEConfigRegistryConfigsConf,
		Mirrors: testClusterV2RKEConfigRegistryMirrorsConf,
	}
	testClusterV2RKEConfigRegistryInterface = []interface{}{
		map[string]interface{}{
			"configs": testClusterV2RKEConfigRegistryConfigsInterface,
			"mirrors": testClusterV2RKEConfigRegistryMirrorsInterface,
		},
	}
}

func TestFlattenClusterV2RKEConfigRegistryConfigs(t *testing.T) {

	cases := []struct {
		Input          map[string]rkev1.RegistryConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterV2RKEConfigRegistryConfigsConf,
			testClusterV2RKEConfigRegistryConfigsInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterV2RKEConfigRegistryConfigs(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterV2RKEConfigRegistryMirrors(t *testing.T) {

	cases := []struct {
		Input          map[string]rkev1.Mirror
		ExpectedOutput []interface{}
	}{
		{
			testClusterV2RKEConfigRegistryMirrorsConf,
			testClusterV2RKEConfigRegistryMirrorsInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterV2RKEConfigRegistryMirrors(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterV2RKEConfigRegistry(t *testing.T) {

	cases := []struct {
		Input          *rkev1.Registry
		ExpectedOutput []interface{}
	}{
		{
			testClusterV2RKEConfigRegistryConf,
			testClusterV2RKEConfigRegistryInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterV2RKEConfigRegistry(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterV2RKEConfigRegistryConfigs(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput map[string]rkev1.RegistryConfig
	}{
		{
			testClusterV2RKEConfigRegistryConfigsInterface,
			testClusterV2RKEConfigRegistryConfigsConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterV2RKEConfigRegistryConfigs(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterV2RKEConfigRegistryMirrors(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput map[string]rkev1.Mirror
	}{
		{
			testClusterV2RKEConfigRegistryMirrorsInterface,
			testClusterV2RKEConfigRegistryMirrorsConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterV2RKEConfigRegistryMirrors(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterV2RKEConfigRegistry(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *rkev1.Registry
	}{
		{
			testClusterV2RKEConfigRegistryInterface,
			testClusterV2RKEConfigRegistryConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterV2RKEConfigRegistry(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
