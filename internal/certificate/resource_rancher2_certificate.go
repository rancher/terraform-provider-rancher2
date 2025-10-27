package rancher2

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRancher2Certificate() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2CertificateCreate,
		Read:   resourceRancher2CertificateRead,
		Update: resourceRancher2CertificateUpdate,
		Delete: resourceRancher2CertificateDelete,

		Schema: certificateFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
	}
}

func resourceRancher2CertificateCreate(d *schema.ResourceData, meta interface{}) error {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	name := d.Get("name").(string)

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		err := meta.(*Config).ProjectExist(projectID)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		certificate, err := expandCertificate(d)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		log.Printf("[INFO] Creating Certificate %s on Project ID %s", name, projectID)

		newCertificate, err := meta.(*Config).CreateCertificate(certificate)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		err = flattenCertificate(d, newCertificate)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		err = resourceRancher2CertificateRead(d, meta)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourceRancher2CertificateRead(d *schema.ResourceData, meta interface{}) error {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	id := d.Id()
	namespaceID := d.Get("namespace_id").(string)

	log.Printf("[INFO] Refreshing Certificate ID %s", id)

	return resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		certificate, err := meta.(*Config).GetCertificate(id, projectID, namespaceID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Certificate ID %s not found.", id)
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(err)
		}

		if err := flattenCertificate(d, certificate); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourceRancher2CertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	id := d.Id()
	namespaceID := d.Get("namespace_id").(string)

	log.Printf("[INFO] Updating Certificate ID %s", id)

	return resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		certificate, err := meta.(*Config).GetCertificate(id, projectID, namespaceID)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		update, err := expandCertificate(d)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		newCertificate, err := meta.(*Config).UpdateCertificate(certificate, update)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		err = flattenCertificate(d, newCertificate)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		err = resourceRancher2CertificateRead(d, meta)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourceRancher2CertificateDelete(d *schema.ResourceData, meta interface{}) error {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	id := d.Id()
	namespaceID := d.Get("namespace_id").(string)

	log.Printf("[INFO] Deleting Certificate ID %s", id)

	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		certificate, err := meta.(*Config).GetCertificate(id, projectID, namespaceID)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Certificate ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(err)
		}

		err = meta.(*Config).DeleteCertificate(certificate)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("[ERROR] Error removing Certificate: %s", err))
		}

		d.SetId("")
		return nil
	})
}
