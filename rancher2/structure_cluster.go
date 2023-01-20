package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenClusterRegistationToken(in *managementClient.ClusterRegistrationToken) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	obj["id"] = in.ID
	obj["cluster_id"] = in.ClusterID
	obj["name"] = in.Name
	obj["command"] = in.Command
	obj["insecure_command"] = in.InsecureCommand
	obj["insecure_node_command"] = in.InsecureNodeCommand
	obj["insecure_windows_node_command"] = in.InsecureWindowsNodeCommand
	obj["manifest_url"] = in.ManifestURL
	obj["node_command"] = in.NodeCommand
	obj["token"] = in.Token
	obj["windows_node_command"] = in.WindowsNodeCommand
	obj["annotations"] = toMapInterface(in.Annotations)
	obj["labels"] = toMapInterface(in.Labels)

	return []interface{}{obj}, nil
}

func flattenClusterAuthEndpoint(in *managementClient.LocalClusterAuthEndpoint) []interface{} {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}
	}

	obj["ca_certs"] = in.CACerts
	obj["enabled"] = in.Enabled
	obj["fqdn"] = in.FQDN

	return []interface{}{obj}
}

func flattenCluster(d *schema.ResourceData, in *Cluster, clusterRegToken *managementClient.ClusterRegistrationToken, kubeConfig *managementClient.GenerateKubeConfigOutput, defaultProjectID, systemProjectID string, monitoringInput *managementClient.MonitoringInput) error {
	if in == nil {
		return fmt.Errorf("[ERROR] flattening cluster: Input cluster is nil")
	}

	if kubeConfig == nil {
		return fmt.Errorf("[ERROR] flattening cluster: Input cluster kube config is nil")
	}

	if in.ID != "" {
		d.SetId(in.ID)
	}

	d.Set("name", in.Name)
	d.Set("description", in.Description)

	if len(in.AgentEnvVars) > 0 {
		d.Set("agent_env_vars", flattenEnvVars(in.AgentEnvVars))
	}

	err := d.Set("cluster_auth_endpoint", flattenClusterAuthEndpoint(in.LocalClusterAuthEndpoint))
	if err != nil {
		return err
	}

	if len(in.ClusterTemplateID) > 0 {
		d.Set("cluster_template_id", in.ClusterTemplateID)
		if len(in.ClusterTemplateRevisionID) > 0 {
			d.Set("cluster_template_revision_id", in.ClusterTemplateRevisionID)
		}
		if len(in.ClusterTemplateQuestions) > 0 {
			d.Set("cluster_template_questions", flattenQuestions(in.ClusterTemplateQuestions))
		}
		if in.ClusterTemplateAnswers != nil {
			for k, v := range readPreservedClusterTemplateAnswers(d) {
				if _, ok := in.ClusterTemplateAnswers.Values[k]; !ok {
					in.ClusterTemplateAnswers.Values[k] = v
				}
			}
			err = d.Set("cluster_template_answers", flattenAnswer(in.ClusterTemplateAnswers))
			if err != nil {
				return err
			}
		}

	}

	if len(in.DefaultPodSecurityPolicyTemplateID) > 0 {
		d.Set("default_pod_security_policy_template_id", in.DefaultPodSecurityPolicyTemplateID)
	}

	if len(in.DesiredAgentImage) > 0 {
		d.Set("desired_agent_image", in.DesiredAgentImage)
	}

	if len(in.DesiredAuthImage) > 0 {
		d.Set("desired_auth_image", in.DesiredAuthImage)
	}

	if len(in.DockerRootDir) > 0 {
		d.Set("docker_root_dir", in.DockerRootDir)
	}
	if len(in.FleetWorkspaceName) > 0 {
		d.Set("fleet_workspace_name", in.FleetWorkspaceName)
	}

	d.Set("enable_cluster_alerting", in.EnableClusterAlerting)
	d.Set("enable_cluster_monitoring", in.EnableClusterMonitoring)
	d.Set("istio_enabled", in.IstioEnabled)

	if in.EnableNetworkPolicy != nil {
		d.Set("enable_network_policy", *in.EnableNetworkPolicy)
	}

	err = d.Set("annotations", toMapInterface(in.Annotations))
	if err != nil {
		return err
	}
	err = d.Set("labels", toMapInterface(in.Labels))
	if err != nil {
		return err
	}
	regToken, err := flattenClusterRegistationToken(clusterRegToken)
	if err != nil {
		return err
	}
	err = d.Set("cluster_registration_token", regToken)
	if err != nil {
		return err
	}

	err = d.Set("cluster_monitoring_input", flattenMonitoringInput(monitoringInput))
	if err != nil {
		return err
	}

	if len(in.CACert) > 0 {
		d.Set("ca_cert", in.CACert)
	}

	d.Set("kube_config", kubeConfig.Config)
	d.Set("default_project_id", defaultProjectID)
	d.Set("system_project_id", systemProjectID)
	d.Set("driver", in.Driver)

	switch driver := ToLower(in.Driver); driver {
	case clusterDriverAKS:
		v, ok := d.Get("aks_config").([]interface{})
		if !ok {
			v = []interface{}{}
		}
		aksConfig, err := flattenClusterAKSConfig(in.AzureKubernetesServiceConfig, v)
		if err != nil {
			return err
		}
		err = d.Set("aks_config", aksConfig)
		if err != nil {
			return err
		}
	case ToLower(clusterDriverAKSV2):
		v, ok := d.Get("aks_config_v2").([]interface{})
		if !ok {
			v = []interface{}{}
		}
		err = d.Set("aks_config_v2", flattenClusterAKSConfigV2(in.AKSConfig, v))
		if err != nil {
			return err
		}
	case clusterDriverEKS:
		v, ok := d.Get("eks_config").([]interface{})
		if !ok {
			v = []interface{}{}
		}
		eksConfig, err := flattenClusterEKSConfig(in.AmazonElasticContainerServiceConfig, v)
		if err != nil {
			return err
		}
		err = d.Set("eks_config", eksConfig)
		if err != nil {
			return err
		}
	case ToLower(clusterDriverEKSV2):
		v, ok := d.Get("eks_config_v2").([]interface{})
		if !ok {
			v = []interface{}{}
		}
		err = d.Set("eks_config_v2", flattenClusterEKSConfigV2(in.EKSConfig, v))
		if err != nil {
			return err
		}
	case clusterDriverGKE:
		v, ok := d.Get("gke_config").([]interface{})
		if !ok {
			v = []interface{}{}
		}
		gkeConfig, err := flattenClusterGKEConfig(in.GoogleKubernetesEngineConfig, v)
		if err != nil {
			return err
		}
		err = d.Set("gke_config", gkeConfig)
		if err != nil {
			return err
		}
	case ToLower(clusterDriverGKEV2):
		v, ok := d.Get("gke_config_v2").([]interface{})
		if !ok {
			v = []interface{}{}
		}
		gkeConfig := flattenClusterGKEConfigV2(in.GKEConfig, v)
		err = d.Set("gke_config_v2", gkeConfig)
		if err != nil {
			return err
		}
	case clusterOKEKind, clusterDriverOKE:
		v, ok := d.Get("oke_config").([]interface{})
		if !ok {
			v = []interface{}{}
		}

		okeConfig, err := flattenClusterOKEConfig(in.OracleKubernetesEngineConfig, v)
		if err != nil {
			return err
		}
		err = d.Set("oke_config", okeConfig)
		if err != nil {
			return err
		}
	}

	// Setting k3s_config, rke2_config and rke_config always as computed
	err = d.Set("k3s_config", flattenClusterK3SConfig(in.K3sConfig))
	if err != nil {
		return err
	}
	err = d.Set("rke2_config", flattenClusterRKE2Config(in.Rke2Config))
	if err != nil {
		return err
	}
	v, ok := d.Get("rke_config").([]interface{})
	if !ok {
		v = []interface{}{}
	}
	rkeConfig, err := flattenClusterRKEConfig(in.RancherKubernetesEngineConfig, v)
	if err != nil {
		return err
	}
	err = d.Set("rke_config", rkeConfig)
	if err != nil {
		return err
	}

	d.Set("windows_prefered_cluster", in.WindowsPreferedCluster)

	return nil
}

func readPreservedClusterTemplateAnswers(d *schema.ResourceData) map[string]string {
	var questions []managementClient.Question
	if q, ok := d.Get("cluster_template_questions").([]interface{}); ok && len(q) > 0 {
		questions = expandQuestions(q)
	}

	var answers *managementClient.Answer
	if a, ok := d.Get("cluster_template_answers").([]interface{}); ok && len(a) > 0 {
		answers = expandAnswer(a)
	}

	preservedAnswers := map[string]string{}
	if questions != nil && answers != nil {
		for _, question := range questions {
			if question.Type == questionTypePassword {
				if answer, ok := answers.Values[question.Variable]; ok {
					preservedAnswers[question.Variable] = answer
				}
			}
		}
	}

	return preservedAnswers
}

func flattenClusterNodes(in []managementClient.Node) []interface{} {
	if len(in) == 0 {
		return []interface{}{}
	}
	out := make([]interface{}, len(in))
	for i, in := range in {
		obj := make(map[string]interface{})

		obj["annotations"] = toMapInterface(in.Annotations)
		obj["capacity"] = toMapInterface(in.Capacity)
		obj["cluster_id"] = in.ClusterID
		obj["external_ip_address"] = in.ExternalIPAddress
		obj["hostname"] = in.Hostname
		obj["id"] = in.ID
		obj["ip_address"] = in.IPAddress
		obj["labels"] = toMapInterface(in.Labels)
		obj["name"] = in.NodeName
		obj["node_pool_id"] = in.NodePoolID
		obj["node_template_id"] = in.NodeTemplateID
		obj["provider_id"] = in.ProviderId
		obj["requested_hostname"] = in.RequestedHostname
		obj["ssh_user"] = in.SshUser
		obj["system_info"] = flattenNodeInfo(in.Info)

		var roles []string
		if in.ControlPlane {
			roles = append(roles, "control_plane")
		}
		if in.Etcd {
			roles = append(roles, "etcd")
		}
		if in.Worker {
			roles = append(roles, "worker")
		}
		obj["roles"] = roles

		out[i] = obj
	}

	return out
}

func flattenNodeInfo(in *managementClient.NodeInfo) map[string]string {
	out := make(map[string]string)

	if in == nil {
		return map[string]string{}
	}

	out["kube_proxy_version"] = in.Kubernetes.KubeProxyVersion
	out["kubelet_version"] = in.Kubernetes.KubeletVersion
	out["container_runtime_version"] = in.OS.DockerVersion
	out["kernel_version"] = in.OS.KernelVersion
	out["operating_system"] = in.OS.OperatingSystem

	return out
}

// Expanders

func expandClusterRegistationToken(p []interface{}, clusterID string) (*managementClient.ClusterRegistrationToken, error) {
	if len(clusterID) == 0 {
		return nil, fmt.Errorf("[ERROR] Expanding Cluster Registration Token: Cluster id is nil")
	}

	obj := &managementClient.ClusterRegistrationToken{}
	obj.ClusterID = clusterID
	obj.Name = clusterRegistrationTokenName

	if len(p) == 0 || p[0] == nil {
		return obj, nil
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["id"].(string); ok && len(v) > 0 {
		obj.ID = v
	}

	if v, ok := in["annotations"].(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in["labels"].(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj, nil
}

func expandClusterAuthEndpoint(p []interface{}) *managementClient.LocalClusterAuthEndpoint {
	obj := &managementClient.LocalClusterAuthEndpoint{}
	if len(p) == 0 || p[0] == nil {
		return obj
	}
	in := p[0].(map[string]interface{})

	if v, ok := in["ca_certs"].(string); ok && len(v) > 0 {
		obj.CACerts = v
	}

	if v, ok := in["enabled"].(bool); ok {
		obj.Enabled = v
	}

	if v, ok := in["fqdn"].(string); ok && len(v) > 0 {
		obj.FQDN = v
	}

	return obj
}

func expandCluster(in *schema.ResourceData) (*Cluster, error) {
	obj := &Cluster{}
	if in == nil {
		return nil, fmt.Errorf("[ERROR] expanding cluster: Input ResourceData is nil")
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.Name = in.Get("name").(string)
	obj.Description = in.Get("description").(string)

	if v, ok := in.Get("agent_env_vars").([]interface{}); ok && len(v) > 0 {
		obj.AgentEnvVars = expandEnvVars(v)
	}

	if v, ok := in.Get("cluster_auth_endpoint").([]interface{}); ok && len(v) > 0 {
		obj.LocalClusterAuthEndpoint = expandClusterAuthEndpoint(v)
	}

	if v, ok := in.Get("cluster_template_id").(string); ok && len(v) > 0 {
		obj.ClusterTemplateID = v
		obj.Driver = clusterDriverRKE
		if v, ok := in.Get("cluster_template_revision_id").(string); ok && len(v) > 0 {
			obj.ClusterTemplateRevisionID = v
			obj.Driver = clusterDriverRKE
		}
		if v, ok := in.Get("cluster_template_answers").([]interface{}); ok && len(v) > 0 {
			obj.ClusterTemplateAnswers = expandAnswer(v)
		}
		if v, ok := in.Get("cluster_template_questions").([]interface{}); ok && len(v) > 0 {
			obj.ClusterTemplateQuestions = expandQuestions(v)
		}
	}

	if v, ok := in.Get("default_pod_security_policy_template_id").(string); ok && len(v) > 0 {
		obj.DefaultPodSecurityPolicyTemplateID = v
	}

	if v, ok := in.Get("desired_agent_image").(string); ok && len(v) > 0 {
		obj.DesiredAgentImage = v
	}

	if v, ok := in.Get("desired_auth_image").(string); ok && len(v) > 0 {
		obj.DesiredAuthImage = v
	}

	if v, ok := in.Get("docker_root_dir").(string); ok && len(v) > 0 {
		obj.DockerRootDir = v
	}

	if v, ok := in.Get("enable_cluster_alerting").(bool); ok {
		obj.EnableClusterAlerting = v
	}

	if v, ok := in.Get("enable_cluster_monitoring").(bool); ok {
		obj.EnableClusterMonitoring = v
	}

	if v, ok := in.Get("enable_network_policy").(bool); ok {
		obj.EnableNetworkPolicy = &v
	}

	if v, ok := in.Get("fleet_workspace_name").(string); ok && len(v) > 0 {
		obj.FleetWorkspaceName = v
	}

	if v, ok := in.Get("aks_config").([]interface{}); ok && len(v) > 0 {
		aksConfig, err := expandClusterAKSConfig(v, obj.Name)
		if err != nil {
			return nil, err
		}
		obj.AzureKubernetesServiceConfig = aksConfig
		obj.Driver = clusterDriverAKS
	}

	if v, ok := in.Get("aks_config_v2").([]interface{}); ok && len(v) > 0 {
		// Setting aks cluster name if empty
		if aksData, ok := v[0].(map[string]interface{}); ok {
			if name, ok := aksData["name"].(string); !ok || len(name) == 0 {
				aksData["name"] = obj.Name
				v[0] = aksData
			}
		}
		obj.AKSConfig = expandClusterAKSConfigV2(v)
		obj.Driver = clusterDriverAKSV2
	}

	if v, ok := in.Get("eks_config").([]interface{}); ok && len(v) > 0 {
		eksConfig, err := expandClusterEKSConfig(v, obj.Name)
		if err != nil {
			return nil, err
		}
		obj.AmazonElasticContainerServiceConfig = eksConfig
		obj.Driver = clusterDriverEKS
	}

	if v, ok := in.Get("eks_config_v2").([]interface{}); ok && len(v) > 0 {
		// Setting eks cluster name if empty
		if eksData, ok := v[0].(map[string]interface{}); ok {
			if name, ok := eksData["name"].(string); !ok || len(name) == 0 {
				eksData["name"] = obj.Name
				v[0] = eksData
			}
		}
		obj.EKSConfig = expandClusterEKSConfigV2(v)
		obj.Driver = clusterDriverEKSV2
	}

	if v, ok := in.Get("gke_config").([]interface{}); ok && len(v) > 0 {
		gkeConfig, err := expandClusterGKEConfig(v, obj.Name)
		if err != nil {
			return nil, err
		}
		obj.GoogleKubernetesEngineConfig = gkeConfig
		obj.Driver = clusterDriverGKE
	}

	if v, ok := in.Get("gke_config_v2").([]interface{}); ok && len(v) > 0 {
		gkeConfig := expandClusterGKEConfigV2(v)
		obj.GKEConfig = gkeConfig
		obj.Driver = clusterDriverGKEV2
	}

	if v, ok := in.Get("oke_config").([]interface{}); ok && len(v) > 0 {
		okeConfig, err := expandClusterOKEConfig(v, obj.Name)
		if err != nil {
			return nil, err
		}
		obj.OracleKubernetesEngineConfig = okeConfig
		obj.Driver = clusterOKEKind
	}

	if v, ok := in.Get("k3s_config").([]interface{}); ok && len(v) > 0 {
		obj.K3sConfig = expandClusterK3SConfig(v)
		obj.Driver = clusterDriverK3S
	}

	if v, ok := in.Get("rke_config").([]interface{}); ok && len(v) > 0 {
		rkeConfig, err := expandClusterRKEConfig(v, obj.Name)
		if err != nil {
			return nil, err
		}
		obj.RancherKubernetesEngineConfig = rkeConfig
		obj.Driver = clusterDriverRKE
	}

	if v, ok := in.Get("rke2_config").([]interface{}); ok && len(v) > 0 {
		obj.Rke2Config = expandClusterRKE2Config(v)
		obj.Driver = clusterDriverRKE2
	}

	if len(obj.Driver) == 0 {
		obj.Driver = clusterDriverImported
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	if v, ok := in.Get("windows_prefered_cluster").(bool); ok {
		obj.WindowsPreferedCluster = v
	}

	return obj, nil
}
