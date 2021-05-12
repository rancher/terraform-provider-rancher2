package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRancher2AppV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2AppV2Create,
		Read:   resourceRancher2AppV2Read,
		Update: resourceRancher2AppV2Update,
		Delete: resourceRancher2AppV2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2AppV2Import,
		},
		Schema: appV2Fields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2AppV2Create(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	repoName := d.Get("repo_name").(string)
	chartName := d.Get("chart_name").(string)
	chartVersion := d.Get("chart_version").(string)

	log.Printf("[INFO] Creating App V2 %s at cluster ID %s", name, clusterID)

	active, cluster, err := meta.(*Config).isClusterActive(clusterID)
	if err != nil {
		return err
	}
	if !active {
		return fmt.Errorf("[ERROR] creating App V2 %s at cluster ID %s: Cluster is not active", name, clusterID)
	}
	d.Set("cluster_name", cluster.Name)

	systemDefaultRegistry, err := meta.(*Config).GetSettingV2ByID(appV2DefaultRegistryID)
	if err != nil {
		return err
	}
	d.Set("system_default_registry", systemDefaultRegistry.Value)

	repo, chartInfo, err := meta.(*Config).InfoAppV2(clusterID, repoName, chartName, chartVersion)
	if err != nil {
		return err
	}

	chartInstallAction, err := expandChartInstallActionV2(d, chartInfo)
	if err != nil {
		return err
	}

	chartOperation, err := meta.(*Config).InstallAppV2(clusterID, repo, chartInstallAction)
	if err != nil {
		return err
	}
	err = appV2OperationWait(meta, clusterID, chartOperation.OperationNamespace+"/"+chartOperation.OperationName)
	if err != nil {
		return fmt.Errorf("[ERROR] installing App V2: %s", err)
	}
	d.SetId(clusterID + appV2ClusterIDsep + chartInstallAction.Namespace + "/" + d.Get("name").(string))

	return resourceRancher2AppV2Read(d, meta)
}

func resourceRancher2AppV2Read(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	log.Printf("[INFO] Refreshing App V2 %s at %s", name, clusterID)

	if clusterName, ok := d.Get("cluster_name").(string); !ok || len(clusterName) == 0 {
		cluster, err := meta.(*Config).GetClusterByID(clusterID)
		if err != nil {
			return err
		}
		d.Set("cluster_name", cluster.Name)
	}
	if systemDefaultRegistry, ok := d.Get("system_default_registry").(string); !ok || len(systemDefaultRegistry) == 0 {
		systemDefaultRegistry, err := meta.(*Config).GetSettingV2ByID(appV2DefaultRegistryID)
		if err != nil {
			return err
		}
		d.Set("system_default_registry", systemDefaultRegistry.Value)
	}
	_, rancherID := splitID(d.Id())
	app, err := meta.(*Config).GetAppV2ByID(clusterID, rancherID)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] App V2 %s not found at %s", name, clusterID)
			d.SetId("")
			return nil
		}
		return err
	}
	return flattenAppV2(d, app)
}

func resourceRancher2AppV2Update(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	repoName := d.Get("repo_name").(string)
	chartName := d.Get("chart_name").(string)
	chartVersion := d.Get("chart_version").(string)
	log.Printf("[INFO] Updating App V2 %s at %s", name, clusterID)

	repo, chartInfo, err := meta.(*Config).InfoAppV2(clusterID, repoName, chartName, chartVersion)
	if err != nil {
		return err
	}
	chartUpgradeAction, err := expandChartUpgradeActionV2(d, chartInfo)
	if err != nil {
		return err
	}

	chartOperation, err := meta.(*Config).UpgradeAppV2(clusterID, repo, chartUpgradeAction)
	if err != nil {
		return err
	}
	err = appV2OperationWait(meta, clusterID, chartOperation.OperationNamespace+"/"+chartOperation.OperationName)
	if err != nil {
		return fmt.Errorf("[ERROR] upgrading App V2: %s", err)
	}
	return resourceRancher2AppV2Read(d, meta)
}

func resourceRancher2AppV2Delete(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	log.Printf("[INFO] Deleting App V2 %s at %s", name, clusterID)

	_, rancherID := splitID(d.Id())
	app, err := meta.(*Config).GetAppV2ByID(clusterID, rancherID)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] App V2 %s not found at %s", name, clusterID)
			d.SetId("")
			return nil
		}
		return err
	}
	err = meta.(*Config).DeleteAppV2(clusterID, app)
	if err != nil {
		return fmt.Errorf("Error removing App V2 %s: %s", name, err)
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"removed"},
		Refresh:    appV2StateRefreshFunc(meta, clusterID, app.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for app (%s) to be deleted: %s", app.ID, waitErr)
	}
	if app.Spec.Chart.Metadata != nil && app.Spec.Chart.Metadata.Annotations != nil && len(app.Spec.Chart.Metadata.Annotations) > 0 && len(app.Spec.Chart.Metadata.Annotations["catalog.cattle.io/auto-install"]) > 0 {
		namespace := d.Get("namespace").(string)
		if len(app.Spec.Chart.Metadata.Annotations["catalog.cattle.io/namespace"]) > 0 {
			namespace = app.Spec.Chart.Metadata.Annotations["catalog.cattle.io/namespace"]
		}
		chartAuto := splitBySep(app.Spec.Chart.Metadata.Annotations["catalog.cattle.io/auto-install"], "=")
		if len(chartAuto) != 2 {
			return fmt.Errorf("bad format on chart annotation catalog.cattle.io/auto-install: %s", app.Spec.Chart.Metadata.Annotations["catalog.cattle.io/auto-install"])
		}
		name := chartAuto[0]
		app, err = meta.(*Config).GetAppV2ByID(clusterID, namespace+"/"+name)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return nil
			}
			return err
		}
		err = meta.(*Config).DeleteAppV2(clusterID, app)
		if err != nil {
			return fmt.Errorf("Error removing App V2 %s: %s", name, err)
		}
		stateConf = &resource.StateChangeConf{
			Pending:    []string{},
			Target:     []string{"removed"},
			Refresh:    appV2StateRefreshFunc(meta, clusterID, app.ID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      1 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		_, waitErr = stateConf.WaitForState()
		if waitErr != nil {
			return fmt.Errorf("[ERROR] waiting for app (%s) to be deleted: %s", app.ID, waitErr)
		}

	}
	return nil
}

// appV2StateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher App.
func appV2StateRefreshFunc(meta interface{}, clusterID, appID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := meta.(*Config).GetAppV2ByID(clusterID, appID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}
		if obj.Status.Summary.Error {
			return nil, "", fmt.Errorf("%s", obj.Status.Summary.State)
		}
		return obj, obj.Status.Summary.State, nil
	}
}

func appV2OperationWait(meta interface{}, clusterID, opID string) error {
	for obj, err := meta.(*Config).GetAppV2OperationByID(clusterID, opID); ; obj, err = meta.(*Config).GetAppV2OperationByID(clusterID, opID) {
		if err != nil {
			return err
		}
		if metadata, ok := obj["metadata"].(map[string]interface{}); ok && len(metadata) > 0 {
			if state, ok := metadata["state"].(map[string]interface{}); ok && len(state) > 0 {
				if transitioning, ok := state["transitioning"].(bool); ok && !transitioning {
					if opError, ok := state["error"].(bool); ok && opError {
						message, err := meta.(*Config).GetAppV2OperationLogs(clusterID, obj)
						if err != nil {
							return fmt.Errorf("%s: %s", state["message"], err)
						}
						return fmt.Errorf("%s", message)
					}
					return nil
				}

			}
		}
		time.Sleep(5 * time.Second)
	}
}
