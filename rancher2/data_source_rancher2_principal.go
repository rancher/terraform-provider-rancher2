package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

const (
	principalTypeGroup = "group"
	principalTypeUser  = "user"
)

var (
	principalTypes = []string{principalTypeGroup, principalTypeUser}
)

func dataSourceRancher2Principal() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2PrincipalRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      principalTypeUser,
				ValidateFunc: validation.StringInSlice(principalTypes, true),
			},
			"exact_match": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceRancher2PrincipalRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	principalType := d.Get("type").(string)
	exactMatch := d.Get("exact_match").(bool)

	collection, err := client.Principal.ListAll(nil)
	if err != nil {
		return err
	}

	principals, err := client.Principal.CollectionActionSearch(collection, &managementClient.SearchPrincipalsInput{
		Name:          name,
		PrincipalType: principalType,
	})
	if err != nil {
		return err
	}

	count := len(principals.Data)
	if count <= 0 {
		return fmt.Errorf("[ERROR] principal \"%s\" of type \"%s\" not found", name, principalType)
	}

	// We always had at least one result here, let's find which can match exactly with the inputted name
	if exactMatch {
		for _, v := range principals.Data {
			if v.Name == name {
				return flattenDataSourcePrincipal(d, &v)
			}
		}
		// This situation will be almost never happened, but we still ensure for the special case
		return fmt.Errorf(
			"[ERROR] principal \"%s\" of type \"%s\" not found. Try again with \"exact_match=false\" for partially matched result",
			name, principalType)
	}

	return flattenDataSourcePrincipal(d, &principals.Data[0])
}

func flattenDataSourcePrincipal(d *schema.ResourceData, in *managementClient.Principal) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)
	d.Set("name", in.Name)
	d.Set("type", in.PrincipalType)

	return nil
}
