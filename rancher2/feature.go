package rancher2

import (
	"fmt"
	"log"

	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

const multipleAuthProvidersFeature = "multiple-auth-providers"

// isFeatureEnabled checks if a given feature is enabled in Rancher.
//
// It returns true if the feature is enabled, false otherwise.
//
// If the feature is not found or access is forbidden, it returns false without an error.
//
// Any other errors encountered during the check are returned.
func isFeatureEnabled(features managementClient.FeatureOperations, name string) (bool, error) {
	feature, err := features.ByID(name)
	if err != nil {
if IsNotFound(err) {
			log.Printf("[INFO] Feature %s not found; treating it as disabled.", name)
			return false, nil
		}
		if IsForbidden(err) {
			log.Printf("[INFO] Access forbidden when reading feature %s; treating it as disabled.", name)
			return false, nil
		}

		return false, fmt.Errorf("getting feature %s: %w", name, err)
	}

	var enabled bool
	if feature.Value != nil {
		enabled = *feature.Value
	}

	return enabled, nil
}

// checkAuthConfigCanBeEnabled returns an error if another non-local AuthConfig
// is already enabled and the multiple-auth-providers feature is disabled.
//
// Without the feature, Rancher only allows the local AuthConfig and one other
// to be enabled at a time.
func checkAuthConfigCanBeEnabled(features managementClient.FeatureOperations, authConfigs managementClient.AuthConfigOperations, name string) error {
	multiEnabled, err := isFeatureEnabled(features, multipleAuthProvidersFeature)
	if err != nil {
		return err
	}

	if multiEnabled {
		return nil
	}

	auths, err := authConfigs.List(NewListOpts(nil))
	if err != nil {
		return fmt.Errorf("listing Auth Configs: %w", err)
	}

	for _, auth := range auths.Data {
		if auth.Enabled && auth.ID != name && auth.ID != "local" {
			return fmt.Errorf("enabling Auth Config %s: Auth Config %s is already enabled, and the %s feature is not enabled in Rancher", name, auth.ID, multipleAuthProvidersFeature)
		}
	}

	return nil
}
