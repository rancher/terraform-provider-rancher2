package rancher2

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rancher/norman/types"
	types2 "github.com/rancher/rancher/pkg/api/steve/catalog/types"
)

func resourceRancher2AppV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2AppV2Create,
		ReadContext:   resourceRancher2AppV2Read,
		UpdateContext: resourceRancher2AppV2Update,
		DeleteContext: resourceRancher2AppV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRancher2AppV2Import,
		},
		Schema: appV2Fields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2AppV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	repoName := d.Get("repo_name").(string)
	chartName := d.Get("chart_name").(string)
	chartVersion := d.Get("chart_version").(string)

	log.Printf("[INFO] Creating App V2 %s at cluster ID %s", name, clusterID)

	active, cluster, err := meta.(*Config).isClusterActive(clusterID)
	if err != nil {
		return diag.FromErr(err)
	}
	if !active {
		return diag.Errorf("[ERROR] creating App V2 %s at cluster ID %s: Cluster is not active", name, clusterID)
	}
	d.Set("cluster_name", cluster.Name)

	systemDefaultRegistry, err := meta.(*Config).GetSettingV2ByID(appV2DefaultRegistryID)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("system_default_registry", systemDefaultRegistry.Value)

	repo, chartInfo, err := infoAppV2(meta.(*Config), clusterID, repoName, chartName, chartVersion)
	if err != nil {
		return diag.FromErr(err)
	}

	chartInstallAction, err := expandChartInstallActionV2(d, chartInfo)
	if err != nil {
		return diag.FromErr(err)
	}

	chartOperation, err := createAppV2(meta.(*Config), clusterID, repo, chartInstallAction)
	if err != nil {
		return diag.FromErr(err)
	}
	err = appV2OperationWait(meta, clusterID, chartOperation.OperationNamespace+"/"+chartOperation.OperationName)
	if err != nil {
		return diag.Errorf("[ERROR] installing App V2: %s", err)
	}
	d.SetId(clusterID + appV2ClusterIDsep + chartInstallAction.Namespace + "/" + d.Get("name").(string))

	return resourceRancher2AppV2Read(ctx, d, meta)
}

func resourceRancher2AppV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	log.Printf("[INFO] Refreshing App V2 %s at %s", name, clusterID)

	return diag.FromErr(retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		if clusterName, ok := d.Get("cluster_name").(string); !ok || len(clusterName) == 0 {
			cluster, err := meta.(*Config).GetClusterByID(clusterID)
			if err != nil {
				return retry.NonRetryableError(err)
			}
			d.Set("cluster_name", cluster.Name)
		}
		if systemDefaultRegistry, ok := d.Get("system_default_registry").(string); !ok || len(systemDefaultRegistry) == 0 {
			systemDefaultRegistry, err := meta.(*Config).GetSettingV2ByID(appV2DefaultRegistryID)
			if err != nil {
				return retry.NonRetryableError(err)
			}
			d.Set("system_default_registry", systemDefaultRegistry.Value)
		}
		_, rancherID := splitID(d.Id())
		app, err := getAppV2ByID(meta.(*Config), clusterID, rancherID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] App V2 %s not found at %s", name, clusterID)
				d.SetId("")
				return nil
			}
			return retry.NonRetryableError(err)
		}

		if err = flattenAppV2(d, app); err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	}))
}

func resourceRancher2AppV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	repoName := d.Get("repo_name").(string)
	chartName := d.Get("chart_name").(string)
	chartVersion := d.Get("chart_version").(string)
	log.Printf("[INFO] Updating App V2 %s at %s", name, clusterID)

	repo, chartInfo, err := infoAppV2(meta.(*Config), clusterID, repoName, chartName, chartVersion)
	if err != nil {
		return diag.FromErr(err)
	}
	chartUpgradeAction, err := expandChartUpgradeActionV2(d, chartInfo)
	if err != nil {
		return diag.FromErr(err)
	}

	chartOperation, err := upgradeAppV2(meta.(*Config), clusterID, repo, chartUpgradeAction)
	if err != nil {
		return diag.FromErr(err)
	}
	err = appV2OperationWait(meta, clusterID, chartOperation.OperationNamespace+"/"+chartOperation.OperationName)
	if err != nil {
		return diag.Errorf("[ERROR] upgrading App V2: %s", err)
	}
	return resourceRancher2AppV2Read(ctx, d, meta)
}

func resourceRancher2AppV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	log.Printf("[INFO] Deleting App V2 %s at %s", name, clusterID)

	_, rancherID := splitID(d.Id())
	app, err := getAppV2ByID(meta.(*Config), clusterID, rancherID)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] App V2 %s not found at %s", name, clusterID)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	err = deleteAppV2(meta.(*Config), clusterID, app)
	if err != nil {
		return diag.Errorf("Error removing App V2 %s: %s", name, err)
	}
	stateConf := &retry.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"removed"},
		Refresh:    appV2StateRefreshFunc(meta, clusterID, app.ID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf("[ERROR] waiting for app (%s) to be deleted: %s", app.ID, waitErr)
	}
	if app.Spec.Chart.Metadata != nil && app.Spec.Chart.Metadata.Annotations != nil && len(app.Spec.Chart.Metadata.Annotations) > 0 && len(app.Spec.Chart.Metadata.Annotations["catalog.cattle.io/auto-install"]) > 0 {
		namespace := d.Get("namespace").(string)
		if len(app.Spec.Chart.Metadata.Annotations["catalog.cattle.io/namespace"]) > 0 {
			namespace = app.Spec.Chart.Metadata.Annotations["catalog.cattle.io/namespace"]
		}
		chartAuto := splitBySep(app.Spec.Chart.Metadata.Annotations["catalog.cattle.io/auto-install"], "=")
		if len(chartAuto) != 2 {
			return diag.Errorf("bad format on chart annotation catalog.cattle.io/auto-install: %s", app.Spec.Chart.Metadata.Annotations["catalog.cattle.io/auto-install"])
		}
		name := chartAuto[0]
		app, err = getAppV2ByID(meta.(*Config), clusterID, namespace+"/"+name)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return nil
			}
			return diag.FromErr(err)
		}
		err = deleteAppV2(meta.(*Config), clusterID, app)
		if err != nil {
			return diag.Errorf("Error removing App V2 %s: %s", name, err)
		}
		stateConf = &retry.StateChangeConf{
			Pending:    []string{},
			Target:     []string{"removed"},
			Refresh:    appV2StateRefreshFunc(meta, clusterID, app.ID),
			Timeout:    d.Timeout(schema.TimeoutDelete),
			Delay:      1 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		_, waitErr = stateConf.WaitForStateContext(ctx)
		if waitErr != nil {
			return diag.Errorf("[ERROR] waiting for app (%s) to be deleted: %s", app.ID, waitErr)
		}

	}
	return nil
}

// appV2StateRefreshFunc returns a retry.StateRefreshFunc, used to watch a Rancher App.
func appV2StateRefreshFunc(meta interface{}, clusterID, appID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := getAppV2ByID(meta.(*Config), clusterID, appID)
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
	for obj, err := getAppV2OperationByID(meta.(*Config), clusterID, opID); ; obj, err = getAppV2OperationByID(meta.(*Config), clusterID, opID) {
		if err != nil {
			return err
		}
		if metadata, ok := obj["metadata"].(map[string]interface{}); ok && len(metadata) > 0 {
			if state, ok := metadata["state"].(map[string]interface{}); ok && len(state) > 0 {
				if transitioning, ok := state["transitioning"].(bool); ok && !transitioning {
					if opError, ok := state["error"].(bool); ok && opError {
						message, err := getAppV2OperationLogs(meta.(*Config), clusterID, obj)
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

// Rancher2 App V2 API CRUD functions
func createAppV2(c *Config, clusterID string, repo *ClusterRepo, chartIntall *types2.ChartInstallAction) (*types2.ChartActionOutput, error) {
	if c == nil {
		return nil, fmt.Errorf("Creating app V2: Provider config is nil")
	}
	if clusterID == "" {
		return nil, fmt.Errorf("Creating app V2: Cluster ID is nil")
	}
	if repo == nil || chartIntall == nil {
		return nil, fmt.Errorf("Creating app V2: Catalog V2 id and chartIntall should be provided")
	}

	client, err := c.CatalogV2Client(clusterID)
	if err != nil {
		return nil, err
	}
	resource := &types.Resource{
		ID:      repo.ID,
		Type:    repo.Type,
		Links:   repo.Links,
		Actions: repo.Actions,
	}
	resp := &types2.ChartActionOutput{}
	err = client.Action(catalogV2APIType, "install", resource, chartIntall, resp)
	if err != nil {
		return nil, fmt.Errorf("failed to install app v2: %v", err)
	}
	return resp, nil
}

func upgradeAppV2(c *Config, clusterID string, repo *ClusterRepo, chartUpgrade *types2.ChartUpgradeAction) (*types2.ChartActionOutput, error) {
	if c == nil {
		return nil, fmt.Errorf("Upgrading app V2: Provider config is nil")
	}
	if clusterID == "" {
		return nil, fmt.Errorf("Upgrading app V2: Cluster ID is nil")
	}
	if repo == nil || chartUpgrade == nil {
		return nil, fmt.Errorf("Upgrading app V2: Catalog V2 id and chartUpgrade should be provided")
	}

	client, err := c.CatalogV2Client(clusterID)
	if err != nil {
		return nil, err
	}
	resource := &types.Resource{
		ID:      repo.ID,
		Type:    repo.Type,
		Links:   repo.Links,
		Actions: repo.Actions,
	}
	resp := &types2.ChartActionOutput{}
	err = client.Action(catalogV2APIType, "upgrade", resource, chartUpgrade, resp)
	if err != nil {
		return nil, fmt.Errorf("failed to upgrade app v2: %v", err)
	}
	return resp, nil
}

func deleteAppV2(c *Config, clusterID string, app *AppV2) error {
	if c == nil {
		return fmt.Errorf("Deleting app V2: Provider config is nil")
	}
	if clusterID == "" {
		return fmt.Errorf("Deleting app V2: Cluster ID is nil")
	}
	if app == nil {
		return fmt.Errorf("App V2 id is nil")
	}

	client, err := c.CatalogV2Client(clusterID)
	if err != nil {
		return err
	}
	resource := &types.Resource{
		ID:      app.ID,
		Type:    app.Type,
		Links:   app.Links,
		Actions: app.Actions,
	}
	var resp interface{}
	return client.Action(appV2APIType, "uninstall", resource, map[string]interface{}{}, resp)
}

func infoAppV2(c *Config, clusterID, repoName, chartName, chartVersion string) (*ClusterRepo, *types2.ChartInfo, error) {
	if c == nil {
		return nil, nil, fmt.Errorf("Getting app V2 info: Provider config is nil")
	}
	if clusterID == "" {
		return nil, nil, fmt.Errorf("Getting app V2 info: Cluster ID is nil")
	}
	if repoName == "" || chartName == "" {
		return nil, nil, fmt.Errorf("Getting app V2 info: Catalog V2 id and chart name should be provided")
	}
	// Waiting for the Catalog V2 is Downloaded
	repo, err := waitCatalogV2Downloaded(c, clusterID, repoName)
	if err != nil {
		return nil, nil, err
	}
	resource := types.Resource{
		ID:      repo.ID,
		Type:    repo.Type,
		Links:   repo.Links,
		Actions: repo.Actions,
	}
	link := "info"
	if resource.Links == nil || len(resource.Links[link]) == 0 {
		return nil, nil, fmt.Errorf("failed to get chart info %s:%s from catalog v2 %s", chartName, chartVersion, repoName)
	}

	resource.Links[link] = resource.Links[link] + "&chartName=" + url.QueryEscape(chartName)
	if len(chartVersion) > 0 {
		resource.Links[link] = resource.Links[link] + "&version=" + url.QueryEscape(chartVersion)
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	client, err := c.CatalogV2Client(clusterID)
	if err != nil {
		return nil, nil, err
	}
	resp := &types2.ChartInfo{}
	for {
		err = client.GetLink(resource, link, resp)
		if err == nil {
			return repo, resp, nil
		}
		if !IsServerError(err) && !IsNotFound(err) {
			return nil, nil, fmt.Errorf("failed to get chart info %s:%s from catalog v2 %s: %v", chartName, chartVersion, repoName, err)
		}
		select {
		case <-time.After(rancher2RetriesWait * time.Second):
		case <-ctx.Done():
			return nil, nil, fmt.Errorf("Timeout getting chart info %s:%s from catalog v2 %s: %v", chartName, chartVersion, repoName, err)
		}
	}
}

func getAppV2ByID(c *Config, clusterID, id string) (*AppV2, error) {
	if c == nil {
		return nil, fmt.Errorf("Getting app V2: Provider config is nil")
	}
	if len(clusterID) == 0 || len(id) == 0 {
		return nil, fmt.Errorf("Getting app V2: Cluster ID and/or App V2 ID is nil")
	}
	resp := &AppV2{}
	err := c.getObjectV2ByID(clusterID, id, appV2APIType, resp)
	if err != nil {
		if !IsServerError(err) && !IsNotFound(err) && !IsForbidden(err) {
			return nil, fmt.Errorf("Getting App V2: %s", err)
		}
		return nil, err
	}
	return resp, nil
}

func getAppV2OperationByID(c *Config, clusterID, id string) (map[string]interface{}, error) {
	if c == nil {
		return nil, fmt.Errorf("Getting app V2 operation: Provider config is nil")
	}
	if len(clusterID) == 0 || len(id) == 0 {
		return nil, fmt.Errorf("Getting app V2 operation: Cluster ID and/or App V2 ID is nil")
	}
	resp := map[string]interface{}{}
	err := c.getObjectV2ByID(clusterID, id, appV2OperationAPIType, &resp)
	if err != nil {
		if !IsServerError(err) && !IsNotFound(err) && !IsForbidden(err) {
			return nil, fmt.Errorf("Getting App V2 logs: %s", err)
		}
		return nil, err
	}
	return resp, nil
}

func getAppV2OperationLogs(c *Config, clusterID string, op map[string]interface{}) (string, error) {
	if c == nil {
		return "", fmt.Errorf("Getting app V2 operation logs: Provider config is nil")
	}
	if len(clusterID) == 0 {
		return "", fmt.Errorf("Getting app V2 operation logs: Cluster ID is nil")
	}
	if op == nil {
		return "", fmt.Errorf("Getting app V2 operation logs: App V2 operation is nil")
	}
	if v, ok := op["id"].(string); !ok || v == "" {
		return "", fmt.Errorf("Getting app V2 operation logs: App V2 operation id is nil")
	}
	opLinks, ok := op["links"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("Getting app V2 operation logs: App V2 operation links are not available")
	}
	links := toMapString(opLinks)
	link := "logs"
	if links == nil || len(links[link]) == 0 {
		return "", fmt.Errorf("failed to get app v2 operation log %s", op["id"])
	}

	resp, err := DoGet(links[link], "", "", c.TokenKey, c.CACerts, c.Insecure)
	if err != nil {
		return "", fmt.Errorf("failed to get app v2 operation log %s: %s", op["id"], err)
	}

	return string(resp), nil
}
