package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2GlobalDNS() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2GlobalDNSCreate,
		Read:   resourceRancher2GlobalDNSRead,
		Update: resourceRancher2GlobalDNSUpdate,
		Delete: resourceRancher2GlobalDNSDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2GlobalDNSImport,
		},

		Schema: GlobalDNSFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceRancher2GlobalDNSCreate(d *schema.ResourceData, meta interface{}) error {
	globalDNS, err := expandGlobalDNS(d)
	if err != nil {
		return err
	}

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating Global DNS registry %s", globalDNS.FQDN)

	newglobalDNS, err := client.GlobalDns.Create(globalDNS)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    globalDNSStateRefreshFunc(client, newglobalDNS.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for global DNS (%s) to be created: %s", newglobalDNS.ID, waitErr)
	}

	err = flattenGlobalDNS(d, newglobalDNS)
	if err != nil {
		return err
	}

	return resourceRancher2GlobalDNSRead(d, meta)
}

func resourceRancher2GlobalDNSRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Global DNS ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		globalDNS, err := client.GlobalDns.ByID(d.Id())
		if err != nil {
			if IsNotFound(err) {
				log.Printf("[INFO] Global DNS ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(err)
		}

		if err = flattenGlobalDNS(d, globalDNS); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourceRancher2GlobalDNSUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Global DNS ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	globalDNS, err := client.GlobalDns.ByID(d.Id())
	if err != nil {
		return err
	}

	update := map[string]interface{}{
		"providerId":  d.Get("provider_id").(string),
		"annotations": toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":      toMapString(d.Get("labels").(map[string]interface{})),
	}

	if d.HasChange("multi_cluster_app_id") {
		if v, ok := d.Get("multi_cluster_app_id").(string); ok && len(v) > 0 {
			update["multiClusterAppId"] = v
		}
	}

	newglobalDNS, err := client.GlobalDns.Update(globalDNS, update)
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"active"},
		Refresh:    globalDNSStateRefreshFunc(client, newglobalDNS.ID),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for global DNS (%s) to be updated: %s", newglobalDNS.ID, waitErr)
	}

	if d.HasChange("project_ids") {
		if v, ok := d.Get("project_ids").([]interface{}); ok {
			projectsIDS := toArrayString(v)
			projectsToAdd := []string{}
			projectsToRemove := []string{}
			for i := range projectsIDS {
				found := false
				for j := range newglobalDNS.ProjectIDs {
					if projectsIDS[i] == newglobalDNS.ProjectIDs[j] {
						found = true
						break
					}
				}
				if !found {
					projectsToAdd = append(projectsToAdd, projectsIDS[i])
				}
			}
			for i := range newglobalDNS.ProjectIDs {
				found := false
				for j := range projectsIDS {
					if projectsIDS[j] == newglobalDNS.ProjectIDs[i] {
						found = true
						break
					}
				}
				if !found {
					projectsToRemove = append(projectsToRemove, newglobalDNS.ProjectIDs[i])
				}
			}
			if len(newglobalDNS.Actions["addProjects"]) > 0 && len(projectsToAdd) > 0 {
				projectAdd := &managementClient.UpdateGlobalDNSTargetsInput{
					ProjectIDs: projectsToAdd,
				}
				err = client.GlobalDns.ActionAddProjects(newglobalDNS, projectAdd)
				if err != nil {
					return err
				}

			}
			if len(newglobalDNS.Actions["removeProjects"]) > 0 && len(projectsToRemove) > 0 {
				projectRemove := &managementClient.UpdateGlobalDNSTargetsInput{
					ProjectIDs: projectsToRemove,
				}
				err = client.GlobalDns.ActionRemoveProjects(newglobalDNS, projectRemove)
				if err != nil {
					return err
				}
			}
		}
	}

	return resourceRancher2GlobalDNSRead(d, meta)
}

func resourceRancher2GlobalDNSDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Global DNS ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	globalDNS, err := client.GlobalDns.ByID(id)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Global DNS ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.GlobalDns.Delete(globalDNS)
	if err != nil {
		return fmt.Errorf("Error removing Global DNS: %s", err)
	}

	log.Printf("[DEBUG] Waiting for global DNS (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"removed"},
		Refresh:    globalDNSStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for global DNS (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// globalDNSStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher Global Role Binding.
func globalDNSStateRefreshFunc(client *managementClient.Client, globalDNSID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		obj, err := client.GlobalDns.ByID(globalDNSID)
		if err != nil {
			if IsNotFound(err) {
				return obj, "removed", nil
			}
			return nil, "", err
		}

		if obj.Removed != "" {
			return obj, "removed", nil
		}

		return obj, "active", nil
	}
}
