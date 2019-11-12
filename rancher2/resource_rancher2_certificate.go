package rancher2

import (
	"fmt"
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

	err := meta.(*Config).ProjectExist(projectID)
	if err != nil {
		return err
	}

	certificate, err := expandCertificate(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating Certificate %s on Project ID %s", name, projectID)

	newCertificate, err := meta.(*Config).CreateCertificate(certificate)
	if err != nil {
		return err
	}

	err = flattenCertificate(d, newCertificate)
	if err != nil {
		return err
	}

	return resourceRancher2CertificateRead(d, meta)
}

func resourceRancher2CertificateRead(d *schema.ResourceData, meta interface{}) error {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	id := d.Id()
	namespaceID := d.Get("namespace_id").(string)

	log.Printf("[INFO] Refreshing Certificate ID %s", id)

	certificate, err := meta.(*Config).GetCertificate(id, projectID, namespaceID)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Certificate ID %s not found.", id)
			d.SetId("")
			return nil
		}
		return err
	}

	return flattenCertificate(d, certificate)
}

func resourceRancher2CertificateUpdate(d *schema.ResourceData, meta interface{}) error {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	id := d.Id()
	namespaceID := d.Get("namespace_id").(string)

	log.Printf("[INFO] Updating Certificate ID %s", id)

	certificate, err := meta.(*Config).GetCertificate(id, projectID, namespaceID)
	if err != nil {
		return err
	}

	certs, err := Base64Decode(d.Get("certs").(string))
	if err != nil {
		return fmt.Errorf("Updating certificate: certs is not base64 encoded")
	}

	key, err := Base64Decode(d.Get("key").(string))
	if err != nil {
		return fmt.Errorf("Updating certificate: key is not base64 encoded")
	}

	update := map[string]interface{}{
		"description": d.Get("description").(string),
		"certs":       certs,
		"keys":        key,
		"annotations": toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":      toMapString(d.Get("labels").(map[string]interface{})),
	}

	newCertificate, err := meta.(*Config).UpdateCertificate(certificate, update)
	if err != nil {
		return err
	}

	err = flattenCertificate(d, newCertificate)
	if err != nil {
		return err
	}

	return resourceRancher2CertificateRead(d, meta)
}

func resourceRancher2CertificateDelete(d *schema.ResourceData, meta interface{}) error {
	_, projectID := splitProjectID(d.Get("project_id").(string))
	id := d.Id()
	namespaceID := d.Get("namespace_id").(string)

	log.Printf("[INFO] Deleting Certificate ID %s", id)

	certificate, err := meta.(*Config).GetCertificate(id, projectID, namespaceID)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Certificate ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = meta.(*Config).DeleteCertificate(certificate)
	if err != nil {
		return fmt.Errorf("Error removing Certificate: %s", err)
	}

	d.SetId("")
	return nil
}
