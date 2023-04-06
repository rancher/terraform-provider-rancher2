package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2PodSecurityPolicyTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2PodSecurityPolicyTemplateCreate,
		Read:   resourceRancher2PodSecurityPolicyTemplateRead,
		Update: resourceRancher2PodSecurityPolicyTemplateUpdate,
		Delete: resourceRancher2PodSecurityPolicyTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2PodSecurityPolicyTemplateImport,
		},

		Schema: podSecurityPolicyTemplateFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2PodSecurityPolicyTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	podSecurityPolicyTemplate := expandPodSecurityPolicyTemplate(d)

	log.Printf("[INFO] Creating PodSecurityPolicyTemplate %s", podSecurityPolicyTemplate.Name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	newPodSecurityPolicyTemplate, err := client.PodSecurityPolicyTemplate.Create(podSecurityPolicyTemplate)
	if err != nil {
		return err
	}

	d.SetId(newPodSecurityPolicyTemplate.ID)

	return resourceRancher2PodSecurityPolicyTemplateRead(d, meta)
}

func resourceRancher2PodSecurityPolicyTemplateRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing PodSecurityPolicyTemplate with ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		pspt, err := client.PodSecurityPolicyTemplate.ByID(d.Id())
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] PodSecurityPolicyTemplate with ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(err)
		}

		if err = flattenPodSecurityPolicyTemplate(d, pspt); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourceRancher2PodSecurityPolicyTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating PodSecurityPolicyTemplate with ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	pspt, err := client.PodSecurityPolicyTemplate.ByID(d.Id())
	if err != nil {
		return err
	}

	update := expandPodSecurityPolicyTemplate(d)

	_, err = client.PodSecurityPolicyTemplate.Update(pspt, update)
	if err != nil {
		return err
	}

	return resourceRancher2PodSecurityPolicyTemplateRead(d, meta)
}

func resourceRancher2PodSecurityPolicyTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	log.Printf("[INFO] Deleting PodSecurityPolicyTemplate with ID %s", id)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	pspt, err := client.PodSecurityPolicyTemplate.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] PodSecurityPolicyTemplate with ID %s not found.", id)
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.PodSecurityPolicyTemplate.Delete(pspt)
	if err != nil {
		return fmt.Errorf("Error removing PodSecurityPolicyTemplate: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"removed"},
		Refresh:    podSecurityPolicyTemplateStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for PodSecurityPolicyTemplate (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// podSecurityPolicyTemplateStateRefreshFunc returns a resource.StateRefreshFunc, used to watch a Rancher PodSecurityPolicyTemplate
func podSecurityPolicyTemplateStateRefreshFunc(client *managementClient.Client, pspID string) resource.StateRefreshFunc {
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
