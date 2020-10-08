package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenNodeDriver(d *schema.ResourceData, in *managementClient.NodeDriver) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)
	d.Set("active", in.Active)
	d.Set("builtin", in.Builtin)
	d.Set("checksum", in.Checksum)
	d.Set("description", in.Description)
	d.Set("name", in.Name)
	d.Set("external_id", in.ExternalID)
	d.Set("ui_url", in.UIURL)
	d.Set("url", in.URL)

	err := d.Set("whitelist_domains", toArrayInterface(in.WhitelistDomains))
	if err != nil {
		return err
	}

	err = d.Set("annotations", toMapInterface(in.Annotations))
	if err != nil {
		return err
	}

	err = d.Set("labels", toMapInterface(in.Labels))
	if err != nil {
		return err
	}

	return nil

}

// Expanders

func expandNodeDriver(in *schema.ResourceData) *managementClient.NodeDriver {
	obj := &managementClient.NodeDriver{}
	if in == nil {
		return nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.Active = in.Get("active").(bool)
	obj.Builtin = in.Get("builtin").(bool)
	obj.Checksum = in.Get("checksum").(string)
	obj.Description = in.Get("description").(string)
	obj.ExternalID = in.Get("external_id").(string)
	obj.Name = in.Get("name").(string)
	obj.UIURL = in.Get("ui_url").(string)
	obj.URL = in.Get("url").(string)

	if v, ok := in.Get("whitelist_domains").([]interface{}); ok && len(v) > 0 {
		obj.WhitelistDomains = toArrayString(v)
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj
}
