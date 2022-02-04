package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testClusterAKSConfigV2NodePoolConf      []managementClient.AKSNodePool
	testClusterAKSConfigV2NodePoolInterface []interface{}
	testClusterAKSConfigV2Conf              *managementClient.AKSClusterConfigSpec
	testClusterAKSConfigV2Interface         []interface{}
)

func init() {
	size := int64(3)
	testClusterAKSConfigV2NodePoolConf = []managementClient.AKSNodePool{
		{
			Name:                newString("name"),
			AvailabilityZones:   &[]string{"zone1", "zone2"},
			Count:               &size,
			EnableAutoScaling:   newTrue(),
			MaxCount:            &size,
			MaxPods:             &size,
			MinCount:            &size,
			Mode:                "test",
			OrchestratorVersion: newString("orchestrator_version"),
			OsDiskSizeGB:        &size,
			OsDiskType:          "os_disk_type",
			OsType:              "os_type",
			VMSize:              "vm_size",
		},
	}
	testClusterAKSConfigV2NodePoolInterface = []interface{}{
		map[string]interface{}{
			"name":                 "name",
			"availability_zones":   []interface{}{"zone1", "zone2"},
			"count":                3,
			"enable_auto_scaling":  true,
			"max_count":            3,
			"max_pods":             3,
			"min_count":            3,
			"mode":                 "test",
			"orchestrator_version": "orchestrator_version",
			"os_disk_size_gb":      3,
			"os_disk_type":         "os_disk_type",
			"os_type":              "os_type",
			"vm_size":              "vm_size",
		},
	}
	testClusterAKSConfigV2Conf = &managementClient.AKSClusterConfigSpec{
		AzureCredentialSecret:      "test",
		AuthBaseURL:                newString("auth_base_url"),
		AuthorizedIPRanges:         &[]string{"type1", "type2"},
		BaseURL:                    newString("base_url"),
		ClusterName:                "name",
		DNSPrefix:                  newString("dns_prefix"),
		HTTPApplicationRouting:     newTrue(),
		Imported:                   false,
		KubernetesVersion:          newString("kubernetes_version"),
		LinuxAdminUsername:         newString("linux_admin_username"),
		LinuxSSHPublicKey:          newString("linux_ssh_public_key"),
		LoadBalancerSKU:            newString("load_balancer_sku"),
		LogAnalyticsWorkspaceGroup: newString("log_analytics_workspace_group"),
		LogAnalyticsWorkspaceName:  newString("log_analytics_workspace_name"),
		Monitoring:                 newTrue(),
		NetworkDNSServiceIP:        newString("network_dns_service_ip"),
		NetworkDockerBridgeCIDR:    newString("network_docker_bridge_cidr"),
		NetworkPlugin:              newString("network_plugin"),
		NetworkPodCIDR:             newString("network_pod_cidr"),
		NetworkPolicy:              newString("network_policy"),
		NetworkServiceCIDR:         newString("network_service_cidr"),
		NodePools:                  testClusterAKSConfigV2NodePoolConf,
		PrivateCluster:             newTrue(),
		ResourceGroup:              "resource_group",
		ResourceLocation:           "resource_location",
		Subnet:                     newString("subnet"),
		Tags: map[string]string{
			"value1": "one",
			"value2": "two",
		},
		VirtualNetwork:              newString("virtual_network"),
		VirtualNetworkResourceGroup: newString("virtual_network_resource_group"),
	}
	testClusterAKSConfigV2Interface = []interface{}{
		map[string]interface{}{
			"cloud_credential_id":           "test",
			"auth_base_url":                 "auth_base_url",
			"authorized_ip_ranges":          []interface{}{"type1", "type2"},
			"base_url":                      "base_url",
			"name":                          "name",
			"dns_prefix":                    "dns_prefix",
			"http_application_routing":      true,
			"imported":                      false,
			"kubernetes_version":            "kubernetes_version",
			"linux_admin_username":          "linux_admin_username",
			"linux_ssh_public_key":          "linux_ssh_public_key",
			"load_balancer_sku":             "load_balancer_sku",
			"log_analytics_workspace_group": "log_analytics_workspace_group",
			"log_analytics_workspace_name":  "log_analytics_workspace_name",
			"monitoring":                    true,
			"network_dns_service_ip":        "network_dns_service_ip",
			"network_docker_bridge_cidr":    "network_docker_bridge_cidr",
			"network_plugin":                "network_plugin",
			"network_pod_cidr":              "network_pod_cidr",
			"network_policy":                "network_policy",
			"network_service_cidr":          "network_service_cidr",
			"node_pools":                    testClusterAKSConfigV2NodePoolInterface,
			"private_cluster":               true,
			"resource_group":                "resource_group",
			"resource_location":             "resource_location",
			"subnet":                        "subnet",
			"tags": map[string]interface{}{
				"value1": "one",
				"value2": "two",
			},
			"virtual_network":                "virtual_network",
			"virtual_network_resource_group": "virtual_network_resource_group",
		},
	}
}

func TestFlattenClusterAKSConfigV2NodePools(t *testing.T) {

	cases := []struct {
		Input          []managementClient.AKSNodePool
		ExpectedOutput []interface{}
	}{
		{
			testClusterAKSConfigV2NodePoolConf,
			testClusterAKSConfigV2NodePoolInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterAKSConfigV2NodePools(tc.Input, tc.ExpectedOutput)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterAKSConfigV2(t *testing.T) {

	cases := []struct {
		Input          *managementClient.AKSClusterConfigSpec
		ExpectedOutput []interface{}
	}{
		{
			testClusterAKSConfigV2Conf,
			testClusterAKSConfigV2Interface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterAKSConfigV2(tc.Input, tc.ExpectedOutput)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterAKSConfigV2NodePools(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []managementClient.AKSNodePool
	}{
		{
			testClusterAKSConfigV2NodePoolInterface,
			testClusterAKSConfigV2NodePoolConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterAKSConfigV2NodePools(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterAKSConfigV2(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.AKSClusterConfigSpec
	}{
		{
			testClusterAKSConfigV2Interface,
			testClusterAKSConfigV2Conf,
		},
	}

	for _, tc := range cases {
		output := expandClusterAKSConfigV2(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
