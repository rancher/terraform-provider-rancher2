package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func dataSourceRancher2Principal() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2PrincipalRead,

		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceRancher2PrincipalRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	email := d.Get("email").(string)
	principalType := d.Get("type").(string)

	collection, err := client.Principal.List(nil)
	if err != nil {
		return err
	}

	principals, err := client.Principal.CollectionActionSearch(collection, &managementClient.SearchPrincipalsInput{
		Name:          email,
		PrincipalType: principalType,
	})
	if err != nil {
		return err
	}

	count := len(principals.Data)
	if count <= 0 {
		return fmt.Errorf("[ERROR] principal \"%s\" of type \"%s\" not found", email, principalType)
	}

	return flattenDataSourcePrincipal(d, &principals.Data[0])
}

func flattenDataSourcePrincipal(d *schema.ResourceData, in *managementClient.Principal) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)
	d.Set("id", in.ID)
	d.Set("type", in.PrincipalType)

	return nil
}
