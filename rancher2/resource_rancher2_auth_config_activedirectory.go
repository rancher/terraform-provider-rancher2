package rancher2

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

const ActiveDirectoryConfigName = "activedirectory"

func authConfigActiveDirectoryFields() map[string]*schema.Schema {
	r := authConfigFields()
	s := map[string]*schema.Schema{
		"servers": {
			Type:     schema.TypeList,
			Required: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
		},
		"service_account_username": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"service_account_password": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"user_search_base": {
			Type:     schema.TypeString,
			Required: true,
		},
		"username": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"password": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"certificate": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
		"connection_timeout": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  5000,
		},
		"default_login_domain": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"group_dn_attribute": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"group_member_mapping_attribute": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"group_member_user_attribute": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"group_name_attribute": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"group_object_class": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"group_search_attribute": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"group_search_base": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"group_search_filter": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"nested_group_membership_enabled": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
		"port": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  389,
		},
		"user_disabled_bit_mask": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  2,
		},
		"user_enabled_attribute": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"user_login_attribute": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"user_name_attribute": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"user_object_class": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"user_search_attribute": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"user_search_filter": {
			Type:     schema.TypeString,
			Optional: true,
			Computed: true,
		},
		"tls": {
			Type:     schema.TypeBool,
			Optional: true,
			Computed: true,
		},
	}

	for k, v := range r {
		s[k] = v
	}

	return s
}

// Flatteners

func flattenAuthConfigActiveDirectory(d *schema.ResourceData, in *managementClient.ActiveDirectoryConfig) error {
	d.SetId(ActiveDirectoryConfigName)

	err := d.Set("name", ActiveDirectoryConfigName)
	if err != nil {
		return err
	}
	err = d.Set("type", managementClient.ActiveDirectoryConfigType)
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
	err = d.Set("servers", toArrayInterface(in.Servers))
	if err != nil {
		return err
	}
	err = d.Set("service_account_username", in.ServiceAccountUsername)
	if err != nil {
		return err
	}
	err = d.Set("user_search_base", in.UserSearchBase)
	if err != nil {
		return err
	}
	err = d.Set("certificate", in.Certificate)
	if err != nil {
		return err
	}
	err = d.Set("connection_timeout", int(in.ConnectionTimeout))
	if err != nil {
		return err
	}
	err = d.Set("default_login_domain", in.DefaultLoginDomain)
	if err != nil {
		return err
	}
	err = d.Set("group_dn_attribute", in.GroupDNAttribute)
	if err != nil {
		return err
	}
	err = d.Set("group_member_mapping_attribute", in.GroupMemberMappingAttribute)
	if err != nil {
		return err
	}
	err = d.Set("group_member_user_attribute", in.GroupMemberUserAttribute)
	if err != nil {
		return err
	}
	err = d.Set("group_name_attribute", in.GroupNameAttribute)
	if err != nil {
		return err
	}
	err = d.Set("group_object_class", in.GroupObjectClass)
	if err != nil {
		return err
	}
	err = d.Set("group_search_attribute", in.GroupSearchAttribute)
	if err != nil {
		return err
	}
	err = d.Set("group_search_base", in.GroupSearchBase)
	if err != nil {
		return err
	}
	err = d.Set("group_search_filter", in.GroupSearchFilter)
	if err != nil {
		return err
	}
	err = d.Set("nested_group_membership_enabled", *in.NestedGroupMembershipEnabled)
	if err != nil {
		return err
	}
	err = d.Set("port", int(in.Port))
	if err != nil {
		return err
	}
	err = d.Set("tls", in.TLS)
	if err != nil {
		return err
	}
	err = d.Set("user_disabled_bit_mask", int(in.UserDisabledBitMask))
	if err != nil {
		return err
	}
	err = d.Set("user_enabled_attribute", in.UserEnabledAttribute)
	if err != nil {
		return err
	}
	err = d.Set("user_login_attribute", in.UserLoginAttribute)
	if err != nil {
		return err
	}
	err = d.Set("user_name_attribute", in.UserNameAttribute)
	if err != nil {
		return err
	}
	err = d.Set("user_object_class", in.UserObjectClass)
	if err != nil {
		return err
	}
	err = d.Set("user_search_attribute", in.UserSearchAttribute)
	if err != nil {
		return err
	}
	err = d.Set("user_search_filter", in.UserSearchFilter)
	if err != nil {
		return err
	}

	return nil
}

// Expanders

func expandAuthConfigActiveDirectory(in *schema.ResourceData) (*managementClient.ActiveDirectoryConfig, error) {
	obj := &managementClient.ActiveDirectoryConfig{}
	if in == nil {
		return nil, fmt.Errorf("expanding %s Auth Config: Input ResourceData is nil", ActiveDirectoryConfigName)
	}

	obj.Name = ActiveDirectoryConfigName
	obj.Type = managementClient.ActiveDirectoryConfigType

	if v, ok := in.Get("access_mode").(string); ok && len(v) > 0 {
		obj.AccessMode = v
	}

	if v, ok := in.Get("allowed_principal_ids").([]interface{}); ok && len(v) > 0 {
		obj.AllowedPrincipalIDs = toArrayString(v)
	}

	if (obj.AccessMode == "required" || obj.AccessMode == "restricted") && len(obj.AllowedPrincipalIDs) == 0 {
		return nil, fmt.Errorf("expanding %s Auth Config: allowed_principal_ids is required on access_mode %s", ActiveDirectoryConfigName, obj.AccessMode)
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

	if v, ok := in.Get("servers").([]interface{}); ok && len(v) > 0 {
		obj.Servers = toArrayString(v)
	}

	if v, ok := in.Get("service_account_username").(string); ok && len(v) > 0 {
		obj.ServiceAccountUsername = v
	}

	if v, ok := in.Get("service_account_password").(string); ok && len(v) > 0 {
		obj.ServiceAccountPassword = v
	}

	if v, ok := in.Get("user_search_base").(string); ok && len(v) > 0 {
		obj.UserSearchBase = v
	}

	if v, ok := in.Get("certificate").(string); ok && len(v) > 0 {
		obj.Certificate = v
	}

	if v, ok := in.Get("connection_timeout").(int); ok && v > 0 {
		obj.ConnectionTimeout = int64(v)
	}

	if v, ok := in.Get("default_login_domain").(string); ok && len(v) > 0 {
		obj.DefaultLoginDomain = v
	}

	if v, ok := in.Get("group_dn_attribute").(string); ok && len(v) > 0 {
		obj.GroupDNAttribute = v
	}

	if v, ok := in.Get("group_member_mapping_attribute").(string); ok && len(v) > 0 {
		obj.GroupMemberMappingAttribute = v
	}

	if v, ok := in.Get("group_member_user_attribute").(string); ok && len(v) > 0 {
		obj.GroupMemberUserAttribute = v
	}

	if v, ok := in.Get("group_name_attribute").(string); ok && len(v) > 0 {
		obj.GroupNameAttribute = v
	}

	if v, ok := in.Get("group_object_class").(string); ok && len(v) > 0 {
		obj.GroupObjectClass = v
	}

	if v, ok := in.Get("group_search_attribute").(string); ok && len(v) > 0 {
		obj.GroupSearchAttribute = v
	}

	if v, ok := in.Get("group_search_base").(string); ok && len(v) > 0 {
		obj.GroupSearchBase = v
	}

	if v, ok := in.Get("group_search_filter").(string); ok && len(v) > 0 {
		obj.GroupSearchFilter = v
	}

	if v, ok := in.Get("nested_group_membership_enabled").(bool); ok {
		obj.NestedGroupMembershipEnabled = &v
	}

	if v, ok := in.Get("port").(int); ok && v > 0 {
		obj.Port = int64(v)
	}

	if v, ok := in.Get("tls").(bool); ok {
		obj.TLS = v
	}

	if v, ok := in.Get("user_disabled_bit_mask").(int); ok && v > 0 {
		obj.UserDisabledBitMask = int64(v)
	}

	if v, ok := in.Get("user_enabled_attribute").(string); ok && len(v) > 0 {
		obj.UserEnabledAttribute = v
	}

	if v, ok := in.Get("user_login_attribute").(string); ok && len(v) > 0 {
		obj.UserLoginAttribute = v
	}

	if v, ok := in.Get("user_name_attribute").(string); ok && len(v) > 0 {
		obj.UserNameAttribute = v
	}

	if v, ok := in.Get("user_object_class").(string); ok && len(v) > 0 {
		obj.UserObjectClass = v
	}

	if v, ok := in.Get("user_search_attribute").(string); ok && len(v) > 0 {
		obj.UserSearchAttribute = v
	}

	if v, ok := in.Get("user_search_filter").(string); ok && len(v) > 0 {
		obj.UserSearchFilter = v
	}

	return obj, nil
}

func resourceRancher2AuthConfigActiveDirectory() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancher2AuthConfigActiveDirectoryCreate,
		Read:   resourceRancher2AuthConfigActiveDirectoryRead,
		Update: resourceRancher2AuthConfigActiveDirectoryUpdate,
		Delete: resourceRancher2AuthConfigActiveDirectoryDelete,

		Schema: authConfigActiveDirectoryFields(),
	}
}

func resourceRancher2AuthConfigActiveDirectoryCreate(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(ActiveDirectoryConfigName)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed to get Auth Config %s: %s", ActiveDirectoryConfigName, err)
	}

	log.Printf("[INFO] Creating Auth Config ActiveDirectory %s", auth.Name)

	authActiveDirectory, err := expandAuthConfigActiveDirectory(d)
	if err != nil {
		return fmt.Errorf("[ERROR] Failed expanding Auth Config %s: %s", ActiveDirectoryConfigName, err)
	}

	authActiveDirectoryTestAndApply := managementClient.ActiveDirectoryTestAndApplyInput{
		ActiveDirectoryConfig: authActiveDirectory,
		Enabled:               authActiveDirectory.Enabled,
		Username:              d.Get("username").(string),
		Password:              d.Get("password").(string),
	}

	err = client.Post(auth.Actions["testAndApply"], authActiveDirectoryTestAndApply, nil)
	if err != nil {
		return fmt.Errorf("[ERROR] Posting Auth Config %s: %s", ActiveDirectoryConfigName, err)
	}

	return resourceRancher2AuthConfigActiveDirectoryRead(d, meta)
}

func resourceRancher2AuthConfigActiveDirectoryRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Auth Config %s", ActiveDirectoryConfigName)
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(ActiveDirectoryConfigName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", ActiveDirectoryConfigName)
			d.SetId("")
			return nil
		}
		return err
	}

	authActiveDirectory, err := meta.(*Config).GetAuthConfig(auth)
	if err != nil {
		return err
	}

	err = flattenAuthConfigActiveDirectory(d, authActiveDirectory.(*managementClient.ActiveDirectoryConfig))
	if err != nil {
		return err
	}

	return nil
}

func resourceRancher2AuthConfigActiveDirectoryUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Auth Config %s", ActiveDirectoryConfigName)

	return resourceRancher2AuthConfigActiveDirectoryCreate(d, meta)
}

func resourceRancher2AuthConfigActiveDirectoryDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Disabling Auth Config %s", ActiveDirectoryConfigName)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	auth, err := client.AuthConfig.ByID(ActiveDirectoryConfigName)
	if err != nil {
		if IsNotFound(err) {
			log.Printf("[INFO] Auth Config %s not found.", ActiveDirectoryConfigName)
			d.SetId("")
			return nil
		}
		return err
	}

	if auth.Enabled == true {
		err = client.Post(auth.Actions["disable"], nil, nil)
		if err != nil {
			return fmt.Errorf("[ERROR] Disabling Auth Config %s: %s", ActiveDirectoryConfigName, err)
		}
	}

	d.SetId("")
	return nil
}
