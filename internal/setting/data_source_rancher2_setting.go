package rancher2

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRancher2Setting() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRancher2SettingRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRancher2SettingRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	log.Printf("[INFO] Refreshing Rancher2 Setting: %s", name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	setting, err := client.Setting.ByID(name)
	if err != nil || setting == nil {
		return err
	}

	d.SetId(name)
	d.Set("value", setting.Value)

	return nil
}
