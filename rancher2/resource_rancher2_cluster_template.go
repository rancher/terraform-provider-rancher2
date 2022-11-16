package rancher2

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
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
		Schema:        clusterTemplateFields(),
		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceRancher2ClusterTemplateResourceV0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceRancher2ClusterTemplateStateUpgradeV0,
				Version: 0,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		CustomizeDiff: customdiff.Sequence(
			customdiff.IfValueChange("template_revisions",
				func(old, new, meta interface{}) bool {
					return true
				},
				func(d *schema.ResourceDiff, meta interface{}) error {
					if !d.HasChange("template_revisions") {
						return nil
					}
					old, new := d.GetChange("template_revisions")
					oldInput := old.([]interface{})
					oldInputLen := len(oldInput)
					newInput := new.([]interface{})
					newInputLen := len(newInput)
					// Indexing old and new inout by ID
					oldInputIndexName := map[string]int{}
					for i := range oldInput {
						if row, ok := oldInput[i].(map[string]interface{}); ok {
							if v, ok := row["name"].(string); ok {
								oldInputIndexName[v] = i
							}
						}
					}
					// Sorting new input
					sortedNewInput := make([]interface{}, len(newInput))
					newRows := []interface{}{}
					lastIndex := 0
					for i := range newInput {
						if row, ok := newInput[i].(map[string]interface{}); ok {
							if name, ok := row["name"].(string); ok {
								if v, ok := oldInputIndexName[name]; ok {
									if oldRow, ok := oldInput[v].(map[string]interface{}); ok {
										oldRow["default"] = row["default"]
										row = oldRow
									}
									if v > i && oldInputLen > newInputLen {
										v = v - (v - i)
									}
									sortedNewInput[v] = row
									lastIndex++
									continue
								}
							}
							row["id"] = ""
							newRows = append(newRows, row)
						}
					}
					for i := range newRows {
						sortedNewInput[lastIndex+i] = newRows[i]
					}
					return d.SetNew("template_revisions", sortedNewInput)
				}),
			customdiff.ValidateValue("template_revisions", func(val, meta interface{}) error {
				hasDefault := false
				names := map[string]int{}
				input := val.([]interface{})
				for i := range input {
					if obj, ok := input[i].(map[string]interface{}); ok {
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
				}
				if !hasDefault {
					return fmt.Errorf("[ERROR] Validating cluster template revisions: NO default defined")
				}
				return nil
			}),
		),
	}
}

func resourceRancher2ClusterTemplateResourceV0() *schema.Resource {
	return &schema.Resource{
		Schema: clusterTemplateFieldsV0(),
	}
}

func resourceRancher2ClusterTemplateStateUpgradeV0(rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
	if tmplRevisions, ok := rawState["template_revisions"].([]interface{}); ok && len(tmplRevisions) > 0 {
		for i1 := range tmplRevisions {
			if tmplRevision, ok := tmplRevisions[i1].(map[string]interface{}); ok && len(tmplRevision) > 0 {
				if clusterConfigs, ok := tmplRevision["cluster_config"].([]interface{}); ok && len(clusterConfigs) > 0 {
					for i2 := range clusterConfigs {
						if clusterConfig, ok := clusterConfigs[i2].(map[string]interface{}); ok && len(clusterConfig) > 0 {
							newValue, err := resourceRancher2ClusterStateUpgradeV0(clusterConfig, meta)
							if err != nil {
								return nil, fmt.Errorf("Upgrading Cluster Template schema V0: %v", err)
							}
							rawState["template_revisions"].([]interface{})[i1].(map[string]interface{})["cluster_config"].([]interface{})[i2] = newValue
						}
					}
				}
			}
		}
	}
	return rawState, nil
}

func resourceRancher2ClusterTemplateCreate(d *schema.ResourceData, meta interface{}) error {
	ctrIndex, clusterTemplate, clusterTemplateRevisions, err := expandClusterTemplate(d)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Creating Cluster Template %s", clusterTemplate.Name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		newClusterTemplate, err := client.ClusterTemplate.Create(clusterTemplate)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		newClusterTemplateRevisions, err := clusterTemplateRevisionsCreate(client, newClusterTemplate.ID, clusterTemplateRevisions)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		d.SetId(newClusterTemplate.ID)

		// Update defaultRevisionId if needed
		if len(newClusterTemplateRevisions) > 0 {
			update := map[string]interface{}{
				"defaultRevisionId": newClusterTemplateRevisions[ctrIndex].ID,
			}

			if _, err = client.ClusterTemplate.Update(newClusterTemplate, update); err != nil {
				return resource.NonRetryableError(err)
			}
		}

		if err = resourceRancher2ClusterTemplateRead(d, meta); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourceRancher2ClusterTemplateRead(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	log.Printf("[INFO] Refreshing Cluster Template ID %s", id)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutRead), func() *resource.RetryError {
		clusterTemplate, err := client.ClusterTemplate.ByID(id)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Cluster Template ID %s not found.", clusterTemplate.ID)
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(err)
		}

		clusterTemplateRevisions, err := clusterTemplateRevisionsRead(client, id)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		if err = flattenClusterTemplate(d, clusterTemplate, clusterTemplateRevisions); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourceRancher2ClusterTemplateUpdate(d *schema.ResourceData, meta interface{}) error {
	id := d.Id()
	log.Printf("[INFO] Updating Cluster Template ID %s", id)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		clusterTemplate, err := client.ClusterTemplate.ByID(id)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		clusterTemplateRevisions := make([]managementClient.ClusterTemplateRevision, 0)
		if d.HasChange("template_revisions") {
			defaultRevisionID, templateRevisions, err := clusterTemplateRevisionsUpdate(client, id, d)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			clusterTemplateRevisions = templateRevisions
			d.Set("default_revision_id", defaultRevisionID)
		}

		update := map[string]interface{}{
			"defaultRevisionId": d.Get("default_revision_id").(string),
			"description":       d.Get("description").(string),
			"members":           expandMembers(d.Get("members").([]interface{})),
			"name":              d.Get("name").(string),
			"annotations":       toMapString(d.Get("annotations").(map[string]interface{})),
			"labels":            toMapString(d.Get("labels").(map[string]interface{})),
		}

		newClusterTemplate, err := client.ClusterTemplate.Update(clusterTemplate, update)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		if len(clusterTemplateRevisions) > 0 {
			err = flattenClusterTemplate(d, newClusterTemplate, clusterTemplateRevisions)
			if err != nil {
				return resource.NonRetryableError(err)
			}
			// Delete removed clusterTemplateRevisions
			err = clusterTemplateRevisionsDelete(client, id, clusterTemplateRevisions)
			if err != nil {
				return resource.NonRetryableError(err)
			}
		}

		if err = resourceRancher2ClusterTemplateRead(d, meta); err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})
}

func resourceRancher2ClusterTemplateDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Cluster Template ID %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		clusterTemplate, err := client.ClusterTemplate.ByID(id)
		if err != nil {
			if IsNotFound(err) || IsForbidden(err) {
				log.Printf("[INFO] Cluster Template ID %s not found.", d.Id())
				d.SetId("")
				return nil
			}
			return resource.NonRetryableError(err)
		}

		err = client.ClusterTemplate.Delete(clusterTemplate)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("[ERROR] Error removing Cluster Template: %s", err))
		}

		d.SetId("")
		return nil
	})
}

func clusterTemplateRevisionsCreate(client *managementClient.Client, ctID string, ctrs []managementClient.ClusterTemplateRevision) ([]managementClient.ClusterTemplateRevision, error) {
	if len(ctID) == 0 {
		return nil, fmt.Errorf("[ERROR] Creating revision: cluster Template ID can't be empty")
	}

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

func clusterTemplateRevisionsRead(client *managementClient.Client, ctID string) ([]managementClient.ClusterTemplateRevision, error) {
	if len(ctID) == 0 {
		return nil, fmt.Errorf("[ERROR] Reading revision: Cluster Template ID can't be empty")
	}
	filters := map[string]interface{}{}
	filters["clusterTemplateId"] = ctID

	clusterTemplateRevisions, err := client.ClusterTemplateRevision.List(NewListOpts(filters))
	return clusterTemplateRevisions.Data, err
}

func clusterTemplateRevisionsUpdate(client *managementClient.Client, ctID string, d *schema.ResourceData) (string, []managementClient.ClusterTemplateRevision, error) {
	if len(ctID) == 0 {
		return "", nil, fmt.Errorf("[ERROR] Updating revision: Cluster Template ID can't be empty")
	}

	old, new := d.GetChange("template_revisions")
	oldData, _ := old.([]interface{})
	data, _ := new.([]interface{})
	indexDefault, ctrs, err := expandClusterTemplateRevisions(data)
	if err != nil {
		return "", nil, err
	}
	clusterTemplateRevisions := make([]managementClient.ClusterTemplateRevision, len(data))
	oldDataIndexName := map[string]int{}
	for i := range oldData {
		if row, ok := oldData[i].(map[string]interface{}); ok {
			if v, ok := row["name"].(string); ok {
				oldDataIndexName[v] = i
			}
		}
	}
	dataIndexName := map[string]int{}
	for i := range data {
		if row, ok := data[i].(map[string]interface{}); ok {
			if v, ok := row["name"].(string); ok {
				dataIndexName[v] = i
			}
		}
	}
	for i := range data {
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
		// Update existing clusterTemplateRevisions if changed
		oldIndex, oldOK := oldDataIndexName[ctrs[i].Name]
		newIndex, newOK := dataIndexName[ctrs[i].Name]
		if oldOK && newOK {
			oldRow, oldOK := oldData[oldIndex].(map[string]interface{})
			in, oldOK := data[newIndex].(map[string]interface{})
			if oldOK && newOK {
				oldRow["default"] = in["default"]
				if !AreEqual(oldRow, in) {
					clusterConfig, err := expandClusterSpecBase(in["cluster_config"].([]interface{}))
					if err != nil {
						return "", nil, err
					}
					enabled := in["enabled"].(bool)
					update := map[string]interface{}{
						"clusterConfig": clusterConfig,
						"enabled":       &enabled,
						"name":          in["name"].(string),
						"questions":     expandQuestions(in["questions"].([]interface{})),
						"annotations":   toMapString(in["annotations"].(map[string]interface{})),
						"labels":        toMapString(in["labels"].(map[string]interface{})),
					}
					ctr, err := client.ClusterTemplateRevision.ByID(ctrs[i].ID)
					if err != nil {
						return "", nil, fmt.Errorf("Getting ClusterTemplateRevision %s: %v", ctrs[i].ID, err)
					}
					newCtr, err := client.ClusterTemplateRevision.Update(ctr, update)
					if err != nil {
						return "", nil, fmt.Errorf("Updating ClusterTemplateRevision %s: %v", ctr.ID, err)
					}
					clusterTemplateRevisions[i] = *newCtr
					continue
				}
			}
		}
		clusterTemplateRevisions[i] = ctrs[i]
	}
	return clusterTemplateRevisions[indexDefault].ID, clusterTemplateRevisions, err
}

func clusterTemplateRevisionsDelete(client *managementClient.Client, ctID string, ctrs []managementClient.ClusterTemplateRevision) error {
	readClusterTemplateRevisions, err := clusterTemplateRevisionsRead(client, ctID)
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
					return fmt.Errorf("Error removing Cluster Template Revision [%s]: %s", readCtr.ID, err)
				}
			}
		}
	}

	return nil
}
