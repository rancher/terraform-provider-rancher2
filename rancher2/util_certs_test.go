package rancher2

import (
	"encoding/base64"
	"testing"
)

func TestDecodeCACertIfBase64(t *testing.T) {
	pem := "-----BEGIN CERTIFICATE-----\nFAKE\n-----END CERTIFICATE-----\n"
	encoded := base64.StdEncoding.EncodeToString([]byte(pem))
	if got := decodeCACertIfBase64(encoded); got != pem {
		t.Fatalf("expected decoded cert, got %q", got)
	}
	if got := decodeCACertIfBase64(pem); got != pem {
		t.Fatalf("expected unchanged cert, got %q", got)
	}
}
