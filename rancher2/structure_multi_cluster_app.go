package rancher2

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

const (
	MultiClusterAppTemplatePrefix = "cattle-global-data:"
)

// Flatteners

func flattenMultiClusterAppTemplateVersionID(d *schema.ResourceData, in string) {
	//Template version ID: cattle-global-data:test-test-1.23.0

	data := strings.TrimPrefix(in, MultiClusterAppTemplatePrefix)
	separator := "-"

	first := strings.Index(data, separator)
	last := strings.LastIndex(data, separator)

	d.Set("template_version_id", in)
	d.Set("catalog_name", data[:first])
	d.Set("template_name", data[first+1:last])
	d.Set("template_version", data[last+1:])
}

func flattenMultiClusterApp(d *schema.ResourceData, in *managementClient.MultiClusterApp) error {
	if in == nil {
		return fmt.Errorf("[ERROR] flattening multi cluster app: Input setting is nil")
	}

	d.SetId(in.ID)
	d.Set("name", in.Name)

	err := d.Set("roles", toArrayInterface(in.Roles))
	if err != nil {
		return err
	}

	err = d.Set("targets", flattenTargets(in.Targets))
	if err != nil {
		return err
	}

	flattenMultiClusterAppTemplateVersionID(d, in.TemplateVersionID)

	err = d.Set("answers", flattenAnswers(in.Answers))
	if err != nil {
		return err
	}

	err = d.Set("members", flattenMembers(in.Members))
	if err != nil {
		return err
	}

	d.Set("revision_history_limit", int(in.RevisionHistoryLimit))

	if in.Status != nil {
		d.Set("revision_id", in.Status.RevisionID)
	}

	err = d.Set("upgrade_strategy", flattenUpgradeStrategy(in.UpgradeStrategy))
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

func expandMultiClusterAppTemplateVersionID(in *schema.ResourceData) string {
	//Template version ID: cattle-global-data:test-test-1.23.0

	catalogName := in.Get("catalog_name").(string)
	appName := in.Get("template_name").(string)
	appVersion := in.Get("template_version").(string)

	return MultiClusterAppTemplatePrefix + catalogName + "-" + appName + "-" + appVersion
}

func expandMultiClusterApp(in *schema.ResourceData) (*managementClient.MultiClusterApp, error) {
	obj := &managementClient.MultiClusterApp{}
	if in == nil {
		return nil, fmt.Errorf("[ERROR] expanding multi cluster app: Input ResourceData is nil")
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.Name = in.Get("name").(string)

	if v, ok := in.Get("roles").([]interface{}); ok && len(v) > 0 {
		obj.Roles = toArrayString(v)
	}

	if v, ok := in.Get("targets").([]interface{}); ok && len(v) > 0 {
		obj.Targets = expandTargets(v)
	}

	obj.TemplateVersionID = expandMultiClusterAppTemplateVersionID(in)

	if v, ok := in.Get("answers").([]interface{}); ok && len(v) > 0 {
		obj.Answers = expandAnswers(v)
	}

	if v, ok := in.Get("members").([]interface{}); ok && len(v) > 0 {
		obj.Members = expandMembers(v)
	}

	if v, ok := in.Get("revision_history_limit").(int); ok && v > 0 {
		obj.RevisionHistoryLimit = int64(v)
	}

	if v, ok := in.Get("revision_id").(string); ok && len(v) > 0 {
		if obj.Status == nil {
			obj.Status = &managementClient.MultiClusterAppStatus{}
		}
		obj.Status.RevisionID = v
	}

	if v, ok := in.Get("upgrade_strategy").([]interface{}); ok && len(v) > 0 {
		obj.UpgradeStrategy = expandUpgradeStrategy(v)
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj, nil
}
