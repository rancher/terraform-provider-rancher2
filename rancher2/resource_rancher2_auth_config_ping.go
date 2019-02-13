package rancher2

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

const PingConfigName = "ping"

//Schemas

func authConfigPingFields() map[string]*schema.Schema {
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
			Type:     schema.TypeString,
			Required: true,
		},
		"sp_key": {
			Type:     schema.TypeString,
			Required: true,
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

func flattenAuthConfigPing(d *schema.ResourceData, in *managementClient.PingConfig) error {
	d.SetId(PingConfigName)

	err := d.Set("name", PingConfigName)
	if err != nil {
		return err
	}
	err = d.Set("type", managementClient.PingConfigType)
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

func expandAuthConfigPing(in *schema.ResourceData) (*managementClient.PingConfig, error) {
	obj := &managementClient.PingConfig{}
	if in == nil {
		return nil, fmt.Errorf("expanding %s Auth Config: Input ResourceData is nil", PingConfigName)
	}

	obj.Name = PingConfigName
	obj.Type = managementClient.PingConfigType

	if v, ok := in.Get("access_mode").(string); ok && len(v) > 0 {
		obj.AccessMode = v
	}

	if v, ok := in.Get("allowed_principal_ids").([]interface{}); ok && len(v) > 0 {
		obj.AllowedPrincipalIDs = toArrayString(v)
	}

	if (obj.AccessMode == "required" || obj.AccessMode == "restricted") && len(obj.AllowedPrincipalIDs) == 0 {
		return nil, fmt.Errorf("expanding %s Auth Config: allowed_principal_ids is required on access_mode %s", PingConfigName, obj.AccessMode)
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

func resourceRancher2AuthConfigPing() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2AuthConfigPingCreate,
		Read:   resourceRancher2AuthConfigPingRead,
		Update: resourceRancher2AuthConfigPingUpdate,
		Delete: resourceRancher2AuthConfigPingDelete,

		Schema: authConfigPingFields(),
	}
}

func resourceRancher2AuthConfigPingCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(PingConfigName)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to get Auth Config %s: %s", PingConfigName, err)
	}

	log.Printf("[INFO] Creating Auth Config %s", PingConfigName)

	authPing, err := expandAuthConfigPing(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", PingConfigName, err)
	}

	newAuth := &managementClient.PingConfig{}
	err = meta.(*Config).UpdateAuthConfig(auth.Links["self"], authPing, newAuth)
	if err != nil {
		return fmt.Errorf("[ERROR] Updating Auth Config %s: %s", PingConfigName, err)
	}

	authPingTestAndEnable := managementClient.SamlConfigTestInput{
		FinalRedirectURL: d.Get("final_redirect_url").(string),
	}

	err = client.Post(auth.Actions["testAndEnable"], authPingTestAndEnable, nil)
	if err != nil {
		return fmt.Errorf("[ERROR] Posting Auth Config %s: %s", PingConfigName, err)
	}

	return resourceRancher2AuthConfigPingRead(d, meta)
}

func resourceRancher2AuthConfigPingRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Auth Config %s", PingConfigName)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(PingConfigName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", PingConfigName)
			d.SetId("")
			return nil
		}
		return err
	}

	authPing, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return err
	}

	err = flattenAuthConfigPing(d, authPing.(*managementClient.PingConfig))
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2AuthConfigPingUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Auth Config %s", PingConfigName)

	return resourceRancher2AuthConfigPingCreate(d, meta)
}

func resourceRancher2AuthConfigPingDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Disabling Auth Config %s", PingConfigName)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(PingConfigName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", PingConfigName)
			d.SetId("")
			return nil
		}
		return err
	}

	if auth.Enabled == true {
		err = client.Post(auth.Actions["disable"], nil, nil)
		if err != nil {
			return fmt.Errorf("[ERROR] Disabling Auth Config %s: %s", PingConfigName, err)
		}
	}

	d.SetId("")
	return nil
}
