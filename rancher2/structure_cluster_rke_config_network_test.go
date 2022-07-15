package rancher2

import (
	"reflect"
	"testing"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

var (
	testClusterRKEConfigNetworkTolerationsConf      []managementClient.Toleration
	testClusterRKEConfigNetworkTolerationsInterface []interface{}
	testClusterRKEConfigNetworkAciConf              *managementClient.AciNetworkProvider
	testClusterRKEConfigNetworkAciInterface         []interface{}
	testClusterRKEConfigNetworkCalicoConf           *managementClient.CalicoNetworkProvider
	testClusterRKEConfigNetworkCalicoInterface      []interface{}
	testClusterRKEConfigNetworkCanalConf            *managementClient.CanalNetworkProvider
	testClusterRKEConfigNetworkCanalInterface       []interface{}
	testClusterRKEConfigNetworkFlannelConf          *managementClient.FlannelNetworkProvider
	testClusterRKEConfigNetworkFlannelInterface     []interface{}
	testClusterRKEConfigNetworkWeaveConf            *managementClient.WeaveNetworkProvider
	testClusterRKEConfigNetworkWeaveInterface       []interface{}
	testClusterRKEConfigNetworkConfAci              *managementClient.NetworkConfig
	testClusterRKEConfigNetworkInterfaceAci         []interface{}
	testClusterRKEConfigNetworkConfCalico           *managementClient.NetworkConfig
	testClusterRKEConfigNetworkInterfaceCalico      []interface{}
	testClusterRKEConfigNetworkConfCanal            *managementClient.NetworkConfig
	testClusterRKEConfigNetworkInterfaceCanal       []interface{}
	testClusterRKEConfigNetworkConfFlannel          *managementClient.NetworkConfig
	testClusterRKEConfigNetworkInterfaceFlannel     []interface{}
	testClusterRKEConfigNetworkConfWeave            *managementClient.NetworkConfig
	testClusterRKEConfigNetworkInterfaceWeave       []interface{}
)

func init() {
	seconds := int64(10)
	testClusterRKEConfigNetworkTolerationsConf = []managementClient.Toleration{
		{
			Key:               "key",
			Value:             "value",
			Effect:            "recipient",
			Operator:          "operator",
			TolerationSeconds: &seconds,
		},
	}
	testClusterRKEConfigNetworkTolerationsInterface = []interface{}{
		map[string]interface{}{
			"key":      "key",
			"value":    "value",
			"effect":   "recipient",
			"operator": "operator",
			"seconds":  10,
		},
	}
	testClusterRKEConfigNetworkAciConf = &managementClient.AciNetworkProvider{
		AEP:                               "RANCHER",
		ApicHosts:                         []string{"192.168.1.10", "192.168.1.11", "192.168.1.12"},
		ApicRefreshTickerAdjust:           "150",
		ApicRefreshTime:                   "1200",
		ApicSubscriptionDelay:             "100",
		ApicUserCrt:                       "cert1",
		ApicUserKey:                       "key1",
		ApicUserName:                      "user1",
		CApic:                             "capic",
		ControllerLogLevel:                "info",
		DisablePeriodicSnatGlobalInfoSync: "false",
		DisableWaitForNetwork:             "false",
		DropLogEnable:                     "false",
		DurationWaitForNetwork:            "300",
		DynamicExternalSubnet:             "172.10.2.1/24",
		EnableEndpointSlice:               "true",
		EncapType:                         "vxlan",
		EpRegistry:                        "registry",
		GbpPodSubnet:                      "172.10.4.1/22",
		HostAgentLogLevel:                 "info",
		ImagePullPolicy:                   "always",
		ImagePullSecret:                   "secret",
		InfraVlan:                         "4093",
		InstallIstio:                      "false",
		IstioProfile:                      "default",
		KafkaBrokers:                      []string{"10.0.0.10", "10.0.0.11", "10.0.0.12"},
		KafkaClientCrt:                    "cert2",
		KafkaClientKey:                    "key2",
		KubeAPIVlan:                       "4001",
		L3Out:                             "l3out",
		L3OutExternalNetworks:             []string{"OUT", "LAN"},
		MaxNodesSvcGraph:                  "100",
		McastRangeEnd:                     "225.20.255.255",
		McastRangeStart:                   "225.20.1.1",
		MTUHeadRoom:                       "100",
		MultusDisable:                     "true",
		NoPriorityClass:                   "false",
		NodePodIfEnable:                   "true",
		NodeSubnet:                        "172.10.0.1/24",
		OVSMemoryLimit:                    "1Gi",
		OpflexAgentLogLevel:               "info",
		OpflexClientSSL:                   "true",
		OpflexDeviceDeleteTimeout:         "600",
		OpflexMode:                        "mode",
		OpflexServerPort:                  "9000",
		OverlayVRFName:                    "vrf1",
		PBRTrackingNonSnat:                "false",
		PodSubnetChunkSize:                "32",
		RunGbpContainer:                   "false",
		RunOpflexServerContainer:          "false",
		ServiceGraphSubnet:                "172.10.3.1/24",
		ServiceMonitorInterval:            "5",
		ServiceVlan:                       "4003",
		SnatContractScope:                 "global",
		SnatNamespace:                     "aci-containers-system",
		SnatPortRangeEnd:                  "65000",
		SnatPortRangeStart:                "5000",
		SnatPortsPerNode:                  "3000",
		SriovEnable:                       "true",
		StaticExternalSubnet:              "172.10.1.1/24",
		SubnetDomainName:                  "domain.tld",
		SystemIdentifier:                  "id",
		Tenant:                            "acitenant",
		Token:                             "acitoken",
		UseAciAnywhereCRD:                 "false",
		UseAciCniPriorityClass:            "true",
		UseClusterRole:                    "true",
		UseHostNetnsVolume:                "false",
		UseOpflexServerVolume:             "false",
		UsePrivilegedContainer:            "false",
		VRFName:                           "vrf",
		VRFTenant:                         "vrf",
		VmmController:                     "controller",
		VmmDomain:                         "domain",
	}
	testClusterRKEConfigNetworkAciInterface = []interface{}{
		map[string]interface{}{
			"aep":                                    "RANCHER",
			"apic_hosts":                             []interface{}{"192.168.1.10", "192.168.1.11", "192.168.1.12"},
			"apic_refresh_ticker_adjust":             "150",
			"apic_refresh_time":                      "1200",
			"apic_subscription_delay":                "100",
			"apic_user_crt":                          "cert1",
			"apic_user_key":                          "key1",
			"apic_user_name":                         "user1",
			"capic":                                  "capic",
			"controller_log_level":                   "info",
			"disable_periodic_snat_global_info_sync": "false",
			"disable_wait_for_network":               "false",
			"drop_log_enable":                        "false",
			"duration_wait_for_network":              "300",
			"extern_dynamic":                         "172.10.2.1/24",
			"enable_endpoint_slice":                  "true",
			"encap_type":                             "vxlan",
			"ep_registry":                            "registry",
			"gbp_pod_subnet":                         "172.10.4.1/22",
			"host_agent_log_level":                   "info",
			"image_pull_policy":                      "always",
			"image_pull_secret":                      "secret",
			"infra_vlan":                             "4093",
			"install_istio":                          "false",
			"istio_profile":                          "default",
			"kafka_brokers":                          []interface{}{"10.0.0.10", "10.0.0.11", "10.0.0.12"},
			"kafka_client_crt":                       "cert2",
			"kafka_client_key":                       "key2",
			"kube_api_vlan":                          "4001",
			"l3out":                                  "l3out",
			"l3out_external_networks":                []interface{}{"OUT", "LAN"},
			"max_nodes_svc_graph":                    "100",
			"mcast_range_end":                        "225.20.255.255",
			"mcast_range_start":                      "225.20.1.1",
			"mtu_head_room":                          "100",
			"multus_disable":                         "true",
			"no_priority_class":                      "false",
			"node_pod_if_enable":                     "true",
			"node_subnet":                            "172.10.0.1/24",
			"ovs_memory_limit":                       "1Gi",
			"opflex_log_level":                       "info",
			"opflex_client_ssl":                      "true",
			"opflex_device_delete_timeout":           "600",
			"opflex_mode":                            "mode",
			"opflex_server_port":                     "9000",
			"overlay_vrf_name":                       "vrf1",
			"pbr_tracking_non_snat":                  "false",
			"pod_subnet_chunk_size":                  "32",
			"run_gbp_container":                      "false",
			"run_opflex_server_container":            "false",
			"node_svc_subnet":                        "172.10.3.1/24",
			"service_monitor_interval":               "5",
			"service_vlan":                           "4003",
			"snat_contract_scope":                    "global",
			"snat_namespace":                         "aci-containers-system",
			"snat_port_range_end":                    "65000",
			"snat_port_range_start":                  "5000",
			"snat_ports_per_node":                    "3000",
			"sriov_enable":                           "true",
			"extern_static":                          "172.10.1.1/24",
			"subnet_domain_name":                     "domain.tld",
			"system_id":                              "id",
			"tenant":                                 "acitenant",
			"token":                                  "acitoken",
			"use_aci_anywhere_crd":                   "false",
			"use_aci_cni_priority_class":             "true",
			"use_cluster_role":                       "true",
			"use_host_netns_volume":                  "false",
			"use_opflex_server_volume":               "false",
			"use_privileged_container":               "false",
			"vrf_name":                               "vrf",
			"vrf_tenant":                             "vrf",
			"vmm_controller":                         "controller",
			"vmm_domain":                             "domain",
		},
	}
	testClusterRKEConfigNetworkCalicoConf = &managementClient.CalicoNetworkProvider{
		CloudProvider: "aws",
	}
	testClusterRKEConfigNetworkCalicoInterface = []interface{}{
		map[string]interface{}{
			"cloud_provider": "aws",
		},
	}
	testClusterRKEConfigNetworkCanalConf = &managementClient.CanalNetworkProvider{
		Iface: "eth0",
	}
	testClusterRKEConfigNetworkCanalInterface = []interface{}{
		map[string]interface{}{
			"iface": "eth0",
		},
	}
	testClusterRKEConfigNetworkFlannelConf = &managementClient.FlannelNetworkProvider{
		Iface: "eth0",
	}
	testClusterRKEConfigNetworkFlannelInterface = []interface{}{
		map[string]interface{}{
			"iface": "eth0",
		},
	}
	testClusterRKEConfigNetworkWeaveConf = &managementClient.WeaveNetworkProvider{
		Password: "password",
	}
	testClusterRKEConfigNetworkWeaveInterface = []interface{}{
		map[string]interface{}{
			"password": "password",
		},
	}
	testClusterRKEConfigNetworkConfAci = &managementClient.NetworkConfig{
		AciNetworkProvider: testClusterRKEConfigNetworkAciConf,
		MTU:                1500,
		Options: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
		Plugin:      networkPluginAciName,
		Tolerations: testClusterRKEConfigNetworkTolerationsConf,
	}
	testClusterRKEConfigNetworkInterfaceAci = []interface{}{
		map[string]interface{}{
			"aci_network_provider": testClusterRKEConfigNetworkAciInterface,
			"mtu":                  1500,
			"options": map[string]interface{}{
				"option1": "value1",
				"option2": "value2",
			},
			"plugin":      networkPluginAciName,
			"tolerations": testClusterRKEConfigNetworkTolerationsInterface,
		},
	}
	testClusterRKEConfigNetworkConfCalico = &managementClient.NetworkConfig{
		CalicoNetworkProvider: testClusterRKEConfigNetworkCalicoConf,
		MTU:                   1500,
		Options: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
		Plugin:      networkPluginCalicoName,
		Tolerations: testClusterRKEConfigNetworkTolerationsConf,
	}
	testClusterRKEConfigNetworkInterfaceCalico = []interface{}{
		map[string]interface{}{
			"calico_network_provider": testClusterRKEConfigNetworkCalicoInterface,
			"mtu":                     1500,
			"options": map[string]interface{}{
				"option1": "value1",
				"option2": "value2",
			},
			"plugin":      networkPluginCalicoName,
			"tolerations": testClusterRKEConfigNetworkTolerationsInterface,
		},
	}
	testClusterRKEConfigNetworkConfCanal = &managementClient.NetworkConfig{
		CanalNetworkProvider: testClusterRKEConfigNetworkCanalConf,
		MTU:                  1500,
		Options: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
		Plugin:      networkPluginCanalName,
		Tolerations: testClusterRKEConfigNetworkTolerationsConf,
	}
	testClusterRKEConfigNetworkInterfaceCanal = []interface{}{
		map[string]interface{}{
			"canal_network_provider": testClusterRKEConfigNetworkCanalInterface,
			"mtu":                    1500,
			"options": map[string]interface{}{
				"option1": "value1",
				"option2": "value2",
			},
			"plugin":      networkPluginCanalName,
			"tolerations": testClusterRKEConfigNetworkTolerationsInterface,
		},
	}
	testClusterRKEConfigNetworkConfFlannel = &managementClient.NetworkConfig{
		FlannelNetworkProvider: testClusterRKEConfigNetworkFlannelConf,
		MTU:                    1500,
		Options: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
		Plugin:      networkPluginFlannelName,
		Tolerations: testClusterRKEConfigNetworkTolerationsConf,
	}
	testClusterRKEConfigNetworkInterfaceFlannel = []interface{}{
		map[string]interface{}{
			"flannel_network_provider": testClusterRKEConfigNetworkFlannelInterface,
			"mtu":                      1500,
			"options": map[string]interface{}{
				"option1": "value1",
				"option2": "value2",
			},
			"plugin":      networkPluginFlannelName,
			"tolerations": testClusterRKEConfigNetworkTolerationsInterface,
		},
	}
	testClusterRKEConfigNetworkConfWeave = &managementClient.NetworkConfig{
		WeaveNetworkProvider: testClusterRKEConfigNetworkWeaveConf,
		MTU:                  1500,
		Options: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
		Plugin:      networkPluginWeaveName,
		Tolerations: testClusterRKEConfigNetworkTolerationsConf,
	}
	testClusterRKEConfigNetworkInterfaceWeave = []interface{}{
		map[string]interface{}{
			"weave_network_provider": testClusterRKEConfigNetworkWeaveInterface,
			"mtu":                    1500,
			"options": map[string]interface{}{
				"option1": "value1",
				"option2": "value2",
			},
			"plugin":      networkPluginWeaveName,
			"tolerations": testClusterRKEConfigNetworkTolerationsInterface,
		},
	}
}

func TestFlattenClusterRKEConfigNetworkAci(t *testing.T) {

	cases := []struct {
		Input          *managementClient.AciNetworkProvider
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigNetworkAciConf,
			testClusterRKEConfigNetworkAciInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigNetworkAci(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterRKEConfigNetworkCalico(t *testing.T) {

	cases := []struct {
		Input          *managementClient.CalicoNetworkProvider
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigNetworkCalicoConf,
			testClusterRKEConfigNetworkCalicoInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigNetworkCalico(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterRKEConfigNetworkCanal(t *testing.T) {

	cases := []struct {
		Input          *managementClient.CanalNetworkProvider
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigNetworkCanalConf,
			testClusterRKEConfigNetworkCanalInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigNetworkCanal(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterRKEConfigNetworkFlannel(t *testing.T) {

	cases := []struct {
		Input          *managementClient.FlannelNetworkProvider
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigNetworkFlannelConf,
			testClusterRKEConfigNetworkFlannelInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigNetworkFlannel(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterRKEConfigNetworkWeave(t *testing.T) {

	cases := []struct {
		Input          *managementClient.WeaveNetworkProvider
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigNetworkWeaveConf,
			testClusterRKEConfigNetworkWeaveInterface,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigNetworkWeave(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestFlattenClusterRKEConfigNetwork(t *testing.T) {

	cases := []struct {
		Input          *managementClient.NetworkConfig
		ExpectedOutput []interface{}
	}{
		{
			testClusterRKEConfigNetworkConfAci,
			testClusterRKEConfigNetworkInterfaceAci,
		},
		{
			testClusterRKEConfigNetworkConfCalico,
			testClusterRKEConfigNetworkInterfaceCalico,
		},
		{
			testClusterRKEConfigNetworkConfCanal,
			testClusterRKEConfigNetworkInterfaceCanal,
		},
		{
			testClusterRKEConfigNetworkConfFlannel,
			testClusterRKEConfigNetworkInterfaceFlannel,
		},
		{
			testClusterRKEConfigNetworkConfWeave,
			testClusterRKEConfigNetworkInterfaceWeave,
		},
	}

	for _, tc := range cases {
		output, err := flattenClusterRKEConfigNetwork(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigNetworkAci(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.AciNetworkProvider
	}{
		{
			testClusterRKEConfigNetworkAciInterface,
			testClusterRKEConfigNetworkAciConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigNetworkAci(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigNetworkCalico(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.CalicoNetworkProvider
	}{
		{
			testClusterRKEConfigNetworkCalicoInterface,
			testClusterRKEConfigNetworkCalicoConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigNetworkCalico(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigNetworkCanal(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.CanalNetworkProvider
	}{
		{
			testClusterRKEConfigNetworkCanalInterface,
			testClusterRKEConfigNetworkCanalConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigNetworkCanal(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigNetworkFlannel(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.FlannelNetworkProvider
	}{
		{
			testClusterRKEConfigNetworkFlannelInterface,
			testClusterRKEConfigNetworkFlannelConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigNetworkFlannel(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigNetworkWeave(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.WeaveNetworkProvider
	}{
		{
			testClusterRKEConfigNetworkWeaveInterface,
			testClusterRKEConfigNetworkWeaveConf,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigNetworkWeave(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}

func TestExpandClusterRKEConfigNetwork(t *testing.T) {

	cases := []struct {
		Input          []interface{}
		ExpectedOutput *managementClient.NetworkConfig
	}{
		{
			testClusterRKEConfigNetworkInterfaceAci,
			testClusterRKEConfigNetworkConfAci,
		},
		{
			testClusterRKEConfigNetworkInterfaceCalico,
			testClusterRKEConfigNetworkConfCalico,
		},
		{
			testClusterRKEConfigNetworkInterfaceCanal,
			testClusterRKEConfigNetworkConfCanal,
		},
		{
			testClusterRKEConfigNetworkInterfaceFlannel,
			testClusterRKEConfigNetworkConfFlannel,
		},
		{
			testClusterRKEConfigNetworkInterfaceWeave,
			testClusterRKEConfigNetworkConfWeave,
		},
	}

	for _, tc := range cases {
		output, err := expandClusterRKEConfigNetwork(tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
