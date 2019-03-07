package rancher2

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

const ADFSConfigName = "adfs"

//Schemas

func authConfigADFSFields() map[string]*schema.Schema {
	r := authConfigFields()
	s := map[string]*schema.Schema{
		"display_name_field": {
			Type:     schema.TypeString,
			Required: true,
		},
		"final_redirect_url": {
			Type:     schema.TypeString,
			Required: true,
		},
		"groups_field": {
			Type:     schema.TypeString,
			Required: true,
		},
		"idp_metadata_content": {
			Type:     schema.TypeString,
			Required: true,
		},
		"rancher_api_host": {
			Type:     schema.TypeString,
			Required: true,
		},
		"sp_cert": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"sp_key": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"uid_field": {
			Type:     schema.TypeString,
			Required: true,
		},
		"user_name_field": {
			Type:     schema.TypeString,
			Required: true,
		},
	}

	for k, v := range r {
		s[k] = v
	}

	return s
}

// Flatteners

func flattenAuthConfigADFS(d *schema.ResourceData, in *managementClient.ADFSConfig) error {
	d.SetId(ADFSConfigName)

	err := d.Set("name", ADFSConfigName)
	if err != nil {
		return err
	}
	err = d.Set("type", managementClient.ADFSConfigType)
	if err != nil {
		return err
	}
	err = d.Set("access_mode", in.AccessMode)
	if err != nil {
		return err
	}
	err = d.Set("allowed_principal_ids", toArrayInterface(in.AllowedPrincipalIDs))
	if err != nil {
		return err
	}
	err = d.Set("enabled", in.Enabled)
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

	err = d.Set("display_name_field", in.DisplayNameField)
	if err != nil {
		return err
	}
	err = d.Set("groups_field", in.GroupsField)
	if err != nil {
		return err
	}
	err = d.Set("idp_metadata_content", in.IDPMetadataContent)
	if err != nil {
		return err
	}
	err = d.Set("rancher_api_host", in.RancherAPIHost)
	if err != nil {
		return err
	}
	err = d.Set("sp_cert", in.SpCert)
	if err != nil {
		return err
	}
	err = d.Set("sp_key", in.SpKey)
	if err != nil {
		return err
	}
	err = d.Set("uid_field", in.UIDField)
	if err != nil {
		return err
	}
	err = d.Set("user_name_field", in.UserNameField)
	if err != nil {
		return err
	}

	return nil
}

// Expanders

func expandAuthConfigADFS(in *schema.ResourceData) (*managementClient.ADFSConfig, error) {
	obj := &managementClient.ADFSConfig{}
	if in == nil {
		return nil, fmt.Errorf("expanding %s Auth Config: Input ResourceData is nil", ADFSConfigName)
	}

	obj.Name = ADFSConfigName
	obj.Type = managementClient.ADFSConfigType

	if v, ok := in.Get("access_mode").(string); ok && len(v) > 0 {
		obj.AccessMode = v
	}

	if v, ok := in.Get("allowed_principal_ids").([]interface{}); ok && len(v) > 0 {
		obj.AllowedPrincipalIDs = toArrayString(v)
	}

	if (obj.AccessMode == "required" || obj.AccessMode == "restricted") && len(obj.AllowedPrincipalIDs) == 0 {
		return nil, fmt.Errorf("expanding %s Auth Config: allowed_principal_ids is required on access_mode %s", ADFSConfigName, obj.AccessMode)
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

func resourceRancher2AuthConfigADFS() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2AuthConfigADFSCreate,
		Read:   resourceRancher2AuthConfigADFSRead,
		Update: resourceRancher2AuthConfigADFSUpdate,
		Delete: resourceRancher2AuthConfigADFSDelete,

		Schema: authConfigADFSFields(),
	}
}

func resourceRancher2AuthConfigADFSCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(ADFSConfigName)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to get Auth Config %s: %s", ADFSConfigName, err)
	}

	log.Printf("[INFO] Creating Auth Config ADFS %s", auth.Name)

	authADFS, err := expandAuthConfigADFS(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", ADFSConfigName, err)
	}

	newAuth := &managementClient.ADFSConfig{}
	err = meta.(*Config).UpdateAuthConfig(auth.Links["self"], authADFS, newAuth)
	if err != nil {
		return fmt.Errorf("[ERROR] Updating Auth Config %s: %s", ADFSConfigName, err)
	}

	authADFSTestAndEnable := managementClient.SamlConfigTestInput{
		FinalRedirectURL: d.Get("final_redirect_url").(string),
	}

	err = client.Post(auth.Actions["testAndEnable"], authADFSTestAndEnable, nil)
	if err != nil {
		return fmt.Errorf("[ERROR] Posting Auth Config %s: %s", ADFSConfigName, err)
	}

	return resourceRancher2AuthConfigADFSRead(d, meta)
}

func resourceRancher2AuthConfigADFSRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Auth Config %s", ADFSConfigName)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(ADFSConfigName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", ADFSConfigName)
			d.SetId("")
			return nil
		}
		return err
	}

	authADFS, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return err
	}

	err = flattenAuthConfigADFS(d, authADFS.(*managementClient.ADFSConfig))
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2AuthConfigADFSUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Auth Config %s", ADFSConfigName)

	return resourceRancher2AuthConfigADFSCreate(d, meta)
}

func resourceRancher2AuthConfigADFSDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Disabling Auth Config %s", ADFSConfigName)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(ADFSConfigName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", ADFSConfigName)
			d.SetId("")
			return nil
		}
		return err
	}

	if auth.Enabled == true {
		err = client.Post(auth.Actions["disable"], nil, nil)
		if err != nil {
			return fmt.Errorf("[ERROR] Disabling Auth Config %s: %s", ADFSConfigName, err)
		}
	}

	d.SetId("")
	return nil
}
