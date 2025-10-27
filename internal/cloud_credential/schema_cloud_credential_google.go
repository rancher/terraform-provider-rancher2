package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const (
	googleConfigDriver = "googlekubernetesengine"
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
