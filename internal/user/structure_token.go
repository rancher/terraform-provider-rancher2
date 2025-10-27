package rancher2

import (
	"fmt"
	"math"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenToken(d *schema.ResourceData, in *managementClient.Token, patch bool) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)

	if len(in.ClusterID) > 0 {
		d.Set("cluster_id", in.ClusterID)
	}

	if len(in.Description) > 0 {
		d.Set("description", in.Description)
	}

	if in.Enabled != nil {
		d.Set("enabled", *in.Enabled)
	}

	d.Set("expired", in.Expired)

	if len(in.Name) > 0 {
		d.Set("name", in.Name)
	}

	if len(in.Token) > 0 {
		d.Set("token", in.Token)
		key := strings.Split(in.Token, ":")
		d.Set("access_key", key[0])
		d.Set("secret_key", key[1])
	}

	if in.TTLMillis >= 1000 {
		if !patch {
			d.Set("ttl", int(in.TTLMillis/1000))
		}
	}

	if len(in.UserID) > 0 {
		d.Set("user_id", in.UserID)
	}

	err := d.Set("annotations", toMapInterface(in.Annotations))
	if err != nil {
		return err
	}

	err = d.Set("labels", toMapInterface(in.Labels))
	if err != nil {
		return err
	}

	return nil

}

// Expanders

func expandToken(in *schema.ResourceData, patch bool) (*managementClient.Token, error) {
	obj := &managementClient.Token{}
	if in == nil {
		return nil, fmt.Errorf("[ERROR] Expanding token: Schema Resource data is nil")
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	if v, ok := in.Get("cluster_id").(string); ok && len(v) > 0 {
		obj.ClusterID = v
	}

	if v, ok := in.Get("description").(string); ok && len(v) > 0 {
		obj.Description = v
	}

	if v, ok := in.Get("ttl").(int); ok && v > 0 {
		if patch {
			// Rancher v2.4.6 ttl is read in minutes from API
			mins := math.Round(float64(v / 60))
			obj.TTLMillis = int64(mins)
		} else {
			obj.TTLMillis = int64(v * 1000)
		}
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj, nil
}
