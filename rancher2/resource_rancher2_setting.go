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

func resourceRancher2Setting() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRancher2SettingCreate,
		ReadContext:   resourceRancher2SettingRead,
		UpdateContext: resourceRancher2SettingUpdate,
		DeleteContext: resourceRancher2SettingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceRancher2SettingImport,
		},
		Schema: settingFields(),
	}
}

func resourceRancher2SettingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	// Checking if setting already exist, updating if already exist. setting id = setting name
	exist, err := client.Setting.ByID(d.Get("name").(string))
	if err == nil {
		d.SetId(exist.ID)
		return resourceRancher2SettingUpdate(ctx, d, meta)
	}
	if err != nil {
		if !IsNotFound(err) || IsForbidden(err) {
			return diag.FromErr(err)
		}
	}

	setting, err := expandSetting(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Creating Setting %s", setting.Name)

	newSetting, err := client.Setting.Create(setting)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(newSetting.ID)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    settingStateRefreshFunc(client, newSetting.ID),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf(
			"[ERROR] waiting for setting (%s) to be created: %s", newSetting.ID, waitErr)
	}

	return resourceRancher2SettingRead(ctx, d, meta)
}

func resourceRancher2SettingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)
	log.Printf("[INFO] Refreshing Rancher2 Setting ID %s", d.Id())

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	setting, err := client.Setting.ByID(name)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Setting ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	err = flattenSetting(d, setting)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceRancher2SettingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Updating Setting ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	setting, err := client.Setting.ByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	update := map[string]interface{}{
		"value":       d.Get("value").(string),
		"annotations": toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":      toMapString(d.Get("labels").(map[string]interface{})),
	}

	newSetting, err := client.Setting.Update(setting, update)
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    settingStateRefreshFunc(client, newSetting.ID),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForStateContext(ctx)
	if waitErr != nil {
		return diag.Errorf(
			"[ERROR] waiting for setting (%s) to be updated: %s", newSetting.ID, waitErr)
	}

	return resourceRancher2SettingRead(ctx, d, meta)
}

func resourceRancher2SettingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] Deleting Setting ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return diag.FromErr(err)
	}

	setting, err := client.Setting.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Setting ID %s not found.", id)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	// Deleting setting if it was cretaed by user
	if setting.CreatorID != "" {
		err = client.Setting.Delete(setting)
		if err != nil {
			return diag.Errorf("Error removing setting: %s", err)
		}

		log.Printf("[DEBUG] Waiting for setting (%s) to be removed", id)

		stateConf := &retry.StateChangeConf{
			Pending:    []string{"active"},
			Target:     []string{"removed"},
			Refresh:    settingStateRefreshFunc(client, id),
			Timeout:    10 * time.Minute,
			Delay:      1 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, waitErr := stateConf.WaitForStateContext(ctx)
		if waitErr != nil {
			return diag.Errorf(
				"[ERROR] waiting for setting (%s) to be removed: %s", id, waitErr)
		}
		// Reseting setting to value = "" if it was cretaed by system
	} else {
		err = d.Set("value", "")
		if err != nil {
			return diag.FromErr(err)
		}

		diagnostics := resourceRancher2SettingUpdate(ctx, d, meta)
		if diagnostics.HasError() {
			return diagnostics
		}
	}

	d.SetId("")
	return nil
}

// settingStateRefreshFunc returns a retry.StateRefreshFunc, used to watch a Rancher Project.
func settingStateRefreshFunc(client *managementClient.Client, settingID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.Setting.ByID(settingID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}
		return obj, "active", nil
	}
}
