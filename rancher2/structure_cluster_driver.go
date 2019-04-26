package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

// Flatteners

func flattenClusterDriver(d *schema.ResourceData, in *managementClient.KontainerDriver) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)

	err := d.Set("active", in.Active)
	if err != nil {
		return err
	}

	err = d.Set("actual_url", in.ActualURL)
	if err != nil {
		return err
	}

	err = d.Set("builtin", in.BuiltIn)
	if err != nil {
		return err
	}

	err = d.Set("checksum", in.Checksum)
	if err != nil {
		return err
	}

	err = d.Set("name", in.Name)
	if err != nil {
		return err
	}

	err = d.Set("ui_url", in.UIURL)
	if err != nil {
		return err
	}

	err = d.Set("url", in.URL)
	if err != nil {
		return err
	}

	err = d.Set("whitelist_domains", toArrayInterface(in.WhitelistDomains))
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

func expandClusterDriver(in *schema.ResourceData) *managementClient.KontainerDriver {
	obj := &managementClient.KontainerDriver{}
	if in == nil {
		return nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.Active = in.Get("active").(bool)
	obj.ActualURL = in.Get("actual_url").(string)
	obj.BuiltIn = in.Get("builtin").(bool)
	obj.Checksum = in.Get("checksum").(string)
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
