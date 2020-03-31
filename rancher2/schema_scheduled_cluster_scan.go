package rancher2

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//Schemas

func scheduledClusterScanConfigFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"cron_schedule": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Crontab schedule",
			ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
				separator := " "
				v, ok := val.(string)
				if !ok || len(v) == 0 {
					errs = append(errs, fmt.Errorf("%q is nil", key))
					return
				}
				result := strings.Split(v, separator)
				if len(result) != 5 {
					errs = append(errs, fmt.Errorf("%q bad format: expected exactly 5 fields <min> <hour> <month_day> <month> <week_day>", key))
				}
				return
			},
		},
		"retention": {
			Type:        schema.TypeInt,
			Optional:    true,
			Computed:    true,
			Description: "Cluster scan retention",
		},
	}

	return s
}

func scheduledClusterScanFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"enabled": {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: "Enable scheduled cluster scan",
		},
		"scan_config": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &schema.Resource{
				Schema: clusterScanConfigFields(),
			},
			Description: "Cluster scan config",
		},
		"schedule_config": {
			Type:     schema.TypeList,
			MaxItems: 1,
			Required: true,
			Elem: &schema.Resource{
				Schema: scheduledClusterScanConfigFields(),
			},
			Description: "Schedule cluster scan config",
		},
	}

	return s
}
