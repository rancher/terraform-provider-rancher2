package cattle

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceCattleSetting() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceCattleSettingRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"value": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceCattleSettingRead(d *schema.ResourceData, meta interface{}) error {
	name := d.Get("name").(string)
	log.Printf("[INFO] Refreshing Cattle Setting: %s", name)

	client, err := meta.(*Config).ManagementClient()
	if err != nil {
		return err
	}

	setting, err := client.Setting.ByID(name)
	if err != nil {
		return err
	}

	d.SetId(name)
	d.Set("value", setting.Value)

	return nil
}
