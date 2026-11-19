package rancher2

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Expanders

// expandAuthConfigName returns the configured AuthConfig name, falling back to
// defaultName when it is not set.
func expandAuthConfigName(in *schema.ResourceData, defaultName string) string {
	if v, ok := in.Get("name").(string); ok && v != "" {
		return v
	}

	return defaultName
}
