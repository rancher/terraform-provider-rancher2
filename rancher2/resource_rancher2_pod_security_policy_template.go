package rancher2

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"log"
	"time"
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

	err = flattenPodSecurityPolicyTemplate(d, newPodSecurityPolicyTemplate)
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

	pspt, err := client.PodSecurityPolicyTemplate.ByID(d.Id())
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] PodSecurityPolicyTemplate with ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = flattenPodSecurityPolicyTemplate(d, pspt)
	if err != nil {
		return err
	}

	return nil
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
	log.Printf("[INFO] Deleting PodSecurityPolicyTemplate with ID %s", d.Id())
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	pspt, err := client.PodSecurityPolicyTemplate.ByID(d.Id())
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] PodSecurityPolicyTemplate with ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.PodSecurityPolicyTemplate.Delete(pspt)
	if err != nil {
		return fmt.Errorf("Error removing PodSecurityPolicyTemplate: %s", err)
	}

	d.SetId("")
	return nil
}
