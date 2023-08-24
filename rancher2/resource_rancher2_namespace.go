package rancher2

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	clusterClient "github.com/rancher/rancher/pkg/client/generated/cluster/v3"
)

func resourceRancher2Namespace() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2NamespaceCreate,
		ReadContext:   resourceRancher2NamespaceRead,
		UpdateContext: resourceRancher2NamespaceUpdate,
		DeleteContext: resourceRancher2NamespaceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRancher2NamespaceImport,
		},

		Schema: namespaceFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2NamespaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clusterID, err := clusterIDFromProjectID(d.Get("project_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	active, _, err := meta.(*Config).isClusterActive(clusterID)
	if err != nil {
		return diag.FromErr(err)
	}
	if !active {
		if v, ok := d.Get("wait_for_cluster").(bool); ok && !v {
			return diag.Errorf("[ERROR] Creating Namespace: Cluster ID %s is not active", clusterID)
		}

		_, err := meta.(*Config).WaitForClusterState(clusterID, clusterActiveCondition, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.Errorf("[ERROR] waiting for cluster ID (%s) to be active: %s", clusterID, err)
		}
	}

	client, err := meta.(*Config).ClusterClient(clusterID)
	if err != nil {
		return diag.FromErr(err)
	}

	ns := expandNamespace(d)

	log.Printf("[INFO] Creating Namespace %s on Cluster ID %s", ns.Name, clusterID)

	newNs, err := client.Namespace.Create(ns)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newNs.ID)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"activating", "forbidden"},
		Target:     []string{"active"},
		Refresh:    namespaceStateRefreshFunc(client, newNs.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf(
			"[ERROR] waiting for namespace (%s) to be created: %s", newNs.ID, waitErr)
	}

	return resourceRancher2NamespaceRead(ctx, d, meta)
}

func resourceRancher2NamespaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clusterID, _ := splitProjectID(d.Get("project_id").(string))

	log.Printf("[INFO] Refreshing Namespace ID %s", d.Id())

	return diag.FromErr(retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		_, _, err := meta.(*Config).isClusterActive(clusterID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Cluster ID %s not found.", clusterID)
				d.SetId("")
				return nil
			}
			return retry.NonRetryableError(err)
		}

		client, err := meta.(*Config).ClusterClient(clusterID)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		ns, err := client.Namespace.ByID(d.Id())
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Namespace ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return retry.NonRetryableError(err)
		}

		if err = flattenNamespace(d, ns); err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	}))
}

func resourceRancher2NamespaceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clusterID, projectID := splitProjectID(d.Get("project_id").(string))

	log.Printf("[INFO] Updating Namespace ID %s", d.Id())

	client, err := meta.(*Config).ClusterClient(clusterID)
	if err != nil {
		return diag.FromErr(err)
	}

	ns, err := client.Namespace.ByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	readClusterID, readProjectID := splitProjectID(ns.ProjectID)

	if projectID != readProjectID && (clusterID == readClusterID || readClusterID == "") {
		log.Printf("[INFO] Moving Namespace ID %s to project %s", d.Id(), projectID)
		nsMove := &clusterClient.NamespaceMove{
			ProjectID: projectID,
		}

		err = client.Namespace.ActionMove(ns, nsMove)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	resourceQuota := expandNamespaceResourceQuota(d.Get("resource_quota").([]interface{}))

	update := map[string]interface{}{
		"description":                   d.Get("description").(string),
		"containerDefaultResourceLimit": expandNamespaceContainerResourceLimit(d.Get("container_resource_limit").([]interface{})),
		"resourceQuota":                 resourceQuota,
		"annotations":                   toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":                        toMapString(d.Get("labels").(map[string]interface{})),
	}

	newNs, err := client.Namespace.Update(ns, update)
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    namespaceStateRefreshFunc(client, newNs.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf(
			"[ERROR] waiting for namespace (%s) to be updated: %s", newNs.ID, waitErr)
	}

	err = flattenNamespace(d, newNs)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRancher2NamespaceRead(ctx, d, meta)
}

func resourceRancher2NamespaceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	clusterID, _ := splitProjectID(d.Get("project_id").(string))

	log.Printf("[INFO] Deleting Namespace ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ClusterClient(clusterID)
	if err != nil {
		return diag.FromErr(err)
	}

	ns, err := client.Namespace.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) || IsServiceUnavailableError(err) {
			log.Printf("[INFO] Namespace ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = client.Namespace.Delete(ns)
	if err != nil {
		return diag.Errorf("Error removing Namespace: %s", err)
	}

	log.Printf("[DEBUG] Waiting for namespace (%s) to be removed", id)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"removing"},
		Target:     []string{"removed", "forbidden"},
		Refresh:    namespaceStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf(
			"[ERROR] waiting for namespace (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// namespaceStateRefreshFunc returns a retry.StateRefreshFunc, used to watch a Rancher Namespace.
func namespaceStateRefreshFunc(client *clusterClient.Client, nsID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.Namespace.ByID(nsID)
		if err != nil {
			if IsNotFound(err) {
				return obj, "removed", nil
			}
			if IsForbidden(err) {
				return obj, "forbidden", nil
			}
			return nil, "", err
		}

		return obj, obj.State, nil
	}
}
