package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testClusterRKEConfigServicesKubeletConf      *managementClient.KubeletService
	testClusterRKEConfigServicesKubeletInterface []interface{}
)

func init() {
	testClusterRKEConfigServicesKubeletConf = &managementClient.KubeletService{
		ClusterDNSServer: "dns.hostname.test",
		ClusterDomain:    "terraform.test",
		ExtraArgs: map[string]string{
			"arg_one": "one",
			"arg_two": "two",
		},
		ExtraBinds:                 []string{"bind_one", "bind_two"},
		ExtraEnv:                   []string{"env_one", "env_two"},
		FailSwapOn:                 true,
		GenerateServingCertificate: true,
		Image:                      "image",
		InfraContainerImage:        "infra_image",
	}
	testClusterRKEConfigServicesKubeletInterface = []interface{}{
		map[string]interface{}{
			"cluster_dns_server": "dns.hostname.test",
			"cluster_domain":     "terraform.test",
			"extra_args": map[string]interface{}{
				"arg_one": "one",
				"arg_two": "two",
			},
			"extra_binds":                  []interface{}{"bind_one", "bind_two"},
			"extra_env":                    []interface{}{"env_one", "env_two"},
			"fail_swap_on":                 true,
			"generate_serving_certificate": true,
			"image":                        "image",
			"infra_container_image":        "infra_image",
		},
	}
}

func TestFlattenClusterRKEConfigServicesKubelet(t *testing.T) {

	cases := []struct {
		Input          *managementClient.KubeletService
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigServicesKubeletConf,
			testClusterRKEConfigServicesKubeletInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigServicesKubelet(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigServicesKubelet(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.KubeletService
	}{
		{
			testClusterRKEConfigServicesKubeletInterface,
			testClusterRKEConfigServicesKubeletConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigServicesKubelet(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
