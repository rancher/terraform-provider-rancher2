package rancher2

import (
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenClusterRKEConfigNetworkAci(in *managementClient.AciNetworkProvider) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.AEP) > 0 {
		obj["aep"] = in.AEP
	}
	if len(in.ApicHosts) > 0 {
		obj["apic_hosts"] = toArrayInterface(in.ApicHosts)
	}
	if len(in.ApicRefreshTickerAdjust) > 0 {
		obj["apic_refresh_ticker_adjust"] = in.ApicRefreshTickerAdjust
	}
	if len(in.ApicRefreshTime) > 0 {
		obj["apic_refresh_time"] = in.ApicRefreshTime
	}
	if len(in.ApicSubscriptionDelay) > 0 {
		obj["apic_subscription_delay"] = in.ApicSubscriptionDelay
	}
	if len(in.ApicUserCrt) > 0 {
		obj["apic_user_crt"] = in.ApicUserCrt
	}
	if len(in.ApicUserKey) > 0 {
		obj["apic_user_key"] = in.ApicUserKey
	}
	if len(in.ApicUserName) > 0 {
		obj["apic_user_name"] = in.ApicUserName
	}
	if len(in.CApic) > 0 {
		obj["capic"] = in.CApic
	}
	if len(in.ControllerLogLevel) > 0 {
		obj["controller_log_level"] = in.ControllerLogLevel
	}
	if len(in.DisablePeriodicSnatGlobalInfoSync) > 0 {
		obj["disable_periodic_snat_global_info_sync"] = in.DisablePeriodicSnatGlobalInfoSync
	}
	if len(in.DisableWaitForNetwork) > 0 {
		obj["disable_wait_for_network"] = in.DisableWaitForNetwork
	}
	if len(in.DropLogEnable) > 0 {
		obj["drop_log_enable"] = in.DropLogEnable
	}
	if len(in.DurationWaitForNetwork) > 0 {
		obj["duration_wait_for_network"] = in.DurationWaitForNetwork
	}
	if len(in.DynamicExternalSubnet) > 0 {
		obj["extern_dynamic"] = in.DynamicExternalSubnet
	}
	if len(in.EnableEndpointSlice) > 0 {
		obj["enable_endpoint_slice"] = in.EnableEndpointSlice
	}
	if len(in.EncapType) > 0 {
		obj["encap_type"] = in.EncapType
	}
	if len(in.EpRegistry) > 0 {
		obj["ep_registry"] = in.EpRegistry
	}
	if len(in.GbpPodSubnet) > 0 {
		obj["gbp_pod_subnet"] = in.GbpPodSubnet
	}
	if len(in.HostAgentLogLevel) > 0 {
		obj["host_agent_log_level"] = in.HostAgentLogLevel
	}
	if len(in.ImagePullPolicy) > 0 {
		obj["image_pull_policy"] = in.ImagePullPolicy
	}
	if len(in.ImagePullSecret) > 0 {
		obj["image_pull_secret"] = in.ImagePullSecret
	}
	if len(in.InfraVlan) > 0 {
		obj["infra_vlan"] = in.InfraVlan
	}
	if len(in.InstallIstio) > 0 {
		obj["install_istio"] = in.InstallIstio
	}
	if len(in.IstioProfile) > 0 {
		obj["istio_profile"] = in.IstioProfile
	}
	if len(in.KafkaBrokers) > 0 {
		obj["kafka_brokers"] = toArrayInterface(in.KafkaBrokers)
	}
	if len(in.KafkaClientCrt) > 0 {
		obj["kafka_client_crt"] = in.KafkaClientCrt
	}
	if len(in.KafkaClientKey) > 0 {
		obj["kafka_client_key"] = in.KafkaClientKey
	}
	if len(in.KubeAPIVlan) > 0 {
		obj["kube_api_vlan"] = in.KubeAPIVlan
	}
	if len(in.L3Out) > 0 {
		obj["l3out"] = in.L3Out
	}
	if len(in.L3OutExternalNetworks) > 0 {
		obj["l3out_external_networks"] = toArrayInterface(in.L3OutExternalNetworks)
	}
	if len(in.MaxNodesSvcGraph) > 0 {
		obj["max_nodes_svc_graph"] = in.MaxNodesSvcGraph
	}
	if len(in.McastRangeEnd) > 0 {
		obj["mcast_range_end"] = in.McastRangeEnd
	}
	if len(in.McastRangeStart) > 0 {
		obj["mcast_range_start"] = in.McastRangeStart
	}
	if len(in.MTUHeadRoom) > 0 {
		obj["mtu_head_room"] = in.MTUHeadRoom
	}
	if len(in.MultusDisable) > 0 {
		obj["multus_disable"] = in.MultusDisable
	}
	if len(in.NoPriorityClass) > 0 {
		obj["no_priority_class"] = in.NoPriorityClass
	}
	if len(in.NodePodIfEnable) > 0 {
		obj["node_pod_if_enable"] = in.NodePodIfEnable
	}
	if len(in.NodeSubnet) > 0 {
		obj["node_subnet"] = in.NodeSubnet
	}
	if len(in.OVSMemoryLimit) > 0 {
		obj["ovs_memory_limit"] = in.OVSMemoryLimit
	}
	if len(in.OpflexAgentLogLevel) > 0 {
		obj["opflex_log_level"] = in.OpflexAgentLogLevel
	}
	if len(in.OpflexClientSSL) > 0 {
		obj["opflex_client_ssl"] = in.OpflexClientSSL
	}
	if len(in.OpflexDeviceDeleteTimeout) > 0 {
		obj["opflex_device_delete_timeout"] = in.OpflexDeviceDeleteTimeout
	}
	if len(in.OpflexMode) > 0 {
		obj["opflex_mode"] = in.OpflexMode
	}
	if len(in.OpflexServerPort) > 0 {
		obj["opflex_server_port"] = in.OpflexServerPort
	}
	if len(in.OverlayVRFName) > 0 {
		obj["overlay_vrf_name"] = in.OverlayVRFName
	}
	if len(in.PBRTrackingNonSnat) > 0 {
		obj["pbr_tracking_non_snat"] = in.PBRTrackingNonSnat
	}
	if len(in.PodSubnetChunkSize) > 0 {
		obj["pod_subnet_chunk_size"] = in.PodSubnetChunkSize
	}
	if len(in.RunGbpContainer) > 0 {
		obj["run_gbp_container"] = in.RunGbpContainer
	}
	if len(in.RunOpflexServerContainer) > 0 {
		obj["run_opflex_server_container"] = in.RunOpflexServerContainer
	}
	if len(in.ServiceGraphSubnet) > 0 {
		obj["node_svc_subnet"] = in.ServiceGraphSubnet
	}
	if len(in.ServiceMonitorInterval) > 0 {
		obj["service_monitor_interval"] = in.ServiceMonitorInterval
	}
	if len(in.ServiceVlan) > 0 {
		obj["service_vlan"] = in.ServiceVlan
	}
	if len(in.SnatContractScope) > 0 {
		obj["snat_contract_scope"] = in.SnatContractScope
	}
	if len(in.SnatNamespace) > 0 {
		obj["snat_namespace"] = in.SnatNamespace
	}
	if len(in.SnatPortRangeEnd) > 0 {
		obj["snat_port_range_end"] = in.SnatPortRangeEnd
	}
	if len(in.SnatPortRangeStart) > 0 {
		obj["snat_port_range_start"] = in.SnatPortRangeStart
	}
	if len(in.SnatPortsPerNode) > 0 {
		obj["snat_ports_per_node"] = in.SnatPortsPerNode
	}
	if len(in.SriovEnable) > 0 {
		obj["sriov_enable"] = in.SriovEnable
	}
	if len(in.StaticExternalSubnet) > 0 {
		obj["extern_static"] = in.StaticExternalSubnet
	}
	if len(in.SubnetDomainName) > 0 {
		obj["subnet_domain_name"] = in.SubnetDomainName
	}
	if len(in.SystemIdentifier) > 0 {
		obj["system_id"] = in.SystemIdentifier
	}
	if len(in.Tenant) > 0 {
		obj["tenant"] = in.Tenant
	}
	if len(in.Token) > 0 {
		obj["token"] = in.Token
	}
	if len(in.UseAciAnywhereCRD) > 0 {
		obj["use_aci_anywhere_crd"] = in.UseAciAnywhereCRD
	}
	if len(in.UseAciCniPriorityClass) > 0 {
		obj["use_aci_cni_priority_class"] = in.UseAciCniPriorityClass
	}
	if len(in.UseClusterRole) > 0 {
		obj["use_cluster_role"] = in.UseClusterRole
	}
	if len(in.UseHostNetnsVolume) > 0 {
		obj["use_host_netns_volume"] = in.UseHostNetnsVolume
	}
	if len(in.UseOpflexServerVolume) > 0 {
		obj["use_opflex_server_volume"] = in.UseOpflexServerVolume
	}
	if len(in.UsePrivilegedContainer) > 0 {
		obj["use_privileged_container"] = in.UsePrivilegedContainer
	}
	if len(in.VRFName) > 0 {
		obj["vrf_name"] = in.VRFName
	}
	if len(in.VRFTenant) > 0 {
		obj["vrf_tenant"] = in.VRFTenant
	}
	if len(in.VmmController) > 0 {
		obj["vmm_controller"] = in.VmmController
	}
	if len(in.VmmDomain) > 0 {
		obj["vmm_domain"] = in.VmmDomain
	}

	return []interface{}{obj}, nil
}

func flattenClusterRKEConfigNetworkCalico(in *managementClient.CalicoNetworkProvider) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.CloudProvider) > 0 {
		obj["cloud_provider"] = in.CloudProvider
	}

	return []interface{}{obj}, nil
}

func flattenClusterRKEConfigNetworkCanal(in *managementClient.CanalNetworkProvider) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.Iface) > 0 {
		obj["iface"] = in.Iface
	}

	return []interface{}{obj}, nil
}

func flattenClusterRKEConfigNetworkFlannel(in *managementClient.FlannelNetworkProvider) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.Iface) > 0 {
		obj["iface"] = in.Iface
	}

	return []interface{}{obj}, nil
}

func flattenClusterRKEConfigNetworkWeave(in *managementClient.WeaveNetworkProvider) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.Password) > 0 {
		obj["password"] = in.Password
	}

	return []interface{}{obj}, nil
}

func flattenClusterRKEConfigNetwork(in *managementClient.NetworkConfig) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	if in.AciNetworkProvider != nil {
		aciNetwork, err := flattenClusterRKEConfigNetworkAci(in.AciNetworkProvider)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["aci_network_provider"] = aciNetwork
	}

	if in.CalicoNetworkProvider != nil {
		calicoNetwork, err := flattenClusterRKEConfigNetworkCalico(in.CalicoNetworkProvider)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["calico_network_provider"] = calicoNetwork
	}

	if in.CanalNetworkProvider != nil {
		canalNetwork, err := flattenClusterRKEConfigNetworkCanal(in.CanalNetworkProvider)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["canal_network_provider"] = canalNetwork
	}

	if in.FlannelNetworkProvider != nil {
		flannelNetwork, err := flattenClusterRKEConfigNetworkFlannel(in.FlannelNetworkProvider)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["flannel_network_provider"] = flannelNetwork
	}

	if in.WeaveNetworkProvider != nil {
		weaveNetwork, err := flattenClusterRKEConfigNetworkWeave(in.WeaveNetworkProvider)
		if err != nil {
			return []interface{}{obj}, err
		}
		obj["weave_network_provider"] = weaveNetwork
	}

	if in.MTU > 0 {
		obj["mtu"] = int(in.MTU)
	}

	if len(in.Options) > 0 {
		obj["options"] = toMapInterface(in.Options)
	}

	if len(in.Plugin) > 0 {
		obj["plugin"] = in.Plugin
	}

	if in.Tolerations != nil && len(in.Tolerations) > 0 {
		obj["tolerations"] = flattenTolerations(in.Tolerations)
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandClusterRKEConfigNetworkAci(p []interface{}) (*managementClient.AciNetworkProvider, error) {
	obj := &managementClient.AciNetworkProvider{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["aep"].(string); ok && len(v) > 0 {
		obj.AEP = v
	}
	if v, ok := in["apic_hosts"].([]interface{}); ok && len(v) > 0 {
		obj.ApicHosts = toArrayString(v)
	}
	if v, ok := in["apic_refresh_ticker_adjust"].(string); ok && len(v) > 0 {
		obj.ApicRefreshTickerAdjust = v
	}
	if v, ok := in["apic_refresh_time"].(string); ok && len(v) > 0 {
		obj.ApicRefreshTime = v
	}
	if v, ok := in["apic_subscription_delay"].(string); ok && len(v) > 0 {
		obj.ApicSubscriptionDelay = v
	}
	if v, ok := in["apic_user_crt"].(string); ok && len(v) > 0 {
		obj.ApicUserCrt = v
	}
	if v, ok := in["apic_user_key"].(string); ok && len(v) > 0 {
		obj.ApicUserKey = v
	}
	if v, ok := in["apic_user_name"].(string); ok && len(v) > 0 {
		obj.ApicUserName = v
	}
	if v, ok := in["capic"].(string); ok && len(v) > 0 {
		obj.CApic = v
	}
	if v, ok := in["controller_log_level"].(string); ok && len(v) > 0 {
		obj.ControllerLogLevel = v
	}
	if v, ok := in["disable_periodic_snat_global_info_sync"].(string); ok && len(v) > 0 {
		obj.DisablePeriodicSnatGlobalInfoSync = v
	}
	if v, ok := in["disable_wait_for_network"].(string); ok && len(v) > 0 {
		obj.DisableWaitForNetwork = v
	}
	if v, ok := in["drop_log_enable"].(string); ok && len(v) > 0 {
		obj.DropLogEnable = v
	}
	if v, ok := in["duration_wait_for_network"].(string); ok && len(v) > 0 {
		obj.DurationWaitForNetwork = v
	}
	if v, ok := in["extern_dynamic"].(string); ok && len(v) > 0 {
		obj.DynamicExternalSubnet = v
	}
	if v, ok := in["enable_endpoint_slice"].(string); ok && len(v) > 0 {
		obj.EnableEndpointSlice = v
	}
	if v, ok := in["encap_type"].(string); ok && len(v) > 0 {
		obj.EncapType = v
	}
	if v, ok := in["ep_registry"].(string); ok && len(v) > 0 {
		obj.EpRegistry = v
	}
	if v, ok := in["gbp_pod_subnet"].(string); ok && len(v) > 0 {
		obj.GbpPodSubnet = v
	}
	if v, ok := in["host_agent_log_level"].(string); ok && len(v) > 0 {
		obj.HostAgentLogLevel = v
	}
	if v, ok := in["image_pull_policy"].(string); ok && len(v) > 0 {
		obj.ImagePullPolicy = v
	}
	if v, ok := in["image_pull_secret"].(string); ok && len(v) > 0 {
		obj.ImagePullSecret = v
	}
	if v, ok := in["infra_vlan"].(string); ok && len(v) > 0 {
		obj.InfraVlan = v
	}
	if v, ok := in["install_istio"].(string); ok && len(v) > 0 {
		obj.InstallIstio = v
	}
	if v, ok := in["istio_profile"].(string); ok && len(v) > 0 {
		obj.IstioProfile = v
	}
	if v, ok := in["kafka_brokers"].([]interface{}); ok && len(v) > 0 {
		obj.KafkaBrokers = toArrayString(v)
	}
	if v, ok := in["kafka_client_crt"].(string); ok && len(v) > 0 {
		obj.KafkaClientCrt = v
	}
	if v, ok := in["kafka_client_key"].(string); ok && len(v) > 0 {
		obj.KafkaClientKey = v
	}
	if v, ok := in["kube_api_vlan"].(string); ok && len(v) > 0 {
		obj.KubeAPIVlan = v
	}
	if v, ok := in["l3out"].(string); ok && len(v) > 0 {
		obj.L3Out = v
	}
	if v, ok := in["l3out_external_networks"].([]interface{}); ok && len(v) > 0 {
		obj.L3OutExternalNetworks = toArrayString(v)
	}
	if v, ok := in["max_nodes_svc_graph"].(string); ok && len(v) > 0 {
		obj.MaxNodesSvcGraph = v
	}
	if v, ok := in["mcast_range_end"].(string); ok && len(v) > 0 {
		obj.McastRangeEnd = v
	}
	if v, ok := in["mcast_range_start"].(string); ok && len(v) > 0 {
		obj.McastRangeStart = v
	}
	if v, ok := in["mtu_head_room"].(string); ok && len(v) > 0 {
		obj.MTUHeadRoom = v
	}
	if v, ok := in["multus_disable"].(string); ok && len(v) > 0 {
		obj.MultusDisable = v
	}
	if v, ok := in["no_priority_class"].(string); ok && len(v) > 0 {
		obj.NoPriorityClass = v
	}
	if v, ok := in["node_pod_if_enable"].(string); ok && len(v) > 0 {
		obj.NodePodIfEnable = v
	}
	if v, ok := in["node_subnet"].(string); ok && len(v) > 0 {
		obj.NodeSubnet = v
	}
	if v, ok := in["ovs_memory_limit"].(string); ok && len(v) > 0 {
		obj.OVSMemoryLimit = v
	}
	if v, ok := in["opflex_log_level"].(string); ok && len(v) > 0 {
		obj.OpflexAgentLogLevel = v
	}
	if v, ok := in["opflex_client_ssl"].(string); ok && len(v) > 0 {
		obj.OpflexClientSSL = v
	}
	if v, ok := in["opflex_device_delete_timeout"].(string); ok && len(v) > 0 {
		obj.OpflexDeviceDeleteTimeout = v
	}
	if v, ok := in["opflex_mode"].(string); ok && len(v) > 0 {
		obj.OpflexMode = v
	}
	if v, ok := in["opflex_server_port"].(string); ok && len(v) > 0 {
		obj.OpflexServerPort = v
	}
	if v, ok := in["overlay_vrf_name"].(string); ok && len(v) > 0 {
		obj.OverlayVRFName = v
	}
	if v, ok := in["pbr_tracking_non_snat"].(string); ok && len(v) > 0 {
		obj.PBRTrackingNonSnat = v
	}
	if v, ok := in["pod_subnet_chunk_size"].(string); ok && len(v) > 0 {
		obj.PodSubnetChunkSize = v
	}
	if v, ok := in["run_gbp_container"].(string); ok && len(v) > 0 {
		obj.RunGbpContainer = v
	}
	if v, ok := in["run_opflex_server_container"].(string); ok && len(v) > 0 {
		obj.RunOpflexServerContainer = v
	}
	if v, ok := in["node_svc_subnet"].(string); ok && len(v) > 0 {
		obj.ServiceGraphSubnet = v
	}
	if v, ok := in["service_monitor_interval"].(string); ok && len(v) > 0 {
		obj.ServiceMonitorInterval = v
	}
	if v, ok := in["service_vlan"].(string); ok && len(v) > 0 {
		obj.ServiceVlan = v
	}
	if v, ok := in["snat_contract_scope"].(string); ok && len(v) > 0 {
		obj.SnatContractScope = v
	}
	if v, ok := in["snat_namespace"].(string); ok && len(v) > 0 {
		obj.SnatNamespace = v
	}
	if v, ok := in["snat_port_range_end"].(string); ok && len(v) > 0 {
		obj.SnatPortRangeEnd = v
	}
	if v, ok := in["snat_port_range_start"].(string); ok && len(v) > 0 {
		obj.SnatPortRangeStart = v
	}
	if v, ok := in["snat_ports_per_node"].(string); ok && len(v) > 0 {
		obj.SnatPortsPerNode = v
	}
	if v, ok := in["sriov_enable"].(string); ok && len(v) > 0 {
		obj.SriovEnable = v
	}
	if v, ok := in["extern_static"].(string); ok && len(v) > 0 {
		obj.StaticExternalSubnet = v
	}
	if v, ok := in["subnet_domain_name"].(string); ok && len(v) > 0 {
		obj.SubnetDomainName = v
	}
	if v, ok := in["system_id"].(string); ok && len(v) > 0 {
		obj.SystemIdentifier = v
	}
	if v, ok := in["tenant"].(string); ok && len(v) > 0 {
		obj.Tenant = v
	}
	if v, ok := in["token"].(string); ok && len(v) > 0 {
		obj.Token = v
	}
	if v, ok := in["use_aci_anywhere_crd"].(string); ok && len(v) > 0 {
		obj.UseAciAnywhereCRD = v
	}
	if v, ok := in["use_aci_cni_priority_class"].(string); ok && len(v) > 0 {
		obj.UseAciCniPriorityClass = v
	}
	if v, ok := in["use_cluster_role"].(string); ok && len(v) > 0 {
		obj.UseClusterRole = v
	}
	if v, ok := in["use_host_netns_volume"].(string); ok && len(v) > 0 {
		obj.UseHostNetnsVolume = v
	}
	if v, ok := in["use_opflex_server_volume"].(string); ok && len(v) > 0 {
		obj.UseOpflexServerVolume = v
	}
	if v, ok := in["use_privileged_container"].(string); ok && len(v) > 0 {
		obj.UsePrivilegedContainer = v
	}
	if v, ok := in["vrf_name"].(string); ok && len(v) > 0 {
		obj.VRFName = v
	}
	if v, ok := in["vrf_tenant"].(string); ok && len(v) > 0 {
		obj.VRFTenant = v
	}
	if v, ok := in["vmm_controller"].(string); ok && len(v) > 0 {
		obj.VmmController = v
	}
	if v, ok := in["vmm_domain"].(string); ok && len(v) > 0 {
		obj.VmmDomain = v
	}

	return obj, nil
}

func expandClusterRKEConfigNetworkCalico(p []interface{}) (*managementClient.CalicoNetworkProvider, error) {
	obj := &managementClient.CalicoNetworkProvider{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["cloud_provider"].(string); ok && len(v) > 0 {
		obj.CloudProvider = v
	}

	return obj, nil
}

func expandClusterRKEConfigNetworkCanal(p []interface{}) (*managementClient.CanalNetworkProvider, error) {
	obj := &managementClient.CanalNetworkProvider{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["iface"].(string); ok && len(v) > 0 {
		obj.Iface = v
	}

	return obj, nil
}

func expandClusterRKEConfigNetworkFlannel(p []interface{}) (*managementClient.FlannelNetworkProvider, error) {
	obj := &managementClient.FlannelNetworkProvider{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["iface"].(string); ok && len(v) > 0 {
		obj.Iface = v
	}

	return obj, nil
}

func expandClusterRKEConfigNetworkWeave(p []interface{}) (*managementClient.WeaveNetworkProvider, error) {
	obj := &managementClient.WeaveNetworkProvider{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["password"].(string); ok && len(v) > 0 {
		obj.Password = v
	}

	return obj, nil
}

func expandClusterRKEConfigNetwork(p []interface{}) (*managementClient.NetworkConfig, error) {
	obj := &managementClient.NetworkConfig{}
	if len(p) == 0 || p[0] == nil {
		obj.Plugin = networkPluginDefault
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["aci_network_provider"].([]interface{}); ok && len(v) > 0 {
		aciNetwork, err := expandClusterRKEConfigNetworkAci(v)
		if err != nil {
			return obj, err
		}
		obj.AciNetworkProvider = aciNetwork
	}

	if v, ok := in["calico_network_provider"].([]interface{}); ok && len(v) > 0 {
		calicoNetwork, err := expandClusterRKEConfigNetworkCalico(v)
		if err != nil {
			return obj, err
		}
		obj.CalicoNetworkProvider = calicoNetwork
	}

	if v, ok := in["canal_network_provider"].([]interface{}); ok && len(v) > 0 {
		canalNetwork, err := expandClusterRKEConfigNetworkCanal(v)
		if err != nil {
			return obj, err
		}
		obj.CanalNetworkProvider = canalNetwork
	}

	if v, ok := in["flannel_network_provider"].([]interface{}); ok && len(v) > 0 {
		flannelNetwork, err := expandClusterRKEConfigNetworkFlannel(v)
		if err != nil {
			return obj, err
		}
		obj.FlannelNetworkProvider = flannelNetwork
	}

	if v, ok := in["weave_network_provider"].([]interface{}); ok && len(v) > 0 {
		weaveNetwork, err := expandClusterRKEConfigNetworkWeave(v)
		if err != nil {
			return obj, err
		}
		obj.WeaveNetworkProvider = weaveNetwork
	}

	if v, ok := in["mtu"].(int); ok && v > 0 {
		obj.MTU = int64(v)
	}

	if v, ok := in["options"].(map[string]interface{}); ok && len(v) > 0 {
		obj.Options = toMapString(v)
	}

	if v, ok := in["plugin"].(string); ok && len(v) > 0 {
		obj.Plugin = v
	}

	if v, ok := in["tolerations"].([]interface{}); ok && len(v) > 0 {
		obj.Tolerations = expandTolerations(v)
	}

	return obj, nil
}
