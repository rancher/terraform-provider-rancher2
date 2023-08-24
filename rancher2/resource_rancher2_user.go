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

func resourceRancher2User() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2UserCreate,
		ReadContext:   resourceRancher2UserRead,
		UpdateContext: resourceRancher2UserUpdate,
		DeleteContext: resourceRancher2UserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRancher2UserImport,
		},

		Schema: userFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceRancher2UserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	user := expandUser(d)

	log.Printf("[INFO] Creating User %s", user.Username)

	newUser, err := client.User.Create(user)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newUser.ID)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    userStateRefreshFunc(client, newUser.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf(
			"[ERROR] waiting for user (%s) to be created: %s", newUser.ID, waitErr)
	}

	return resourceRancher2UserRead(ctx, d, meta)
}

func resourceRancher2UserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Refreshing User ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}
	err = retry.RetryContext(ctx, d.Timeout(schema.TimeoutRead), func() *retry.RetryError {
		user, err := client.User.ByID(d.Id())
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] User ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return retry.NonRetryableError(err)

		}

		if err = flattenUser(d, user); err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	})
	return diag.FromErr(err)
}

func resourceRancher2UserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating User ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	user, err := client.User.ByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Update user password if needed
	_, user, err = meta.(*Config).SetUserPassword(user, d.Get("password").(string))
	if err != nil {
		return diag.Errorf("[ERROR] Updating Admin password: %s", err)
	}

	update := map[string]interface{}{
		"name":        d.Get("name").(string),
		"enabled":     d.Get("enabled").(bool),
		"annotations": toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":      toMapString(d.Get("labels").(map[string]interface{})),
	}

	newUser, err := client.User.Update(user, update)
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    userStateRefreshFunc(client, newUser.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf(
			"[ERROR] waiting for user (%s) to be updated: %s", newUser.ID, waitErr)
	}

	return resourceRancher2UserRead(ctx, d, meta)
}

func resourceRancher2UserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting User ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	user, err := client.User.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] User ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = client.User.Delete(user)
	if err != nil {
		return diag.Errorf("Error removing User: %s", err)
	}

	log.Printf("[DEBUG] Waiting for user (%s) to be removed", id)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"removing"},
		Target:     []string{"removed"},
		Refresh:    userStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf(
			"[ERROR] waiting for user (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// userStateRefreshFunc returns a retry.StateRefreshFunc, used to watch a Rancher User.
func userStateRefreshFunc(client *managementClient.Client, userID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.User.ByID(userID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}
		return obj, obj.State, nil
	}
}
