package rancher2

import (
	"github.com/hashicorp/terraform/helper/schema"
	managementClient "github.com/rancher/types/client/management/v3"
)

// Flatteners

func flattenAuthConfigFreeIpa(d *schema.ResourceData, in *managementClient.LdapConfig) error {
	err := flattenAuthConfigLdap(d, in)
	if err != nil {
		return err
	}

	d.SetId(AuthConfigFreeIpaName)

	err = d.Set("name", AuthConfigFreeIpaName)
	if err != nil {
		return err
	}
	err = d.Set("type", managementClient.FreeIpaConfigType)
	if err != nil {
		return err
	}

	return nil
}

// Expanders

func expandAuthConfigFreeIpa(in *schema.ResourceData) (*managementClient.LdapConfig, error) {
	obj, err := expandAuthConfigLdap(in)
	if err != nil {
		return nil, err
	}

	obj.Name = AuthConfigFreeIpaName
	obj.Type = managementClient.FreeIpaConfigType

	return obj, nil
}
