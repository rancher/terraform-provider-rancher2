package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testClusterRKEConfigServicesKubeproxyConf           *managementClient.KubeproxyService
	testClusterRKEConfigServicesKubeproxyInterface      []interface{}
	testClusterRKEConfigServicesKubeletConf             *managementClient.KubeletService
	testClusterRKEConfigServicesKubeletInterface        []interface{}
	testClusterRKEConfigServicesKubeControllerConf      *managementClient.KubeControllerService
	testClusterRKEConfigServicesKubeControllerInterface []interface{}
	testClusterRKEConfigServicesKubeAPIConf             *managementClient.KubeAPIService
	testClusterRKEConfigServicesKubeAPIInterface        []interface{}
	testClusterRKEConfigServicesETCDConf                *managementClient.ETCDService
	testClusterRKEConfigServicesETCDInterface           []interface{}
	testClusterRKEConfigServicesConf                    *managementClient.RKEConfigServices
	testClusterRKEConfigServicesInterface               []interface{}
)

func init() {
	testClusterRKEConfigServicesKubeproxyConf = &managementClient.KubeproxyService{
		ExtraArgs: map[string]string{
			"arg_one": "one",
			"arg_two": "two",
		},
		ExtraBinds: []string{"bind_one", "bind_two"},
		ExtraEnv:   []string{"env_one", "env_two"},
		Image:      "image",
	}
	testClusterRKEConfigServicesKubeproxyInterface = []interface{}{
		map[string]interface{}{
			"extra_args": map[string]interface{}{
				"arg_one": "one",
				"arg_two": "two",
			},
			"extra_binds": []interface{}{"bind_one", "bind_two"},
			"extra_env":   []interface{}{"env_one", "env_two"},
			"image":       "image",
		},
	}
	testClusterRKEConfigServicesKubeletConf = &managementClient.KubeletService{
		ClusterDNSServer: "dns.hostname.test",
		ClusterDomain:    "terraform.test",
		ExtraArgs: map[string]string{
			"arg_one": "one",
			"arg_two": "two",
		},
		ExtraBinds:          []string{"bind_one", "bind_two"},
		ExtraEnv:            []string{"env_one", "env_two"},
		FailSwapOn:          true,
		Image:               "image",
		InfraContainerImage: "infra_image",
	}
	testClusterRKEConfigServicesKubeletInterface = []interface{}{
		map[string]interface{}{
			"cluster_dns_server": "dns.hostname.test",
			"cluster_domain":     "terraform.test",
			"extra_args": map[string]interface{}{
				"arg_one": "one",
				"arg_two": "two",
			},
			"extra_binds":           []interface{}{"bind_one", "bind_two"},
			"extra_env":             []interface{}{"env_one", "env_two"},
			"fail_swap_on":          true,
			"image":                 "image",
			"infra_container_image": "infra_image",
		},
	}
	testClusterRKEConfigServicesKubeControllerConf = &managementClient.KubeControllerService{
		ClusterCIDR: "10.42.0.0/16",
		ExtraArgs: map[string]string{
			"arg_one": "one",
			"arg_two": "two",
		},
		ExtraBinds:            []string{"bind_one", "bind_two"},
		ExtraEnv:              []string{"env_one", "env_two"},
		Image:                 "image",
		ServiceClusterIPRange: "10.43.0.0/16",
	}
	testClusterRKEConfigServicesKubeControllerInterface = []interface{}{
		map[string]interface{}{
			"cluster_cidr": "10.42.0.0/16",
			"extra_args": map[string]interface{}{
				"arg_one": "one",
				"arg_two": "two",
			},
			"extra_binds":              []interface{}{"bind_one", "bind_two"},
			"extra_env":                []interface{}{"env_one", "env_two"},
			"image":                    "image",
			"service_cluster_ip_range": "10.43.0.0/16",
		},
	}
	testClusterRKEConfigServicesKubeAPIConf = &managementClient.KubeAPIService{
		ExtraArgs: map[string]string{
			"arg_one": "one",
			"arg_two": "two",
		},
		ExtraBinds:            []string{"bind_one", "bind_two"},
		ExtraEnv:              []string{"env_one", "env_two"},
		Image:                 "image",
		PodSecurityPolicy:     true,
		ServiceClusterIPRange: "10.43.0.0/16",
		ServiceNodePortRange:  "30000-32000",
	}
	testClusterRKEConfigServicesKubeAPIInterface = []interface{}{
		map[string]interface{}{
			"extra_args": map[string]interface{}{
				"arg_one": "one",
				"arg_two": "two",
			},
			"extra_binds":              []interface{}{"bind_one", "bind_two"},
			"extra_env":                []interface{}{"env_one", "env_two"},
			"image":                    "image",
			"pod_security_policy":      true,
			"service_cluster_ip_range": "10.43.0.0/16",
			"service_node_port_range":  "30000-32000",
		},
	}
	testClusterRKEConfigServicesETCDConf = &managementClient.ETCDService{
		CACert:       "XXXXXXXX",
		Cert:         "YYYYYYYY",
		Creation:     "creation",
		ExternalURLs: []string{"url_one", "url_two"},
		ExtraArgs: map[string]string{
			"arg_one": "one",
			"arg_two": "two",
		},
		ExtraBinds: []string{"bind_one", "bind_two"},
		ExtraEnv:   []string{"env_one", "env_two"},
		Image:      "image",
		Key:        "ZZZZZZZZ",
		Path:       "/etcd",
		Retention:  "6h",
		Snapshot:   true,
	}
	testClusterRKEConfigServicesETCDInterface = []interface{}{
		map[string]interface{}{
			"ca_cert":       "XXXXXXXX",
			"cert":          "YYYYYYYY",
			"creation":      "creation",
			"external_urls": []interface{}{"url_one", "url_two"},
			"extra_args": map[string]interface{}{
				"arg_one": "one",
				"arg_two": "two",
			},
			"extra_binds": []interface{}{"bind_one", "bind_two"},
			"extra_env":   []interface{}{"env_one", "env_two"},
			"image":       "image",
			"key":         "ZZZZZZZZ",
			"path":        "/etcd",
			"retention":   "6h",
			"snapshot":    true,
		},
	}
	testClusterRKEConfigServicesConf = &managementClient.RKEConfigServices{
		Etcd:           testClusterRKEConfigServicesETCDConf,
		KubeAPI:        testClusterRKEConfigServicesKubeAPIConf,
		KubeController: testClusterRKEConfigServicesKubeControllerConf,
		Kubelet:        testClusterRKEConfigServicesKubeletConf,
		Kubeproxy:      testClusterRKEConfigServicesKubeproxyConf,
	}
	testClusterRKEConfigServicesInterface = []interface{}{
		map[string]interface{}{
			"etcd":            testClusterRKEConfigServicesETCDInterface,
			"kube_api":        testClusterRKEConfigServicesKubeAPIInterface,
			"kube_controller": testClusterRKEConfigServicesKubeControllerInterface,
			"kubelet":         testClusterRKEConfigServicesKubeletInterface,
			"kubeproxy":       testClusterRKEConfigServicesKubeproxyInterface,
		},
	}
}

func TestFlattenClusterRKEConfigServicesKubeproxy(t *testing.T) {

	cases := []struct {
		Input          *managementClient.KubeproxyService
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigServicesKubeproxyConf,
			testClusterRKEConfigServicesKubeproxyInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigServicesKubeproxy(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
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

func TestFlattenClusterRKEConfigServicesKubeController(t *testing.T) {

	cases := []struct {
		Input          *managementClient.KubeControllerService
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigServicesKubeControllerConf,
			testClusterRKEConfigServicesKubeControllerInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigServicesKubeController(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterRKEConfigServicesKubeAPI(t *testing.T) {

	cases := []struct {
		Input          *managementClient.KubeAPIService
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigServicesKubeAPIConf,
			testClusterRKEConfigServicesKubeAPIInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigServicesKubeAPI(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterRKEConfigServicesEtcd(t *testing.T) {

	cases := []struct {
		Input          *managementClient.ETCDService
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigServicesETCDConf,
			testClusterRKEConfigServicesETCDInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigServicesEtcd(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
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
		output, err := flattenClusterRKEConfigServices(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigServicesKubeproxy(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.KubeproxyService
	}{
		{
			testClusterRKEConfigServicesKubeproxyInterface,
			testClusterRKEConfigServicesKubeproxyConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigServicesKubeproxy(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
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

func TestExpandClusterRKEConfigServicesKubeController(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.KubeControllerService
	}{
		{
			testClusterRKEConfigServicesKubeControllerInterface,
			testClusterRKEConfigServicesKubeControllerConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigServicesKubeController(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigServicesKubeAPI(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.KubeAPIService
	}{
		{
			testClusterRKEConfigServicesKubeAPIInterface,
			testClusterRKEConfigServicesKubeAPIConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigServicesKubeAPI(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigServicesEtcd(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.ETCDService
	}{
		{
			testClusterRKEConfigServicesETCDInterface,
			testClusterRKEConfigServicesETCDConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigServicesEtcd(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
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
