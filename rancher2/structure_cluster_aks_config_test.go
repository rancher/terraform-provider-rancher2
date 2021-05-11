package rancher2

import (
	"reflect"
	"testing"
)

var (
	testClusterAKSConfigConf      *AzureKubernetesServiceConfig
	testClusterAKSConfigInterface []interface{}
)

func init() {
	testClusterAKSConfigConf = &AzureKubernetesServiceConfig{
		AADClientAppID:                     "add_client_app_id",
		AADServerAppID:                     "add_server_app_id",
		AADServerAppSecret:                 "aad_server_app_secret",
		AADTenantID:                        "aad_tenant_id",
		AdminUsername:                      "admin",
		AgentDNSPrefix:                     "dns",
		AgentPoolName:                      "agent",
		AgentStorageProfile:                "agent_storage_profile",
		AgentVMSize:                        "size",
		AuthBaseURL:                        "auth_base_url",
		BaseURL:                            "url",
		ClientID:                           "client_id",
		ClientSecret:                       "client_secret",
		Count:                              3,
		DisplayName:                        "test",
		DriverName:                         clusterDriverAKS,
		DNSServiceIP:                       "dns_ip",
		DockerBridgeCIDR:                   "192.168.1.0/16",
		EnableHTTPApplicationRouting:       true,
		EnableMonitoring:                   newTrue(),
		KubernetesVersion:                  "version",
		LoadBalancerSku:                    clusterAKSLoadBalancerSkuStandard,
		Location:                           "location",
		LogAnalyticsWorkspace:              "log_analytics_workspace",
		LogAnalyticsWorkspaceResourceGroup: "log_analytics_workspace_resource_group",
		MasterDNSPrefix:                    "dns_prefix",
		Name:                               "test",
		MaxPods:                            100,
		NetworkPlugin:                      "network_plugin",
		NetworkPolicy:                      "network_policy",
		AgentOsdiskSizeGB:                  16,
		PodCIDR:                            "pod_cidr",
		ResourceGroup:                      "resource_group",
		SSHPublicKeyContents:               "key",
		ServiceCIDR:                        "service_cidr",
		Subnet:                             "subnet",
		SubscriptionID:                     "subscription_id",
		Tags: []string{
			"tag1=value1",
			"tag2=value2",
		},
		TenantID:                    "tenant_id",
		VirtualNetwork:              "virtual_network",
		VirtualNetworkResourceGroup: "network_resource_group",
	}
	testClusterAKSConfigInterface = []interface{}{
		map[string]interface{}{
			"add_client_app_id":                      "add_client_app_id",
			"add_server_app_id":                      "add_server_app_id",
			"aad_server_app_secret":                  "aad_server_app_secret",
			"aad_tenant_id":                          "aad_tenant_id",
			"admin_username":                         "admin",
			"agent_dns_prefix":                       "dns",
			"agent_os_disk_size":                     16,
			"agent_pool_name":                        "agent",
			"agent_storage_profile":                  "agent_storage_profile",
			"agent_vm_size":                          "size",
			"auth_base_url":                          "auth_base_url",
			"base_url":                               "url",
			"client_id":                              "client_id",
			"client_secret":                          "client_secret",
			"count":                                  3,
			"dns_service_ip":                         "dns_ip",
			"docker_bridge_cidr":                     "192.168.1.0/16",
			"enable_http_application_routing":        true,
			"enable_monitoring":                      true,
			"kubernetes_version":                     "version",
			"load_balancer_sku":                      clusterAKSLoadBalancerSkuStandard,
			"location":                               "location",
			"log_analytics_workspace":                "log_analytics_workspace",
			"log_analytics_workspace_resource_group": "log_analytics_workspace_resource_group",
			"master_dns_prefix":                      "dns_prefix",
			"max_pods":                               100,
			"network_plugin":                         "network_plugin",
			"network_policy":                         "network_policy",
			"pod_cidr":                               "pod_cidr",
			"resource_group":                         "resource_group",
			"ssh_public_key_contents":                "key",
			"service_cidr":                           "service_cidr",
			"subnet":                                 "subnet",
			"subscription_id":                        "subscription_id",
			"tags": []interface{}{
				"tag1=value1",
				"tag2=value2",
			},
			"tenant_id":                      "tenant_id",
			"virtual_network":                "virtual_network",
			"virtual_network_resource_group": "network_resource_group",
		},
	}
}

func TestFlattenClusterAKSConfig(t *testing.T) {

	cases := []struct {
		Input          *AzureKubernetesServiceConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterAKSConfigConf,
			testClusterAKSConfigInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterAKSConfig(tc.Input, testClusterAKSConfigInterface)
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
		ExpectedOutput *AzureKubernetesServiceConfig
	}{
		{
			testClusterAKSConfigInterface,
			testClusterAKSConfigConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterAKSConfig(tc.Input, "test")
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
