package rancher2

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRancher2Registry() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2RegistryCreate,
		ReadContext:   resourceRancher2RegistryRead,
		UpdateContext: resourceRancher2RegistryUpdate,
		DeleteContext: resourceRancher2RegistryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRancher2RegistryImport,
		},

		Schema: registryFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2RegistryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	name := d.Get("name").(string)

	return diag.FromErr(retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		err := meta.(*Config).ProjectExist(projectID)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		registry := expandRegistry(d)

		log.Printf("[INFO] Creating Registry %s on Project ID %s", name, projectID)

		newRegistry, err := meta.(*Config).CreateRegistry(registry)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		err = flattenRegistry(d, newRegistry)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		diagnostics := resourceRancher2RegistryRead(ctx, d, meta)
		if diagnostics.HasError() {
			return retry.NonRetryableError(errors.New(diagnostics[0].Summary))
		}

		return nil
	}))
}

func resourceRancher2RegistryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	id := d.Id()
	namespaceID := d.Get("namespace_id").(string)

	log.Printf("[INFO] Refreshing Registry ID %s", id)

	return diag.FromErr(retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		registry, err := meta.(*Config).GetRegistry(id, projectID, namespaceID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Registry ID %s not found.", id)
				d.SetId("")
				return nil
			}
			return retry.NonRetryableError(err)
		}

		if err = flattenRegistry(d, registry); err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	}))
}

func resourceRancher2RegistryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	id := d.Id()
	namespaceID := d.Get("namespace_id").(string)

	log.Printf("[INFO] Updating Registry ID %s", id)

	return diag.FromErr(retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		registry, err := meta.(*Config).GetRegistry(id, projectID, namespaceID)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		update := map[string]interface{}{
			"description": d.Get("description").(string),
			"registries":  expandRegistryCredential(d.Get("registries").([]interface{})),
			"annotations": toMapString(d.Get("annotations").(map[string]interface{})),
			"labels":      toMapString(d.Get("labels").(map[string]interface{})),
		}

		newRegistry, err := meta.(*Config).UpdateRegistry(registry, update)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		err = flattenRegistry(d, newRegistry)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		diagnostics := resourceRancher2RegistryRead(ctx, d, meta)
		if diagnostics.HasError() {
			return retry.NonRetryableError(errors.New(diagnostics[0].Summary))
		}

		return nil
	}))
}

func resourceRancher2RegistryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	id := d.Id()
	namespaceID := d.Get("namespace_id").(string)

	log.Printf("[INFO] Deleting Registry ID %s", id)

	return diag.FromErr(retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		registry, err := meta.(*Config).GetRegistry(id, projectID, namespaceID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Registry ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return retry.NonRetryableError(err)
		}

		err = meta.(*Config).DeleteRegistry(registry)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("[ERROR] Error removing Registry: %s", err))
		}

		d.SetId("")
		return nil
	}))
}
