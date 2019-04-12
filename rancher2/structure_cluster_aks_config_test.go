package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/types/client/management/v3"
)

var (
	testClusterAKSConfigConf      *managementClient.AzureKubernetesServiceConfig
	testClusterAKSConfigInterface []interface{}
)

func init() {
	testClusterAKSConfigConf = &managementClient.AzureKubernetesServiceConfig{
		AdminUsername:        "admin",
		AgentDNSPrefix:       "dns",
		AgentPoolName:        "agent",
		AgentVMSize:          "size",
		BaseURL:              "url",
		ClientID:             "client_id",
		ClientSecret:         "client_secret",
		Count:                3,
		DNSServiceIP:         "dns_ip",
		DockerBridgeCIDR:     "192.168.1.0/16",
		KubernetesVersion:    "version",
		Location:             "location",
		MasterDNSPrefix:      "dns_prefix",
		OsDiskSizeGB:         16,
		ResourceGroup:        "resource_group",
		SSHPublicKeyContents: "key",
		ServiceCIDR:          "service_cidr",
		Subnet:               "subnet",
		SubscriptionID:       "subscription_id",
		Tag: map[string]string{
			"tag1": "value1",
			"tag2": "value2",
		},
		TenantID:                    "tenant_id",
		VirtualNetwork:              "virtual_network",
		VirtualNetworkResourceGroup: "network_resource_group",
	}
	testClusterAKSConfigInterface = []interface{}{
		map[string]interface{}{
			"admin_username":          "admin",
			"agent_dns_prefix":        "dns",
			"agent_pool_name":         "agent",
			"agent_vm_size":           "size",
			"base_url":                "url",
			"client_id":               "client_id",
			"client_secret":           "client_secret",
			"count":                   3,
			"dns_service_ip":          "dns_ip",
			"docker_bridge_cidr":      "192.168.1.0/16",
			"kubernetes_version":      "version",
			"location":                "location",
			"master_dns_prefix":       "dns_prefix",
			"os_disk_size_gb":         16,
			"resource_group":          "resource_group",
			"ssh_public_key_contents": "key",
			"service_cidr":            "service_cidr",
			"subnet":                  "subnet",
			"subscription_id":         "subscription_id",
			"tag": map[string]interface{}{
				"tag1": "value1",
				"tag2": "value2",
			},
			"tenant_id":                      "tenant_id",
			"virtual_network":                "virtual_network",
			"virtual_network_resource_group": "network_resource_group",
		},
	}
}

func TestFlattenClusterAKSConfig(t *testing.T) {

	cases := []struct {
		Input          *managementClient.AzureKubernetesServiceConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterAKSConfigConf,
			testClusterAKSConfigInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterAKSConfig(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterAKSConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.AzureKubernetesServiceConfig
	}{
		{
			testClusterAKSConfigInterface,
			testClusterAKSConfigConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterAKSConfig(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
