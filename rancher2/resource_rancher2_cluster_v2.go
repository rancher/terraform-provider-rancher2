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
)

func resourceRancher2ClusterV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2ClusterV2Create,
		Read:   resourceRancher2ClusterV2Read,
		Update: resourceRancher2ClusterV2Update,
		Delete: resourceRancher2ClusterV2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2ClusterV2Import,
		},
		Schema:        clusterV2Fields(),
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceRancher2ClusterV2Resource().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceRancher2ClusterV2StateUpgradeV0,
				Version: 0,
			},
		},
		CustomizeDiff: func(d *schema.ResourceDiff, i interface{}) error {
			if d.HasChange("rke_config") {
				oldObj, newObj := d.GetChange("rke_config")
				oldInterface, oldOk := oldObj.([]interface{})
				newInterface, newOk := newObj.([]interface{})
				if oldOk && newOk && len(newInterface) > 0 {
					oldConfig := expandClusterV2RKEConfig(oldInterface)
					newConfig := expandClusterV2RKEConfig(newInterface)
					if reflect.DeepEqual(oldConfig, newConfig) {
						d.Clear("rke_config")
					} else {
						d.SetNew("rke_config", flattenClusterV2RKEConfig(newConfig))
					}
				}
			}
			if d.HasChange("local_auth_endpoint") {
				oldObj, newObj := d.GetChange("local_auth_endpoint")
				oldInterface, oldOk := oldObj.([]interface{})
				newInterface, newOk := newObj.([]interface{})
				if oldOk && newOk && len(newInterface) > 0 {
					if m, ok := newInterface[0].(map[string]interface{}); ok {
						if ca, ok := m["ca_certs"].(string); ok && ca != "" {
							if use, ok := m["use_internal_ca_certs"].(bool); ok && use {
								return fmt.Errorf("only one of \"ca_certs\" or \"use_internal_ca_certs\" can be set")
							}
						}
					}
					oldConfig := expandClusterV2LocalAuthEndpoint(oldInterface)
					newConfig := expandClusterV2LocalAuthEndpoint(newInterface)
					oldUse := false
					newUse := false
					if len(oldInterface) > 0 {
						if m, ok := oldInterface[0].(map[string]interface{}); ok {
							if v, ok := m["use_internal_ca_certs"].(bool); ok {
								oldUse = v
							}
						}
					}
					if len(newInterface) > 0 {
						if m, ok := newInterface[0].(map[string]interface{}); ok {
							if v, ok := m["use_internal_ca_certs"].(bool); ok {
								newUse = v
							}
						}
					}
					if reflect.DeepEqual(oldConfig, newConfig) && oldUse == newUse {
						d.Clear("local_auth_endpoint")
					} else {
						d.SetNew("local_auth_endpoint", newObj)
					}
				}
			}
			return nil
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},
	}
}

func resourceRancher2ClusterV2Resource() *schema.Resource {
	return &schema.Resource{
		Schema: clusterV2FieldsV0(),
	}
}

func resourceRancher2ClusterV2StateUpgradeV0(rawState map[string]any, meta interface{}) (map[string]any, error) {
	if rkeConfigs, ok := rawState["rke_config"].([]any); ok && len(rkeConfigs) > 0 {
		for i := range rkeConfigs {
			if rkeConfig, ok := rkeConfigs[i].(map[string]any); ok && len(rkeConfig) > 0 {
				if machineSelectorConfigs, ok := rkeConfig["machine_selector_config"].([]any); ok && len(machineSelectorConfigs) > 0 {

					// upgrade all machine selector configs
					for m := range machineSelectorConfigs {
						if machineSelectorConfig, ok := machineSelectorConfigs[m].(map[string]any); ok && len(machineSelectorConfig) > 0 {

							// machine selector config data found. Migrate state from map -> string
							if config, ok := machineSelectorConfig["config"].(map[string]any); ok {
								newValue := ""
								if conf, err := mapInterfaceToYAML(config); err == nil {
									newValue = conf
								}
								rawState["rke_config"].([]interface{})[i].(map[string]any)["machine_selector_config"].([]any)[m].(map[string]any)["config"] = newValue
							}
						}
					}
				}
			}
		}
	}
	return rawState, nil
}

func resourceRancher2ClusterV2Create(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	cluster, err := expandClusterV2(d)
	if err != nil {
		return err
	}

	useInternal := false
	if v, ok := d.Get("local_auth_endpoint").([]interface{}); ok && len(v) > 0 {
		if m, ok := v[0].(map[string]interface{}); ok {
			if b, ok := m["use_internal_ca_certs"].(bool); ok {
				useInternal = b
			}
		}
	}

	log.Printf("[INFO] Creating Cluster V2 %s", name)

	newCluster, err := createClusterV2(meta.(*Config), cluster)
	if err != nil {
		return err
	}
	d.SetId(newCluster.ID)
	newCluster, err = waitForClusterV2State(meta.(*Config), newCluster.ID, clusterV2CreatedCondition, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return err
	}

	// Waiting for cluster v2 active if it has machine pools defined
	if newCluster.Spec.RKEConfig != nil && newCluster.Spec.RKEConfig.MachinePools != nil && len(newCluster.Spec.RKEConfig.MachinePools) > 0 {
		newCluster, err = waitForClusterV2State(meta.(*Config), newCluster.ID, clusterV2ActiveCondition, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return err
		}
	}

	if useInternal {
		caCert, err := getClusterCACert(meta.(*Config), newCluster.Status.ClusterName)
		if err != nil {
			return err
		}
		newCluster.Spec.LocalClusterAuthEndpoint.CACerts = caCert
		_, err = updateClusterV2(meta.(*Config), newCluster.ID, newCluster)
		if err != nil {
			return err
		}
	}

	return resourceRancher2ClusterV2Read(d, meta)
}

func resourceRancher2ClusterV2Read(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Cluster V2 %s", d.Id())

	cluster, err := getClusterV2ByID(meta.(*Config), d.Id())
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) || IsNotAccessibleByID(err) {
			log.Printf("[INFO] Cluster V2 %s not found", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}
	d.Set("cluster_v1_id", cluster.Status.ClusterName)
	err = setClusterV2LegacyData(d, meta.(*Config))
	if err != nil {
		return err
	}
	if err := flattenClusterV2(d, cluster); err != nil {
		return err
	}
	return setClusterV2LocalAuthEndpointInternalFlag(d, meta.(*Config), cluster)
}

func resourceRancher2ClusterV2Update(d *schema.ResourceData, meta interface{}) error {
	cluster, err := expandClusterV2(d)
	if err != nil {
		return err
	}

	if v, ok := d.Get("local_auth_endpoint").([]interface{}); ok && len(v) > 0 {
		if m, ok := v[0].(map[string]interface{}); ok {
			if b, ok := m["use_internal_ca_certs"].(bool); ok && b {
				caCert, err := getClusterCACert(meta.(*Config), d.Get("cluster_v1_id").(string))
				if err != nil {
					return err
				}
				cluster.Spec.LocalClusterAuthEndpoint.CACerts = caCert
			}
		}
	}

	log.Printf("[INFO] Updating Cluster V2 %s", d.Id())

	newCluster, err := updateClusterV2(meta.(*Config), d.Id(), cluster)
	if err != nil {
		return err
	}
	// Waiting for cluster v2 active if it has machine pools defined
	if newCluster.Spec.RKEConfig != nil && newCluster.Spec.RKEConfig.MachinePools != nil && len(newCluster.Spec.RKEConfig.MachinePools) > 0 {
		newCluster, err = waitForClusterV2State(meta.(*Config), newCluster.ID, clusterV2ActiveCondition, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return err
		}
	}
	return resourceRancher2ClusterV2Read(d, meta)
}

func resourceRancher2ClusterV2Delete(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	log.Printf("[INFO] Deleting Cluster V2 %s", name)

	cluster, err := getClusterV2ByID(meta.(*Config), d.Id())
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			d.SetId("")
			return nil
		}
	}
	err = deleteClusterV2(meta.(*Config), cluster)
	if err != nil {
		return err
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"removed"},
		Refresh:    clusterV2StateRefreshFunc(meta, cluster.ID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for cluster (%s) to be removed: %w", cluster.ID, waitErr)
	}

	// Rancher deletes the Management v3 Cluster under the hook, we should wait for the deletion to success
	v1ClusterName := cluster.Status.ClusterName
	if v1ClusterName != "" {
		client, err := meta.(*Config).ManagementClient()
		if err != nil {
			return err
		}
		stateConf = &resource.StateChangeConf{
			Pending:    []string{"removing"},
			Target:     []string{"removed"},
			Refresh:    clusterStateRefreshFunc(client, v1ClusterName),
			Timeout:    d.Timeout(schema.TimeoutDelete),
			Delay:      1 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		_, waitErr = stateConf.WaitForState()
		if waitErr != nil {
			return fmt.Errorf("[ERROR] waiting for cluster (%s) to be removed: %w", cluster.ID, waitErr)
		}
	}

	d.SetId("")
	return nil
}

// clusterV2StateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher Cluster v2.
func clusterV2StateRefreshFunc(meta interface{}, objID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := getClusterV2ByID(meta.(*Config), objID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) || IsNotAccessibleByID(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}
		return obj, "active", nil
	}
}

// Rancher2 Cluster V2 API CRUD functions
func createClusterV2(c *Config, obj *ClusterV2) (*ClusterV2, error) {
	if c == nil {
		return nil, fmt.Errorf("Creating cluster V2: Provider config is nil")
	}
	if obj == nil {
		return nil, fmt.Errorf("Creating cluster V2: Cluster V2 is nil")
	}
	resp := &ClusterV2{}
	err := c.createObjectV2(rancher2DefaultLocalClusterID, clusterV2APIType, obj, resp)
	if err != nil {
		return nil, fmt.Errorf("Creating cluster V2: %s", err)
	}
	return resp, nil
}

func deleteClusterV2(c *Config, obj *ClusterV2) error {
	if c == nil {
		return fmt.Errorf("Deleting cluster V2: Provider config is nil")
	}
	if obj == nil {
		return fmt.Errorf("Deleting cluster V2: Cluster V2 is nil")
	}
	resource := &norman.Resource{
		ID:      obj.ID,
		Type:    clusterV2APIType,
		Links:   obj.Links,
		Actions: obj.Actions,
	}
	return c.deleteObjectV2(rancher2DefaultLocalClusterID, resource)
}

func getClusterV2ByID(c *Config, id string) (*ClusterV2, error) {
	if c == nil {
		return nil, fmt.Errorf("Getting cluster V2: Provider config is nil")
	}
	if len(id) == 0 {
		return nil, fmt.Errorf("Getting cluster V2: Cluster V2 ID is empty")
	}
	resp := &ClusterV2{}
	err := c.getObjectV2ByID(rancher2DefaultLocalClusterID, id, clusterV2APIType, resp)
	if err != nil {
		if !IsServerError(err) && !IsNotFound(err) && !IsForbidden(err) {
			return nil, fmt.Errorf("Getting cluster V2: %w", err)
		}
		return nil, err
	}
	return resp, nil
}

func updateClusterV2(c *Config, id string, obj *ClusterV2) (*ClusterV2, error) {
	if c == nil {
		return nil, fmt.Errorf("Updating cluster V2: Provider config is nil")
	}
	if len(id) == 0 {
		return nil, fmt.Errorf("Updating cluster V2: Cluster V2 ID is empty")
	}
	if obj == nil {
		return nil, fmt.Errorf("Updating cluster V2: Cluster V2 is nil")
	}
	resp := &ClusterV2{}
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	for {
		err := c.updateObjectV2(rancher2DefaultLocalClusterID, id, clusterV2APIType, obj, resp)
		if err == nil {
			return resp, err
		}
		if !IsServerError(err) && !IsUnknownSchemaType(err) && !IsConflict(err) {
			return nil, err
		}
		if IsConflict(err) {
			// Read cluster again and update ObjectMeta.ResourceVersion before retry
			newObj := &ClusterV2{}
			err = c.getObjectV2ByID(rancher2DefaultLocalClusterID, id, clusterV2APIType, newObj)
			if err != nil {
				return nil, err
			}
			obj.ObjectMeta.ResourceVersion = newObj.ObjectMeta.ResourceVersion
		}
		select {
		case <-time.After(rancher2RetriesWait * time.Second):
		case <-ctx.Done():
			return nil, fmt.Errorf("Timeout updating cluster V2 ID %s: %w", id, err)
		}
	}
}

func waitForClusterV2State(c *Config, id, state string, interval time.Duration) (*ClusterV2, error) {
	if id == "" || state == "" {
		return nil, fmt.Errorf("Cluster V2 ID and/or condition is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), interval)
	defer cancel()
	for {
		obj, err := getClusterV2ByID(c, id)
		if err != nil {
			log.Printf("[DEBUG] Retrying on error Refreshing Cluster V2 %s: %v", id, err)
			if !IsNotFound(err) && !IsForbidden(err) && !IsNotAccessibleByID(err) {
				return nil, fmt.Errorf("Getting cluster V2 ID (%s): %w", id, err)
			}
			if IsNotAccessibleByID(err) {
				// Restarting clients to update RBAC
				c.RestartClients()
			}
		}
		if obj != nil {
			for i := range obj.Status.Conditions {
				if obj.Status.Conditions[i].Type == state {
					// Status of the condition, one of True, False, Unknown.
					if obj.Status.Conditions[i].Status == "Unknown" {
						break
					}
					if obj.Status.Conditions[i].Status == "True" {
						return obj, nil
					}
					// When cluster condition is false, retrying if it has been updated for last rancher2WaitFalseCond seconds
					lastUpdate, err := time.Parse(time.RFC3339, obj.Status.Conditions[i].LastUpdateTime)
					if err == nil && time.Since(lastUpdate) < rancher2WaitFalseCond*time.Second {
						break
					}
					return nil, fmt.Errorf("Cluster V2 ID %s: %s", id, obj.Status.Conditions[i].Message)
				}
			}
		}
		select {
		case <-time.After(rancher2RetriesWait * time.Second):
		case <-ctx.Done():
			return nil, fmt.Errorf("Timeout waiting for cluster V2 ID %s", id)
		}
	}
}

func setClusterV2LegacyData(d *schema.ResourceData, c *Config) error {
	format := "Setting cluster V2 legacy data: %w"

	if c == nil {
		return fmt.Errorf("Setting cluster V2 legacy data: Provider config is nil")
	}
	clusterV1ID := d.Get("cluster_v1_id").(string)
	if len(clusterV1ID) == 0 {
		return fmt.Errorf("Setting cluster V2 legacy data: cluster_v1_id is empty")
	}

	client, err := c.ManagementClient()
	if err != nil {
		return fmt.Errorf(format, err)
	}

	cluster := &Cluster{}
	err = client.APIBaseClient.ByID(managementClient.ClusterType, clusterV1ID, cluster)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Cluster ID %s not found.", cluster.ID)
			return nil
		}
		return fmt.Errorf(format, err)
	}

	clusterRegistrationToken, err := findClusterRegistrationToken(client, cluster.ID)
	if err != nil && !IsForbidden(err) {
		return fmt.Errorf(format, err)
	}
	regToken, _ := flattenClusterRegistrationToken(clusterRegistrationToken)
	err = d.Set("cluster_registration_token", regToken)
	if err != nil {
		return fmt.Errorf(format, err)
	}

	kubeConfig, err := getClusterKubeconfig(c, cluster.ID, d.Get("kube_config").(string))
	if err != nil {
		return fmt.Errorf(format, err)
	}
	d.Set("kube_config", kubeConfig.Config)

	return nil
}

func getClusterCACert(c *Config, clusterV1ID string) (string, error) {
	if c == nil {
		return "", fmt.Errorf("provider config is nil")
	}
	if clusterV1ID == "" {
		return "", fmt.Errorf("cluster_v1_id is empty")
	}
	client, err := c.ManagementClient()
	if err != nil {
		return "", err
	}
	cluster := &Cluster{}
	err = client.APIBaseClient.ByID(managementClient.ClusterType, clusterV1ID, cluster)
	if err != nil {
		return "", err
	}
	return decodeCACertIfBase64(cluster.CACert), nil
}

func setClusterV2LocalAuthEndpointInternalFlag(d *schema.ResourceData, c *Config, cluster *ClusterV2) error {
	if cluster == nil || c == nil {
		return fmt.Errorf("setting local auth endpoint internal flag: missing data")
	}
	lae := cluster.Spec.LocalClusterAuthEndpoint
	useInternal := false
	if cluster.Status.ClusterName != "" && lae.CACerts != "" {
		caCert, err := getClusterCACert(c, cluster.Status.ClusterName)
		if err != nil {
			return err
		}
		if lae.CACerts == caCert {
			useInternal = true
		}
	}
	if v, ok := d.Get("local_auth_endpoint").([]interface{}); ok && len(v) > 0 {
		if m, ok := v[0].(map[string]interface{}); ok {
			m["use_internal_ca_certs"] = useInternal
			d.Set("local_auth_endpoint", []interface{}{m})
		}
	}
	return nil
}
