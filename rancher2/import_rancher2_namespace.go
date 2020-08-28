package rancher2

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	clusterClient "github.com/rancher/rancher/pkg/client/generated/cluster/v3"
)

func resourceRancher2NamespaceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	projectID, resourceID := splitID(d.Id())

	clusterID, projectID := splitProjectID(projectID)

	client, err := meta.(*Config).ClusterClient(clusterID)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	ns, err := client.Namespace.ByID(resourceID)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	d.Set("project_id", clusterID)
	if projectID != "" {
		log.Printf("[INFO] Moving Namespace ID %s to project %s", d.Id(), projectID)
		nsMove := &clusterClient.NamespaceMove{
			ProjectID: projectID,
		}

		err = client.Namespace.ActionMove(ns, nsMove)
		if err != nil {
			return []*schema.ResourceData{}, err
		}
		d.Set("project_id", projectID)
	}

	err = flattenNamespace(d, ns)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
