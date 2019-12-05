package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

func resourceRancher2Notifier() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2NotifierCreate,
		Read:   resourceRancher2NotifierRead,
		Update: resourceRancher2NotifierUpdate,
		Delete: resourceRancher2NotifierDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2NotifierImport,
		},

		Schema: notifierFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2NotifierCreate(d *schema.ResourceData, meta interface{}) error {
	notifier, err := expandNotifier(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating Notifier %s", notifier.Name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	newNotifier, err := client.Notifier.Create(notifier)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    notifierStateRefreshFunc(client, newNotifier.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for notifier (%s) to be created: %s", newNotifier.ID, waitErr)
	}

	err = flattenNotifier(d, newNotifier)
	if err != nil {
		return err
	}

	return resourceRancher2NotifierRead(d, meta)
}

func resourceRancher2NotifierRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Notifier ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	notifier, err := client.Notifier.ByID(d.Id())
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Notifier ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = flattenNotifier(d, notifier)
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2NotifierUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Notifier ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	notifier, err := client.Notifier.ByID(d.Id())
	if err != nil {
		return err
	}

	update := map[string]interface{}{
		"name":         d.Get("name").(string),
		"description":  d.Get("description").(string),
		"clusterId":    d.Get("cluster_id").(string),
		"sendResolved": d.Get("send_resolved").(bool),
		"annotations":  toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":       toMapString(d.Get("labels").(map[string]interface{})),
	}

	if notifier.PagerdutyConfig != nil && d.HasChange("pagerduty_config") {
		update["pagerdutyConfig"] = expandNotifierPagerdutyConfig(d.Get("pagerduty_config").([]interface{}))
	}

	if notifier.SlackConfig != nil && d.HasChange("slack_config") {
		update["slackConfig"] = expandNotifierSlackConfig(d.Get("slack_config").([]interface{}))
	}

	if notifier.SMTPConfig != nil && d.HasChange("smtp_config") {
		update["smtpConfig"] = expandNotifierSMTPConfig(d.Get("smtp_config").([]interface{}))
	}

	if notifier.WebhookConfig != nil && d.HasChange("webhook_config") {
		update["webhookConfig"] = expandNotifierWebhookConfig(d.Get("webhook_config").([]interface{}))
	}

	if notifier.WechatConfig != nil && d.HasChange("wechat_config") {
		update["wechatConfig"] = expandNotifierWechatConfig(d.Get("wechat_config").([]interface{}))
	}

	newNotifier, err := client.Notifier.Update(notifier, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    notifierStateRefreshFunc(client, newNotifier.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for notifier (%s) to be updated: %s", newNotifier.ID, waitErr)
	}

	return resourceRancher2NotifierRead(d, meta)
}

func resourceRancher2NotifierDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Notifier ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	notifier, err := client.Notifier.ByID(id)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Notifier ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.Notifier.Delete(notifier)
	if err != nil {
		return fmt.Errorf("Error removing Notifier: %s", err)
	}

	log.Printf("[DEBUG] Waiting for notifier (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"removing"},
		Target:     []string{"removed"},
		Refresh:    notifierStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for notifier (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// notifierStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher Notifier.
func notifierStateRefreshFunc(client *managementClient.Client, notifierID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.Notifier.ByID(notifierID)
		if err != nil {
			if IsNotFound(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		if obj.Removed != "" {
			return obj, "removed", nil
		}

		return obj, obj.State, nil
	}
}
