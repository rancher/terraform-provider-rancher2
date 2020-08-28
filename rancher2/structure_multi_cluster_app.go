package rancher2

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

const (
	MultiClusterAppTemplatePrefix = "cattle-global-data:"
)

// Flatteners

func flattenMultiClusterAppTemplateVersionID(d *schema.ResourceData, externalID string) string {
	//Global catalog url: catalog://?catalog=demo&template=test&version=1.23.0

	str := strings.TrimPrefix(externalID, AppTemplateExternalIDPrefix)
	values := strings.Split(str, "&")
	out := make(map[string]string, len(values))
	for _, v := range values {
		pair := strings.Split(v, "=")
		if len(pair) != 2 {
			continue
		}
		out[pair[0]] = pair[1]
	}

	d.Set("catalog_name", out["catalog"])
	d.Set("template_name", out["template"])
	d.Set("template_version", out["version"])

	//Template version ID: cattle-global-data:test-test-1.23.0
	return MultiClusterAppTemplatePrefix + out["catalog"] + "-" + out["template"] + "-" + out["version"]
}

func flattenMultiClusterApp(d *schema.ResourceData, in *managementClient.MultiClusterApp, externalID string) error {
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

	d.Set("template_version_id", flattenMultiClusterAppTemplateVersionID(d, externalID))

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
