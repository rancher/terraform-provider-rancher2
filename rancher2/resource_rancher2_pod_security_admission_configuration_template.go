package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func resourceRancher2PodSecurityAdmissionConfigurationTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2PodSecurityAdmissionConfigurationTemplateCreate,
		Read:   resourceRancher2PodSecurityAdmissionConfigurationTemplateRead,
		Update: resourceRancher2PodSecurityAdmissionConfigurationTemplateUpdate,
		Delete: resourceRancher2PodSecurityAdmissionConfigurationTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2PodSecurityAdmissionConfigurationTemplateImport,
		},

		Schema: podSecurityAdmissionConfigurationTemplateFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2PodSecurityAdmissionConfigurationTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	podSecurityAdmissionConfigurationTemplate, err := expandPodSecurityAdmissionConfigurationTemplate(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating PodSecurityAdmissionConfigurationTemplate %s", podSecurityAdmissionConfigurationTemplate.Name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	newPodSecurityAdmissionConfigurationTemplate, err := client.PodSecurityAdmissionConfigurationTemplate.Create(podSecurityAdmissionConfigurationTemplate)
	if err != nil {
		return err
	}

	d.SetId(newPodSecurityAdmissionConfigurationTemplate.ID)

	return resourceRancher2PodSecurityAdmissionConfigurationTemplateRead(d, meta)
}

func resourceRancher2PodSecurityAdmissionConfigurationTemplateRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing PodSecurityAdmissionConfigurationTemplate with ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		psact, err := client.PodSecurityAdmissionConfigurationTemplate.ByID(d.Id())
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] PodSecurityAdmissionConfigurationTemplate with ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(err)
		}

		if err = flattenPodSecurityAdmissionConfigurationTemplate(d, psact); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourceRancher2PodSecurityAdmissionConfigurationTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating PodSecurityAdmissionConfigurationTemplate with ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	psact, err := client.PodSecurityAdmissionConfigurationTemplate.ByID(d.Id())
	if err != nil {
		return err
	}

	update, err := expandPodSecurityAdmissionConfigurationTemplate(d)
	if err != nil {
		return err
	}

	_, err = client.PodSecurityAdmissionConfigurationTemplate.Update(psact, update)
	if err != nil {
		return err
	}

	return resourceRancher2PodSecurityAdmissionConfigurationTemplateRead(d, meta)
}

func resourceRancher2PodSecurityAdmissionConfigurationTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	log.Printf("[INFO] Deleting PodSecurityAdmissionConfigurationTemplate with ID %s", id)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	psact, err := client.PodSecurityAdmissionConfigurationTemplate.ByID(id)
	if err != nil {
		if IsNotFound(err) || IsForbidden(err) {
			log.Printf("[INFO] PodSecurityAdmissionConfigurationTemplate with ID %s not found.", id)
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.PodSecurityAdmissionConfigurationTemplate.Delete(psact)
	if err != nil {
		return fmt.Errorf("[ERROR] removing PodSecurityAdmissionConfigurationTemplate: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active"},
		Target:     []string{"removed"},
		Refresh:    podSecurityAdmissionConfigurationTemplateStateRefreshFunc(client, id),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"[ERROR] waiting for PodSecurityAdmissionConfigurationTemplate (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// podSecurityAdmissionConfigurationTemplateStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// a Rancher PodSecurityAdmissionConfiguration Template
func podSecurityAdmissionConfigurationTemplateStateRefreshFunc(client *managementClient.Client, pspID string) resource.StateRefreshFunc {
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
