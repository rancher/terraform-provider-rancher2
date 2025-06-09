package rancher2

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRancher2NamespaceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	// in this context the clusterID is rancher2_cluster_v2.cluster_v1_id or rancher2_cluster.id
	// the input is expected to be "project_id.namespace_id"
	// project_id is a concatenation of "cluster_id:project_id" or data.rancher2_cluster.downstream_cluster.default_project_id
	log.Printf("[INFO] Splitting given id %s", d.Id())
	projectID, resourceID := splitID(d.Id())
	log.Printf("[INFO] Splitting project %s", projectID)
	clusterID, projectID := splitProjectID(projectID)
	log.Printf("[INFO] Using cluster id %s", clusterID)

	client, err := meta.(*Config).ClusterClient(clusterID)
	if err != nil {
		log.Printf("[INFO] Problem getting cluster client for cluster with id \"%s\"", clusterID)
		return []*schema.ResourceData{}, err
	}

	ns, err := client.Namespace.ByID(resourceID)
	if err != nil {
		return []*schema.ResourceData{}, err
	}
	if ns.ProjectID != projectID {
		return []*schema.ResourceData{}, fmt.Errorf("[ERROR] Project ID \"%s\" in import object doesn't match resource to import (\"%s\").", projectID, ns.ProjectID)
	}

	err = flattenNamespace(d, ns)
	if err != nil {
		return []*schema.ResourceData{}, err
	}

	return []*schema.ResourceData{d}, nil
}
