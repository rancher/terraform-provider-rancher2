package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testClusterGKEConfigV2NodeTaintsConfigConf                              []managementClient.GKENodeTaintConfig
	testClusterGKEConfigV2NodeTaintsConfigInterface                         []interface{}
	testClusterGKEConfigV2ClusterAddonsConf                                 *managementClient.GKEClusterAddons
	testClusterGKEConfigV2ClusterAddonsInterface                            []interface{}
	testClusterGKEConfigV2IPAllocationPolicyConf                            *managementClient.GKEIPAllocationPolicy
	testClusterGKEConfigV2IPAllocationPolicyInterface                       []interface{}
	testClusterGKEConfigV2MasterAuthorizedNetworksConfigCidrBlocksConf      []managementClient.GKECidrBlock
	testClusterGKEConfigV2MasterAuthorizedNetworksConfigCidrBlocksInterface []interface{}
	testClusterGKEConfigV2MasterAuthorizedNetworksConfigConf                *managementClient.GKEMasterAuthorizedNetworksConfig
	testClusterGKEConfigV2MasterAuthorizedNetworksConfigInterface           []interface{}
	testClusterGKEConfigV2NodePoolsAutoscalingConf                          *managementClient.GKENodePoolAutoscaling
	testClusterGKEConfigV2NodePoolsAutoscalingInterface                     []interface{}
	testClusterGKEConfigV2NodeConfigConf                                    *managementClient.GKENodeConfig
	testClusterGKEConfigV2NodeConfigInterface                               []interface{}
	testClusterGKEConfigV2NodePoolsManagementConf                           *managementClient.GKENodePoolManagement
	testClusterGKEConfigV2NodePoolsManagementInterface                      []interface{}
	testClusterGKEConfigV2NodePoolsConfigConf                               []managementClient.GKENodePoolConfig
	testClusterGKEConfigV2NodePoolsConfigInterface                          []interface{}
	testClusterGKEConfigV2PrivateClusterConfigConf                          *managementClient.GKEPrivateClusterConfig
	testClusterGKEConfigV2PrivateClusterConfigInterface                     []interface{}
	testClusterGKEConfigV2Conf                                              *managementClient.GKEClusterConfigSpec
	testClusterGKEConfigV2Interface                                         []interface{}
)

func init() {
	testClusterGKEConfigV2NodeTaintsConfigConf = []managementClient.GKENodeTaintConfig{
		{
			Key:    "key",
			Value:  "value",
			Effect: "recipient",
		},
	}
	testClusterGKEConfigV2NodeTaintsConfigInterface = []interface{}{
		map[string]interface{}{
			"key":    "key",
			"value":  "value",
			"effect": "recipient",
		},
	}
	testClusterGKEConfigV2ClusterAddonsConf = &managementClient.GKEClusterAddons{
		HTTPLoadBalancing:        true,
		HorizontalPodAutoscaling: true,
		NetworkPolicyConfig:      true,
	}
	testClusterGKEConfigV2ClusterAddonsInterface = []interface{}{
		map[string]interface{}{
			"http_load_balancing":        true,
			"horizontal_pod_autoscaling": true,
			"network_policy_config":      true,
		},
	}
	testClusterGKEConfigV2IPAllocationPolicyConf = &managementClient.GKEIPAllocationPolicy{
		ClusterIpv4CidrBlock:       "cluster_ipv4_cidr_block",
		ClusterSecondaryRangeName:  "cluster_secondary_range_name",
		CreateSubnetwork:           true,
		NodeIpv4CidrBlock:          "node_ipv4_cidr_block",
		ServicesIpv4CidrBlock:      "services_ipv4_cidr_block",
		ServicesSecondaryRangeName: "services_secondary_range_name",
		SubnetworkName:             "subnetwork_name",
		UseIPAliases:               true,
	}
	testClusterGKEConfigV2IPAllocationPolicyInterface = []interface{}{
		map[string]interface{}{
			"cluster_ipv4_cidr_block":       "cluster_ipv4_cidr_block",
			"cluster_secondary_range_name":  "cluster_secondary_range_name",
			"create_subnetwork":             true,
			"node_ipv4_cidr_block":          "node_ipv4_cidr_block",
			"services_ipv4_cidr_block":      "services_ipv4_cidr_block",
			"services_secondary_range_name": "services_secondary_range_name",
			"subnetwork_name":               "subnetwork_name",
			"use_ip_aliases":                true,
		},
	}
	testClusterGKEConfigV2MasterAuthorizedNetworksConfigCidrBlocksConf = []managementClient.GKECidrBlock{
		{
			CidrBlock:   "cidr_block",
			DisplayName: "display_name",
		},
	}
	testClusterGKEConfigV2MasterAuthorizedNetworksConfigCidrBlocksInterface = []interface{}{
		map[string]interface{}{
			"cidr_block":   "cidr_block",
			"display_name": "display_name",
		},
	}
	testClusterGKEConfigV2MasterAuthorizedNetworksConfigConf = &managementClient.GKEMasterAuthorizedNetworksConfig{
		CidrBlocks: testClusterGKEConfigV2MasterAuthorizedNetworksConfigCidrBlocksConf,
		Enabled:    true,
	}
	testClusterGKEConfigV2MasterAuthorizedNetworksConfigInterface = []interface{}{
		map[string]interface{}{
			"cidr_blocks": testClusterGKEConfigV2MasterAuthorizedNetworksConfigCidrBlocksInterface,
			"enabled":     true,
		},
	}
	testClusterGKEConfigV2NodePoolsAutoscalingConf = &managementClient.GKENodePoolAutoscaling{
		Enabled:      true,
		MaxNodeCount: int64(20),
		MinNodeCount: int64(10),
	}
	testClusterGKEConfigV2NodePoolsAutoscalingInterface = []interface{}{
		map[string]interface{}{
			"enabled":        true,
			"max_node_count": 20,
			"min_node_count": 10,
		},
	}
	testClusterGKEConfigV2NodeConfigConf = &managementClient.GKENodeConfig{
		DiskSizeGb: int64(20),
		DiskType:   "disk_type",
		ImageType:  "image_type",
		Labels: map[string]string{
			"label1": "value1",
			"label2": "value2",
		},
		LocalSsdCount: int64(1),
		MachineType:   "machine_type",
		OauthScopes:   []string{"oauth1", "oauth1"},
		Preemptible:   true,
		Tags:          []string{"tags1", "tags2"},
		Taints:        testClusterGKEConfigV2NodeTaintsConfigConf,
	}
	testClusterGKEConfigV2NodeConfigInterface = []interface{}{
		map[string]interface{}{
			"disk_size_gb": 20,
			"disk_type":    "disk_type",
			"image_type":   "image_type",
			"labels": map[string]interface{}{
				"label1": "value1",
				"label2": "value2",
			},
			"local_ssd_count": 1,
			"machine_type":    "machine_type",
			"oauth_scopes":    []interface{}{"oauth1", "oauth1"},
			"preemptible":     true,
			"tags":            []interface{}{"tags1", "tags2"},
			"taints":          testClusterGKEConfigV2NodeTaintsConfigInterface,
		},
	}
	testClusterGKEConfigV2NodePoolsManagementConf = &managementClient.GKENodePoolManagement{
		AutoRepair:  true,
		AutoUpgrade: true,
	}
	testClusterGKEConfigV2NodePoolsManagementInterface = []interface{}{
		map[string]interface{}{
			"auto_repair":  true,
			"auto_upgrade": true,
		},
	}
	initialNodeCount := int64(3)
	maxPodsConstraint := int64(10)
	testClusterGKEConfigV2NodePoolsConfigConf = []managementClient.GKENodePoolConfig{
		{
			Autoscaling:       testClusterGKEConfigV2NodePoolsAutoscalingConf,
			Config:            testClusterGKEConfigV2NodeConfigConf,
			InitialNodeCount:  &initialNodeCount,
			Management:        testClusterGKEConfigV2NodePoolsManagementConf,
			MaxPodsConstraint: &maxPodsConstraint,
			Name:              newString("name"),
			Version:           newString("version"),
		},
	}
	testClusterGKEConfigV2NodePoolsConfigInterface = []interface{}{
		map[string]interface{}{
			"autoscaling":         testClusterGKEConfigV2NodePoolsAutoscalingInterface,
			"config":              testClusterGKEConfigV2NodeConfigInterface,
			"initial_node_count":  3,
			"management":          testClusterGKEConfigV2NodePoolsManagementInterface,
			"max_pods_constraint": 10,
			"name":                "name",
			"version":             "version",
		},
	}

	testClusterGKEConfigV2PrivateClusterConfigConf = &managementClient.GKEPrivateClusterConfig{
		EnablePrivateEndpoint: true,
		EnablePrivateNodes:    true,
		MasterIpv4CidrBlock:   "master_ipv4_cidr_block",
	}
	testClusterGKEConfigV2PrivateClusterConfigInterface = []interface{}{
		map[string]interface{}{
			"enable_private_endpoint": true,
			"enable_private_nodes":    true,
			"master_ipv4_cidr_block":  "master_ipv4_cidr_block",
		},
	}

	testClusterGKEConfigV2Conf = &managementClient.GKEClusterConfigSpec{
		ClusterAddons:          testClusterGKEConfigV2ClusterAddonsConf,
		ClusterIpv4CidrBlock:   newString("cluster_ipv4_cidr_block"),
		ClusterName:            "name",
		Description:            "description",
		EnableKubernetesAlpha:  newTrue(),
		GoogleCredentialSecret: "google_credential_secret",
		IPAllocationPolicy:     testClusterGKEConfigV2IPAllocationPolicyConf,
		Imported:               false,
		KubernetesVersion:      newString("kubernetes_version"),
		Labels: map[string]string{
			"label1": "value1",
			"label2": "value2",
		},
		Locations:                      []string{"access1", "access2"},
		LoggingService:                 newString("logging_service"),
		MaintenanceWindow:              newString("maintenance_window"),
		MasterAuthorizedNetworksConfig: testClusterGKEConfigV2MasterAuthorizedNetworksConfigConf,
		MonitoringService:              newString("monitoring_service"),
		Network:                        newString("network"),
		NetworkPolicyEnabled:           newTrue(),
		NodePools:                      testClusterGKEConfigV2NodePoolsConfigConf,
		PrivateClusterConfig:           testClusterGKEConfigV2PrivateClusterConfigConf,
		ProjectID:                      "project_id",
		Region:                         "region",
		Subnetwork:                     newString("subnetwork"),
		Zone:                           "zone",
	}
	testClusterGKEConfigV2Interface = []interface{}{
		map[string]interface{}{
			"cluster_addons":           testClusterGKEConfigV2ClusterAddonsInterface,
			"cluster_ipv4_cidr_block":  "cluster_ipv4_cidr_block",
			"name":                     "name",
			"description":              "description",
			"enable_kubernetes_alpha":  true,
			"google_credential_secret": "google_credential_secret",
			"ip_allocation_policy":     testClusterGKEConfigV2IPAllocationPolicyInterface,
			"imported":                 false,
			"kubernetes_version":       "kubernetes_version",
			"labels": map[string]interface{}{
				"label1": "value1",
				"label2": "value2",
			},
			"locations":                         []interface{}{"access1", "access2"},
			"logging_service":                   "logging_service",
			"maintenance_window":                "maintenance_window",
			"master_authorized_networks_config": testClusterGKEConfigV2MasterAuthorizedNetworksConfigInterface,
			"monitoring_service":                "monitoring_service",
			"network":                           "network",
			"network_policy_enabled":            true,
			"node_pools":                        testClusterGKEConfigV2NodePoolsConfigInterface,
			"private_cluster_config":            testClusterGKEConfigV2PrivateClusterConfigInterface,
			"project_id":                        "project_id",
			"region":                            "region",
			"subnetwork":                        "subnetwork",
			"zone":                              "zone",
		},
	}
}

func TestFlattenClusterGKEConfigV2ClusterAddons(t *testing.T) {

	cases := []struct {
		Input          *managementClient.GKEClusterAddons
		ExpectedOutput []interface{}
	}{
		{
			testClusterGKEConfigV2ClusterAddonsConf,
			testClusterGKEConfigV2ClusterAddonsInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterGKEConfigV2ClusterAddons(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterGKEConfigV2IPAllocationPolicy(t *testing.T) {

	cases := []struct {
		Input          *managementClient.GKEIPAllocationPolicy
		ExpectedOutput []interface{}
	}{
		{
			testClusterGKEConfigV2IPAllocationPolicyConf,
			testClusterGKEConfigV2IPAllocationPolicyInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterGKEConfigV2IPAllocationPolicy(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterGKEConfigV2MasterAuthorizedNetworksConfig(t *testing.T) {

	cases := []struct {
		Input          *managementClient.GKEMasterAuthorizedNetworksConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterGKEConfigV2MasterAuthorizedNetworksConfigConf,
			testClusterGKEConfigV2MasterAuthorizedNetworksConfigInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterGKEConfigV2MasterAuthorizedNetworksConfig(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterGKEConfigV2NodePoolsConfig(t *testing.T) {

	cases := []struct {
		Input          []managementClient.GKENodePoolConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterGKEConfigV2NodePoolsConfigConf,
			testClusterGKEConfigV2NodePoolsConfigInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterGKEConfigV2NodePoolsConfig(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterGKEConfigV2PrivateClusterConfig(t *testing.T) {

	cases := []struct {
		Input          *managementClient.GKEPrivateClusterConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterGKEConfigV2PrivateClusterConfigConf,
			testClusterGKEConfigV2PrivateClusterConfigInterface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterGKEConfigV2PrivateClusterConfig(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterGKEConfigV2(t *testing.T) {

	cases := []struct {
		Input          *managementClient.GKEClusterConfigSpec
		ExpectedOutput []interface{}
	}{
		{
			testClusterGKEConfigV2Conf,
			testClusterGKEConfigV2Interface,
		},
	}

	for _, tc := range cases {
		output := flattenClusterGKEConfigV2(tc.Input, tc.ExpectedOutput)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterGKEConfigV2ClusterAddons(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.GKEClusterAddons
	}{
		{
			testClusterGKEConfigV2ClusterAddonsInterface,
			testClusterGKEConfigV2ClusterAddonsConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterGKEConfigV2ClusterAddons(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterGKEConfigV2IPAllocationPolicy(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.GKEIPAllocationPolicy
	}{
		{
			testClusterGKEConfigV2IPAllocationPolicyInterface,
			testClusterGKEConfigV2IPAllocationPolicyConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterGKEConfigV2IPAllocationPolicy(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterGKEConfigV2MasterAuthorizedNetworksConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.GKEMasterAuthorizedNetworksConfig
	}{
		{
			testClusterGKEConfigV2MasterAuthorizedNetworksConfigInterface,
			testClusterGKEConfigV2MasterAuthorizedNetworksConfigConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterGKEConfigV2MasterAuthorizedNetworksConfig(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterGKEConfigV2NodePoolsConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput []managementClient.GKENodePoolConfig
	}{
		{
			testClusterGKEConfigV2NodePoolsConfigInterface,
			testClusterGKEConfigV2NodePoolsConfigConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterGKEConfigV2NodePoolsConfig(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterGKEConfigV2PrivateClusterConfig(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.GKEPrivateClusterConfig
	}{
		{
			testClusterGKEConfigV2PrivateClusterConfigInterface,
			testClusterGKEConfigV2PrivateClusterConfigConf,
		},
	}

	for _, tc := range cases {
		output := expandClusterGKEConfigV2PrivateClusterConfig(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterGKEConfigV2(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.GKEClusterConfigSpec
	}{
		{
			testClusterGKEConfigV2Interface,
			testClusterGKEConfigV2Conf,
		},
	}

	for _, tc := range cases {
		output := expandClusterGKEConfigV2(tc.Input)
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
