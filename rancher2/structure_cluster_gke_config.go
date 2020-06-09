package rancher2

// Flatteners

func flattenClusterGKEConfig(in *GoogleKubernetesEngineConfig, p []interface{}) ([]interface{}, error) {
	var obj map[string]interface{}
	if len(p) == 0 || p[0] == nil {
		obj = make(map[string]interface{})
	} else {
		obj = p[0].(map[string]interface{})
	}

	if in == nil {
		return []interface{}{}, nil
	}

	if len(in.ClusterIpv4Cidr) > 0 {
		obj["cluster_ipv4_cidr"] = in.ClusterIpv4Cidr
	}

	if len(in.Credential) > 0 {
		obj["credential"] = in.Credential
	}

	if len(in.Description) > 0 {
		obj["description"] = in.Description
	}

	if in.DiskSizeGb > 0 {
		obj["disk_size_gb"] = int(in.DiskSizeGb)
	}

	if len(in.DiskType) > 0 {
		obj["disk_type"] = in.DiskType
	}

	obj["enable_alpha_feature"] = in.EnableAlphaFeature
	obj["enable_auto_repair"] = in.EnableAutoRepair
	obj["enable_auto_upgrade"] = in.EnableAutoUpgrade

	if in.EnableHTTPLoadBalancing != nil {
		obj["enable_http_load_balancing"] = *in.EnableHTTPLoadBalancing
	}

	if in.EnableHorizontalPodAutoscaling != nil {
		obj["enable_horizontal_pod_autoscaling"] = *in.EnableHorizontalPodAutoscaling
	}

	obj["enable_kubernetes_dashboard"] = in.EnableKubernetesDashboard
	obj["enable_legacy_abac"] = in.EnableLegacyAbac
	obj["enable_master_authorized_network"] = in.EnableMasterAuthorizedNetwork

	if in.EnableNetworkPolicyConfig != nil {
		obj["enable_network_policy_config"] = *in.EnableNetworkPolicyConfig
	}

	obj["enable_nodepool_autoscaling"] = in.EnableNodepoolAutoscaling
	obj["enable_private_endpoint"] = in.EnablePrivateEndpoint
	obj["enable_private_nodes"] = in.EnablePrivateNodes

	if in.EnableStackdriverLogging != nil {
		obj["enable_stackdriver_logging"] = *in.EnableStackdriverLogging
	}

	if in.EnableHorizontalPodAutoscaling != nil {
		obj["enable_stackdriver_monitoring"] = *in.EnableStackdriverMonitoring
	}

	if len(in.ImageType) > 0 {
		obj["image_type"] = in.ImageType
	}

	if len(in.IPPolicyClusterIpv4CidrBlock) > 0 {
		obj["ip_policy_cluster_ipv4_cidr_block"] = in.IPPolicyClusterIpv4CidrBlock
	}

	if len(in.IPPolicyClusterSecondaryRangeName) > 0 {
		obj["ip_policy_cluster_secondary_range_name"] = in.IPPolicyClusterSecondaryRangeName
	}

	obj["ip_policy_create_subnetwork"] = in.IPPolicyCreateSubnetwork

	if len(in.IPPolicyNodeIpv4CidrBlock) > 0 {
		obj["ip_policy_node_ipv4_cidr_block"] = in.IPPolicyNodeIpv4CidrBlock
	}

	if len(in.IPPolicyServicesIpv4CidrBlock) > 0 {
		obj["ip_policy_services_ipv4_cidr_block"] = in.IPPolicyServicesIpv4CidrBlock
	}

	if len(in.IPPolicyServicesSecondaryRangeName) > 0 {
		obj["ip_policy_services_secondary_range_name"] = in.IPPolicyServicesSecondaryRangeName
	}

	if len(in.IPPolicySubnetworkName) > 0 {
		obj["ip_policy_subnetwork_name"] = in.IPPolicySubnetworkName
	}

	obj["issue_client_certificate"] = in.IssueClientCertificate
	obj["kubernetes_dashboard"] = in.KubernetesDashboard

	if len(in.Labels) > 0 {
		obj["labels"] = toMapInterface(in.Labels)
	}

	if in.LocalSsdCount > 0 {
		obj["local_ssd_count"] = int(in.LocalSsdCount)
	}

	if len(in.Locations) > 0 {
		obj["locations"] = toArrayInterface(in.Locations)
	}

	if len(in.MachineType) > 0 {
		obj["machine_type"] = in.MachineType
	}

	if len(in.MaintenanceWindow) > 0 {
		obj["maintenance_window"] = in.MaintenanceWindow
	}

	if len(in.MasterAuthorizedNetworkCidrBlocks) > 0 {
		obj["master_authorized_network_cidr_blocks"] = toArrayInterface(in.MasterAuthorizedNetworkCidrBlocks)
	}

	if len(in.MasterIpv4CidrBlock) > 0 {
		obj["master_ipv4_cidr_block"] = in.MasterIpv4CidrBlock
	}

	if len(in.MasterVersion) > 0 {
		obj["master_version"] = in.MasterVersion
	}

	if in.MaxNodeCount > 0 {
		obj["max_node_count"] = int(in.MaxNodeCount)
	}

	if in.MinNodeCount > 0 {
		obj["min_node_count"] = int(in.MinNodeCount)
	}

	if len(in.Network) > 0 {
		obj["network"] = in.Network
	}

	if in.NodeCount > 0 {
		obj["node_count"] = int(in.NodeCount)
	}

	if len(in.NodePool) > 0 {
		obj["node_pool"] = in.NodePool
	}

	if len(in.NodeVersion) > 0 {
		obj["node_version"] = in.NodeVersion
	}

	if len(in.OauthScopes) > 0 {
		obj["oauth_scopes"] = toArrayInterface(in.OauthScopes)
	}

	obj["preemptible"] = in.Preemptible

	if len(in.ProjectID) > 0 {
		obj["project_id"] = in.ProjectID
	}

	if len(in.Region) > 0 {
		obj["region"] = in.Region
	}

	if len(in.ResourceLabels) > 0 {
		obj["resource_labels"] = toMapInterface(in.ResourceLabels)
	}

	if len(in.ServiceAccount) > 0 {
		obj["service_account"] = in.ServiceAccount
	}

	if len(in.SubNetwork) > 0 {
		obj["sub_network"] = in.SubNetwork
	}

	obj["use_ip_aliases"] = in.UseIPAliases

	if len(in.Taints) > 0 {
		obj["taints"] = toArrayInterface(in.Taints)
	}

	if len(in.Zone) > 0 {
		obj["zone"] = in.Zone
	}

	return []interface{}{obj}, nil
}

// Expanders

func expandClusterGKEConfig(p []interface{}, name string) (*GoogleKubernetesEngineConfig, error) {
	obj := &GoogleKubernetesEngineConfig{}
	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	obj.DisplayName = name
	obj.Name = name
	obj.DriverName = clusterDriverGKE

	if v, ok := in["cluster_ipv4_cidr"].(string); ok && len(v) > 0 {
		obj.ClusterIpv4Cidr = v
	}

	if v, ok := in["credential"].(string); ok && len(v) > 0 {
		obj.Credential = v
	}

	if v, ok := in["description"].(string); ok && len(v) > 0 {
		obj.Description = v
	}

	if v, ok := in["disk_size_gb"].(int); ok && v > 0 {
		obj.DiskSizeGb = int64(v)
	}

	if v, ok := in["disk_type"].(string); ok && len(v) > 0 {
		obj.DiskType = v
	}

	if v, ok := in["enable_alpha_feature"].(bool); ok {
		obj.EnableAlphaFeature = v
	}

	if v, ok := in["enable_auto_repair"].(bool); ok {
		obj.EnableAutoRepair = v
	}

	if v, ok := in["enable_auto_upgrade"].(bool); ok {
		obj.EnableAutoUpgrade = v
	}

	if v, ok := in["enable_http_load_balancing"].(bool); ok {
		obj.EnableHTTPLoadBalancing = &v
	}

	if v, ok := in["enable_horizontal_pod_autoscaling"].(bool); ok {
		obj.EnableHorizontalPodAutoscaling = &v
	}

	if v, ok := in["enable_kubernetes_dashboard"].(bool); ok {
		obj.EnableKubernetesDashboard = v
	}

	if v, ok := in["enable_legacy_abac"].(bool); ok {
		obj.EnableLegacyAbac = v
	}

	if v, ok := in["enable_master_authorized_network"].(bool); ok {
		obj.EnableMasterAuthorizedNetwork = v
	}

	if v, ok := in["enable_network_policy_config"].(bool); ok {
		obj.EnableNetworkPolicyConfig = &v
	}

	if v, ok := in["enable_nodepool_autoscaling"].(bool); ok {
		obj.EnableNodepoolAutoscaling = v
	}

	if v, ok := in["enable_private_endpoint"].(bool); ok {
		obj.EnablePrivateEndpoint = v
	}

	if v, ok := in["enable_private_nodes"].(bool); ok {
		obj.EnablePrivateNodes = v
	}

	if v, ok := in["enable_stackdriver_logging"].(bool); ok {
		obj.EnableStackdriverLogging = &v
	}

	if v, ok := in["enable_stackdriver_monitoring"].(bool); ok {
		obj.EnableStackdriverMonitoring = &v
	}

	if v, ok := in["image_type"].(string); ok && len(v) > 0 {
		obj.ImageType = v
	}

	if v, ok := in["ip_policy_cluster_ipv4_cidr_block"].(string); ok && len(v) > 0 {
		obj.IPPolicyClusterIpv4CidrBlock = v
	}

	if v, ok := in["ip_policy_cluster_secondary_range_name"].(string); ok && len(v) > 0 {
		obj.IPPolicyClusterSecondaryRangeName = v
	}

	if v, ok := in["ip_policy_create_subnetwork"].(bool); ok {
		obj.IPPolicyCreateSubnetwork = v
	}

	if v, ok := in["ip_policy_node_ipv4_cidr_block"].(string); ok && len(v) > 0 {
		obj.IPPolicyNodeIpv4CidrBlock = v
	}

	if v, ok := in["ip_policy_services_ipv4_cidr_block"].(string); ok && len(v) > 0 {
		obj.IPPolicyServicesIpv4CidrBlock = v
	}

	if v, ok := in["ip_policy_services_secondary_range_name"].(string); ok && len(v) > 0 {
		obj.IPPolicyServicesSecondaryRangeName = v
	}

	if v, ok := in["ip_policy_subnetwork_name"].(string); ok && len(v) > 0 {
		obj.IPPolicySubnetworkName = v
	}

	if v, ok := in["issue_client_certificate"].(bool); ok {
		obj.IssueClientCertificate = v
	}

	if v, ok := in["kubernetes_dashboard"].(bool); ok {
		obj.KubernetesDashboard = v
	}

	if v, ok := in["labels"].(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	if v, ok := in["local_ssd_count"].(int); ok && v > 0 {
		obj.LocalSsdCount = int64(v)
	}

	if v, ok := in["locations"].([]interface{}); ok && len(v) > 0 {
		obj.Locations = toArrayString(v)
	}

	if v, ok := in["machine_type"].(string); ok && len(v) > 0 {
		obj.MachineType = v
	}

	if v, ok := in["maintenance_window"].(string); ok && len(v) > 0 {
		obj.MaintenanceWindow = v
	}

	if v, ok := in["master_authorized_network_cidr_blocks"].([]interface{}); ok && len(v) > 0 {
		obj.MasterAuthorizedNetworkCidrBlocks = toArrayString(v)
	}

	if v, ok := in["master_ipv4_cidr_block"].(string); ok && len(v) > 0 {
		obj.MasterIpv4CidrBlock = v
	}

	if v, ok := in["master_version"].(string); ok && len(v) > 0 {
		obj.MasterVersion = v
	}

	if v, ok := in["max_node_count"].(int); ok && v > 0 {
		obj.MaxNodeCount = int64(v)
	}

	if v, ok := in["min_node_count"].(int); ok && v > 0 {
		obj.MinNodeCount = int64(v)
	}

	if v, ok := in["network"].(string); ok && len(v) > 0 {
		obj.Network = v
	}

	if v, ok := in["node_count"].(int); ok && v > 0 {
		obj.NodeCount = int64(v)
	}

	if v, ok := in["node_pool"].(string); ok && len(v) > 0 {
		obj.NodePool = v
	}

	if v, ok := in["node_version"].(string); ok && len(v) > 0 {
		obj.NodeVersion = v
	}

	if v, ok := in["oauth_scopes"].([]interface{}); ok && len(v) > 0 {
		obj.OauthScopes = toArrayString(v)
	}

	if v, ok := in["preemptible"].(bool); ok {
		obj.Preemptible = v
	}

	if v, ok := in["project_id"].(string); ok && len(v) > 0 {
		obj.ProjectID = v
	}

	if v, ok := in["region"].(string); ok && len(v) > 0 {
		obj.Region = v
	}

	if v, ok := in["resource_labels"].(map[string]interface{}); ok && len(v) > 0 {
		obj.ResourceLabels = toMapString(v)
	}

	if v, ok := in["service_account"].(string); ok && len(v) > 0 {
		obj.ServiceAccount = v
	}

	if v, ok := in["sub_network"].(string); ok && len(v) > 0 {
		obj.SubNetwork = v
	}

	if v, ok := in["use_ip_aliases"].(bool); ok {
		obj.UseIPAliases = v
	}

	if v, ok := in["taints"].([]interface{}); ok && len(v) > 0 {
		obj.Taints = toArrayString(v)
	}

	if v, ok := in["zone"].(string); ok && len(v) > 0 {
		obj.Zone = v
	}

	return obj, nil
}
