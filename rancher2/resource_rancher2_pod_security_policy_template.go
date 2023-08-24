package rancher2

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2PodSecurityPolicyTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2PodSecurityPolicyTemplateCreate,
		ReadContext:   resourceRancher2PodSecurityPolicyTemplateRead,
		UpdateContext: resourceRancher2PodSecurityPolicyTemplateUpdate,
		DeleteContext: resourceRancher2PodSecurityPolicyTemplateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRancher2PodSecurityPolicyTemplateImport,
		},

		Schema: podSecurityPolicyTemplateFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2PodSecurityPolicyTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	podSecurityPolicyTemplate := expandPodSecurityPolicyTemplate(d)

	log.Printf("[INFO] Creating PodSecurityPolicyTemplate %s", podSecurityPolicyTemplate.Name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	newPodSecurityPolicyTemplate, err := client.PodSecurityPolicyTemplate.Create(podSecurityPolicyTemplate)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newPodSecurityPolicyTemplate.ID)

	return resourceRancher2PodSecurityPolicyTemplateRead(ctx, d, meta)
}

func resourceRancher2PodSecurityPolicyTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Refreshing PodSecurityPolicyTemplate with ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		pspt, err := client.PodSecurityPolicyTemplate.ByID(d.Id())
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] PodSecurityPolicyTemplate with ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return retry.NonRetryableError(err)
		}

		if err = flattenPodSecurityPolicyTemplate(d, pspt); err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	}))
}

func resourceRancher2PodSecurityPolicyTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating PodSecurityPolicyTemplate with ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	pspt, err := client.PodSecurityPolicyTemplate.ByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	update := expandPodSecurityPolicyTemplate(d)

	_, err = client.PodSecurityPolicyTemplate.Update(pspt, update)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRancher2PodSecurityPolicyTemplateRead(ctx, d, meta)
}

func resourceRancher2PodSecurityPolicyTemplateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := d.Id()
	log.Printf("[INFO] Deleting PodSecurityPolicyTemplate with ID %s", id)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	pspt, err := client.PodSecurityPolicyTemplate.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] PodSecurityPolicyTemplate with ID %s not found.", id)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = client.PodSecurityPolicyTemplate.Delete(pspt)
	if err != nil {
		return diag.Errorf("Error removing PodSecurityPolicyTemplate: %s", err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"removed"},
		Refresh:    podSecurityPolicyTemplateStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf(
			"[ERROR] waiting for PodSecurityPolicyTemplate (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// podSecurityPolicyTemplateStateRefreshFunc returns a retry.StateRefreshFunc, used to watch a Rancher PodSecurityPolicyTemplate
func podSecurityPolicyTemplateStateRefreshFunc(client *managementClient.Client, pspID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.PodSecurityPolicyTemplate.ByID(pspID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		return obj, "active", nil
	}
}
