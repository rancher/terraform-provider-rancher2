package rancher2

import (
	"encoding/json"
	"reflect"
	"testing"
)

var (
	testClusterGKEConfigConf       *GoogleKubernetesEngineConfig
	testClusterGKEConfigLegacyConf *GoogleKubernetesEngineConfig
	testClusterGKEConfigInterface  []interface{}
)

func init() {
	nodePool := GoogleKubernetesEngineNodePool{
		BaseNodePool: BaseNodePool{
			AddDefaultLabel: false,
			AddDefaultTaint: false,
			AdditionalLabels: map[string]string{
				"label1": "value1",
				"label2": "value2",
			},
			AdditionalTaints: []K8sTaint{
				{Effect: "NoSchedule", Key: "taint", Value: "value"},
			},
			Name: "node_pool",
		},
		DiskSizeGb:         16,
		DiskType:           "disk_type",
		EnableAutoRepair:   true,
		EnableAutoUpgrade:  true,
		ImageType:          "image",
		LocalSsdCount:      1,
		MachineType:        "machine",
		MaximumNodeCount:   3,
		MinimumNodeCount:   1,
		MinimumCpuPlatform: "Intel Skylake",
		OauthScopes:        []string{"scope1", "scope2"},
		Preemptible:        true,
		ServiceAccount:     "service_account",
		Version:            "node_version",
	}

	nodePoolBytes, _ := json.Marshal(nodePool)

	testClusterGKEConfigConf = &GoogleKubernetesEngineConfig{
		ClusterIpv4Cidr:                    "192.168.1.0/16",
		Credential:                         "credential",
		Description:                        "description",
		DisplayName:                        "test",
		EnableAlphaFeature:                 true,
		EnableHTTPLoadBalancing:            newTrue(),
		EnableHorizontalPodAutoscaling:     newTrue(),
		EnableKubernetesDashboard:          true,
		EnableLegacyAbac:                   true,
		EnableMasterAuthorizedNetwork:      true,
		EnableNetworkPolicyConfig:          newTrue(),
		EnableNodepoolAutoscaling:          true,
		EnablePrivateEndpoint:              true,
		EnablePrivateNodes:                 true,
		EnableStackdriverLogging:           newTrue(),
		EnableStackdriverMonitoring:        newTrue(),
		IPPolicyClusterIpv4CidrBlock:       "ip_policy_cluster_ipv4_cidr_block",
		IPPolicyClusterSecondaryRangeName:  "ip_policy_cluster_secondary_range_name",
		IPPolicyCreateSubnetwork:           true,
		IPPolicyNodeIpv4CidrBlock:          "ip_policy_node_ipv4_cidr_block",
		IPPolicyServicesIpv4CidrBlock:      "ip_policy_services_ipv4_cidr_block",
		IPPolicyServicesSecondaryRangeName: "ip_policy_services_secondary_range_name",
		IPPolicySubnetworkName:             "ip_policy_subnetwork_name",
		IssueClientCertificate:             true,
		KubernetesDashboard:                true,
		Locations:                          []string{"location1", "location2"},
		MaintenanceWindow:                  "maintenance",
		MasterAuthorizedNetworkCidrBlocks:  []string{"master1", "master2"},
		MasterIpv4CidrBlock:                "master_ipv4_cidr_block",
		MasterVersion:                      "version",
		Name:                               "test",
		Network:                            "network",
		NodePools:                          []string{string(nodePoolBytes)},
		ProjectID:                          "project-test",
		ResourceLabels: map[string]string{
			"Rlabel1": "value1",
			"Rlabel2": "value2",
		},
		SubNetwork:               "subnetwork",
		UseIPAliases:             true,
		Zone:                     "zone",
		DefaultMaxPodsConstraint: 32,
		Region:                   "region",
	}

	testClusterGKEConfigLegacyConf = &GoogleKubernetesEngineConfig{
		ClusterIpv4Cidr:                    "192.168.1.0/16",
		Credential:                         "credential",
		Description:                        "description",
		DiskType:                           "disk_type",
		DiskSizeGb:                         16,
		DisplayName:                        "test",
		EnableAlphaFeature:                 true,
		EnableAutoRepair:                   true,
		EnableAutoUpgrade:                  true,
		EnableHTTPLoadBalancing:            newTrue(),
		EnableHorizontalPodAutoscaling:     newTrue(),
		EnableKubernetesDashboard:          true,
		EnableLegacyAbac:                   true,
		EnableMasterAuthorizedNetwork:      true,
		EnableNetworkPolicyConfig:          newTrue(),
		EnableNodepoolAutoscaling:          true,
		EnablePrivateEndpoint:              true,
		EnablePrivateNodes:                 true,
		EnableStackdriverLogging:           newTrue(),
		EnableStackdriverMonitoring:        newTrue(),
		ImageType:                          "image",
		IPPolicyClusterIpv4CidrBlock:       "ip_policy_cluster_ipv4_cidr_block",
		IPPolicyClusterSecondaryRangeName:  "ip_policy_cluster_secondary_range_name",
		IPPolicyCreateSubnetwork:           true,
		IPPolicyNodeIpv4CidrBlock:          "ip_policy_node_ipv4_cidr_block",
		IPPolicyServicesIpv4CidrBlock:      "ip_policy_services_ipv4_cidr_block",
		IPPolicyServicesSecondaryRangeName: "ip_policy_services_secondary_range_name",
		IPPolicySubnetworkName:             "ip_policy_subnetwork_name",
		IssueClientCertificate:             true,
		KubernetesDashboard:                true,
		Labels: map[string]string{
			"label1": "value1",
			"label2": "value2",
		},
		LocalSsdCount:                     1,
		Locations:                         []string{"location1", "location2"},
		MachineType:                       "machine",
		MaintenanceWindow:                 "maintenance",
		MasterAuthorizedNetworkCidrBlocks: []string{"master1", "master2"},
		MasterIpv4CidrBlock:               "master_ipv4_cidr_block",
		MasterVersion:                     "version",
		MaxNodeCount:                      3,
		MinNodeCount:                      1,
		MinCpuPlatform:                    "Intel Skylake",
		Name:                              "test",
		Network:                           "network",
		NodePool:                          "node_pool",
		NodeVersion:                       "node_version",
		OauthScopes:                       []string{"scope1", "scope2"},
		Preemptible:                       true,
		ProjectID:                         "project-test",
		ResourceLabels: map[string]string{
			"Rlabel1": "value1",
			"Rlabel2": "value2",
		},
		ServiceAccount:           "service_account",
		SubNetwork:               "subnetwork",
		UseIPAliases:             true,
		Taints:                   []string{"NoSchedule:taint=value"},
		Zone:                     "zone",
		DefaultMaxPodsConstraint: 32,
		Region:                   "region",
	}
	testClusterGKEConfigInterface = []interface{}{
		map[string]interface{}{
			"cluster_ipv4_cidr":                       "192.168.1.0/16",
			"credential":                              "credential",
			"description":                             "description",
			"enable_alpha_feature":                    true,
			"enable_http_load_balancing":              true,
			"enable_horizontal_pod_autoscaling":       true,
			"enable_kubernetes_dashboard":             true,
			"enable_legacy_abac":                      true,
			"enable_master_authorized_network":        true,
			"enable_network_policy_config":            true,
			"enable_nodepool_autoscaling":             true,
			"enable_private_endpoint":                 true,
			"enable_private_nodes":                    true,
			"enable_stackdriver_logging":              true,
			"enable_stackdriver_monitoring":           true,
			"ip_policy_cluster_ipv4_cidr_block":       "ip_policy_cluster_ipv4_cidr_block",
			"ip_policy_cluster_secondary_range_name":  "ip_policy_cluster_secondary_range_name",
			"ip_policy_create_subnetwork":             true,
			"ip_policy_node_ipv4_cidr_block":          "ip_policy_node_ipv4_cidr_block",
			"ip_policy_services_ipv4_cidr_block":      "ip_policy_services_ipv4_cidr_block",
			"ip_policy_services_secondary_range_name": "ip_policy_services_secondary_range_name",
			"ip_policy_subnetwork_name":               "ip_policy_subnetwork_name",
			"issue_client_certificate":                true,
			"kubernetes_dashboard":                    true,
			"locations":                               []interface{}{"location1", "location2"},
			"maintenance_window":                      "maintenance",
			"master_authorized_network_cidr_blocks":   []interface{}{"master1", "master2"},
			"master_ipv4_cidr_block":                  "master_ipv4_cidr_block",
			"master_version":                          "version",
			"network":                                 "network",
			"node_pools": []interface{}{
				map[string]interface{}{
					"add_default_label": false,
					"add_default_taint": false,
					"additional_labels": map[string]interface{}{
						"label1": "value1",
						"label2": "value2",
					},
					"additional_taints": []interface{}{
						map[string]interface{}{"effect": "NoSchedule", "key": "taint", "value": "value"},
					},
					"disk_size_gb":        16,
					"disk_type":           "disk_type",
					"enable_auto_repair":  true,
					"enable_auto_upgrade": true,
					"image_type":          "image",
					"local_ssd_count":     1,
					"machine_type":        "machine",
					"max_node_count":      3,
					"min_node_count":      1,
					"min_cpu_platform":    "Intel Skylake",
					"name":                "node_pool",
					"oauth_scopes":        []interface{}{"scope1", "scope2"},
					"preemptible":         true,
					"service_account":     "service_account",
					"version":             "node_version",
				},
			},
			"project_id": "project-test",
			"resource_labels": map[string]interface{}{
				"Rlabel1": "value1",
				"Rlabel2": "value2",
			},
			"sub_network":                 "subnetwork",
			"use_ip_aliases":              true,
			"zone":                        "zone",
			"default_max_pods_constraint": 32,
			"region":                      "region",
		},
	}
}

func TestFlattenClusterGKEConfig(t *testing.T) {

	cases := map[string]struct {
		Input          *GoogleKubernetesEngineConfig
		ExpectedOutput []interface{}
	}{
		"GKECluster": {
			testClusterGKEConfigConf,
			testClusterGKEConfigInterface,
		},
		"LegacyGKECluster": {
			testClusterGKEConfigLegacyConf,
			testClusterGKEConfigInterface,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			output, err := flattenClusterGKEConfig(tc.Input)
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

func TestExpandClusterGKEConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *GoogleKubernetesEngineConfig
	}{
		{
			testClusterGKEConfigInterface,
			testClusterGKEConfigConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterGKEConfig(&GoogleKubernetesEngineConfig{}, tc.Input, "test")
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
