package rancher2

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	projectClient "github.com/rancher/rancher/pkg/client/generated/project/v3"
)

const (
	AppTemplateExternalIDPrefix  = "catalog://?"
	AppCatalogClusterLocalPrefix = "local"
	AppCatalogClusterPrefix      = "c-"
	AppCatalogProjectPrefix      = "p-"
)

// Flatteners

func flattenAppExternalID(d *schema.ResourceData, in string) {
	//Global catalog url: catalog://?catalog=demo&template=test&version=1.23.0
	//Cluster catalog url: catalog://?catalog=c-XXXXX/test&type=clusterCatalog&template=test&version=1.23.0
	//Project catalog url: catalog://?catalog=p-XXXXX/test&type=projectCatalog&template=test&version=1.23.0

	str := strings.TrimPrefix(in, AppTemplateExternalIDPrefix)
	values := strings.Split(str, "&")
	out := make(map[string]string, len(values))
	for _, v := range values {
		pair := strings.Split(v, "=")
		if len(pair) != 2 {
			continue
		}
		if pair[0] == "catalog" && (strings.HasPrefix(pair[1], AppCatalogClusterLocalPrefix) || strings.HasPrefix(pair[1], AppCatalogClusterPrefix) || strings.HasPrefix(pair[1], AppCatalogProjectPrefix)) {
			pair[1] = strings.Replace(pair[1], "/", ":", -1)
		}
		out[pair[0]] = pair[1]
	}

	d.Set("external_id", in)
	d.Set("catalog_name", out["catalog"])
	d.Set("template_name", out["template"])
	d.Set("template_version", out["version"])
}

func flattenApp(d *schema.ResourceData, in *projectClient.App) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)

	flattenAppExternalID(d, in.ExternalID)

	d.Set("name", in.Name)
	d.Set("project_id", in.ProjectID)
	d.Set("target_namespace", in.TargetNamespace)

	if len(in.Answers) > 0 {
		err := d.Set("answers", toMapInterface(in.Answers))
		if err != nil {
			return err
		}
	}

	if len(in.Description) > 0 {
		d.Set("description", in.Description)
	}

	if len(in.AppRevisionID) > 0 {
		d.Set("revision_id", in.AppRevisionID)
	}

	if len(in.ValuesYaml) > 0 {
		d.Set("values_yaml", Base64Encode(in.ValuesYaml))
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

func expandAppExternalID(in *schema.ResourceData) string {
	//Global catalog url: catalog://?catalog=test&template=test&version=1.23.0
	//Cluster catalog url: catalog://?catalog=c-XXXXX/test&type=clusterCatalog&template=test&version=1.23.0
	//Project catalog url: catalog://?catalog=p-XXXXX/test&type=projectCatalog&template=test&version=1.23.0

	catalogName := in.Get("catalog_name").(string)
	appName := in.Get("template_name").(string)
	appVersion := in.Get("template_version").(string)

	if strings.HasPrefix(catalogName, AppCatalogClusterLocalPrefix) || strings.HasPrefix(catalogName, AppCatalogClusterPrefix) {
		catalogName = strings.Replace(catalogName, ":", "/", -1)
		catalogName = catalogName + "&type=clusterCatalog"
	}
	if strings.HasPrefix(catalogName, AppCatalogProjectPrefix) {
		catalogName = strings.Replace(catalogName, ":", "/", -1)
		catalogName = catalogName + "&type=projectCatalog"
	}

	catalogPart := "catalog=" + catalogName
	appNamePart := "&template=" + appName
	appVersionPart := "&version=" + appVersion

	return AppTemplateExternalIDPrefix + catalogPart + appNamePart + appVersionPart
}

func expandApp(in *schema.ResourceData) (*projectClient.App, error) {
	obj := &projectClient.App{}
	if in == nil {
		return nil, nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.ExternalID = expandAppExternalID(in)
	obj.Name = in.Get("name").(string)
	obj.ProjectID = in.Get("project_id").(string)
	obj.TargetNamespace = in.Get("target_namespace").(string)

	if v, ok := in.Get("answers").(map[string]interface{}); ok && len(v) > 0 {
		obj.Answers = toMapString(v)
	}

	if v, ok := in.Get("description").(string); ok && len(v) > 0 {
		obj.Description = v
	}

	if v, ok := in.Get("revision_id").(string); ok && len(v) > 0 {
		obj.AppRevisionID = v
	}

	if v, ok := in.Get("values_yaml").(string); ok && len(v) > 0 {
		values, err := Base64Decode(v)
		if err != nil {
			return nil, fmt.Errorf("expanding app: values_yaml is not base64 encoded: %s", v)
		}
		obj.ValuesYaml = values
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	obj.Timeout = int64(in.Timeout(schema.TimeoutCreate).Seconds())

	return obj, nil
}
