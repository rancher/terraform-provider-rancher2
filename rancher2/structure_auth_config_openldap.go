package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenAuthConfigOpenLdap(d *schema.ResourceData, in *managementClient.LdapConfig) error {
	d.SetId(AuthConfigOpenLdapName)
	d.Set("name", AuthConfigOpenLdapName)
	d.Set("type", managementClient.OpenLdapConfigType)

	err := flattenAuthConfigLdap(d, in)
	if err != nil {
		return err
	}

	return nil
}

// Expanders

func expandAuthConfigOpenLdap(in *schema.ResourceData) (*managementClient.LdapConfig, error) {
	obj, err := expandAuthConfigLdap(in)
	if err != nil {
		return nil, err
	}

	obj.Name = AuthConfigOpenLdapName
	obj.Type = managementClient.OpenLdapConfigType

	return obj, nil
}
