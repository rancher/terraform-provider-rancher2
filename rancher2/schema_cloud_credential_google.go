package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	// Before Rancher v2.12 only the GKE Kontainer driver was supported,
	// and the provider would automatically enable that driver when
	// a google cloud credential was created. Rancher v2.13 dropped support
	// for automatically enabling kontainer drivers, so we should migrate existing
	// google cloud credentials to use the google node driver instead.
	//
	// This is done to ensure that terraform state files have a uniform driver value
	// across all google cloud credentials, both old and new.

	oldGoogleConfigKontainerDriver = "googlekubernetesengine"
	googleConfigNodeDriver         = "google"
)

//Types

type googleCredentialConfig struct {
	AuthEncodedJSON string `json:"authEncodedJson,omitempty" yaml:"authEncodedJson,omitempty"`
}

//Schemas

func cloudCredentialGoogleFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"auth_encoded_json": {
			Type:        schema.TypeString,
			Required:    true,
			Sensitive:   true,
			Description: "Google auth encoded json",
			ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
				v, ok := val.(string)
				if !ok || len(v) == 0 {
					return
				}
				_, err := jsonToMapInterface(v)
				if err != nil {
					errs = append(errs, fmt.Errorf("%q must be in json format, error: %v", key, err))
					return
				}
				return
			},
		},
	}

	return s
}
