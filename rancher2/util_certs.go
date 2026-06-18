package rancher2

import "encoding/base64"

func decodeCACertIfBase64(cert string) string {
	if cert == "" {
		return cert
	}
	if b, err := base64.StdEncoding.DecodeString(cert); err == nil {
		return string(b)
	}
	return cert
}
