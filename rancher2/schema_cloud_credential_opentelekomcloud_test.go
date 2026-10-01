package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/stretchr/testify/require"
	"github.com/zclconf/go-cty/cty"
)

func TestTCloudPublicCredentialSchema(t *testing.T) {
	r := resourceRancher2CloudCredential()
	config := terraform.NewResourceConfigRaw(map[string]interface{}{
		"name": "tcloud-public",
		"tcloud_public_credential_config": []interface{}{map[string]interface{}{
			"access_key": "test-access-key",
			"secret_key": "test-secret-key",
			"region":     "eu-de",
		}},
	})
	_, errs := r.Validate(config)
	require.Empty(t, errs)
	for _, field := range []string{"access_key", "secret_key", "password"} {
		require.True(t, cloudCredentialOpentelekomcloudFields()[field].Sensitive, field)
	}
	require.Len(t, cloudCredentialOpentelekomcloudFields(), 9)
	require.Contains(t, r.Schema["tcloud_public_credential_config"].ConflictsWith, "nutanix_credential_config")
	for _, field := range allCloudCredentialDriverConfigFields {
		if field != "tcloud_public_credential_config" {
			require.Contains(t, r.Schema[field].ConflictsWith, "tcloud_public_credential_config")
		}
	}
	conflicting := terraform.NewResourceConfigRaw(map[string]interface{}{
		"name":                            "conflict",
		"tcloud_public_credential_config": []interface{}{map[string]interface{}{"access_key": "test", "secret_key": "test"}},
		"openstack_credential_config":     []interface{}{map[string]interface{}{"password": "test"}},
	})
	_, errs = r.Validate(conflicting)
	require.NotEmpty(t, errs)
}

func TestTCloudPublicCredentialAuthenticationDiff(t *testing.T) {
	cases := []struct {
		name    string
		values  map[string]interface{}
		wantErr bool
	}{
		{"access keys", map[string]interface{}{"access_key": "test-ak", "secret_key": "test-sk"}, false},
		{"password", map[string]interface{}{"username": "test-user", "password": "test-password", "domain_name": "test-domain"}, false},
		{"both methods", map[string]interface{}{
			"access_key": "test-ak", "secret_key": "test-sk", "username": "test-user", "password": "test-password", "domain_name": "test-domain",
		}, false},
		{"empty", map[string]interface{}{}, true},
		{"missing secret key", map[string]interface{}{"access_key": "test-ak"}, true},
		{"missing access key", map[string]interface{}{"secret_key": "test-sk"}, true},
		{"missing password", map[string]interface{}{"username": "test-user", "domain_name": "test-domain"}, true},
		{"missing domain", map[string]interface{}{"username": "test-user", "password": "test-password"}, true},
		{"partial keys with password method", map[string]interface{}{
			"access_key": "test-ak", "username": "test-user", "password": "test-password", "domain_name": "test-domain",
		}, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			config := terraform.NewResourceConfigRaw(map[string]interface{}{
				"name": "tcloud-public", "tcloud_public_credential_config": []interface{}{tc.values},
			})
			_, err := resourceRancher2CloudCredential().Diff(nil, config, nil)
			if tc.wantErr {
				require.ErrorContains(t, err, "tcloud_public_credential_config requires")
			} else {
				require.NoError(t, err)
			}
		})
	}
	t.Run("other drivers unchanged", func(t *testing.T) {
		config := terraform.NewResourceConfigRaw(map[string]interface{}{
			"name": "openstack", "openstack_credential_config": []interface{}{map[string]interface{}{"password": "test-password"}},
		})
		_, err := resourceRancher2CloudCredential().Diff(nil, config, nil)
		require.NoError(t, err)
	})
	t.Run("unknown secret key", func(t *testing.T) {
		r := resourceRancher2CloudCredential()
		core := r.CoreConfigSchema()
		values := make(map[string]cty.Value)
		for key, valueType := range core.ImpliedType().AttributeTypes() {
			values[key] = cty.NullVal(valueType)
		}
		credential := make(map[string]cty.Value)
		for key := range cloudCredentialOpentelekomcloudFields() {
			credential[key] = cty.NullVal(cty.String)
		}
		credential["access_key"] = cty.StringVal("test-ak")
		credential["secret_key"] = cty.UnknownVal(cty.String)
		values["name"] = cty.StringVal("tcloud-public")
		values["tcloud_public_credential_config"] = cty.ListVal([]cty.Value{cty.ObjectVal(credential)})
		config := terraform.NewResourceConfigShimmed(cty.ObjectVal(values), core)
		_, err := r.Diff(nil, config, nil)
		require.NoError(t, err)
	})
}

func TestTCloudPublicCredentialRejectInvalidCreate(t *testing.T) {
	d := schema.TestResourceDataRaw(t, cloudCredentialFields(), map[string]interface{}{
		"name": "invalid", "tcloud_public_credential_config": []interface{}{map[string]interface{}{"access_key": "test-ak"}},
	})
	require.ErrorContains(t, resourceRancher2CloudCredentialCreate(d, nil), "tcloud_public_credential_config requires")
}
