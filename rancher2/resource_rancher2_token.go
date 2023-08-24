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

func resourceRancher2Token() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2TokenCreate,
		ReadContext:   resourceRancher2TokenRead,
		UpdateContext: resourceRancher2TokenUpdate,
		DeleteContext: resourceRancher2TokenDelete,

		Schema: tokenFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceRancher2TokenCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Creating Token")
	patch, err := meta.(*Config).IsRancherVersionGreaterThanOrEqualAndLessThan(rancher2TokeTTLMinutesVersion, rancher2TokeTTLMilisVersion)
	if err != nil {
		return diag.FromErr(err)
	}
	token, err := expandToken(d, patch)
	if err != nil {
		return diag.FromErr(err)
	}

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *retry.RetryError {
		newToken, err := client.Token.Create(token)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		err = flattenToken(d, newToken, patch)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		diag2 := resourceRancher2TokenRead(ctx, d, meta)
		if diag2.HasError() {
			return retry.NonRetryableError(errors.New(diag2[0].Summary))
		}

		return nil
	})
	return diag.FromErr(err)
}

func resourceRancher2TokenRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Refreshing Token ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		token, err := client.Token.ByID(d.Id())
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Token ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return retry.NonRetryableError(err)
		}

		renew := d.Get("renew").(bool)
		if (!*token.Enabled || token.Expired) && renew {
			d.Set("renew", false)
		}

		patch, err := meta.(*Config).IsRancherVersionGreaterThanOrEqualAndLessThan(rancher2TokeTTLMinutesVersion, rancher2TokeTTLMilisVersion)
		if err != nil {
			return retry.NonRetryableError(err)
		}
		err = flattenToken(d, token, patch)
		if err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	})
	return diag.FromErr(err)
}

func resourceRancher2TokenUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceRancher2TokenRead(ctx, d, meta)
}

func resourceRancher2TokenDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting Token ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *retry.RetryError {
		token, err := client.Token.ByID(id)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Token ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return retry.NonRetryableError(err)
		}

		err = client.Token.Delete(token)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("[ERROR] Error removing Token: %s", err))
		}

		d.SetId("")
		return nil
	})
	return diag.FromErr(err)
}

func isTokenValid(c *Config, id string) (bool, error) {
	if len(id) == 0 {
		return false, nil
	}

	client, err := c.ManagementClient()
	if err != nil {
		return false, err
	}

	token, err := client.Token.ByID(id)
	if err != nil {
		if !IsNotFound(err) && !IsForbidden(err) {
			return false, err
		}
		return false, nil
	}
	// Token is valid if it's enabled and not expired
	return (token.Enabled != nil && *token.Enabled && !token.Expired), nil
}
