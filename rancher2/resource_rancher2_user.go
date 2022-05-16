package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	norman "github.com/rancher/norman/types"
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
		CustomizeDiff: userTokenComputedIf(userTokenFieldsList),
		Schema:        userFields(),
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

	if v, ok := d.Get("token_config").([]interface{}); ok && len(v) > 0 {
		patch, err := meta.(*Config).IsRancherVersionGreaterThanOrEqualAndLessThan(rancher2TokeTTLMinutesVersion, rancher2TokeTTLMilisVersion)
		if err != nil {
			return err
		}

		if token := expandUserToken(d, patch); token != nil {
			err = resourceRancher2UserCreateLoginBinding(d, client)
			if err != nil {
				return err
			}

			userClient, err := getClientForUser(d, meta)
			if err != nil {
				return err
			}
			defer doUserLogout(d, userClient)

			err = resourceRancher2UserRecreateToken(d, meta, userClient, token)
			if err != nil {
				return err
			}
		}
	}

	return resourceRancher2UserRead(d, meta)
}

func resourceRancher2UserCreateLoginBinding(d *schema.ResourceData, client *managementClient.Client) error {
	log.Printf("[INFO] Creating login Role Binding for User %s", d.Id())
	newLoginRoleBinding, err := client.GlobalRoleBinding.Create(&managementClient.GlobalRoleBinding{
		GlobalRoleID: "user-base",
		UserID:       d.Id(),
	})
	if err != nil {
		return err
	}

	d.Set("login_role_binding_id", newLoginRoleBinding.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    globalRoleBindingStateRefreshFunc(client, newLoginRoleBinding.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for login role binding (%s) to be created: %s", newLoginRoleBinding.ID, waitErr)
	}

	return nil
}

func resourceRancher2UserRecreateToken(d *schema.ResourceData, meta interface{}, client *managementClient.Client, token *managementClient.Token) error {
	// Delete current token, if exists
	if v, ok := d.Get("token_id").(string); ok && len(v) > 0 {
		log.Printf("[DEBUG] Removing token %s from user %s", v, d.Id())
		currentToken, err := client.Token.ByID(v)
		if err != nil {
			if IsNotFound(err) {
				clearUserTokenFields(d)
			} else {
				return err
			}
		}

		err = client.Token.Delete(currentToken)
		if err != nil {
			return nil
		}

		clearUserTokenFields(d)
	}

	newToken, err := client.Token.Create(token)
	if err != nil {
		return fmt.Errorf("[ERROR] Creating User token: %s", err)
	}

	d.Set("token_id", newToken.ID)

	patch, err := meta.(*Config).IsRancherVersionGreaterThanOrEqualAndLessThan(rancher2TokeTTLMinutesVersion, rancher2TokeTTLMilisVersion)
	if err != nil {
		return err
	}

	return flattenUserToken(d, newToken, patch)
}

func resourceRancher2UserRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing User ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	user, err := client.User.ByID(d.Id())
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] User ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	if v, ok := d.Get("token_id").(string); ok && len(v) > 0 {
		userClient, err := getClientForUser(d, meta)
		if err != nil {
			return err
		}
		defer doUserLogout(d, userClient)

		err = resourceRancher2UserTokenRead(d, meta, userClient)
		if err != nil {
			return err
		}
	}

	return flattenUser(d, user)
}

func resourceRancher2UserTokenRead(d *schema.ResourceData, meta interface{}, client *managementClient.Client) error {
	tokenID := d.Get("token_id").(string)
	log.Printf("[INFO] Refreshing Token ID %s for User %s", tokenID, d.Id())

	token, err := client.Token.ByID(tokenID)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Token ID %s not found.", tokenID)
			clearUserTokenFields(d)
			return nil
		}
		return err
	}

	patch, err := meta.(*Config).IsRancherVersionGreaterThanOrEqualAndLessThan(rancher2TokeTTLMinutesVersion, rancher2TokeTTLMilisVersion)
	if err != nil {
		return err
	}

	return flattenUserToken(d, token, patch)
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

	if v, ok := d.Get("token_config").([]interface{}); ok && len(v) > 0 {
		err := recreateUserTokenIfNecessary(d, meta)
		if err != nil {
			return err
		}
	} else if v, ok := d.Get("token_id").(string); ok && len(v) > 0 {
		err := deleteUserToken(d, meta)
		if err != nil {
			return err
		}
	}

	return resourceRancher2UserRead(d, meta)
}

func recreateUserTokenIfNecessary(d *schema.ResourceData, meta interface{}) error {
	hasChange := d.HasChange("token_config")
	renewTokenIfNecessary := d.Get("token_config.0.renew").(bool)

	if hasChange || renewTokenIfNecessary {
		client, err := getClientForUser(d, meta)
		if err != nil {
			return err
		}
		defer doUserLogout(d, client)

		expiredToken := false
		if !hasChange && renewTokenIfNecessary {
			expiredToken, err = isUserTokenExpired(client, d.Get("token_id").(string))
			if err != nil {
				return err
			}
		}

		if hasChange || expiredToken {
			log.Printf("[INFO] Recreating Token for User ID %s (change=%v, renew=%v, expired=%v)", d.Id(), hasChange, renewTokenIfNecessary, expiredToken)

			patch, err := meta.(*Config).IsRancherVersionGreaterThanOrEqualAndLessThan(rancher2TokeTTLMinutesVersion, rancher2TokeTTLMilisVersion)
			if err != nil {
				return err
			}

			token := expandUserToken(d, patch)

			err = resourceRancher2UserRecreateToken(d, meta, client, token)
			if err != nil {
				return err
			}

			return resourceRancher2UserTokenRead(d, meta, client)
		}
	}

	return nil
}

func isUserTokenExpired(client *managementClient.Client, id string) (bool, error) {
	if len(id) == 0 {
		return true, nil
	}

	token, err := client.Token.ByID(id)
	if err != nil {
		if IsNotFound(err) {
			return true, nil
		}
		return false, err
	}

	return token.Expired, nil
}

func deleteUserToken(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Token ID %s", d.Get("token_id").(string))
	client, err := getClientForUser(d, meta)
	if err != nil {
		return err
	}
	defer doUserLogout(d, client)

	token, err := client.Token.ByID(d.Get("token_id").(string))
	if err != nil {
		if IsNotFound(err) {
			clearUserTokenFields(d)
			return nil
		}
		return err
	}

	err = client.Token.Delete(token)
	if err != nil {
		return err
	}

	clearUserTokenFields(d)

	return nil
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

func getClientForUser(d *schema.ResourceData, meta interface{}) (*managementClient.Client, error) {
	if v, ok := d.Get("temp_token").(string); ok && len(v) > 0 {
		client, err := newManagementClient(meta, v)
		if err != nil {
			if IsNotFound(err) || IsUnauthorized(err) || IsForbidden(err) {
				d.Set("temp_token_id", "")
				d.Set("temp_token", "")
				return doUserLogin(d, meta)
			}
			return nil, err
		}

		return client, nil
	}

	return doUserLogin(d, meta)
}

func newManagementClient(meta interface{}, token string) (*managementClient.Client, error) {
	options := meta.(*Config).CreateClientOpts()
	options.URL = options.URL + rancher2ClientAPIVersion
	options.TokenKey = token
	return managementClient.NewClient(options)
}

func doUserLogin(d *schema.ResourceData, meta interface{}) (*managementClient.Client, error) {
	log.Printf("[DEBUG] Creating Temp Token for User %s", d.Id())
	tempTokenID, tempTokenValue, err := DoUserLogin(meta.(*Config).URL, d.Get("username").(string), d.Get("password").(string), "0", "Temp Terraform API token", meta.(*Config).CACerts, meta.(*Config).Insecure)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Login with %s user: %v", d.Id(), err)
	}

	d.Set("temp_token_id", tempTokenID)
	d.Set("temp_token", tempTokenValue)

	return newManagementClient(meta, tempTokenValue)
}

func doUserLogout(d *schema.ResourceData, client *managementClient.Client) error {
	if v, ok := d.Get("temp_token_id").(string); ok && len(v) > 0 {
		log.Printf("[DEBUG] Deleting Temp Token for User %s", d.Id())
		existing := &norman.Resource{
			ID: v,
			Actions: map[string]string{
				"logout": client.Opts.URL + "/tokens?action=logout",
			},
		}

		err := client.APIBaseClient.Action("token", "logout", existing, map[string]interface{}{}, nil)
		if err != nil {
			return err
		}

		d.Set("temp_token_id", "")
		d.Set("temp_token", "")
	}

	return nil
}

func clearUserTokenFields(d *schema.ResourceData) {
	d.Set("token_id", "")
	d.Set("token_name", "")
	d.Set("token_enabled", false)
	d.Set("token_expired", true)
	d.Set("auth_token", "")
	d.Set("access_key", "")
	d.Set("secret_key", "")
}
