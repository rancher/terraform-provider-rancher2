package rancher2

import (
	"encoding/json"
	"reflect"
	"testing"
)

var (
	testClusterAKSConfigConf      *AzureKubernetesServiceConfig
	testClusterAKSConfigInterface []interface{}

	testClusterAKSConfigLegacyConf      *AzureKubernetesServiceConfig
	testClusterAKSConfigLegacyInterface []interface{}
)

func init() {
	nodePool := AzureKubernetesServiceNodePool{
		BaseNodePool: BaseNodePool{
			Name:             "agent",
			AdditionalLabels: map[string]string{"label": "value"},
			AdditionalTaints: []K8sTaint{
				{Effect: "NoSchedule", Key: "taint_key", Value: "taint_value"},
			},
		},
		AvailabilityZones: []string{"az1", "az2"},
		CreatePoolPerZone: true,
		EnableAutoScaling: newTrue(),
		MaxCount:          33,
		MaxPods:           100,
		MinCount:          11,
		OsDiskSizeGB:      16,
		Type:              "agent_pool_type",
		Version:           "agent_version",
		VMSize:            "size",
	}

	nodePoolBytes, _ := json.Marshal(nodePool)

	testClusterAKSConfigConf = &AzureKubernetesServiceConfig{
		AADClientAppID:                     "add_client_app_id",
		AADServerAppID:                     "add_server_app_id",
		AADServerAppSecret:                 "aad_server_app_secret",
		AADTenantID:                        "aad_tenant_id",
		AdminUsername:                      "admin",
		AgentDNSPrefix:                     "dns",
		AgentStorageProfile:                "agent_storage_profile",
		AuthBaseURL:                        "auth_base_url",
		BaseURL:                            "url",
		ClientID:                           "client_id",
		ClientSecret:                       "client_secret",
		Count:                              3,
		DisplayName:                        "test",
		DNSServiceIP:                       "dns_ip",
		DockerBridgeCIDR:                   "192.168.1.0/16",
		EnableHTTPApplicationRouting:       true,
		EnableMonitoring:                   newTrue(),
		KubernetesVersion:                  "version",
		LoadBalancerSku:                    "load_balancer_sku",
		Location:                           "location",
		LogAnalyticsWorkspace:              "log_analytics_workspace",
		LogAnalyticsWorkspaceResourceGroup: "log_analytics_workspace_resource_group",
		MasterDNSPrefix:                    "dns_prefix",
		Name:                               "test",
		NodePools:                          []string{string(nodePoolBytes)},
		NetworkPlugin:                      "network_plugin",
		NetworkPolicy:                      "network_policy",
		PodCIDR:                            "pod_cidr",
		ResourceGroup:                      "resource_group",
		SSHPublicKeyContents:               "key",
		ServiceCIDR:                        "service_cidr",
		Subnet:                             "subnet",
		SubscriptionID:                     "subscription_id",
		Tags: map[string]string{
			"tag1": "value1",
			"tag2": "value2",
		},
		TenantID:                    "tenant_id",
		VirtualNetwork:              "virtual_network",
		VirtualNetworkResourceGroup: "network_resource_group",
	}

	testClusterAKSConfigLegacyConf = &AzureKubernetesServiceConfig{
		AADClientAppID:                     "add_client_app_id",
		AADServerAppID:                     "add_server_app_id",
		AADServerAppSecret:                 "aad_server_app_secret",
		AADTenantID:                        "aad_tenant_id",
		AdminUsername:                      "admin",
		AgentDNSPrefix:                     "dns",
		AgentPoolName:                      "agent",
		AgentPoolType:                      "agent_pool_type",
		AgentStorageProfile:                "agent_storage_profile",
		AgentVMSize:                        "size",
		AuthBaseURL:                        "auth_base_url",
		AvailabilityZones:                  []string{"az1", "az2"},
		BaseURL:                            "url",
		ClientID:                           "client_id",
		ClientSecret:                       "client_secret",
		Count:                              3,
		DisplayName:                        "test",
		DNSServiceIP:                       "dns_ip",
		DockerBridgeCIDR:                   "192.168.1.0/16",
		EnableAutoScaling:                  newTrue(),
		EnableHTTPApplicationRouting:       true,
		EnableMonitoring:                   newTrue(),
		KubernetesVersion:                  "version",
		LoadBalancerSku:                    "load_balancer_sku",
		Location:                           "location",
		LogAnalyticsWorkspace:              "log_analytics_workspace",
		LogAnalyticsWorkspaceResourceGroup: "log_analytics_workspace_resource_group",
		MasterDNSPrefix:                    "dns_prefix",
		MaxCount:                           33,
		MinCount:                           11,
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
		Tags: map[string]string{
			"tag1": "value1",
			"tag2": "value2",
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
			"agent_storage_profile":                  "agent_storage_profile",
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
			"load_balancer_sku":                      "load_balancer_sku",
			"location":                               "location",
			"log_analytics_workspace":                "log_analytics_workspace",
			"log_analytics_workspace_resource_group": "log_analytics_workspace_resource_group",
			"master_dns_prefix":                      "dns_prefix",
			"network_plugin":                         "network_plugin",
			"network_policy":                         "network_policy",
			"node_pools": []interface{}{
				map[string]interface{}{
					"add_default_label": false,
					"add_default_taint": false,
					"additional_labels": map[string]interface{}{"label": "value"},
					"additional_taints": []interface{}{
						map[string]interface{}{
							"effect": "NoSchedule",
							"key":    "taint_key",
							"value":  "taint_value",
						},
					},
					"availability_zones":   []interface{}{"az1", "az2"},
					"create_pool_per_zone": true,
					"enable_auto_scaling":  true,
					"max_count":            33,
					"max_pods":             100,
					"min_count":            11,
					"name":                 "agent",
					"os_disk_size":         16,
					"version":              "agent_version",
					"vm_size":              "size",
					"type":                 "agent_pool_type",
				},
			},
			"pod_cidr":                "pod_cidr",
			"resource_group":          "resource_group",
			"ssh_public_key_contents": "key",
			"service_cidr":            "service_cidr",
			"subnet":                  "subnet",
			"subscription_id":         "subscription_id",
			"tags": map[string]interface{}{
				"tag1": "value1",
				"tag2": "value2",
			},
			"tenant_id":                      "tenant_id",
			"virtual_network":                "virtual_network",
			"virtual_network_resource_group": "network_resource_group",
		},
	}

	testClusterAKSConfigLegacyInterface = []interface{}{
		map[string]interface{}{
			"add_client_app_id":                      "add_client_app_id",
			"add_server_app_id":                      "add_server_app_id",
			"aad_server_app_secret":                  "aad_server_app_secret",
			"aad_tenant_id":                          "aad_tenant_id",
			"admin_username":                         "admin",
			"agent_dns_prefix":                       "dns",
			"agent_storage_profile":                  "agent_storage_profile",
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
			"load_balancer_sku":                      "load_balancer_sku",
			"location":                               "location",
			"log_analytics_workspace":                "log_analytics_workspace",
			"log_analytics_workspace_resource_group": "log_analytics_workspace_resource_group",
			"master_dns_prefix":                      "dns_prefix",
			"network_plugin":                         "network_plugin",
			"network_policy":                         "network_policy",
			"node_pools": []interface{}{
				map[string]interface{}{
					"add_default_label":    false,
					"add_default_taint":    false,
					"availability_zones":   []interface{}{"az1", "az2"},
					"create_pool_per_zone": true,
					"enable_auto_scaling":  true,
					"max_count":            33,
					"max_pods":             100,
					"min_count":            11,
					"name":                 "agent",
					"os_disk_size":         16,
					"vm_size":              "size",
					"type":                 "agent_pool_type",
				},
			},
			"pod_cidr":                "pod_cidr",
			"resource_group":          "resource_group",
			"ssh_public_key_contents": "key",
			"service_cidr":            "service_cidr",
			"subnet":                  "subnet",
			"subscription_id":         "subscription_id",
			"tags": map[string]interface{}{
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

	cases := map[string]struct {
		Input          *AzureKubernetesServiceConfig
		ExpectedOutput []interface{}
	}{
		"AKSCluster": {
			testClusterAKSConfigConf,
			testClusterAKSConfigInterface,
		},
		"AKSLegacyCluster": {
			testClusterAKSConfigLegacyConf,
			testClusterAKSConfigLegacyInterface,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			output, err := flattenClusterAKSConfig(tc.Input)
			if err != nil {
				t.Fatalf("[ERROR] on flattener: %#v", err)
			}
			if !reflect.DeepEqual(output, tc.ExpectedOutput) {
				t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
					tc.ExpectedOutput, output)
			}
		})
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
		output, err := expandClusterAKSConfig(&AzureKubernetesServiceConfig{}, tc.Input, "test")
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
