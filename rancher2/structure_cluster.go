package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
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
	obj["windows_node_command"] = in.WindowsNodeCommand
	obj["annotations"] = toMapInterface(in.Annotations)
	obj["labels"] = toMapInterface(in.Labels)

	return []interface{}{obj}, nil
}

func flattenCluster(d *schema.ResourceData, in *Cluster, clusterRegToken *managementClient.ClusterRegistrationToken, kubeConfig *managementClient.GenerateKubeConfigOutput) error {
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

	err := d.Set("name", in.Name)
	if err != nil {
		return err
	}
	err = d.Set("description", in.Description)
	if err != nil {
		return err
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

	err = d.Set("kube_config", kubeConfig.Config)
	if err != nil {
		return err
	}

	err = d.Set("driver", in.Driver)
	if err != nil {
		return err
	}

	switch in.Driver {
	case clusterDriverAKS:
		aksConfig, err := flattenClusterAKSConfig(in.AzureKubernetesServiceConfig)
		if err != nil {
			return err
		}
		d.Set("aks_config", aksConfig)
		if err != nil {
			return err
		}
	case clusterDriverEKS:
		eksConfig, err := flattenClusterEKSConfig(in.AmazonElasticContainerServiceConfig)
		if err != nil {
			return err
		}
		d.Set("eks_config", eksConfig)
		if err != nil {
			return err
		}
	case clusterDriverGKE:
		gkeConfig, err := flattenClusterGKEConfig(in.GoogleKubernetesEngineConfig)
		if err != nil {
			return err
		}
		d.Set("gke_config", gkeConfig)
		if err != nil {
			return err
		}
	case clusterDriverRKE:
		rkeConfig, err := flattenClusterRKEConfig(in.RancherKubernetesEngineConfig)
		if err != nil {
			return err
		}
		err = d.Set("rke_config", rkeConfig)
		if err != nil {
			return err
		}
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
