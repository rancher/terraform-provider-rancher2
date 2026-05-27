package rancher2

import (
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
)

func expandOIDCClient(in *schema.ResourceData) (*managementClient.OIDCClient, error) {
	obj := &managementClient.OIDCClient{}
	if in == nil {
		return nil, errors.New("expanding OIDC Client: Input ResourceData is nil")
	}

	if v := in.Id(); len(v) > 0 {
		obj.ID = v
	}

	if v, ok := in.GetOk("description"); ok {
		obj.Description = v.(string)
	}

	v := in.Get("redirect_uris")
	obj.RedirectURIs = toArrayString(v.([]any))

	v, ok := in.GetOk("token_expiration_seconds")
	if ok {
		// This should be safe because the field is declared as an integer.
		tokenExpirationSeconds := v.(int)
		obj.TokenExpirationSeconds = int64(tokenExpirationSeconds)
	} else {
		obj.TokenExpirationSeconds = 0
	}

	v, ok = in.GetOk("refresh_token_expiration_seconds")
	if ok {
		// This should be safe because the field is declared as an integer.
		refreshTokenExpirationSeconds := v.(int)
		obj.RefreshTokenExpirationSeconds = int64(refreshTokenExpirationSeconds)
	} else {
		obj.RefreshTokenExpirationSeconds = 0
	}

	if v, ok := in.Get("annotations").(map[string]any); ok && len(v) > 0 {
		obj.Annotations = toMapString(v)
	}

	if v, ok := in.Get("labels").(map[string]any); ok && len(v) > 0 {
		obj.Labels = toMapString(v)
	}

	return obj, nil
}

func flattenOIDCClient(d *schema.ResourceData, in *managementClient.OIDCClient) error {
	if in == nil {
		return errors.New("flattening OIDC Client: Input config is nil")
	}

	if in.ID != "" {
		d.SetId(in.ID)
	}

	return errors.Join(
		d.Set("description", in.Description),
		d.Set("redirect_uris", in.RedirectURIs),
		d.Set("token_expiration_seconds", int(in.TokenExpirationSeconds)),
		d.Set("refresh_token_expiration_seconds", int(in.RefreshTokenExpirationSeconds)),
		d.Set("client_id", in.Status.ClientID),
		d.Set("annotations", toMapInterface(in.Annotations)),
		d.Set("labels", toMapInterface(in.Labels)),
	)
}
