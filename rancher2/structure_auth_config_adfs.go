package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenAuthConfigADFS(d *schema.ResourceData, in *managementClient.ADFSConfig) error {
	d.SetId(AuthConfigADFSName)
	d.Set("name", AuthConfigADFSName)
	d.Set("type", managementClient.ADFSConfigType)
	d.Set("access_mode", in.AccessMode)

	err := d.Set("allowed_principal_ids", toArrayInterface(in.AllowedPrincipalIDs))
	if err != nil {
		return err
	}

	d.Set("enabled", in.Enabled)

	err = d.Set("annotations", toMapInterface(in.Annotations))
	if err != nil {
		return err
	}
	err = d.Set("labels", toMapInterface(in.Labels))
	if err != nil {
		return err
	}

	d.Set("display_name_field", in.DisplayNameField)
	d.Set("groups_field", in.GroupsField)
	d.Set("idp_metadata_content", in.IDPMetadataContent)
	d.Set("rancher_api_host", in.RancherAPIHost)
	d.Set("sp_cert", in.SpCert)
	d.Set("uid_field", in.UIDField)
	d.Set("user_name_field", in.UserNameField)

	return nil
}

// Expanders

func expandAuthConfigADFS(in *schema.ResourceData) (*managementClient.ADFSConfig, error) {
	obj := &managementClient.ADFSConfig{}
	if in == nil {
		return nil, fmt.Errorf("expanding %s Auth Config: Input ResourceData is nil", AuthConfigADFSName)
	}

	obj.Name = AuthConfigADFSName
	obj.Type = managementClient.ADFSConfigType

	if v, ok := in.Get("access_mode").(string); ok && len(v) > 0 {
		obj.AccessMode = v
	}

	if v, ok := in.Get("allowed_principal_ids").([]interface{}); ok && len(v) > 0 {
		obj.AllowedPrincipalIDs = toArrayString(v)
	}

	if (obj.AccessMode == "required" || obj.AccessMode == "restricted") && len(obj.AllowedPrincipalIDs) == 0 {
		return nil, fmt.Errorf("expanding %s Auth Config: allowed_principal_ids is required on access_mode %s", AuthConfigADFSName, obj.AccessMode)
	}

	if v, ok := in.Get("enabled").(bool); ok {
		obj.Enabled = v
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	if v, ok := in.Get("display_name_field").(string); ok && len(v) > 0 {
		obj.DisplayNameField = v
	}

	if v, ok := in.Get("groups_field").(string); ok && len(v) > 0 {
		obj.GroupsField = v
	}

	if v, ok := in.Get("idp_metadata_content").(string); ok && len(v) > 0 {
		obj.IDPMetadataContent = v
	}

	if v, ok := in.Get("rancher_api_host").(string); ok && len(v) > 0 {
		obj.RancherAPIHost = v
	}

	if v, ok := in.Get("sp_cert").(string); ok && len(v) > 0 {
		obj.SpCert = v
	}

	if v, ok := in.Get("sp_key").(string); ok && len(v) > 0 {
		obj.SpKey = v
	}

	if v, ok := in.Get("uid_field").(string); ok && len(v) > 0 {
		obj.UIDField = v
	}

	if v, ok := in.Get("user_name_field").(string); ok && len(v) > 0 {
		obj.UserNameField = v
	}

	return obj, nil
}
