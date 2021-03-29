package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

	newSecret, err := meta.(*Config).CreateSecretV2(clusterID, secret)
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

	secret, err := meta.(*Config).GetSecretV2ByID(clusterID, rancherID)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] Secret V2 %s not found at cluster ID %s", rancherID, clusterID)
			d.SetId("")
			return nil
		}
		return err
	}
	return flattenSecretV2(d, secret)
}

func resourceRancher2SecretV2Update(d *schema.ResourceData, meta interface{}) error {
	clusterID, rancherID := splitID(d.Id())
	secret := expandSecretV2(d)
	log.Printf("[INFO] Updating Secret V2 %s at Cluster ID %s", rancherID, clusterID)

	newSecret, err := meta.(*Config).UpdateSecretV2(clusterID, rancherID, secret)
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
	secret, err := meta.(*Config).GetSecretV2ByID(clusterID, rancherID)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			d.SetId("")
			return nil
		}
	}
	err = meta.(*Config).DeleteSecretV2(clusterID, secret)
	if err != nil {
		return err
	}
	stateConf := &resource.StateChangeConf{
		Pending:    []string{},
		Target:     []string{"removed"},
		Refresh:    secretV2StateRefreshFunc(meta, clusterID, secret.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
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
		obj, err := meta.(*Config).GetSecretV2ByID(clusterID, secretID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}
		return obj, "active", nil
	}
}
