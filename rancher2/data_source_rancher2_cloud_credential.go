package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRancher2CloudCredential() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2CloudCredentialRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"annotations": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceRancher2CloudCredentialRead(d *schema.ResourceData, meta interface{}) error {
	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	filters := map[string]interface{}{"name": name}
	listOpts := NewListOpts(filters)

	credentials, err := client.CloudCredential.List(listOpts)
	if err != nil {
		return err
	}

	count := len(credentials.Data)
	if count <= 0 {
		return fmt.Errorf("[ERROR] catalog with name \"%s\" not found", name)
	}
	if count > 1 {
		return fmt.Errorf("[ERROR] found %d catalogs with name \"%s\"", count, name)
	}

	credential := credentials.Data[0]

	d.SetId(credential.ID)
	d.Set("name", credential.Name)
	err = d.Set("annotations", toMapInterface(credential.Annotations))
	if err != nil {
		return err
	}
	err = d.Set("labels", toMapInterface(credential.Labels))
	if err != nil {
		return err
	}

	return nil
}
