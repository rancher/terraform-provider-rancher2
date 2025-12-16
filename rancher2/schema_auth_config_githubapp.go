package rancher2

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"maps"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

const AuthConfigGithubAppName = "githubapp"

func authConfigGithubAppFields() map[string]*schema.Schema {
	s := map[string]*schema.Schema{
		"client_id": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"client_secret": {
			Type:      schema.TypeString,
			Required:  true,
			Sensitive: true,
		},
		"hostname": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "github.com",
		},
		"tls": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		"app_id": {
			Type:         schema.TypeString,
			Description:  "The GitHub App ID is provided on the GitHub apps page.",
			Required:     true,
			ValidateFunc: isValidIntegerString,
		},
		"private_key": {
			Type:         schema.TypeString,
			Required:     true,
			Description:  "PEM format private key for signing requests.",
			ValidateFunc: isPEMEncodedPrivateKey,
		},
		"installation_id": {
			Type:         schema.TypeString,
			Description:  "If the Installation ID is not provided, all installations for the App will be queried.",
			Optional:     true,
			ValidateFunc: isValidIntegerString,
		},
	}

	maps.Copy(s, authConfigFields())

	return s
}

func isValidIntegerString(i any, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	if _, err := strconv.ParseInt(v, 10, 64); err != nil {
		errors = append(errors, fmt.Errorf("expected %q to be a valid integer, got %v", k, v))
	}

	return warnings, errors
}

func isPEMEncodedPrivateKey(i any, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	var block *pem.Block
	if block, _ = pem.Decode([]byte(v)); block == nil {
		errors = append(errors, fmt.Errorf("expected %q to be PEM encoded", k))
		return warnings, errors
	}

	var key any
	parsedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		parsedKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			errors = append(errors, fmt.Errorf("expected %q to be an RSA Private Key", k))
			return warnings, errors
		} else {
			key = parsedKey
		}
	} else {
		key = parsedKey
	}

	if _, ok := key.(*rsa.PrivateKey); !ok {
		errors = append(errors, fmt.Errorf("expected %q to be an RSA Private Key", k))
		return warnings, errors
	}

	return warnings, errors
}
