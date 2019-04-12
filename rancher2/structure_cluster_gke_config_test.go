package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testClusterGKEConfigConf      *managementClient.GoogleKubernetesEngineConfig
	testClusterGKEConfigInterface []interface{}
)

func init() {
	testClusterGKEConfigConf = &managementClient.GoogleKubernetesEngineConfig{
		ClusterIpv4Cidr:                "192.168.1.0/16",
		Credential:                     "credential",
		Description:                    "description",
		DiskSizeGb:                     16,
		EnableAlphaFeature:             true,
		EnableHTTPLoadBalancing:        newTrue(),
		EnableHorizontalPodAutoscaling: newTrue(),
		EnableKubernetesDashboard:      true,
		EnableLegacyAbac:               true,
		EnableNetworkPolicyConfig:      newTrue(),
		EnableStackdriverLogging:       newTrue(),
		EnableStackdriverMonitoring:    newTrue(),
		ImageType:                      "image",
		Labels: map[string]string{
			"label1": "value1",
			"label2": "value2",
		},
		Locations:         []string{"location1", "location2"},
		MachineType:       "machine",
		MaintenanceWindow: "maintenance",
		MasterVersion:     "version",
		Network:           "network",
		NodeCount:         3,
		NodeVersion:       "node_version",
		ProjectID:         "project-test",
		SubNetwork:        "subnetwork",
		Zone:              "zone",
	}
	testClusterGKEConfigInterface = []interface{}{
		map[string]interface{}{
			"cluster_ipv4_cidr":                 "192.168.1.0/16",
			"credential":                        "credential",
			"description":                       "description",
			"disk_size_gb":                      16,
			"enable_alpha_feature":              true,
			"enable_http_load_balancing":        true,
			"enable_horizontal_pod_autoscaling": true,
			"enable_kubernetes_dashboard":       true,
			"enable_legacy_abac":                true,
			"enable_network_policy_config":      true,
			"enable_stackdriver_logging":        true,
			"enable_stackdriver_monitoring":     true,
			"image_type":                        "image",
			"labels": map[string]interface{}{
				"label1": "value1",
				"label2": "value2",
			},
			"locations":          []interface{}{"location1", "location2"},
			"machine_type":       "machine",
			"maintenance_window": "maintenance",
			"master_version":     "version",
			"network":            "network",
			"node_count":         3,
			"node_version":       "node_version",
			"project_id":         "project-test",
			"sub_network":        "subnetwork",
			"zone":               "zone",
		},
	}
}

func TestFlattenClusterGKEConfig(t *testing.T) {

	cases := []struct {
		Input          *managementClient.GoogleKubernetesEngineConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterGKEConfigConf,
			testClusterGKEConfigInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterGKEConfig(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterGKEConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.GoogleKubernetesEngineConfig
	}{
		{
			testClusterGKEConfigInterface,
			testClusterGKEConfigConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterGKEConfig(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
