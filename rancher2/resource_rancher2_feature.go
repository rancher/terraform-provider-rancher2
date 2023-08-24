package rancher2

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRancher2Feature() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2FeatureCreate,
		ReadContext:   resourceRancher2FeatureRead,
		UpdateContext: resourceRancher2FeatureUpdate,
		DeleteContext: resourceRancher2FeatureDelete,
		Schema:        featureFields(),
	}
}

func resourceRancher2FeatureCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// New features can't be created from the Rancher API just update existing
	d.SetId(d.Get("name").(string))

	return resourceRancher2FeatureUpdate(ctx, d, meta)
}

func resourceRancher2FeatureRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	log.Printf("[INFO] Refreshing Rancher2 Feature ID %s", d.Id())

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	feature, err := client.Feature.ByID(name)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Feature ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("[ERROR] refreshing feature %s: %v", d.Id(), err)
	}

	err = flattenFeature(d, feature)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceRancher2FeatureUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating Feature ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	feature, err := client.Feature.ByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	featValue := d.Get("value").(bool)
	update := map[string]interface{}{
		"value":       &featValue,
		"annotations": toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":      toMapString(d.Get("labels").(map[string]interface{})),
	}

	newFeature, err := client.Feature.Update(feature, update)
	if err != nil {
		// Bad gateway or service unavailable error may be fine if Rancher is restarted
		if !IsBadGatewayError(err) && !IsServiceUnavailableError(err) {
			return diag.Errorf("[ERROR] updating feature %s: %v", d.Id(), err)
		}
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"active", "rebooting"},
		Target:     []string{"active"},
		Refresh:    featureStateRefreshFunc(meta, newFeature.ID),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf(
			"[ERROR] waiting for feature (%s) to be updated: %s", newFeature.ID, waitErr)
	}

	return resourceRancher2FeatureRead(ctx, d, meta)
}

func resourceRancher2FeatureDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting Feature ID %s", d.Id())
	// Not removing feature from Rancher just from tfstate
	d.SetId("")
	return nil
}

// featureStateRefreshFunc returns a retry.StateRefreshFunc, used to watch a Rancher Project.
func featureStateRefreshFunc(meta interface{}, featureID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		err := meta.(*Config).isRancherReady()
		if err != nil {
			return nil, "rebooting", nil
		}
		client, err := meta.(*Config).ManagementClient()
		if err != nil {
			return nil, "", err
		}
		obj, err := client.Feature.ByID(featureID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}
		return obj, "active", nil
	}
}
