package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	clusterImportedKind = "imported"
)

var (
	clusterKinds = []string{clusterImportedKind, clusterEksKind, clusterAksKind, clusterGkeKind, clusterRkeKind}
)

// Schema

func clusterFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"name": &schema.Schema{
			Type:     schema.TypeString,
			Required: true,
		},
		"kind": &schema.Schema{
			Type:         schema.TypeString,
			Optional:     true,
			Default:      clusterRkeKind,
			ValidateFunc: validation.StringInSlice(clusterKinds, true),
		},
		"kube_config": &schema.Schema{
			Type:     schema.TypeString,
			Computed: true,
		},
		"rke_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"aks_config", "eks_config", "gke_config"},
			Elem: &schema.Resource{
				Schema: rkeConfigFields(),
			},
		},
		"eks_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"aks_config", "gke_config", "rke_config"},
			Elem: &schema.Resource{
				Schema: eksConfigFields(),
			},
		},
		"aks_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"eks_config", "gke_config", "rke_config"},
			Elem: &schema.Resource{
				Schema: aksConfigFields(),
			},
		},
		"gke_config": &schema.Schema{
			Type:          schema.TypeList,
			MaxItems:      1,
			Optional:      true,
			ConflictsWith: []string{"aks_config", "eks_config", "rke_config"},
			Elem: &schema.Resource{
				Schema: gkeConfigFields(),
			},
		},
		"description": &schema.Schema{
			Type:     schema.TypeString,
			Optional: true,
		},
		"cluster_registration_token": &schema.Schema{
			Type:     schema.TypeList,
			MaxItems: 1,
			Computed: true,
			Elem: &schema.Resource{
				Schema: clusterRegistationTokenFields(),
			},
		},
		"annotations": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
		"labels": &schema.Schema{
			Type:     schema.TypeMap,
			Optional: true,
			Computed: true,
		},
	}

	return s
}

// Flatteners

func flattenCluster(d *schema.ResourceData, in *managementClient.Cluster, clusterRegToken *managementClient.ClusterRegistrationToken, kubeConfig *managementClient.GenerateKubeConfigOutput) error {
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

	kind := d.Get("kind").(string)
	if kind == "" {
		if in.AzureKubernetesServiceConfig != nil {
			kind = clusterAksKind
		}
		if in.AmazonElasticContainerServiceConfig != nil {
			kind = clusterEksKind
		}
		if in.GoogleKubernetesEngineConfig != nil {
			kind = clusterGkeKind
		}
		if in.RancherKubernetesEngineConfig != nil {
			kind = clusterRkeKind
		}
		if kind == "" {
			kind = clusterImportedKind
		}

		err = d.Set("kind", kind)
		if err != nil {
			return err
		}
	}

	switch kind {
	case clusterAksKind:
		aksConfig, err := flattenAksConfig(in.AzureKubernetesServiceConfig)
		if err != nil {
			return err
		}
		d.Set("aks_config", aksConfig)
		if err != nil {
			return err
		}
	case clusterEksKind:
		eksConfig, err := flattenEksConfig(in.AmazonElasticContainerServiceConfig)
		if err != nil {
			return err
		}
		d.Set("eks_config", eksConfig)
		if err != nil {
			return err
		}
	case clusterGkeKind:
		gkeConfig, err := flattenGkeConfig(in.GoogleKubernetesEngineConfig)
		if err != nil {
			return err
		}
		d.Set("gke_config", gkeConfig)
		if err != nil {
			return err
		}
	case clusterRkeKind:
		rkeConfig, err := flattenRkeConfig(in.RancherKubernetesEngineConfig)
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

func expandCluster(in *schema.ResourceData) (*managementClient.Cluster, error) {
	obj := &managementClient.Cluster{}
	if in == nil {
		return nil, fmt.Errorf("[ERROR] expanding cluster: Input ResourceData is nil")
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.Name = in.Get("name").(string)
	obj.Description = in.Get("description").(string)

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	switch kind := in.Get("kind").(string); kind {
	case clusterRkeKind:
		rkeConfig, err := expandRkeConfig(in.Get("rke_config").([]interface{}))
		if err != nil {
			return nil, err
		}
		obj.RancherKubernetesEngineConfig = rkeConfig
	case clusterEksKind:
		eksConfig, err := expandEksConfig(in.Get("eks_config").([]interface{}))
		if err != nil {
			return nil, err
		}
		obj.AmazonElasticContainerServiceConfig = eksConfig
	case clusterAksKind:
		aksConfig, err := expandAksConfig(in.Get("aks_config").([]interface{}))
		if err != nil {
			return nil, err
		}
		obj.AzureKubernetesServiceConfig = aksConfig
	case clusterGkeKind:
		gkeConfig, err := expandGkeConfig(in.Get("gke_config").([]interface{}))
		if err != nil {
			return nil, err
		}
		obj.GoogleKubernetesEngineConfig = gkeConfig
	}

	return obj, nil
}

func resourceRancher2Cluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2ClusterCreate,
		Read:   resourceRancher2ClusterRead,
		Update: resourceRancher2ClusterUpdate,
		Delete: resourceRancher2ClusterDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2ClusterImport,
		},
		Schema: clusterFields(),
		// Setting default timeouts to be liberal in order to accommodate managed Kubernetes providers like EKS, GKE, and AKS
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
	}
}

func resourceRancher2ClusterCreate(d *schema.ResourceData, meta interface{}) error {
	cluster, err := expandCluster(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating Cluster %s", cluster.Name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	expectedState := "active"

	kind := d.Get("kind").(string)

	if kind == clusterImportedKind {
		expectedState = "pending"
	}

	if kind == clusterRkeKind {
		expectedState = "provisioning"
	}

	newCluster, err := client.Cluster.Create(cluster)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{expectedState},
		Refresh:    clusterStateRefreshFunc(client, newCluster.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for cluster (%s) to be created: %s", newCluster.ID, waitErr)
	}

	d.SetId(newCluster.ID)

	return resourceRancher2ClusterRead(d, meta)
}

func resourceRancher2ClusterRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Cluster ID %s", d.Id())

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	cluster, err := client.Cluster.ByID(d.Id())
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Cluster ID %s not found.", cluster.ID)
			d.SetId("")
			return nil
		}
		return err
	}

	clusterRegistrationToken, err := findClusterRegistrationToken(client, cluster.ID)
	if err != nil {
		return err
	}

	kubeConfig, err := client.Cluster.ActionGenerateKubeconfig(cluster)
	if err != nil {
		return err
	}

	err = flattenCluster(d, cluster, clusterRegistrationToken, kubeConfig)
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2ClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Cluster ID %s", d.Id())

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	cluster, err := client.Cluster.ByID(d.Id())
	if err != nil {
		return err
	}

	update := map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"annotations": toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":      toMapString(d.Get("labels").(map[string]interface{})),
	}

	switch kind := d.Get("kind").(string); kind {
	case clusterRkeKind:
		rkeConfig, err := expandRkeConfig(d.Get("rke_config").([]interface{}))
		if err != nil {
			return err
		}
		update["rancherKubernetesEngineConfig"] = rkeConfig
	case clusterEksKind:
		eksConfig, err := expandEksConfig(d.Get("eks_config").([]interface{}))
		if err != nil {
			return err
		}
		update["amazonElasticContainerServiceConfig"] = eksConfig
	case clusterAksKind:
		aksConfig, err := expandAksConfig(d.Get("aks_config").([]interface{}))
		if err != nil {
			return err
		}
		update["azureKubernetesServiceConfig"] = aksConfig
	case clusterGkeKind:
		gkeConfig, err := expandGkeConfig(d.Get("gke_config").([]interface{}))
		if err != nil {
			return err
		}
		update["googleKubernetesEngineConfig"] = gkeConfig
	}

	newCluster, err := meta.(*Config).UpdateClusterByID(cluster, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active", "provisioning", "pending", "updating"},
		Target:     []string{"active", "provisioning", "pending"},
		Refresh:    clusterStateRefreshFunc(client, newCluster.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for cluster (%s) to be updated: %s", newCluster.ID, waitErr)
	}

	d.SetId(newCluster.ID)

	return resourceRancher2ClusterRead(d, meta)
}

func resourceRancher2ClusterDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Cluster ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	cluster, err := client.Cluster.ByID(id)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Cluster ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.Cluster.Delete(cluster)
	if err != nil {
		return fmt.Errorf("Error removing Cluster: %s", err)
	}

	log.Printf("[DEBUG] Waiting for cluster (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"removing"},
		Target:     []string{"removed"},
		Refresh:    clusterStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for cluster (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

func resourceRancher2ClusterImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	err := resourceRancher2ClusterRead(d, meta)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}

// ClusterStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher Cluster.
func clusterStateRefreshFunc(client *managementClient.Client, clusterID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.Cluster.ByID(clusterID)
		if err != nil {
			if IsNotFound(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		if obj.Removed != "" {
			return obj, "removed", nil
		}

		return obj, obj.State, nil
	}
}
