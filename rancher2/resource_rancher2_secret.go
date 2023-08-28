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

func resourceRancher2Secret() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2SecretCreate,
		ReadContext:   resourceRancher2SecretRead,
		UpdateContext: resourceRancher2SecretUpdate,
		DeleteContext: resourceRancher2SecretDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRancher2SecretImport,
		},

		Schema: secretFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2SecretCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	name := d.Get("name").(string)

	err := retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		err := meta.(*Config).ProjectExist(projectID)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		secret := expandSecret(d)

		log.Printf("[INFO] Creating Secret %s on Project ID %s", name, projectID)

		newSecret, err := meta.(*Config).CreateSecret(secret)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		err = flattenSecret(d, newSecret)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		diagnostics := resourceRancher2SecretRead(ctx, d, meta)
		if diagnostics.HasError() {
			return retry.NonRetryableError(errors.New(diagnostics[0].Summary))
		}

		return nil
	})
	return diag.FromErr(err)
}

func resourceRancher2SecretRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	id := d.Id()
	namespaceID := d.Get("namespace_id").(string)

	log.Printf("[INFO] Refreshing Secret ID %s", id)

	secret, err := meta.(*Config).GetSecret(id, projectID, namespaceID)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Secret ID %s not found.", id)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	return diag.FromErr(flattenSecret(d, secret))
}

func resourceRancher2SecretUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	id := d.Id()
	namespaceID := d.Get("namespace_id").(string)

	log.Printf("[INFO] Updating Secret ID %s", id)

	return diag.FromErr(retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		secret, err := meta.(*Config).GetSecret(id, projectID, namespaceID)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		update := map[string]interface{}{
			"description": d.Get("description").(string),
			"data":        toMapString(d.Get("data").(map[string]interface{})),
			"annotations": toMapString(d.Get("annotations").(map[string]interface{})),
			"labels":      toMapString(d.Get("labels").(map[string]interface{})),
		}

		newSecret, err := meta.(*Config).UpdateSecret(secret, update)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		err = flattenSecret(d, newSecret)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		diagnostics := resourceRancher2SecretRead(ctx, d, meta)
		if diagnostics.HasError() {
			return retry.NonRetryableError(errors.New(diagnostics[0].Summary))
		}

		return nil
	}))
}

func resourceRancher2SecretDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	id := d.Id()
	namespaceID := d.Get("namespace_id").(string)

	log.Printf("[INFO] Deleting Secret ID %s", id)

	return diag.FromErr(retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		secret, err := meta.(*Config).GetSecret(id, projectID, namespaceID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Secret ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return retry.NonRetryableError(err)
		}

		err = meta.(*Config).DeleteSecret(secret)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("[ERROR] Error removing Secret: %s", err))
		}

		d.SetId("")
		return nil
	}))
}
