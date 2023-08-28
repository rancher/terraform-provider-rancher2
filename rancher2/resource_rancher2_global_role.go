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

func resourceRancher2GlobalRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2GlobalRoleCreate,
		ReadContext:   resourceRancher2GlobalRoleRead,
		UpdateContext: resourceRancher2GlobalRoleUpdate,
		DeleteContext: resourceRancher2GlobalRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRancher2GlobalRoleImport,
		},

		Schema: globalRoleFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2GlobalRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		globalRole := expandGlobalRole(d)

		log.Printf("[INFO] Creating global role")

		newGlobalRole, err := client.GlobalRole.Create(globalRole)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		d.SetId(newGlobalRole.ID)

		diagnostics := resourceRancher2GlobalRoleRead(ctx, d, meta)
		if diagnostics.HasError() {
			return retry.NonRetryableError(errors.New(diagnostics[0].Summary))
		}

		return nil
	}))
}

func resourceRancher2GlobalRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Refreshing global role ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		globalRole, err := client.GlobalRole.ByID(d.Id())
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] global role ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return retry.NonRetryableError(err)
		}

		if err = flattenGlobalRole(d, globalRole); err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	}))
}

func resourceRancher2GlobalRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating global role ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		globalRole, err := client.GlobalRole.ByID(d.Id())
		if err != nil {
			return retry.NonRetryableError(err)
		}

		update := map[string]interface{}{
			"description":    d.Get("description").(string),
			"name":           d.Get("name").(string),
			"newUserDefault": d.Get("new_user_default").(bool),
			"rules":          expandPolicyRules(d.Get("rules").([]interface{})),
			"annotations":    toMapString(d.Get("annotations").(map[string]interface{})),
			"labels":         toMapString(d.Get("labels").(map[string]interface{})),
		}

		if _, err = client.GlobalRole.Update(globalRole, update); err != nil {
			return retry.NonRetryableError(err)
		}

		if diagnostics := resourceRancher2GlobalRoleRead(ctx, d, meta); diagnostics.HasError() {
			return retry.NonRetryableError(errors.New(diagnostics[0].Summary))
		}

		return nil
	}))
}

func resourceRancher2GlobalRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting global role ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		globalRole, err := client.GlobalRole.ByID(id)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Global role ID %s not found.", id)
				d.SetId("")
				return nil
			}
			return retry.NonRetryableError(err)
		}

		if !globalRole.Builtin {
			if err = client.GlobalRole.Delete(globalRole); err != nil {
				return retry.NonRetryableError(fmt.Errorf("[ERROR] Error removing global role: %s", err))
			}
		}

		d.SetId("")
		return nil
	}))
}
