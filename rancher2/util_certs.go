package rancher2

import "encoding/base64"

// decodeCACertIfBase64 decodes cert if it is base64 encoded, returning the
// original string unchanged when it is not (or is empty). Rancher's v3
// management API returns the cluster CA certificate base64-encoded; the
// provisioning v1 API's local_auth_endpoint.ca_certs expects raw PEM.
func decodeCACertIfBase64(cert string) string {
	if cert == "" {
		return cert
	}
	if b, err := base64.StdEncoding.DecodeString(cert); err == nil {
		return string(b)
	}
	return cert
}
