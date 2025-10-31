package rancher2

import (
	"reflect"
	"testing"
)

const testPublicKey = `
-----BEGIN RSA PUBLIC KEY-----
MIIBCgKCAQEA+xGZ/wcz9ugFpP07Nspo6U17l0YhFiFpxxU4pTk3Lifz9R3zsIsu
ERwta7+fWIfxOo208ett/jhskiVodSEt3QBGh4XBipyWopKwZ93HHaDVZAALi/2A
+xTBtWdEo7XGUujKDvC2/aZKukfjpOiUI8AhLAfjmlcD/UZ1QPh0mHsglRNCmpCw
mwSXA9VNmhz+PiB+Dml4WWnKW/VHo2ujTXxq7+efMU4H2fny3Se3KYOsFPFGZ1TN
QSYlFuShWrHPtiLmUdPoP6CV2mML1tk+l7DIIqXrQhLUKDACeM5roMx0kLhUWB8P
+0uj1CNlNN4JRZlC7xFfqiMbFRU9Z4N6YwIDAQAB
-----END RSA PUBLIC KEY-----
`

func TestIsValidIntegerString(t *testing.T) {
	validationTests := map[string]struct {
		value      any
		wantErrors []string
	}{
		"valid integer value": {
			value: "1",
		},
		"invalid integer value": {
			value:      12345,
			wantErrors: []string{`expected type of "testing" to be string`},
		},
		"invalid string": {
			value:      "a",
			wantErrors: []string{`expected "testing" to be a valid integer, got a`},
		},
	}

	assertErrorsEqual := func(t *testing.T, errors []error, want []string) {
		se := func(e []error) []string {
			var se []string
			for _, v := range e {
				se = append(se, v.Error())
			}
			return se
		}(errors)
		if !reflect.DeepEqual(se, want) {
			t.Errorf("errors do not match, got %#v, want %#v", se, want)
		}
	}

	for k, tt := range validationTests {
		t.Run(k, func(t *testing.T) {
			_, errors := isValidIntegerString(tt.value, "testing")
			assertErrorsEqual(t, errors, tt.wantErrors)
		})
	}
}

func TestIsPEMEncodedPrivateKey(t *testing.T) {
	validationTests := map[string]struct {
		value      any
		wantErrors []string
	}{
		"non-string value": {
			value:      12345,
			wantErrors: []string{`expected type of "testing" to be string`},
		},
		"non-PEM encoded string": {
			value:      "CERT",
			wantErrors: []string{`expected "testing" to be PEM encoded`},
		},
		"invalid PEM encoded cert": {
			value:      "----BEGIN PRIVATE KEY-----\nTESTING\n-----END PRIVATE KEY-----\n",
			wantErrors: []string{`expected "testing" to be PEM encoded`},
		},
		"invalid PEM public key": {
			value:      testPublicKey,
			wantErrors: []string{`expected "testing" to be an RSA Private Key`},
		},
		"valid PEM Private Key PKCS1 encoded": {
			value: testPrivateKeyPKCS1,
		},
		"valid PEM Private Key PKCS8 encoded": {
			value: testPrivateKeyPKCS8,
		},
	}

	assertErrorsEqual := func(t *testing.T, errors []error, want []string) {
		se := func(e []error) []string {
			var se []string
			for _, v := range e {
				se = append(se, v.Error())
			}
			return se
		}(errors)
		if !reflect.DeepEqual(se, want) {
			t.Errorf("errors do not match, got %#v, want %#v", se, want)
		}
	}

	for k, tt := range validationTests {
		t.Run(k, func(t *testing.T) {
			_, errors := isPEMEncodedPrivateKey(tt.value, "testing")
			assertErrorsEqual(t, errors, tt.wantErrors)
		})
	}
}
