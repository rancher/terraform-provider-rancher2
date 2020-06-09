package rancher2

import (
	"reflect"
	"testing"
)

var (
	testClusterGKEConfigConf      *GoogleKubernetesEngineConfig
	testClusterGKEConfigInterface []interface{}
)

func init() {
	testClusterGKEConfigConf = &GoogleKubernetesEngineConfig{
		ClusterIpv4Cidr:                    "192.168.1.0/16",
		Credential:                         "credential",
		Description:                        "description",
		DiskType:                           "disk_type",
		DiskSizeGb:                         16,
		DisplayName:                        "test",
		DriverName:                         clusterDriverGKE,
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
		Name:                              "test",
		Network:                           "network",
		NodeCount:                         3,
		NodePool:                          "node_pool",
		NodeVersion:                       "node_version",
		OauthScopes:                       []string{"scope1", "scope2"},
		Preemptible:                       true,
		ProjectID:                         "project-test",
		Region:                            "region",
		ResourceLabels: map[string]string{
			"Rlabel1": "value1",
			"Rlabel2": "value2",
		},
		ServiceAccount: "service_account",
		SubNetwork:     "subnetwork",
		UseIPAliases:   true,
		Taints:         []string{"taint1", "taint2"},
		Zone:           "zone",
	}
	testClusterGKEConfigInterface = []interface{}{
		map[string]interface{}{
			"cluster_ipv4_cidr":                       "192.168.1.0/16",
			"credential":                              "credential",
			"description":                             "description",
			"disk_type":                               "disk_type",
			"disk_size_gb":                            16,
			"enable_alpha_feature":                    true,
			"enable_auto_repair":                      true,
			"enable_auto_upgrade":                     true,
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
			"image_type":                              "image",
			"ip_policy_cluster_ipv4_cidr_block":       "ip_policy_cluster_ipv4_cidr_block",
			"ip_policy_cluster_secondary_range_name":  "ip_policy_cluster_secondary_range_name",
			"ip_policy_create_subnetwork":             true,
			"ip_policy_node_ipv4_cidr_block":          "ip_policy_node_ipv4_cidr_block",
			"ip_policy_services_ipv4_cidr_block":      "ip_policy_services_ipv4_cidr_block",
			"ip_policy_services_secondary_range_name": "ip_policy_services_secondary_range_name",
			"ip_policy_subnetwork_name":               "ip_policy_subnetwork_name",
			"issue_client_certificate":                true,
			"kubernetes_dashboard":                    true,
			"labels": map[string]interface{}{
				"label1": "value1",
				"label2": "value2",
			},
			"local_ssd_count":                       1,
			"locations":                             []interface{}{"location1", "location2"},
			"machine_type":                          "machine",
			"maintenance_window":                    "maintenance",
			"master_authorized_network_cidr_blocks": []interface{}{"master1", "master2"},
			"master_ipv4_cidr_block":                "master_ipv4_cidr_block",
			"master_version":                        "version",
			"max_node_count":                        3,
			"min_node_count":                        1,
			"network":                               "network",
			"node_count":                            3,
			"node_pool":                             "node_pool",
			"node_version":                          "node_version",
			"oauth_scopes":                          []interface{}{"scope1", "scope2"},
			"preemptible":                           true,
			"project_id":                            "project-test",
			"region":                                "region",
			"resource_labels": map[string]interface{}{
				"Rlabel1": "value1",
				"Rlabel2": "value2",
			},
			"service_account": "service_account",
			"sub_network":     "subnetwork",
			"use_ip_aliases":  true,
			"taints":          []interface{}{"taint1", "taint2"},
			"zone":            "zone",
		},
	}
}

func TestFlattenClusterGKEConfig(t *testing.T) {

	cases := []struct {
		Input          *GoogleKubernetesEngineConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterGKEConfigConf,
			testClusterGKEConfigInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterGKEConfig(tc.Input, testClusterGKEConfigInterface)
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
		ExpectedOutput *GoogleKubernetesEngineConfig
	}{
		{
			testClusterGKEConfigInterface,
			testClusterGKEConfigConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterGKEConfig(tc.Input, "test")
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
