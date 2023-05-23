package rancher2

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	norman "github.com/rancher/norman/types"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	projectClient "github.com/rancher/rancher/pkg/client/generated/project/v3"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
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
		CustomizeDiff: func(d *schema.ResourceDiff, i interface{}) error {
			if d.Get("driver") == clusterDriverEKSV2 && d.HasChange("eks_config_v2") {
				old, new := d.GetChange("eks_config_v2")
				oldObj := expandClusterEKSConfigV2(old.([]interface{}))
				newObj := expandClusterEKSConfigV2(new.([]interface{}))
				if reflect.DeepEqual(oldObj, newObj) {
					d.Clear("eks_config_v2")
				} else {
					d.SetNew("eks_config_v2", flattenClusterEKSConfigV2(newObj, []interface{}{}))
				}
			}
			return nil
		},
		Schema:        clusterFields(),
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceRancher2ClusterResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceRancher2ClusterStateUpgradeV0,
				Version: 0,
			},
		},
		// Setting default timeouts to be liberal in order to accommodate managed Kubernetes providers like EKS, GKE, and AKS
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
	}
}

func resourceRancher2ClusterResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: clusterFieldsV0(),
	}
}

func resourceRancher2ClusterStateUpgradeV0(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if rkeConfigs, ok := rawState["rke_config"].([]interface{}); ok && len(rkeConfigs) > 0 {
		for i1 := range rkeConfigs {
			if rkeConfig, ok := rkeConfigs[i1].(map[string]interface{}); ok && len(rkeConfig) > 0 {
				if services, ok := rkeConfig["services"].([]interface{}); ok && len(services) > 0 {
					for i2 := range services {
						if service, ok := services[i2].(map[string]interface{}); ok && len(service) > 0 {
							if kubeApis, ok := service["kube_api"].([]interface{}); ok && len(kubeApis) > 0 {
								for i3 := range kubeApis {
									if kubeAPI, ok := kubeApis[i3].(map[string]interface{}); ok && len(kubeAPI) > 0 {
										if eventRates, ok := kubeAPI["event_rate_limit"].([]interface{}); ok && len(eventRates) > 0 {
											for i4 := range eventRates {
												if eventRate, ok := eventRates[i4].(map[string]interface{}); ok && len(eventRate) > 0 {
													if config, ok := eventRate["configuration"].(map[string]interface{}); ok {
														newValue := ""
														if len(config) > 0 {
															conf, err := mapInterfaceToYAML(config)
															if err == nil {
																newValue = conf
															}
														}
														rawState["rke_config"].([]interface{})[i1].(map[string]interface{})["services"].([]interface{})[i2].(map[string]interface{})["kube_api"].([]interface{})[i3].(map[string]interface{})["event_rate_limit"].([]interface{})[i4].(map[string]interface{})["configuration"] = newValue
													}
												}
											}
										}
										if secretEncs, ok := kubeAPI["secrets_encryption_config"].([]interface{}); ok && len(secretEncs) > 0 {
											for i4 := range secretEncs {
												if secretEnc, ok := secretEncs[i4].(map[string]interface{}); ok && len(secretEnc) > 0 {
													if config, ok := secretEnc["custom_config"].(map[string]interface{}); ok {
														newValue := ""
														if len(config) > 0 {
															conf, err := mapInterfaceToYAML(config)
															if err == nil {
																newValue = conf
															}
														}
														rawState["rke_config"].([]interface{})[i1].(map[string]interface{})["services"].([]interface{})[i2].(map[string]interface{})["kube_api"].([]interface{})[i3].(map[string]interface{})["secrets_encryption_config"].([]interface{})[i4].(map[string]interface{})["custom_config"] = newValue
													}
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return rawState, nil
}

func resourceRancher2ClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	cluster, err := expandCluster(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating Cluster %s", cluster.Name)

	expectedState := "active"

	if cluster.Driver == clusterDriverImported || (cluster.Driver == clusterDriverEKSV2 && cluster.EKSConfig.Imported) {
		expectedState = "pending"
	}

	if cluster.Driver == clusterDriverRKE || cluster.Driver == clusterDriverK3S || cluster.Driver == clusterDriverRKE2 {
		expectedState = "provisioning"
	}

	// Creating cluster with monitoring disabled
	cluster.EnableClusterMonitoring = false
	newCluster := &Cluster{}
	if cluster.EKSConfig != nil && !cluster.EKSConfig.Imported {
		clusterStr, _ := interfaceToJSON(cluster)
		clusterMap, _ := jsonToMapInterface(clusterStr)
		clusterMap["eksConfig"] = fixClusterEKSConfigV2(d.Get("eks_config_v2").([]interface{}), structToMap(cluster.EKSConfig))
		err = client.APIBaseClient.Create(managementClient.ClusterType, clusterMap, newCluster)
	} else if cluster.GKEConfig != nil && !cluster.GKEConfig.Imported {
		clusterStr, _ := interfaceToJSON(cluster)
		clusterMap, _ := jsonToMapInterface(clusterStr)
		clusterMap["gkeConfig"] = fixClusterGKEConfigV2(structToMap(cluster.GKEConfig))
		err = client.APIBaseClient.Create(managementClient.ClusterType, clusterMap, newCluster)
	} else {
		err = client.APIBaseClient.Create(managementClient.ClusterType, cluster, newCluster)
	}
	if err != nil {
		return err
	}

	newCluster.EnableClusterMonitoring = d.Get("enable_cluster_monitoring").(bool)
	d.SetId(newCluster.ID)

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
	if newCluster.EnableClusterMonitoring {
		if len(newCluster.Actions[monitoringActionEnable]) == 0 {
			err = client.APIBaseClient.ByID(managementClient.ClusterType, newCluster.ID, newCluster)
			if err != nil {
				return err
			}
		}
		clusterResource := &norman.Resource{
			ID:      newCluster.ID,
			Type:    newCluster.Type,
			Links:   newCluster.Links,
			Actions: newCluster.Actions,
		}
		// Retry enabling monitoring until timeout if got apierr 500 https://github.com/rancher/rancher/issues/30188
		ctx, cancel := context.WithTimeout(context.Background(), meta.(*Config).Timeout)
		defer cancel()
		for {
			err = client.APIBaseClient.Action(managementClient.ClusterType, monitoringActionEnable, clusterResource, monitoringInput, nil)
			if err == nil {
				return resourceRancher2ClusterRead(d, meta)
			}
			if !IsServerError(err) {
				return err
			}
			select {
			case <-time.After(rancher2RetriesWait * time.Second):
			case <-ctx.Done():
				break
			}
		}
	}

	return resourceRancher2ClusterRead(d, meta)
}

func resourceRancher2ClusterRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Cluster ID %s", d.Id())

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		cluster := &Cluster{}
		err = client.APIBaseClient.ByID(managementClient.ClusterType, d.Id(), cluster)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Cluster ID %s not found.", cluster.ID)
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(err)
		}

		clusterRegistrationToken, err := findClusterRegistrationToken(client, cluster.ID)
		if err != nil && !IsForbidden(err) {
			return resource.NonRetryableError(err)
		}

		defaultProjectID, systemProjectID, err := meta.(*Config).GetClusterSpecialProjectsID(cluster.ID)
		if err != nil && !IsForbidden(err) {
			return resource.NonRetryableError(err)
		}

		kubeConfig, err := getClusterKubeconfig(meta.(*Config), cluster.ID, d.Get("kube_config").(string))
		if err != nil && !IsForbidden(err) {
			return resource.NonRetryableError(err)
		}

		var monitoringInput *managementClient.MonitoringInput
		if len(cluster.Annotations[monitoringInputAnnotation]) > 0 {
			monitoringInput = &managementClient.MonitoringInput{}
			err = jsonToInterface(cluster.Annotations[monitoringInputAnnotation], monitoringInput)
			if err != nil {
				return resource.NonRetryableError(err)
			}

		}

		if err = flattenCluster(
			d,
			cluster,
			clusterRegistrationToken,
			kubeConfig,
			defaultProjectID,
			systemProjectID,
			monitoringInput); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
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
		"agentEnvVars":                       expandEnvVars(d.Get("agent_env_vars").([]interface{})),
		"description":                        d.Get("description").(string),
		"defaultPodSecurityPolicyTemplateId": d.Get("default_pod_security_policy_template_id").(string),
		"desiredAgentImage":                  d.Get("desired_agent_image").(string),
		"desiredAuthImage":                   d.Get("desired_auth_image").(string),
		"dockerRootDir":                      d.Get("docker_root_dir").(string),
		"fleetWorkspaceName":                 d.Get("fleet_workspace_name").(string),
		"enableClusterAlerting":              d.Get("enable_cluster_alerting").(bool),
		"enableClusterMonitoring":            d.Get("enable_cluster_monitoring").(bool),
		"enableNetworkPolicy":                &enableNetworkPolicy,
		"istioEnabled":                       d.Get("enable_cluster_istio").(bool),
		"localClusterAuthEndpoint":           expandClusterAuthEndpoint(d.Get("cluster_auth_endpoint").([]interface{})),
		"annotations":                        toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":                             toMapString(d.Get("labels").(map[string]interface{})),
	}

	// cluster_monitoring is not updated here. Setting old `enable_cluster_monitoring` value if it was updated
	if d.HasChange("enable_cluster_monitoring") {
		update["enableClusterMonitoring"] = !d.Get("enable_cluster_monitoring").(bool)
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

	replace := false
	switch driver := ToLower(d.Get("driver").(string)); driver {
	case clusterDriverAKS:
		aksConfig, err := expandClusterAKSConfig(d.Get("aks_config").([]interface{}), d.Get("name").(string))
		if err != nil {
			return err
		}
		update["azureKubernetesServiceConfig"] = aksConfig
	case ToLower(clusterDriverAKSV2):
		aksConfigV2 := expandClusterAKSConfigV2(d.Get("aks_config_v2").([]interface{}))
		update["aksConfig"] = aksConfigV2
	case clusterDriverEKS:
		eksConfig, err := expandClusterEKSConfig(d.Get("eks_config").([]interface{}), d.Get("name").(string))
		if err != nil {
			return err
		}
		update["amazonElasticContainerServiceConfig"] = eksConfig
	case ToLower(clusterDriverEKSV2):
		eksConfigV2 := expandClusterEKSConfigV2(d.Get("eks_config_v2").([]interface{}))
		update["eksConfig"] = fixClusterEKSConfigV2(d.Get("eks_config_v2").([]interface{}), structToMap(eksConfigV2))
	case clusterDriverGKE:
		gkeConfig, err := expandClusterGKEConfig(d.Get("gke_config").([]interface{}), d.Get("name").(string))
		if err != nil {
			return err
		}
		update["googleKubernetesEngineConfig"] = gkeConfig
	case ToLower(clusterDriverGKEV2):
		gkeConfig := expandClusterGKEConfigV2(d.Get("gke_config_v2").([]interface{}))
		update["gkeConfig"] = fixClusterGKEConfigV2(structToMap(gkeConfig))
	case clusterOKEKind:
		okeConfig, err := expandClusterOKEConfig(d.Get("oke_config").([]interface{}), d.Get("name").(string))
		if err != nil {
			return err
		}
		update["okeEngineConfig"] = okeConfig
	case ToLower(clusterDriverRKE):
		rkeConfig, err := expandClusterRKEConfig(d.Get("rke_config").([]interface{}), d.Get("name").(string))
		if err != nil {
			return err
		}
		update["rancherKubernetesEngineConfig"] = rkeConfig
		replace = d.HasChange("rke_config")
	case clusterDriverK3S:
		update["k3sConfig"] = expandClusterK3SConfig(d.Get("k3s_config").([]interface{}))
	case clusterDriverRKE2:
		update["rke2Config"] = expandClusterRKE2Config(d.Get("rke2_config").([]interface{}))
	}

	newCluster := &Cluster{}
	if replace {
		err = client.APIBaseClient.Replace(managementClient.ClusterType, cluster, update, newCluster)
	} else {
		err = client.APIBaseClient.Update(managementClient.ClusterType, cluster, update, newCluster)
	}
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active", "provisioning", "pending", "updating", "upgrading"},
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

	// cluster_monitoring is updated here
	if d.HasChange("enable_cluster_monitoring") || d.HasChange("cluster_monitoring_input") {
		clusterResource := &norman.Resource{
			ID:      newCluster.ID,
			Type:    newCluster.Type,
			Links:   newCluster.Links,
			Actions: newCluster.Actions,
		}
		enableMonitoring := d.Get("enable_cluster_monitoring").(bool)
		if enableMonitoring {
			monitoringInput := expandMonitoringInput(d.Get("cluster_monitoring_input").([]interface{}))
			if len(newCluster.Actions[monitoringActionEnable]) > 0 {
				err = client.APIBaseClient.Action(managementClient.ClusterType, monitoringActionEnable, clusterResource, monitoringInput, nil)
				if err != nil {
					return err
				}
			} else {
				monitorVersionChanged := false
				if d.HasChange("cluster_monitoring_input") {
					old, new := d.GetChange("cluster_monitoring_input")
					oldInput := old.([]interface{})
					oldInputLen := len(oldInput)
					oldVersion := ""
					if oldInputLen > 0 {
						oldRow, oldOK := oldInput[0].(map[string]interface{})
						if oldOK {
							oldVersion = oldRow["version"].(string)
						}
					}
					newInput := new.([]interface{})
					newInputLen := len(newInput)
					newVersion := ""
					if newInputLen > 0 {
						newRow, newOK := newInput[0].(map[string]interface{})
						if newOK {
							newVersion = newRow["version"].(string)
						}
					}
					if oldVersion != newVersion {
						monitorVersionChanged = true
					}
				}
				if monitorVersionChanged && monitoringInput != nil {
					err = updateClusterMonitoringApps(meta, d.Get("system_project_id").(string), monitoringInput.Version)
					if err != nil {
						return err
					}
				}
				err = client.APIBaseClient.Action(managementClient.ClusterType, monitoringActionEdit, clusterResource, monitoringInput, nil)
				if err != nil {
					return err
				}
			}
		} else if len(newCluster.Actions[monitoringActionDisable]) > 0 {
			err = client.APIBaseClient.Action(managementClient.ClusterType, monitoringActionDisable, clusterResource, nil, nil)
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
		if IsNotFound(err) || IsForbidden(err) {
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
	log.Printf("[TRACE] Finding cluster registration token for %s", clusterID)
	for i := range clusterRegistrationTokenNames {
		regTokenID := clusterID + ":" + clusterRegistrationTokenNames[i]
		for retry, retries := 1, 10; retry <= retries; retry++ {
			regToken, err := client.ClusterRegistrationToken.ByID(regTokenID)
			if err != nil {
				if !IsNotFound(err) {
					return nil, err
				}
				log.Printf("[TRACE] Cluster registration token %s not found for %s", regTokenID, clusterID)
				break
			}
			if (len(regToken.Command) > 0 && len(regToken.NodeCommand) > 0) || retry == retries {
				log.Printf("[INFO] Found existing cluster registration token for %s", clusterID)
				return regToken, nil
			}
			log.Printf("[DEBUG] Sleeping for 3 seconds before checking cluster registration token for %s", clusterID)
			time.Sleep(3 * time.Second)
		}
	}
	log.Printf("[TRACE] Cluster registration token not found for %s", clusterID)
	return createClusterRegistrationToken(client, clusterID)
}

func createClusterRegistrationToken(client *managementClient.Client, clusterID string) (*managementClient.ClusterRegistrationToken, error) {
	log.Printf("[DEBUG] Creating cluster registration token for %s", clusterID)

	regToken, err := expandClusterRegistationToken([]interface{}{}, clusterID)
	if err != nil {
		return nil, err
	}

	newRegToken, err := client.ClusterRegistrationToken.Create(regToken)
	if err != nil {
		if IsConflict(err) {
			log.Printf("[INFO] Found existing cluster registration token for %s", clusterID)
			regTokenID := clusterID + ":" + clusterRegistrationTokenName
			return client.ClusterRegistrationToken.ByID(regTokenID)
		}
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
	newRegToken, err = client.ClusterRegistrationToken.ByID(newRegToken.ID)
	if err != nil {
		return nil, err
	}
	log.Printf("[INFO] Created cluster registration token %s for %s", newRegToken.ID, clusterID)
	return newRegToken, nil
}

func isKubeConfigValid(c *Config, config string) (string, bool, error) {
	token, tokenValid, err := isKubeConfigTokenValid(c, config)
	if err != nil {
		return "", false, err
	}
	if !tokenValid {
		return "", false, nil
	}
	kubeconfig, err := clientcmd.RESTConfigFromKubeConfig([]byte(config))
	if err != nil {
		return "", false, fmt.Errorf("Checking Kubeconfig: %v", err)
	}
	_, err = kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		return token, false, nil
	}

	return token, true, nil
}

func isKubeConfigTokenValid(c *Config, config string) (string, bool, error) {
	token, err := getTokenFromKubeConfig(config)
	if err != nil {
		return "", false, fmt.Errorf("Getting Kubeconfig token: %v", err)
	}
	isValid, err := isTokenValid(c, splitTokenID(token))
	if err != nil {
		return "", false, fmt.Errorf("Checking Kubeconfig token: %v", err)
	}
	return token, isValid, nil
}

func replaceKubeConfigToken(c *Config, config, token string) (string, error) {
	if len(token) == 0 {
		return config, nil
	}
	kubeconfig, err := getObjFromKubeConfig(config)
	if err != nil {
		return "", fmt.Errorf("Getting K8s config object: %v", err)
	}
	if kubeconfig == nil || kubeconfig.AuthInfos == nil || len(kubeconfig.AuthInfos) == 0 {
		return config, nil
	}

	client, err := c.ManagementClient()
	if err != nil {
		return "", fmt.Errorf("Replacing cluster Kubeconfig token: %v", err)
	}
	removeToken, err := client.Token.ByID(splitTokenID(kubeconfig.AuthInfos[0].AuthInfo.Token))
	if err != nil {
		if !IsNotFound(err) && !IsForbidden(err) {
			return "", err
		}
	}

	err = client.Token.Delete(removeToken)
	if err != nil {
		return "", fmt.Errorf("Error removing Token: %s", err)
	}
	kubeconfig.AuthInfos[0].AuthInfo.Token = token
	return getKubeConfigFromObj(kubeconfig)
}

func getClusterKubeconfig(c *Config, id, origconfig string) (*managementClient.GenerateKubeConfigOutput, error) {
	action := "generateKubeconfig"
	cluster := &Cluster{}

	token, kubeValid, err := isKubeConfigValid(c, origconfig)
	if err != nil {
		return nil, fmt.Errorf("Getting cluster Kubeconfig: %v", err)
	}
	if kubeValid {
		return &managementClient.GenerateKubeConfigOutput{Config: origconfig}, nil
	}

	client, err := c.ManagementClient()
	if err != nil {
		return nil, fmt.Errorf("Getting cluster Kubeconfig: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	for {
		err = client.APIBaseClient.ByID(managementClient.ClusterType, id, cluster)
		if err != nil {
			if !IsNotFound(err) && !IsForbidden(err) && !IsServiceUnavailableError(err) {
				return nil, fmt.Errorf("Getting cluster Kubeconfig: %v", err)
			}
		} else if len(cluster.Actions[action]) > 0 {
			isRancher26, err := c.IsRancherVersionGreaterThanOrEqual("2.6.0")
			if err != nil {
				return nil, err
			}
			kubeConfig := &managementClient.GenerateKubeConfigOutput{}
			if isRancher26 && cluster.LocalClusterAuthEndpoint != nil && cluster.LocalClusterAuthEndpoint.Enabled {
				if connected, _, _ := c.isClusterConnected(cluster.ID); !connected {
					log.Printf("[WARN] Getting cluster Kubeconfig: kubeconfig is not yet available for cluster %s", cluster.Name)
					return kubeConfig, nil
				}
			}
			clusterResource := &norman.Resource{
				ID:      cluster.ID,
				Type:    cluster.Type,
				Links:   cluster.Links,
				Actions: cluster.Actions,
			}
			err = client.APIBaseClient.Action(managementClient.ClusterType, action, clusterResource, nil, kubeConfig)
			if err == nil {
				if isRancher26 && len(token) > 0 {
					newConfig, err := replaceKubeConfigToken(c, kubeConfig.Config, token)
					if err != nil {
						return nil, err
					}
					kubeConfig.Config = newConfig
				}
				return kubeConfig, nil
			}
			if !IsNotFound(err) && !IsForbidden(err) && !IsServiceUnavailableError(err) {
				return nil, fmt.Errorf("Getting cluster Kubeconfig: %v", err)
			}
		}
		select {
		case <-time.After(rancher2RetriesWait * time.Second):
		case <-ctx.Done():
			return nil, fmt.Errorf("Timeout getting cluster Kubeconfig: %v", err)
		}
	}
}

func updateClusterMonitoringApps(meta interface{}, systemProjectID, version string) error {
	cliProject, err := meta.(*Config).ProjectClient(systemProjectID)
	if err != nil {
		return err
	}

	filters := map[string]interface{}{
		"targetNamespace": clusterMonitoringV1Namespace,
	}

	listOpts := NewListOpts(filters)

	apps, err := cliProject.App.List(listOpts)
	if err != nil {
		return err
	}

	for _, a := range apps.Data {
		if a.Name == "cluster-monitoring" || a.Name == "monitoring-operator" {
			externalID := updateVersionExternalID(a.ExternalID, version)
			upgrade := &projectClient.AppUpgradeConfig{
				Answers:      a.Answers,
				ExternalID:   externalID,
				ForceUpgrade: true,
			}

			err = cliProject.App.ActionUpgrade(&a, upgrade)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
