package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceRancher2NamespaceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	clusterID, resourceID := splitID(d.Id())

	client, err := meta.(*Config).ClusterClient(clusterID)
	if err != nil {
		return []*schema.ResourceData{}, err
	}
	err = d.Set("project_id", clusterID)
	if err != nil {
		return []*schema.ResourceData{}, err
	}
	ns, err := client.Namespace.ByID(resourceID)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	err = flattenNamespace(d, ns)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
