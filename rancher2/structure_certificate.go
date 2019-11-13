package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	projectClient "github.com/rancher/types/client/project/v3"
)

// Flatteners

func flattenProjectCertificate(d *schema.ResourceData, in *projectClient.Certificate) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)

	d.Set("certs", Base64Encode(in.Certs))
	d.Set("project_id", in.ProjectID)

	if len(in.Description) > 0 {
		d.Set("description", in.Description)
	}

	if len(in.Name) > 0 {
		d.Set("name", in.Name)
	}

	err := d.Set("annotations", toMapInterface(in.Annotations))
	if err != nil {
		return err
	}

	err = d.Set("labels", toMapInterface(in.Labels))
	if err != nil {
		return err
	}

	return nil

}

func flattenNamespacedCertificate(d *schema.ResourceData, in *projectClient.NamespacedCertificate) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)

	d.Set("certs", Base64Encode(in.Certs))
	d.Set("project_id", in.ProjectID)
	d.Set("namespace_id", in.NamespaceId)

	if len(in.Description) > 0 {
		d.Set("description", in.Description)
	}

	if len(in.Name) > 0 {
		d.Set("name", in.Name)
	}

	err := d.Set("annotations", toMapInterface(in.Annotations))
	if err != nil {
		return err
	}

	err = d.Set("labels", toMapInterface(in.Labels))
	if err != nil {
		return err
	}

	return nil

}

func flattenCertificate(d *schema.ResourceData, in interface{}) error {
	namespaceID := d.Get("namespace_id").(string)
	if len(namespaceID) > 0 {
		return flattenNamespacedCertificate(d, in.(*projectClient.NamespacedCertificate))
	}

	return flattenProjectCertificate(d, in.(*projectClient.Certificate))

}

// Expanders

func expandProjectCertificate(in *schema.ResourceData) (*projectClient.Certificate, error) {
	obj := &projectClient.Certificate{}
	if in == nil {
		return nil, nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	if v, ok := in.Get("certs").(string); ok && len(v) > 0 {
		certs, err := Base64Decode(v)
		if err != nil {
			return nil, fmt.Errorf("expanding certificate: certs is not base64 encoded: %s", v)
		}
		obj.Certs = certs
	}

	if v, ok := in.Get("key").(string); ok && len(v) > 0 {
		key, err := Base64Decode(v)
		if err != nil {
			return nil, fmt.Errorf("expanding certificate: key is not base64 encoded: %s", v)
		}
		obj.Key = key
	}

	_, projectID := splitProjectID(in.Get("project_id").(string))
	obj.ProjectID = projectID

	if v, ok := in.Get("description").(string); ok && len(v) > 0 {
		obj.Description = v
	}

	if v, ok := in.Get("name").(string); ok && len(v) > 0 {
		obj.Name = v
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj, nil
}

func expandNamespacedCertificate(in *schema.ResourceData) (*projectClient.NamespacedCertificate, error) {
	obj := &projectClient.NamespacedCertificate{}
	if in == nil {
		return nil, nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	if v, ok := in.Get("certs").(string); ok && len(v) > 0 {
		certs, err := Base64Decode(v)
		if err != nil {
			return nil, fmt.Errorf("expanding certificate: certs is not base64 encoded: %s", v)
		}
		obj.Certs = certs
	}

	if v, ok := in.Get("key").(string); ok && len(v) > 0 {
		key, err := Base64Decode(v)
		if err != nil {
			return nil, fmt.Errorf("expanding certificate: key is not base64 encoded: %s", v)
		}
		obj.Key = key
	}

	_, projectID := splitProjectID(in.Get("project_id").(string))
	obj.ProjectID = projectID

	if v, ok := in.Get("description").(string); ok && len(v) > 0 {
		obj.Description = v
	}

	if v, ok := in.Get("name").(string); ok && len(v) > 0 {
		obj.Name = v
	}

	obj.NamespaceId = in.Get("namespace_id").(string)

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj, nil
}

func expandCertificate(in *schema.ResourceData) (interface{}, error) {
	namespaceID := in.Get("namespace_id").(string)
	if len(namespaceID) > 0 {
		return expandNamespacedCertificate(in)
	}

	return expandProjectCertificate(in)
}
