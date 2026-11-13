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
					if err := validateClusterV2LocalAuthEndpointResourceDiff(d, newInterface); err != nil {
						return err
					}
					if clusterV2LocalAuthEndpointDiffEqual(oldInterface, newInterface) {
						d.Clear("local_auth_endpoint")
					} else {
						d.SetNew("local_auth_endpoint", newObj)
					}
				}
			}
			return customizeClusterV2LocalAuthEndpointCASync(d, i)
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

	if clusterV2LocalAuthEndpointShouldUseInternalCACerts(d) {
		caCert, err := getClusterV2LocalAuthEndpointCACert(newCluster.Status.ClusterName, func(clusterV1ID string) (string, error) {
			return getClusterCACert(meta.(*Config), clusterV1ID)
		})
		if err != nil {
			return err
		}
		if caCert != "" {
			newCluster, err = updateClusterV2LocalAuthEndpointCACerts(
				func() (*ClusterV2, error) {
					return getClusterV2ByID(meta.(*Config), newCluster.ID)
				},
				func(latest *ClusterV2) (*ClusterV2, error) {
					return updateClusterV2Once(meta.(*Config), latest.ID, latest)
				},
				caCert, rancher2RetriesWait*time.Second, meta.(*Config).Timeout)
			if err != nil {
				return err
			}
		}
	}

	return resourceRancher2ClusterV2ReadInternal(d, meta, false)
}

func resourceRancher2ClusterV2Read(d *schema.ResourceData, meta interface{}) error {
	return resourceRancher2ClusterV2ReadInternal(d, meta, true)
}

func resourceRancher2ClusterV2ReadInternal(d *schema.ResourceData, meta interface{}, detectInternalCADrift bool) error {
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
	err = setClusterV2LegacyData(d, meta.(*Config), true)
	if err != nil {
		return err
	}

	// use_internal_ca_certs has no server-side representation, so it must be
	// captured before flattenClusterV2 rewrites the local_auth_endpoint block
	// and reinjected afterward. It is never derived from cluster CA content
	// (see clusterV2LocalAuthEndpointUseInternalCACerts).
	useInternalCACerts := clusterV2LocalAuthEndpointUseInternalCACerts(d)
	syncRequired := false
	if detectInternalCADrift && useInternalCACerts && cluster.Spec.LocalClusterAuthEndpoint.Enabled {
		caCert, err := getClusterV2LocalAuthEndpointCACert(cluster.Status.ClusterName, func(clusterV1ID string) (string, error) {
			return getClusterCACert(meta.(*Config), clusterV1ID)
		})
		if err != nil {
			return err
		}
		syncRequired = clusterV2LocalAuthEndpointCASyncRequired(useInternalCACerts, cluster, caCert)
	}
	if err := flattenClusterV2(d, cluster); err != nil {
		return err
	}
	if err := setClusterV2LocalAuthEndpointUseInternalCACerts(d, useInternalCACerts); err != nil {
		return err
	}
	return d.Set("local_auth_endpoint_ca_sync_required", syncRequired)
}

func resourceRancher2ClusterV2Update(d *schema.ResourceData, meta interface{}) error {
	cluster, err := expandClusterV2(d)
	if err != nil {
		return err
	}

	if clusterV2LocalAuthEndpointShouldUseInternalCACerts(d) {
		clusterV1ID := d.Get("cluster_v1_id").(string)
		caCert, err := getClusterV2LocalAuthEndpointCACert(clusterV1ID, func(clusterV1ID string) (string, error) {
			return getClusterCACert(meta.(*Config), clusterV1ID)
		})
		if err != nil {
			return err
		}
		if caCert != "" {
			cluster.Spec.LocalClusterAuthEndpoint.CACerts = caCert
		} else {
			current, err := getClusterV2ByID(meta.(*Config), d.Id())
			if err != nil {
				return err
			}
			cluster.Spec.LocalClusterAuthEndpoint.CACerts = current.Spec.LocalClusterAuthEndpoint.CACerts
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
	return resourceRancher2ClusterV2ReadInternal(d, meta, false)
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

// updateClusterV2Once updates a cluster once and returns API conflicts to the
// caller. This lets CA reconciliation refetch the complete cluster before retrying.
func updateClusterV2Once(c *Config, id string, obj *ClusterV2) (*ClusterV2, error) {
	if c == nil {
		return nil, fmt.Errorf("Updating cluster V2: Provider config is nil")
	}
	if id == "" {
		return nil, fmt.Errorf("Updating cluster V2: Cluster V2 ID is empty")
	}
	if obj == nil {
		return nil, fmt.Errorf("Updating cluster V2: Cluster V2 is nil")
	}
	resp := &ClusterV2{}
	if err := c.updateObjectV2(rancher2DefaultLocalClusterID, id, clusterV2APIType, obj, resp); err != nil {
		return nil, err
	}
	return resp, nil
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

func setClusterV2LegacyData(d *schema.ResourceData, c *Config, generateKubeConfig bool) error {
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

	if generateKubeConfig {
		kubeConfig, err := getClusterKubeconfig(c, cluster.ID, d.Get("kube_config").(string))
		if err != nil {
			return fmt.Errorf(format, err)
		}
		d.Set("kube_config", kubeConfig.Config)
	}

	return nil
}

// clusterV2LocalAuthEndpointRawMap returns the single element of a
// local_auth_endpoint list ([]interface{} as read from *schema.ResourceData
// or a diff) as a map, or (nil, false) if the block is empty/absent.
func clusterV2LocalAuthEndpointRawMap(v []interface{}) (map[string]interface{}, bool) {
	if len(v) == 0 || v[0] == nil {
		return nil, false
	}
	m, ok := v[0].(map[string]interface{})
	return m, ok
}

// clusterV2LocalAuthEndpointUseInternalCACertsValue reads use_internal_ca_certs
// out of a raw local_auth_endpoint list. It never inspects ca_certs, enabled,
// or fqdn: use_internal_ca_certs has no server-side representation, so its
// value is only ever what was explicitly stored here.
func clusterV2LocalAuthEndpointUseInternalCACertsValue(v []interface{}) bool {
	m, ok := clusterV2LocalAuthEndpointRawMap(v)
	if !ok {
		return false
	}
	b, _ := m["use_internal_ca_certs"].(bool)
	return b
}

// clusterV2LocalAuthEndpointUseInternalCACerts reads the current
// use_internal_ca_certs value out of d. It is used both to read the planned
// value (Create/Update) and to capture the prior value before Read overwrites
// the local_auth_endpoint block (see setClusterV2LocalAuthEndpointUseInternalCACerts).
func clusterV2LocalAuthEndpointUseInternalCACerts(d *schema.ResourceData) bool {
	v, _ := d.Get("local_auth_endpoint").([]interface{})
	return clusterV2LocalAuthEndpointUseInternalCACertsValue(v)
}

// clusterV2LocalAuthEndpointShouldUseInternalCACerts reports whether the
// endpoint is enabled and configured to use the internal CA.
func clusterV2LocalAuthEndpointShouldUseInternalCACerts(d *schema.ResourceData) bool {
	v, _ := d.Get("local_auth_endpoint").([]interface{})
	m, ok := clusterV2LocalAuthEndpointRawMap(v)
	if !ok {
		return false
	}
	enabled, _ := m["enabled"].(bool)
	return enabled && clusterV2LocalAuthEndpointUseInternalCACertsValue(v)
}

// clusterV2LocalAuthEndpointCASyncRequired reports whether the endpoint CA
// differs from the available internal CA.
func clusterV2LocalAuthEndpointCASyncRequired(useInternalCACerts bool, cluster *ClusterV2, caCert string) bool {
	if !useInternalCACerts || cluster == nil || !cluster.Spec.LocalClusterAuthEndpoint.Enabled || caCert == "" {
		return false
	}
	return caCert != cluster.Spec.LocalClusterAuthEndpoint.CACerts
}

// getClusterV2LocalAuthEndpointCACert skips the CA fetch until Rancher assigns
// the management cluster ID.
func getClusterV2LocalAuthEndpointCACert(clusterV1ID string, fetch func(string) (string, error)) (string, error) {
	if clusterV1ID == "" {
		log.Printf("[INFO] Skipping local_auth_endpoint CA certificate fetch because cluster_v1_id is not available")
		return "", nil
	}
	log.Printf("[INFO] Fetching cluster %s CA certificate for local_auth_endpoint", clusterV1ID)
	return fetch(clusterV1ID)
}

// customizeClusterV2LocalAuthEndpointCASync schedules an update when Read
// detects internal CA drift for an enabled endpoint.
func customizeClusterV2LocalAuthEndpointCASync(d *schema.ResourceDiff, _ interface{}) error {
	syncRequired, _ := d.Get("local_auth_endpoint_ca_sync_required").(bool)
	if !syncRequired || !d.NewValueKnown("local_auth_endpoint.0.enabled") ||
		!d.NewValueKnown("local_auth_endpoint.0.use_internal_ca_certs") {
		return nil
	}
	v, _ := d.Get("local_auth_endpoint").([]interface{})
	m, ok := clusterV2LocalAuthEndpointRawMap(v)
	if !ok {
		return nil
	}
	enabled, _ := m["enabled"].(bool)
	useInternalCACerts, _ := m["use_internal_ca_certs"].(bool)
	if !enabled || !useInternalCACerts {
		return nil
	}
	return d.SetNewComputed("local_auth_endpoint_ca_sync_required")
}

// updateClusterV2LocalAuthEndpointCACerts refetches the complete cluster before
// each CA update attempt, so conflict retries preserve concurrent changes.
func updateClusterV2LocalAuthEndpointCACerts(fetch func() (*ClusterV2, error), update func(*ClusterV2) (*ClusterV2, error),
	caCert string, retryInterval, timeout time.Duration) (*ClusterV2, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	for {
		cluster, err := fetch()
		if err != nil {
			return nil, err
		}
		cluster.Spec.LocalClusterAuthEndpoint.CACerts = caCert
		updated, err := update(cluster)
		if err == nil {
			return updated, nil
		}
		if !IsConflict(err) {
			return nil, err
		}
		select {
		case <-time.After(retryInterval):
		case <-ctx.Done():
			return nil, fmt.Errorf("Timeout updating cluster V2 ID %s: %w", cluster.ID, err)
		}
	}
}

// setClusterV2LocalAuthEndpointUseInternalCACerts reinjects useInternalCACerts
// into the current local_auth_endpoint block. It must be called after
// flattenClusterV2 has repopulated enabled/fqdn/ca_certs from the API object,
// since d.Set on a nested list block resets any sub-field not present in the
// new value to its zero value. If Rancher omits the block, a preserved true
// value recreates it so the next apply can fetch and apply the internal CA.
//
// When useInternalCACerts is true, it also blanks ca_certs back to empty.
// ca_certs is Optional but not Computed, and the user's config never sets it
// while use_internal_ca_certs is true (they are mutually exclusive), so if
// the real CA value flattenClusterV2 just wrote were left in state, every
// subsequent plan would show a diff trying to clear it back to what config
// says (empty), and the plan would never converge. The real CA is still sent
// to Rancher when available; only its reflection in Terraform state is
// intentionally not tracked.
func setClusterV2LocalAuthEndpointUseInternalCACerts(d *schema.ResourceData, useInternalCACerts bool) error {
	v, _ := d.Get("local_auth_endpoint").([]interface{})
	m, ok := clusterV2LocalAuthEndpointRawMap(v)
	if !ok {
		if !useInternalCACerts {
			return nil
		}
		m = map[string]interface{}{}
	}
	m["use_internal_ca_certs"] = useInternalCACerts
	if useInternalCACerts {
		m["ca_certs"] = ""
	}
	return d.Set("local_auth_endpoint", []interface{}{m})
}

// validateClusterV2LocalAuthEndpointResourceDiff validates each known endpoint
// field independently and defers checks that depend on unknown values.
func validateClusterV2LocalAuthEndpointResourceDiff(d *schema.ResourceDiff, newInterface []interface{}) error {
	if !d.NewValueKnown("local_auth_endpoint.0.use_internal_ca_certs") {
		return nil
	}
	m, ok := clusterV2LocalAuthEndpointRawMap(newInterface)
	if !ok {
		return nil
	}
	useInternal, _ := m["use_internal_ca_certs"].(bool)
	if !useInternal {
		return nil
	}
	if d.NewValueKnown("local_auth_endpoint.0.ca_certs") {
		if caCerts, _ := m["ca_certs"].(string); caCerts != "" {
			return fmt.Errorf(`only one of "ca_certs" or "use_internal_ca_certs" can be set`)
		}
	}
	if d.NewValueKnown("local_auth_endpoint.0.fqdn") {
		if fqdn, _ := m["fqdn"].(string); fqdn == "" {
			return fmt.Errorf(`"fqdn" is required in "local_auth_endpoint" when "use_internal_ca_certs" is true`)
		}
	}
	return nil
}

// clusterV2LocalAuthEndpointDiffEqual reports whether old and new
// local_auth_endpoint values are equivalent, including use_internal_ca_certs.
// CustomizeDiff uses this to decide whether to clear a spurious diff; a
// change to use_internal_ca_certs alone must always be reported as unequal.
func clusterV2LocalAuthEndpointDiffEqual(oldInterface, newInterface []interface{}) bool {
	oldConfig := expandClusterV2LocalAuthEndpoint(oldInterface)
	newConfig := expandClusterV2LocalAuthEndpoint(newInterface)
	if !reflect.DeepEqual(oldConfig, newConfig) {
		return false
	}

	oldVal := clusterV2LocalAuthEndpointUseInternalCACertsValue(oldInterface)
	newVal := clusterV2LocalAuthEndpointUseInternalCACertsValue(newInterface)
	return oldVal == newVal
}

// getClusterCACert fetches and returns the v3 management cluster's CA
// certificate, base64-decoded if necessary. It performs no retries and does
// not treat an empty result as an error.
func getClusterCACert(c *Config, clusterV1ID string) (string, error) {
	if c == nil {
		return "", fmt.Errorf("Getting cluster CA cert: Provider config is nil")
	}
	if clusterV1ID == "" {
		return "", fmt.Errorf("Getting cluster CA cert: cluster_v1_id is empty")
	}
	client, err := c.ManagementClient()
	if err != nil {
		return "", fmt.Errorf("Getting cluster CA cert: %w", err)
	}
	cluster := &Cluster{}
	err = client.APIBaseClient.ByID(managementClient.ClusterType, clusterV1ID, cluster)
	if err != nil {
		return "", fmt.Errorf("Getting cluster CA cert: %w", err)
	}
	return decodeCACertIfBase64(cluster.CACert), nil
}
