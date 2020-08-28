package rancher2

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	projectClient "github.com/rancher/rancher/pkg/client/generated/project/v3"
)

var (
	testProjectCertificateConf         *projectClient.Certificate
	testProjectCertificateInterface    map[string]interface{}
	testNamespacedCertificateConf      *projectClient.NamespacedCertificate
	testNamespacedCertificateInterface map[string]interface{}
)

func init() {
	testCertificateCert := `-----BEGIN CERTIFICATE-----
MIIGUTCCBDmgAwIBAgIJAPnyJ5qP2EJ/MA0GCSqGSIb3DQEBBQUAMHgxCzAJBgNV
BAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRIwEAYDVQQHEwlDdXBlcnRpbm8x
FTATBgNVBAoTDFJhbmNoZXIgTGFiczENMAsGA1UECxMEVGVzdDEaMBgGA1UEAxMR
dGVzdC50ZXJyYWZvcm0uaW8wHhcNMTkwOTAyMTcyNzAyWhcNMjkwODMwMTcyNzAy
WjB4MQswCQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTESMBAGA1UEBxMJ
Q3VwZXJ0aW5vMRUwEwYDVQQKEwxSYW5jaGVyIExhYnMxDTALBgNVBAsTBFRlc3Qx
GjAYBgNVBAMTEXRlc3QudGVycmFmb3JtLmlvMIICIjANBgkqhkiG9w0BAQEFAAOC
Ag8AMIICCgKCAgEAqpPA9mpY3U0kSwwXPDNhrHA6HMpdzZw2xtJ5nzJA99m5+RsR
ZJH61x4Ujsd4kOSCcMFDPGK21aPRpRqi65qkyWPeUSGXBAQIS8MbZjllv2hXSfir
aFoRmnOE9UMOFeXlD+Y1yVxT8oT1aCDrL4ma6tQDGYOZQ6jISssUI7oNsHCtm6st
rcuKsxHuj8SKRp2x2ES4vDN7fAYbBl+CTIqht4RqGjX76cFQIk4SmyuO2qkdt/fU
jvRKy2ORt61yn7K4HP8UNXEGY7YJhlLMRwlLdic02ONAqoiBT9yorpO0P078xTQN
b7nuyY0g5408zn5+kDRVXCf9m4oKOSaZu5tRmcNG7oGmy0+MTsJlaAqiM67stljc
5AUE0xEmG1+OSe1yvMbpQwmWd/l+kV0/LDVRKCWauKiJi3oI0VSCRrcb6dOIVqDs
zjSYL4HzLxSXuuO3x7es6MEhkLM9I/VXAaxc4h4SuHTH2Zg0bXl5+PGQbqzqo1iL
a8Xs87IuRgMRRKoVbTw8Hb3XmcHW/Jv/2rzvPqGLineaffuGkmAvibnbwXUkXGcj
gTcIduWGSq/EfcclC1LszEszvymYNTV4XtsFOK+y4WDH/1U0cDfi+6wOjEvC3g1J
DI8jXtFlAJSqc/gIIbTQ42hkPfzjghc5WsQTM8mM7Zfhrd8DFcaRInhoTYUCAwEA
AaOB3TCB2jAdBgNVHQ4EFgQUovUCgqoRazEq0IJidzVk6FSwgwswgaoGA1UdIwSB
ojCBn4AUovUCgqoRazEq0IJidzVk6FSwgwuhfKR6MHgxCzAJBgNVBAYTAlVTMRMw
EQYDVQQIEwpDYWxpZm9ybmlhMRIwEAYDVQQHEwlDdXBlcnRpbm8xFTATBgNVBAoT
DFJhbmNoZXIgTGFiczENMAsGA1UECxMEVGVzdDEaMBgGA1UEAxMRdGVzdC50ZXJy
YWZvcm0uaW+CCQD58ieaj9hCfzAMBgNVHRMEBTADAQH/MA0GCSqGSIb3DQEBBQUA
A4ICAQAiiuP3sIiC+49MD6Ge1l1Rgn+m9lQCW0dPVb02zaxVh4/Qv4vK0RixJR8h
8qrBSpbC97fZ6CRizCSQgd4GL9ilVb8E6/IKla4ddvDorYoZHC2DYWiXSoIiZisg
j31KLLeUy+/2fr9Mc0jFEIq9inUjgwe9YNgQq0DQFs7udhyFavWg9q/fqP5dSQNR
nLNeG175ygGCeM19aIlGXAeFvoROuWNeDF7PHwkGPgpM0O29vQe8U84w0j1KOkyQ
tErI3iRtB1aV04JfItRNmSMNsAL0bQ37+4Jud/c3lDKa2dT8aTHW8SqRE63YBech
cXRRi8zG/8TdMsSauKwY0U+yQ3Ywa0XemmzC6DNj9pnwMn4eiXCT9+PXw2wyP7E5
fM0PWcGyMXpB4MslUqy9Kcu76BjKuQ50IiZ5zniCQZ+Fvm86pIQmrzlfy+ye7+TQ
sncnFoJvJ1Oe3o65QEiBZI4/aVkT/W45Rci6lTm3voh/kICebX7NTh+BCnPoKC+T
W3/N8M8YphOubjk+hRBLeDNuLJ9MXLoo/5WmHgtMvgagk/lg/TKyMv6dTwAVXLR6
c9CpbQhim0aocB3pMip59YJL+wxm2w+zE5g22VFoyEM7/vkE2fMByzSW3adMreoP
7DTWlp6wx9Lnlmg1VP9WFRQwy6+sQXY1klNyzupwr2kuWynrjQ==
-----END CERTIFICATE-----`
	testCertificateKey := `-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`
	testProjectCertificateConf = &projectClient.Certificate{
		Algorithm:       "RSA",
		CN:              "test.terraform.io",
		CertFingerprint: "30:B7:8C:E6:D8:69:53:A4:C3:5C:55:D1:BD:98:18:34:23:AE:B7:F2",
		Certs:           testCertificateCert,
		Description:     "description",
		ExpiresAt:       "2029-08-30 17:27:02 +0000 UTC",
		IssuedAt:        "2019-09-02 17:27:02 +0000 UTC",
		Issuer:          "test.terraform.io",
		Key:             testCertificateKey,
		KeySize:         "512",
		SerialNumber:    "18010501404599796351",
		ProjectID:       "project:test",
		Name:            "name",
		Version:         "3",
		Annotations: map[string]string{
			"node_one": "one",
			"node_two": "two",
		},
		Labels: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
	}
	testProjectCertificateInterface = map[string]interface{}{
		"certs":       Base64Encode(testCertificateCert),
		"key":         Base64Encode(testCertificateKey),
		"project_id":  "project:test",
		"name":        "name",
		"description": "description",
		"annotations": map[string]interface{}{
			"node_one": "one",
			"node_two": "two",
		},
		"labels": map[string]interface{}{
			"option1": "value1",
			"option2": "value2",
		},
	}
	testNamespacedCertificateConf = &projectClient.NamespacedCertificate{
		Algorithm:       "RSA",
		CN:              "test.terraform.io",
		CertFingerprint: "30:B7:8C:E6:D8:69:53:A4:C3:5C:55:D1:BD:98:18:34:23:AE:B7:F2",
		Certs:           testCertificateCert,
		Description:     "description",
		ExpiresAt:       "2029-08-30 17:27:02 +0000 UTC",
		IssuedAt:        "2019-09-02 17:27:02 +0000 UTC",
		Issuer:          "test.terraform.io",
		Key:             testCertificateKey,
		KeySize:         "512",
		SerialNumber:    "18010501404599796351",
		ProjectID:       "project:test",
		Name:            "name",
		Version:         "3",
		NamespaceId:     "namespace_id",
		Annotations: map[string]string{
			"node_one": "one",
			"node_two": "two",
		},
		Labels: map[string]string{
			"option1": "value1",
			"option2": "value2",
		},
	}
	testNamespacedCertificateInterface = map[string]interface{}{
		"certs":        Base64Encode(testCertificateCert),
		"key":          Base64Encode(testCertificateKey),
		"project_id":   "project:test",
		"name":         "name",
		"description":  "description",
		"namespace_id": "namespace_id",
		"annotations": map[string]interface{}{
			"node_one": "one",
			"node_two": "two",
		},
		"labels": map[string]interface{}{
			"option1": "value1",
			"option2": "value2",
		},
	}
}

func TestFlattenCertificate(t *testing.T) {

	cases := []struct {
		Input          interface{}
		ExpectedOutput map[string]interface{}
	}{
		{
			testProjectCertificateConf,
			testProjectCertificateInterface,
		},
		{
			testNamespacedCertificateConf,
			testNamespacedCertificateInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, certificateFields(), tc.ExpectedOutput)
		err := flattenCertificate(output, tc.Input)
		if err != nil {
			t.Fatalf("[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.ExpectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		if !reflect.DeepEqual(expectedOutput, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from flattener.\nExpected: %#v\nGiven:    %#v",
				expectedOutput, tc.ExpectedOutput)
		}
	}
}

func TestExpandCertificate(t *testing.T) {

	cases := []struct {
		Input          map[string]interface{}
		ExpectedOutput interface{}
	}{
		{
			testProjectCertificateInterface,
			testProjectCertificateConf,
		},
		{
			testNamespacedCertificateInterface,
			testNamespacedCertificateConf,
		},
	}

	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, certificateFields(), tc.Input)
		output, err := expandCertificate(inputResourceData)
		if err != nil {
			t.Fatalf("[ERROR] on expander: %#v", err)
		}
		if !reflect.DeepEqual(output, tc.ExpectedOutput) {
			t.Fatalf("Unexpected output from expander.\nExpected: %#v\nGiven:    %#v",
				tc.ExpectedOutput, output)
		}
	}
}
