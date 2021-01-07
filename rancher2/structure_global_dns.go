package rancher2

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners
func flattenGlobalDNS(d *schema.ResourceData, in *managementClient.GlobalDns) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)
	d.Set("fqdn", in.FQDN)
	d.Set("provider_id", in.ProviderID)

	if len(in.Name) > 0 {
		d.Set("name", in.Name)
	}
	if len(in.MultiClusterAppID) > 0 {
		d.Set("multi_cluster_app_id", in.MultiClusterAppID)
	}
	if in.ProjectIDs != nil && len(in.ProjectIDs) > 0 {
		d.Set("project_ids", toArrayInterface(in.ProjectIDs))
	}
	if in.TTL > 0 {
		d.Set("ttl", int(in.TTL))
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

// Expanders
func expandGlobalDNS(in *schema.ResourceData) (*managementClient.GlobalDns, error) {
	obj := &managementClient.GlobalDns{}
	if in == nil {
		return nil, fmt.Errorf("resource rancher2_global_dns_provider data cannot be nil")
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.Name = in.Get("name").(string)
	obj.FQDN = in.Get("fqdn").(string)
	obj.ProviderID = in.Get("provider_id").(string)

	if v, ok := in.Get("project_ids").([]interface{}); ok && len(v) > 0 {
		obj.ProjectIDs = toArrayString(v)
	}

	if v, ok := in.Get("multi_cluster_app_id").(string); ok && len(v) > 0 {
		obj.MultiClusterAppID = v
	}

	if v, ok := in.Get("ttl").(int); ok && v > 0 {
		obj.TTL = int64(v)
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj, nil
}
