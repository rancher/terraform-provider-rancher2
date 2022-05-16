package rancher2

import (
	"math"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

// Flatteners

func flattenUser(d *schema.ResourceData, in *managementClient.User) error {
	if in == nil {
		return nil
	}

	d.SetId(in.ID)
	d.Set("username", in.Username)
	d.Set("enabled", in.Enabled)

	if len(in.Name) > 0 {
		d.Set("name", in.Name)
	}

	err := d.Set("principal_ids", toArrayInterface(in.PrincipalIDs))
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

	return nil
}

func flattenUserToken(d *schema.ResourceData, in *managementClient.Token, patch bool) error {
	if in == nil {
		return nil
	}

	d.Set("token_id", in.ID)

	if len(in.ClusterID) > 0 {
		d.Set("token_config.0.cluster_id", in.ClusterID)
	}

	if len(in.Description) > 0 {
		d.Set("token_config.0.description", in.Description)
	}

	if in.Enabled != nil {
		d.Set("token_enabled", *in.Enabled)
	}

	d.Set("token_expired", in.Expired)

	if len(in.Name) > 0 {
		d.Set("token_name", in.Name)
	}

	if len(in.Token) > 0 {
		d.Set("auth_token", in.Token)
		key := strings.Split(in.Token, ":")
		d.Set("access_key", key[0])
		d.Set("secret_key", key[1])
	}

	if in.TTLMillis >= 1000 {
		if !patch {
			d.Set("token_config.0.ttl", int(in.TTLMillis/1000))
		}
	}

	return nil
}

// Expanders

func expandUser(in *schema.ResourceData) *managementClient.User {
	obj := &managementClient.User{}
	if in == nil {
		return nil
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	obj.Password = in.Get("password").(string)
	obj.Username = in.Get("username").(string)
	enabled := in.Get("enabled").(bool)
	obj.Enabled = &enabled

	if v, ok := in.Get("name").(string); ok && len(v) > 0 {
		obj.Name = v
	}

	if v, ok := in.Get("principal_ids").([]interface{}); ok && len(v) > 0 {
		obj.PrincipalIDs = toArrayString(v)
	}

	if v, ok := in.Get("annotations").(map[string]interface{}); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]interface{}); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj
}

func expandUserToken(in *schema.ResourceData, patch bool) *managementClient.Token {
	if in == nil {
		return nil
	}

	if v, ok := in.Get("token_config").([]interface{}); ok && len(v) > 0 {
		obj := &managementClient.Token{}

		if v, ok := in.Get("token_id").(string); ok && len(v) > 0 {
			obj.ID = v
		}

		if v, ok := in.Get("token_config.0.cluster_id").(string); ok && len(v) > 0 {
			obj.ClusterID = v
		}

		if v, ok := in.Get("token_config.0.description").(string); ok && len(v) > 0 {
			obj.Description = v
		}

		if v, ok := in.Get("token_config.0.ttl").(int); ok && v > 0 {
			if patch {
				// Rancher v2.4.6 ttl is read in minutes from API
				mins := math.Round(float64(v / 60))
				obj.TTLMillis = int64(mins)
			} else {
				obj.TTLMillis = int64(v * 1000)
			}
		}

		return obj
	}

	return nil
}
