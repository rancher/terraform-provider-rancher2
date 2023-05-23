package rancher2

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/rancher/norman/types"
)

func resourceRancher2SecretV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2SecretV2Create,
		Read:   resourceRancher2SecretV2Read,
		Update: resourceRancher2SecretV2Update,
		Delete: resourceRancher2SecretV2Delete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2SecretV2Import,
		},
		Schema: secretV2Fields(),
		CustomizeDiff: customdiff.ForceNewIf("immutable", func(d *schema.ResourceDiff, m interface{}) bool {
			if d.HasChange("immutable") {
				return !d.Get("immutable").(bool)
			}
			return d.Get("immutable").(bool)
		}),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2SecretV2Create(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	secret := expandSecretV2(d)

	log.Printf("[INFO] Creating Secret V2 %s", name)

	newSecret, err := createSecretV2(meta.(*Config), clusterID, secret)
	if err != nil {
		return err
	}
	d.SetId(clusterID + secretV2ClusterIDsep + newSecret.ID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    secretV2StateRefreshFunc(meta, clusterID, newSecret.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for secret (%s) to be active: %s", newSecret.ID, waitErr)
	}
	return resourceRancher2SecretV2Read(d, meta)
}

func resourceRancher2SecretV2Read(d *schema.ResourceData, meta interface{}) error {
	clusterID, rancherID := splitID(d.Id())
	log.Printf("[INFO] Refreshing Secret V2 %s at Cluster ID %s", rancherID, clusterID)

	return resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		secret, err := getSecretV2ByID(meta.(*Config), clusterID, rancherID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Secret V2 %s not found at cluster ID %s", rancherID, clusterID)
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(err)
		}

		if err = flattenSecretV2(d, secret); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourceRancher2SecretV2Update(d *schema.ResourceData, meta interface{}) error {
	clusterID, rancherID := splitID(d.Id())
	secret := expandSecretV2(d)
	log.Printf("[INFO] Updating Secret V2 %s at Cluster ID %s", rancherID, clusterID)

	newSecret, err := updateSecretV2(meta.(*Config), clusterID, rancherID, secret)
	if err != nil {
		return err
	}
	d.SetId(clusterID + secretV2ClusterIDsep + newSecret.ID)
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"active"},
		Refresh:    secretV2StateRefreshFunc(meta, clusterID, newSecret.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for secret (%s) to be active: %s", newSecret.ID, waitErr)
	}
	return resourceRancher2SecretV2Read(d, meta)
}

func resourceRancher2SecretV2Delete(d *schema.ResourceData, meta interface{}) error {
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	log.Printf("[INFO] Deleting Secret V2 %s", name)

	_, rancherID := splitID(d.Id())
	secret, err := getSecretV2ByID(meta.(*Config), clusterID, rancherID)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			d.SetId("")
			return nil
		}
	}
	err = deleteSecretV2(meta.(*Config), clusterID, secret)
	if err != nil {
		return err
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"removed"},
		Refresh:    secretV2StateRefreshFunc(meta, clusterID, secret.ID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf("[ERROR] waiting for secret (%s) to be active: %s", secret.ID, waitErr)
	}
	d.SetId("")
	return nil
}

// secretV2StateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher Secret v2.
func secretV2StateRefreshFunc(meta interface{}, clusterID, secretID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := getSecretV2ByID(meta.(*Config), clusterID, secretID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}
		return obj, "active", nil
	}
}

// Rancher2 Secret V2 API CRUD functions
func createSecretV2(c *Config, clusterID string, secret *SecretV2) (*SecretV2, error) {
	if c == nil {
		return nil, fmt.Errorf("Creating secret V2: Provider config is nil")
	}
	if clusterID == "" {
		return nil, fmt.Errorf("Creating secret V2: Cluster ID is nil")
	}
	if secret == nil {
		return nil, fmt.Errorf("Creating secret V2: object is nil")
	}
	// Converting secret V2 object to map[string]interface{} as type fields is duplicated
	secret2, err := interfaceToMap(secret)
	if err != nil {
		return nil, err
	}
	secret2["type"] = secret2["_type"]
	resp := &SecretV2{}
	err = c.createObjectV2(clusterID, secretV2APIType, secret2, resp)
	if err != nil {
		return nil, fmt.Errorf("Creating secret V2: %s", err)
	}
	return resp, nil
}

func deleteSecretV2(c *Config, clusterID string, obj *SecretV2) error {
	if c == nil {
		return fmt.Errorf("Deleting secret V2: Provider config is nil")
	}
	if clusterID == "" {
		return fmt.Errorf("Deleting secret V2: Cluster ID is nil")
	}
	if obj == nil {
		return fmt.Errorf("Deleting secret V2: object is nil")
	}
	resource := &types.Resource{
		ID:      obj.ID,
		Type:    secretV2APIType,
		Links:   obj.Links,
		Actions: obj.Actions,
	}
	return c.deleteObjectV2(clusterID, resource)
}

func getSecretV2ByID(c *Config, clusterID, id string) (*SecretV2, error) {
	if c == nil {
		return nil, fmt.Errorf("Getting secret V2: Provider config is nil")
	}
	if len(clusterID) == 0 || len(id) == 0 {
		return nil, fmt.Errorf("Getting secret V2: Cluster ID and/or Secret V2 ID is nil")
	}
	resp := &SecretV2{}
	err := c.getObjectV2ByID(clusterID, id, secretV2APIType, resp)
	if err != nil {
		if !IsServerError(err) && !IsNotFound(err) && !IsForbidden(err) {
			return nil, fmt.Errorf("Getting secret V2: %s", err)
		}
		return nil, err
	}
	return resp, nil
}

func updateSecretV2(c *Config, clusterID, id string, obj *SecretV2) (*SecretV2, error) {
	if c == nil {
		return nil, fmt.Errorf("Updating secret V2: Provider config is nil")
	}
	if len(clusterID) == 0 || len(id) == 0 {
		return nil, fmt.Errorf("Updating secret V2: Cluster ID and/or Secret V2 ID is nil")
	}
	if obj == nil {
		return nil, fmt.Errorf("Updating secret V2: Cluster V2 is nil")
	}
	// Converting secret V2 object to map[string]interface{} as type fields is duplicated
	updateMap, err := interfaceToMap(obj)
	if err != nil {
		return nil, err
	}
	updateMap["type"] = updateMap["_type"]
	resp := &SecretV2{}
	ctx, cancel := context.WithTimeout(context.Background(), c.Timeout)
	defer cancel()
	for {
		err := c.updateObjectV2(clusterID, id, secretV2APIType, updateMap, resp)
		if err == nil {
			return resp, err
		}
		if !IsServerError(err) && !IsUnknownSchemaType(err) && !IsConflict(err) {
			return nil, err
		}
		if IsConflict(err) {
			// Read secret again and update ObjectMeta.ResourceVersion before retry
			newObj := &SecretV2{}
			err = c.getObjectV2ByID(clusterID, id, secretV2APIType, newObj)
			if err != nil {
				return nil, err
			}
			obj.ObjectMeta.ResourceVersion = newObj.ObjectMeta.ResourceVersion
			// Converting secret V2 object to map[string]interface{} as type fields is duplicated
			updateMap, err := interfaceToMap(obj)
			if err != nil {
				return nil, err
			}
			updateMap["type"] = updateMap["_type"]
		}
		select {
		case <-time.After(rancher2RetriesWait * time.Second):
		case <-ctx.Done():
			return nil, fmt.Errorf("Timeout updating secret V2 ID %s: %v", id, err)
		}
	}
}
