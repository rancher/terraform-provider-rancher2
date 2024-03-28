package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2User() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2UserCreate,
		Read:   resourceRancher2UserRead,
		Update: resourceRancher2UserUpdate,
		Delete: resourceRancher2UserDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2UserImport,
		},

		Schema: userFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceRancher2UserCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	user := expandUser(d)

	log.Printf("[INFO] Creating User %s", user.Username)

	newUser, err := client.User.Create(user)
	if err != nil {
		return err
	}

	d.SetId(newUser.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    userStateRefreshFunc(client, newUser.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for user (%s) to be created: %s", newUser.ID, waitErr)
	}

	return resourceRancher2UserRead(d, meta)
}

func resourceRancher2UserRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing User ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		user, err := client.User.ByID(d.Id())
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] User ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(err)
		}

		if err = flattenUser(d, user); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourceRancher2UserUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating User ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	user, err := client.User.ByID(d.Id())
	if err != nil {
		return err
	}

	// Update user password if needed
	_, user, err = meta.(*Config).SetUserPassword(user, d.Get("password").(string))
	if err != nil {
		return fmt.Errorf("[ERROR] Updating Admin password: %s", err)
	}

	update := map[string]interface{}{
		"name":        d.Get("name").(string),
		"enabled":     d.Get("enabled").(bool),
		"annotations": toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":      toMapString(d.Get("labels").(map[string]interface{})),
	}

	newUser, err := client.User.Update(user, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    userStateRefreshFunc(client, newUser.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for user (%s) to be updated: %s", newUser.ID, waitErr)
	}

	return resourceRancher2UserRead(d, meta)
}

func resourceRancher2UserDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting User ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	user, err := client.User.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] User ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.User.Delete(user)
	if err != nil {
		return fmt.Errorf("Error removing User: %s", err)
	}

	log.Printf("[DEBUG] Waiting for user (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"removing"},
		Target:     []string{"removed"},
		Refresh:    userStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for user (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// userStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher User.
func userStateRefreshFunc(client *managementClient.Client, userID string) resource.StateRefreshFunc {
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
