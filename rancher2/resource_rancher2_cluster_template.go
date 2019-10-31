package rancher2

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

func resourceRancher2ClusterTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2ClusterTemplateCreate,
		Read:   resourceRancher2ClusterTemplateRead,
		Update: resourceRancher2ClusterTemplateUpdate,
		Delete: resourceRancher2ClusterTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: resourceRancher2ClusterTemplateImport,
		},
		Schema: clusterTemplateFields(),
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		CustomizeDiff: customdiff.Sequence(
			customdiff.IfValueChange("template_revisions",
				func(old, new, meta interface{}) bool {
					return len(new.([]interface{})) < len(old.([]interface{}))
				},
				func(d *schema.ResourceDiff, meta interface{}) error {
					old, new := d.GetChange("template_revisions")
					oldInput := old.([]interface{})
					newInput := new.([]interface{})
					// Setting template_revisions order if is removed
					if len(newInput) != len(oldInput) {
						for i := range newInput {
							oldObj := oldInput[i].(map[string]interface{})
							newObj := newInput[i].(map[string]interface{})
							oldName := oldObj["name"].(string)
							newName := newObj["name"].(string)
							if oldName != newName {
								for j := range oldInput {
									oldObj = oldInput[j].(map[string]interface{})
									oldName = oldObj["name"].(string)
									if oldName == newName {
										newInput[i] = oldInput[j]
										break
									}
								}
							}
						}
						err := d.SetNew("template_revisions", newInput)
						if err != nil {
							return err
						}
					}
					return nil
				}),
			customdiff.ValidateValue("template_revisions", func(val, meta interface{}) error {
				hasDefault := false
				names := map[string]int{}
				input := val.([]interface{})
				for i := range input {
					obj := input[i].(map[string]interface{})
					if v, ok := obj["default"].(bool); ok && v {
						if hasDefault {
							return fmt.Errorf("[ERROR] Validating cluster template revisions: more than one default defined")
						}
						hasDefault = true
					}
					if v, ok := obj["name"].(string); ok && len(v) > 0 {
						names[v]++
						if names[v] > 1 {
							return fmt.Errorf("[ERROR] Validating cluster template revisions: name \"%s\" is repeated", v)
						}
					}
				}
				if !hasDefault {
					return fmt.Errorf("[ERROR] Validating cluster template revisions: NO default defined")
				}
				return nil
			}),
		),
	}
}

func resourceRancher2ClusterTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	clusterTemplate, clusterTemplateRevisions, err := expandClusterTemplate(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating Cluster Template %s", clusterTemplate.Name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	newClusterTemplate, err := client.ClusterTemplate.Create(clusterTemplate)
	if err != nil {
		return err
	}

	newClusterTemplateRevisions, err := clusterTemplateRevisionsCreate(client, newClusterTemplate.ID, clusterTemplateRevisions)
	if err != nil {
		return err
	}

	err = flattenClusterTemplate(d, newClusterTemplate, newClusterTemplateRevisions)
	if err != nil {
		return err
	}

	// Update defaultRevisionId if needed
	if len(newClusterTemplateRevisions) > 1 {
		update := map[string]interface{}{
			"defaultRevisionId": d.Get("default_revision_id").(string),
		}

		_, err = client.ClusterTemplate.Update(newClusterTemplate, update)
		if err != nil {
			return err
		}
	}

	return resourceRancher2ClusterTemplateRead(d, meta)
}

func resourceRancher2ClusterTemplateRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	log.Printf("[INFO] Refreshing Cluster Template ID %s", id)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	clusterTemplate, err := client.ClusterTemplate.ByID(id)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Cluster Template ID %s not found.", clusterTemplate.ID)
			d.SetId("")
			return nil
		}
		return err
	}

	clusterTemplateRevisions, err := clusterTemplateRevisionsRead(client, clusterTemplate.ID, d.Get("template_revisions").([]interface{}))
	if err != nil {
		return err
	}

	return flattenClusterTemplate(d, clusterTemplate, clusterTemplateRevisions)
}

func resourceRancher2ClusterTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	log.Printf("[INFO] Updating Cluster Template ID %s", id)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	clusterTemplate, err := client.ClusterTemplate.ByID(id)
	if err != nil {
		return err
	}

	defaultRevisionId := d.Get("default_revision_id").(string)
	clusterTemplateRevisions := []managementClient.ClusterTemplateRevision{}
	if d.HasChange("template_revisions") {
		defaultRevisionId, clusterTemplateRevisions, err = clusterTemplateRevisionsUpdate(client, id, d.Get("template_revisions").([]interface{}))
		if err != nil {
			return err
		}
	}

	update := map[string]interface{}{
		"defaultRevisionId": defaultRevisionId,
		"description":       d.Get("description").(string),
		"members":           expandMembers(d.Get("members").([]interface{})),
		"name":              d.Get("name").(string),
		"annotations":       toMapString(d.Get("annotations").(map[string]interface{})),
		"labels":            toMapString(d.Get("labels").(map[string]interface{})),
	}

	_, err = client.ClusterTemplate.Update(clusterTemplate, update)
	if err != nil {
		return err
	}

	if len(clusterTemplateRevisions) > 0 {
		err = flattenClusterTemplate(d, clusterTemplate, clusterTemplateRevisions)
		if err != nil {
			return err
		}
		// Delete removed clusterTemplateRevisions
		err = clusterTemplateRevisionsDelete(client, id, clusterTemplateRevisions, d.Get("template_revisions").([]interface{}))
		if err != nil {
			return err
		}
	}

	return resourceRancher2ClusterTemplateRead(d, meta)
}

func resourceRancher2ClusterTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Cluster Template ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	clusterTemplate, err := client.ClusterTemplate.ByID(id)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Cluster Template ID %s not found.", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	err = client.ClusterTemplate.Delete(clusterTemplate)
	if err != nil {
		return fmt.Errorf("Error removing Cluster Template: %s", err)
	}

	d.SetId("")
	return nil
}

func clusterTemplateRevisionsCreate(client *managementClient.Client, ctID string, ctrs []managementClient.ClusterTemplateRevision) ([]managementClient.ClusterTemplateRevision, error) {
	clusterTemplateRevisions := make([]managementClient.ClusterTemplateRevision, len(ctrs))

	for i := range ctrs {
		ctrs[i].ClusterTemplateID = ctID
		clusterTemplateRevision, err := client.ClusterTemplateRevision.Create(&ctrs[i])
		if err != nil {
			return nil, err
		}
		clusterTemplateRevisions[i] = *clusterTemplateRevision
	}

	return clusterTemplateRevisions, nil
}

func clusterTemplateRevisionsRead(client *managementClient.Client, ctID string, data []interface{}) ([]managementClient.ClusterTemplateRevision, error) {
	filters := map[string]interface{}{}
	filters["clusterTemplateId"] = ctID

	clusterTemplateRevisions, err := client.ClusterTemplateRevision.List(NewListOpts(filters))
	if err != nil {
		return nil, err
	}

	// Sorting clusterTemplateRevisions.Data array by data interface
	sorted := make([]managementClient.ClusterTemplateRevision, len(clusterTemplateRevisions.Data))
	newCTR := []managementClient.ClusterTemplateRevision{}
	lastIndex := 0
	for i := range sorted {
		found := false
		for j := range data {
			row := data[j].(map[string]interface{})
			if (row["id"] != "" && clusterTemplateRevisions.Data[i].ID == row["id"]) || (row["id"] == "" && clusterTemplateRevisions.Data[i].Name == row["name"]) {
				sorted[j] = clusterTemplateRevisions.Data[i]
				found = true
				lastIndex++
				break
			}
		}
		if !found {
			newCTR = append(newCTR, clusterTemplateRevisions.Data[i])
		}
	}

	for i := range newCTR {
		sorted[lastIndex+i] = newCTR[i]
	}

	return sorted, nil
}

func clusterTemplateRevisionsUpdate(client *managementClient.Client, ctID string, data []interface{}) (string, []managementClient.ClusterTemplateRevision, error) {
	ctrs, err := expandClusterTemplateRevisions(data)
	if err != nil {
		return "", nil, err
	}

	clusterTemplateRevisions := make([]managementClient.ClusterTemplateRevision, len(data))
	for i := range data {
		in := data[i].(map[string]interface{})
		// Create new clusterTemplateRevisions
		if len(ctrs[i].ID) == 0 {
			ctrs[i].ClusterTemplateID = ctID
			newCtr, err := client.ClusterTemplateRevision.Create(&ctrs[i])
			if err != nil {
				return "", nil, err
			}
			clusterTemplateRevisions[i] = *newCtr
			continue
		}
		// Update existing clusterTemplateRevisions
		ctr, err := client.ClusterTemplateRevision.ByID(ctrs[i].ID)
		if err != nil {
			return "", nil, err
		}
		enabled := in["enabled"].(bool)
		update := map[string]interface{}{
			"clusterConfig": expandClusterSpecBase(in["cluster_config"].([]interface{})),
			"enabled":       &enabled,
			"name":          in["name"].(string),
			"questions":     expandQuestions(in["questions"].([]interface{})),
			"annotations":   toMapString(in["annotations"].(map[string]interface{})),
			"labels":        toMapString(in["labels"].(map[string]interface{})),
		}

		newCtr, err := client.ClusterTemplateRevision.Update(ctr, update)
		if err != nil {
			return "", nil, err
		}

		clusterTemplateRevisions[i] = *newCtr
	}

	ctrID, _, err := flattenClusterTemplateRevisions(clusterTemplateRevisions, "", data)

	return ctrID, clusterTemplateRevisions, nil
}

func clusterTemplateRevisionsDelete(client *managementClient.Client, ctID string, ctrs []managementClient.ClusterTemplateRevision, data []interface{}) error {
	readClusterTemplateRevisions, err := clusterTemplateRevisionsRead(client, ctID, data)
	if err != nil {
		return err
	}

	if len(readClusterTemplateRevisions) != len(ctrs) {
		for _, readCtr := range readClusterTemplateRevisions {
			found := false
			for _, ctr := range ctrs {
				if readCtr.ID == ctr.ID {
					found = true
					break
				}
			}
			if !found {
				err = client.ClusterTemplateRevision.Delete(&readCtr)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
