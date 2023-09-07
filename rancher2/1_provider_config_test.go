package rancher2

import (
	"errors"
	"testing"
)

type ProviderConfigTestCase struct {
	Config         *Config
	ExpectedOutput *Config
	ExpectedError  error
}

func TestProviderConfig_ProxyURL(t *testing.T) {
	cases := map[string]ProviderConfigTestCase{
		"configuring valid http proxy url": {
			&Config{
				ProxyURL: "http://localhost:8080",
			},
			nil,
			nil,
		},
		"configuring valid https proxy url": {
			&Config{
				ProxyURL: "https://localhost:8443",
			},
			nil,
			nil,
		},
		"configuring valid socks5 proxy url": {
			&Config{
				ProxyURL: "socks5://localhost:2000",
			},
			nil,
			nil,
		},
		"configuring invalid proxy url": {
			&Config{
				ProxyURL: "http://local host:8080",
			},
			nil,
			errors.New("[ERROR] invalid proxy address \"http://local host:8080\": parse \"http://local host:8080\": invalid character \" \" in host name"),
		},
		"configuring invalid proxy url scheme": {
			&Config{
				ProxyURL: "foo://localhost:80",
			},
			nil,
			errors.New("[ERROR] invalid proxy scheme \"foo\" (should be one of http, https or socks5)"),
		},
	}
	runTestCases(t, cases)
}

func runTestCases(t *testing.T, cases map[string]ProviderConfigTestCase) {
	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			expectedOutput := tc.ExpectedOutput
			if expectedOutput == nil {
				expectedOutput = tc.Config
			}
			output, err := providerValidateConfig(tc.Config)
			if err == nil && tc.ExpectedError != nil {
				t.Fatalf("expected error \"%v\", got none", tc.ExpectedError)
			}
			if err != nil && tc.ExpectedError != nil && err.Error() != tc.ExpectedError.Error() {
				t.Fatalf("expected error '%v', got '%v'", tc.ExpectedError, err)
			}
			if err != nil && tc.ExpectedError == nil {
				t.Fatalf("invalid provider configuration: %v", err)
			}
			if err == nil && expectedOutput != output {
				t.Fatalf("invalid provider configuration\nExpected: %#v\nGiven:    %#v", expectedOutput, output)
			}
		})
	}
}
