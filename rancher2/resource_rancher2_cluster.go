package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	norman "github.com/rancher/norman/types"
	managementClient "github.com/rancher/types/client/management/v3"
)

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

	if cluster.Driver == clusterDriverImported {
		expectedState = "pending"
	}

	if cluster.Driver == clusterDriverRKE {
		expectedState = "provisioning"
	}

	newCluster := &Cluster{}
	err = client.APIBaseClient.Create(managementClient.ClusterType, cluster, newCluster)
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

	monitoringInput := expandMonitoringInput(d.Get("cluster_monitoring_input").([]interface{}))
	if newCluster.EnableClusterMonitoring && monitoringInput != nil {
		clusterResource := &norman.Resource{
			ID:      newCluster.ID,
			Type:    newCluster.Type,
			Links:   newCluster.Links,
			Actions: newCluster.Actions,
		}
		err = client.APIBaseClient.Action(managementClient.ClusterType, "editMonitoring", clusterResource, monitoringInput, nil)
		if err != nil {
			return err
		}
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

	cluster := &Cluster{}
	err = client.APIBaseClient.ByID(managementClient.ClusterType, d.Id(), cluster)
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

	clusterResource := &norman.Resource{
		ID:      cluster.ID,
		Type:    cluster.Type,
		Links:   cluster.Links,
		Actions: cluster.Actions,
	}
	kubeConfig := &managementClient.GenerateKubeConfigOutput{}
	err = client.APIBaseClient.Action(managementClient.ClusterType, "generateKubeconfig", clusterResource, nil, kubeConfig)
	if err != nil {
		return err
	}

	defaultProjectID, systemProjectID, err := meta.(*Config).GetClusterSpecialProjectsID(cluster.ID)
	if err != nil {
		return err
	}

	monitoringInput := &managementClient.MonitoringInput{}
	if cluster.EnableClusterMonitoring {
		monitoringOutput := &managementClient.MonitoringOutput{}
		err = client.APIBaseClient.Action(managementClient.ClusterType, "viewMonitoring", clusterResource, nil, monitoringOutput)
		if err != nil {
			return err
		}

		if monitoringOutput != nil && len(monitoringOutput.Answers) > 0 {
			monitoringInput.Answers = monitoringOutput.Answers
		}
	}

	err = flattenCluster(d, cluster, clusterRegistrationToken, kubeConfig, defaultProjectID, systemProjectID, monitoringInput)
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

	cluster := &norman.Resource{}
	err = client.APIBaseClient.ByID(managementClient.ClusterType, d.Id(), cluster)
	if err != nil {
		return err
	}

	enableNetworkPolicy := d.Get("enable_network_policy").(bool)
	update := map[string]interface{}{
		"name":                               d.Get("name").(string),
		"description":                        d.Get("description").(string),
		"defaultPodSecurityPolicyTemplateId": d.Get("default_pod_security_policy_template_id").(string),
		"desiredAgentImage":                  d.Get("desired_agent_image").(string),
		"desiredAuthImage":                   d.Get("desired_auth_image").(string),
		"dockerRootDir":                      d.Get("docker_root_dir").(string),
		"enableClusterAlerting":              d.Get("enable_cluster_alerting").(bool),
		"enableClusterMonitoring":            d.Get("enable_cluster_monitoring").(bool),
		"enableNetworkPolicy":                &enableNetworkPolicy,
		"istioEnabled":                       d.Get("enable_cluster_istio").(bool),
		"localClusterAuthEndpoint":           expandClusterAuthEndpoint(d.Get("cluster_auth_endpoint").([]interface{})),
		"annotations":                        toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":                             toMapString(d.Get("labels").(map[string]interface{})),
	}

	if clusterTemplateID, ok := d.Get("cluster_template_id").(string); ok && len(clusterTemplateID) > 0 {
		update["clusterTemplateId"] = clusterTemplateID
		if clusterTemplateRevisionID, ok := d.Get("cluster_template_revision_id").(string); ok && len(clusterTemplateRevisionID) > 0 {
			update["clusterTemplateRevisionId"] = clusterTemplateRevisionID
		}
		if answers, ok := d.Get("cluster_template_answers").([]interface{}); ok && len(answers) > 0 {
			update["answers"] = expandAnswer(answers)
		}
		if questions, ok := d.Get("cluster_template_questions").([]interface{}); ok && len(questions) > 0 {
			update["questions"] = expandQuestions(questions)
		}
	}

	switch driver := d.Get("driver").(string); driver {
	case clusterDriverAKS:
		aksConfig, err := expandClusterAKSConfig(d.Get("aks_config").([]interface{}), d.Get("name").(string))
		if err != nil {
			return err
		}
		update["azureKubernetesServiceConfig"] = aksConfig
	case clusterDriverEKS:
		eksConfig, err := expandClusterEKSConfig(d.Get("eks_config").([]interface{}), d.Get("name").(string))
		if err != nil {
			return err
		}
		update["amazonElasticContainerServiceConfig"] = eksConfig
	case clusterDriverGKE:
		gkeConfig, err := expandClusterGKEConfig(d.Get("gke_config").([]interface{}), d.Get("name").(string))
		if err != nil {
			return err
		}
		update["googleKubernetesEngineConfig"] = gkeConfig
	case clusterDriverRKE:
		rkeConfig, err := expandClusterRKEConfig(d.Get("rke_config").([]interface{}), d.Get("name").(string))
		if err != nil {
			return err
		}
		update["rancherKubernetesEngineConfig"] = rkeConfig
	}

	newCluster := &Cluster{}
	err = client.APIBaseClient.Update(managementClient.ClusterType, cluster, update, newCluster)
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

	if d.HasChange("enable_cluster_monitoring") || d.HasChange("cluster_monitoring_input") {
		monitoringInput := expandMonitoringInput(d.Get("cluster_monitoring_input").([]interface{}))
		if newCluster.EnableClusterMonitoring && monitoringInput != nil {
			clusterResource := &norman.Resource{
				ID:      newCluster.ID,
				Type:    newCluster.Type,
				Links:   newCluster.Links,
				Actions: newCluster.Actions,
			}
			err = client.APIBaseClient.Action(managementClient.ClusterType, "editMonitoring", clusterResource, monitoringInput, nil)
			if err != nil {
				return err
			}
		}
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

	cluster := &norman.Resource{}
	err = client.APIBaseClient.ByID(managementClient.ClusterType, d.Id(), cluster)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Cluster ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.APIBaseClient.Delete(cluster)
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

// clusterStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher Cluster.
func clusterStateRefreshFunc(client *managementClient.Client, clusterID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj := &Cluster{}
		err := client.APIBaseClient.ByID(managementClient.ClusterType, clusterID, obj)
		if err != nil {
			// The IsForbidden check is used in the case the user performing the action does not have the
			// right to retrieve the full list of clusters. If the user tries to retrieve the cluster that
			// just got deleted, instead of getting a 404 not found response it will get a 403 forbidden
			// eventhough it had the right to access the cluster before it was deleted. If we reach this
			// code path, it means that the user had the right to access the cluster, delete it, hence
			// meaning that the delete was successful.
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		return obj, obj.State, nil
	}
}

// clusterRegistrationTokenStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher ClusterRegistrationToken.
func clusterRegistrationTokenStateRefreshFunc(client *managementClient.Client, clusterID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.ClusterRegistrationToken.ByID(clusterID)
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

func findFlattenClusterRegistrationToken(client *managementClient.Client, clusterID string) ([]interface{}, error) {
	clusterReg, err := findClusterRegistrationToken(client, clusterID)
	if err != nil {
		return []interface{}{}, err
	}

	return flattenClusterRegistationToken(clusterReg)
}

func findClusterRegistrationToken(client *managementClient.Client, clusterID string) (*managementClient.ClusterRegistrationToken, error) {
	regTokenID := clusterID + ":" + clusterRegistrationTokenName
	regToken, err := client.ClusterRegistrationToken.ByID(regTokenID)

	if err != nil {
		if IsNotFound(err) {
			return createClusterRegistrationToken(client, clusterID)
		}
		return nil, err
	}

	return regToken, nil
}

func createClusterRegistrationToken(client *managementClient.Client, clusterID string) (*managementClient.ClusterRegistrationToken, error) {
	regToken, err := expandClusterRegistationToken([]interface{}{}, clusterID)
	if err != nil {
		return nil, err
	}

	newRegToken, err := client.ClusterRegistrationToken.Create(regToken)
	if err != nil {
		return nil, err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    clusterRegistrationTokenStateRefreshFunc(client, newRegToken.ID),
		Timeout:    5 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return nil, fmt.Errorf("[ERROR] waiting for cluster registration token (%s) to be created: %s", newRegToken.ID, waitErr)
	}
	return newRegToken, nil
}
