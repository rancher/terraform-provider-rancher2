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

func resourceRancher2RoleTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2RoleTemplateCreate,
		ReadContext:   resourceRancher2RoleTemplateRead,
		UpdateContext: resourceRancher2RoleTemplateUpdate,
		DeleteContext: resourceRancher2RoleTemplateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRancher2RoleTemplateImport,
		},

		Schema: roleTemplateFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2RoleTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		roleTemplate := expandRoleTemplate(d)
		if roleTemplate == nil {
			log.Printf("[INFO] Expanded role template was empty")
			return nil
		}

		log.Printf("[INFO] Creating role template")

		newRoleTemplate, err := client.RoleTemplate.Create(roleTemplate)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		d.SetId(newRoleTemplate.ID)

		diagnostics := resourceRancher2RoleTemplateRead(ctx, d, meta)
		if diagnostics.HasError() {
			return retry.NonRetryableError(errors.New(diagnostics[0].Summary))
		}

		return nil
	}))
}

func resourceRancher2RoleTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Refreshing role template ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		roleTemplate, err := client.RoleTemplate.ByID(d.Id())
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] role template ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return retry.NonRetryableError(err)
		}

		if err = flattenRoleTemplate(d, roleTemplate); err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	}))
}

func resourceRancher2RoleTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating role template ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(retry.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *retry.RetryError {
		roleTemplate, err := client.RoleTemplate.ByID(d.Id())
		if err != nil {
			return retry.NonRetryableError(err)
		}

		update := map[string]interface{}{
			"administrative":  d.Get("administrative").(bool),
			"context":         d.Get("context").(string),
			"description":     d.Get("description").(string),
			"external":        d.Get("external").(bool),
			"hidden":          d.Get("hidden").(bool),
			"locked":          d.Get("locked").(bool),
			"name":            d.Get("name").(string),
			"roleTemplateIds": toArrayString(d.Get("role_template_ids").([]interface{})),
			"rules":           expandPolicyRules(d.Get("rules").([]interface{})),
			"annotations":     toMapString(d.Get("annotations").(map[string]interface{})),
			"labels":          toMapString(d.Get("labels").(map[string]interface{})),
		}

		switch update["context"] {
		case roleTemplateContextCluster:
			update["clusterCreatorDefault"] = d.Get("default_role").(bool)
			update["projectCreatorDefault"] = false
		case roleTemplateContextProject:
			update["clusterCreatorDefault"] = false
			update["projectCreatorDefault"] = d.Get("default_role").(bool)
		}

		_, err = client.RoleTemplate.Update(roleTemplate, update)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		diagnostics := resourceRancher2RoleTemplateRead(ctx, d, meta)
		if diagnostics.HasError() {
			return retry.NonRetryableError(errors.New(diagnostics[0].Summary))
		}

		return nil
	}))
}

func resourceRancher2RoleTemplateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting role template ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		roleTemplate, err := client.RoleTemplate.ByID(id)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Role template ID %s not found.", id)
				d.SetId("")
				return nil
			}
			return retry.NonRetryableError(err)
		}

		if !roleTemplate.Builtin {
			err = client.RoleTemplate.Delete(roleTemplate)
			if err != nil {
				return retry.NonRetryableError(fmt.Errorf("[ERROR] Error removing role template: %s", err))
			}
		}

		d.SetId("")
		return nil
	}))
}
