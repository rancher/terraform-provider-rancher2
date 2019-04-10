package rancher2

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

const AzureADConfigName = "azuread"

//Schemas

func authConfigAzureADFields() map[string]*schema.Schema {
	r := authConfigFields()
	s := map[string]*schema.Schema{
		"application_id": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"application_secret": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"auth_endpoint": {
			Type:     schema.TypeString,
			Required: true,
		},
		"endpoint": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "https://login.microsoftonline.com/",
		},
		"graph_endpoint": {
			Type:     schema.TypeString,
			Required: true,
		},
		"rancher_url": {
			Type:     schema.TypeString,
			Required: true,
		},
		"tenant_id": {
			Type:     schema.TypeString,
			Required: true,
		},
		"token_endpoint": {
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

func flattenAuthConfigAzureAD(d *schema.ResourceData, in *managementClient.AzureADConfig) error {
	d.SetId(AzureADConfigName)

	err := d.Set("name", AzureADConfigName)
	if err != nil {
		return err
	}
	err = d.Set("type", managementClient.AzureADConfigType)
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

	err = d.Set("application_id", in.ApplicationID)
	if err != nil {
		return err
	}
	err = d.Set("auth_endpoint", in.AuthEndpoint)
	if err != nil {
		return err
	}
	err = d.Set("endpoint", in.Endpoint)
	if err != nil {
		return err
	}
	err = d.Set("graph_endpoint", in.GraphEndpoint)
	if err != nil {
		return err
	}
	err = d.Set("rancher_url", in.RancherURL)
	if err != nil {
		return err
	}
	err = d.Set("tenant_id", in.TenantID)
	if err != nil {
		return err
	}
	err = d.Set("token_endpoint", in.TokenEndpoint)
	if err != nil {
		return err
	}

	return nil
}

// Expanders

func expandAuthConfigAzureAD(in *schema.ResourceData) (*managementClient.AzureADConfig, error) {
	obj := &managementClient.AzureADConfig{}
	if in == nil {
		return nil, fmt.Errorf("expanding %s Auth Config: Input ResourceData is nil", AzureADConfigName)
	}

	obj.Name = AzureADConfigName
	obj.Type = managementClient.AzureADConfigType

	if v, ok := in.Get("access_mode").(string); ok && len(v) > 0 {
		obj.AccessMode = v
	}

	if v, ok := in.Get("allowed_principal_ids").([]interface{}); ok && len(v) > 0 {
		obj.AllowedPrincipalIDs = toArrayString(v)
	}

	if (obj.AccessMode == "required" || obj.AccessMode == "restricted") && len(obj.AllowedPrincipalIDs) == 0 {
		return nil, fmt.Errorf("expanding %s Auth Config: allowed_principal_ids is required on access_mode %s", AzureADConfigName, obj.AccessMode)
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

	if v, ok := in.Get("application_id").(string); ok && len(v) > 0 {
		obj.ApplicationID = v
	}

	if v, ok := in.Get("application_secret").(string); ok && len(v) > 0 {
		obj.ApplicationSecret = v
	}

	if v, ok := in.Get("auth_endpoint").(string); ok && len(v) > 0 {
		obj.AuthEndpoint = v
	}

	if v, ok := in.Get("endpoint").(string); ok {
		obj.Endpoint = v
	}

	if v, ok := in.Get("graph_endpoint").(string); ok && len(v) > 0 {
		obj.GraphEndpoint = v
	}

	if v, ok := in.Get("rancher_url").(string); ok && len(v) > 0 {
		obj.RancherURL = v
	}

	if v, ok := in.Get("tenant_id").(string); ok && len(v) > 0 {
		obj.TenantID = v
	}

	if v, ok := in.Get("token_endpoint").(string); ok {
		obj.TokenEndpoint = v
	}

	return obj, nil
}

func resourceRancher2AuthConfigAzureAD() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2AuthConfigAzureADCreate,
		Read:   resourceRancher2AuthConfigAzureADRead,
		Update: resourceRancher2AuthConfigAzureADUpdate,
		Delete: resourceRancher2AuthConfigAzureADDelete,

		Schema: authConfigAzureADFields(),
	}
}

func resourceRancher2AuthConfigAzureADCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AzureADConfigName)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to get Auth Config %s: %s", AzureADConfigName, err)
	}

	log.Printf("[INFO] Creating Auth Config AzureAD %s", auth.Name)

	authAzureAD, err := expandAuthConfigAzureAD(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", AzureADConfigName, err)
	}

	// Checking if other auth config is enabled
	if authAzureAD.Enabled {
		err = meta.(*Config).CheckAuthConfigEnabled(AzureADConfigName)
		if err != nil {
			return fmt.Errorf("[ERROR] Checking to enable Auth Config %s: %s", AzureADConfigName, err)
		}
	}

	// Updated auth config
	newAuth := &managementClient.AzureADConfig{}
	err = meta.(*Config).UpdateAuthConfig(auth.Links["self"], authAzureAD, newAuth)
	if err != nil {
		return fmt.Errorf("[ERROR] Updating Auth Config %s: %s", AzureADConfigName, err)
	}

	return resourceRancher2AuthConfigAzureADRead(d, meta)
}

func resourceRancher2AuthConfigAzureADRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Auth Config %s", AzureADConfigName)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AzureADConfigName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AzureADConfigName)
			d.SetId("")
			return nil
		}
		return err
	}

	authAzureAD, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return err
	}

	err = flattenAuthConfigAzureAD(d, authAzureAD.(*managementClient.AzureADConfig))
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2AuthConfigAzureADUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Auth Config %s", AzureADConfigName)

	return resourceRancher2AuthConfigAzureADCreate(d, meta)
}

func resourceRancher2AuthConfigAzureADDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Disabling Auth Config %s", AzureADConfigName)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(AzureADConfigName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", AzureADConfigName)
			d.SetId("")
			return nil
		}
		return err
	}

	if auth.Enabled == true {
		err = client.Post(auth.Actions["disable"], nil, nil)
		if err != nil {
			return fmt.Errorf("[ERROR] Disabling Auth Config %s: %s", AzureADConfigName, err)
		}
	}

	d.SetId("")
	return nil
}
