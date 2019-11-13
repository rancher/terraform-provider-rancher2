package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

// Flatteners

func flattenClusterRegistationToken(in *managementClient.ClusterRegistrationToken) ([]interface{}, error) {
	obj := make(map[string]interface{})
	if in == nil {
		return []interface{}{}, nil
	}

	obj["id"] = in.ID
	obj["cluster_id"] = in.ClusterID
	obj["name"] = clusterRegistrationTokenName
	obj["command"] = in.Command
	obj["insecure_command"] = in.InsecureCommand
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

	if clusterRegToken == nil {
		return fmt.Errorf("[ERROR] flattening cluster: Input cluster registration token is nil")
	}

	if kubeConfig == nil {
		return fmt.Errorf("[ERROR] flattening cluster: Input cluster kube config is nil")
	}

	if in.ID != "" {
		d.SetId(in.ID)
	}

	d.Set("name", in.Name)
	d.Set("description", in.Description)

	err := d.Set("cluster_auth_endpoint", flattenClusterAuthEndpoint(in.LocalClusterAuthEndpoint))
	if err != nil {
		return err
	}

	if len(in.ClusterTemplateID) > 0 {
		d.Set("cluster_template_id", in.ClusterTemplateID)
		if len(in.ClusterTemplateRevisionID) > 0 {
			d.Set("cluster_template_revision_id", in.ClusterTemplateRevisionID)
		}
		if in.ClusterTemplateAnswers != nil {
			err = d.Set("cluster_template_answers", flattenAnswer(in.ClusterTemplateAnswers))
			if err != nil {
				return err
			}
		}
		if len(in.ClusterTemplateQuestions) > 0 {
			d.Set("cluster_template_questions", flattenQuestions(in.ClusterTemplateQuestions))
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

	d.Set("enable_cluster_alerting", in.EnableClusterAlerting)
	d.Set("enable_cluster_monitoring", in.EnableClusterMonitoring)
	d.Set("enable_cluster_istio", in.IstioEnabled)

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

	d.Set("kube_config", kubeConfig.Config)
	d.Set("default_project_id", defaultProjectID)
	d.Set("system_project_id", systemProjectID)
	d.Set("driver", in.Driver)

	switch in.Driver {
	case clusterDriverAKS:
		aksConfig, err := flattenClusterAKSConfig(in.AzureKubernetesServiceConfig)
		if err != nil {
			return err
		}
		err = d.Set("aks_config", aksConfig)
		if err != nil {
			return err
		}
	case clusterDriverEKS:
		eksConfig, err := flattenClusterEKSConfig(in.AmazonElasticContainerServiceConfig)
		if err != nil {
			return err
		}
		err = d.Set("eks_config", eksConfig)
		if err != nil {
			return err
		}
	case clusterDriverGKE:
		gkeConfig, err := flattenClusterGKEConfig(in.GoogleKubernetesEngineConfig)
		if err != nil {
			return err
		}
		err = d.Set("gke_config", gkeConfig)
		if err != nil {
			return err
		}
	}

	// Setting rke_config always as computed
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

	return nil
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

	if v, ok := in.Get("enable_cluster_istio").(bool); ok {
		obj.IstioEnabled = v
	}

	if v, ok := in.Get("enable_network_policy").(bool); ok {
		obj.EnableNetworkPolicy = &v
	}

	if v, ok := in.Get("aks_config").([]interface{}); ok && len(v) > 0 {
		aksConfig, err := expandClusterAKSConfig(v, obj.Name)
		if err != nil {
			return nil, err
		}
		obj.AzureKubernetesServiceConfig = aksConfig
		obj.Driver = clusterDriverAKS
	}

	if v, ok := in.Get("eks_config").([]interface{}); ok && len(v) > 0 {
		eksConfig, err := expandClusterEKSConfig(v, obj.Name)
		if err != nil {
			return nil, err
		}
		obj.AmazonElasticContainerServiceConfig = eksConfig
		obj.Driver = clusterDriverEKS
	}

	if v, ok := in.Get("gke_config").([]interface{}); ok && len(v) > 0 {
		gkeConfig, err := expandClusterGKEConfig(v, obj.Name)
		if err != nil {
			return nil, err
		}
		obj.GoogleKubernetesEngineConfig = gkeConfig
		obj.Driver = clusterDriverGKE
	}

	if v, ok := in.Get("rke_config").([]interface{}); ok && len(v) > 0 {
		rkeConfig, err := expandClusterRKEConfig(v, obj.Name)
		if err != nil {
			return nil, err
		}
		obj.RancherKubernetesEngineConfig = rkeConfig
		obj.Driver = clusterDriverRKE
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

	return obj, nil
}
