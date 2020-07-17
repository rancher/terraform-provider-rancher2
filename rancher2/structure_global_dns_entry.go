package rancher2

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

// Flatteners
func flattenGlobalDNSEntry(d *schema.ResourceData, in *managementClient.GlobalDNS) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)
	d.Set("fqdn", in.FQDN)
	d.Set("name", in.Name)
	d.Set("multi_cluster_app_id", in.MultiClusterAppID)
	d.Set("project_ids", in.ProjectIDs)
	d.Set("provider_id", in.ProviderID)

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

// Expanders
func expandGlobalDNSEntry(in *schema.ResourceData) (*managementClient.GlobalDNS, error) {
	obj := &managementClient.GlobalDNS{}
	if in == nil {
		return nil, fmt.Errorf("resource rancher2_global_dns_provider data cannot be nil")
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.Name = in.Get("name").(string)
	obj.FQDN = in.Get("fqdn").(string)
	obj.ProviderID = in.Get("provider_id").(string)

	if len(in.Get("project_ids").([]interface{})) > 0 {
		for _, s := range in.Get("project_ids").([]interface{}) {
			obj.ProjectIDs = append(obj.ProjectIDs, s.(string))
		}
	}

	if in.Get("multi_cluster_app_id").(string) != "" {
		obj.MultiClusterAppID = in.Get("multi_cluster_app_id").(string)
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj, nil
}
