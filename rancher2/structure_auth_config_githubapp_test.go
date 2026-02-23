package rancher2

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	"github.com/stretchr/testify/assert"
)

const testPrivateKeyPKCS1 = `
-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQCjcGqTkOq0CR3rTx0ZSQSIdTrDrFAYl29611xN8aVgMQIWtDB/
lD0W5TpKPuU9iaiG/sSn/VYt6EzN7Sr332jj7cyl2WrrHI6ujRswNy4HojMuqtfa
b5FFDpRmCuvl35fge18OvoQTJELhhJ1EvJ5KUeZiuJ3u3YyMnxxXzLuKbQIDAQAB
AoGAPrNDz7TKtaLBvaIuMaMXgBopHyQd3jFKbT/tg2Fu5kYm3PrnmCoQfZYXFKCo
ZUFIS/G1FBVWWGpD/MQ9tbYZkKpwuH+t2rGndMnLXiTC296/s9uix7gsjnT4Naci
5N6EN9pVUBwQmGrYUTHFc58ThtelSiPARX7LSU2ibtJSv8ECQQDWBRrrAYmbCUN7
ra0DFT6SppaDtvvuKtb+mUeKbg0B8U4y4wCIK5GH8EyQSwUWcXnNBO05rlUPbifs
DLv/u82lAkEAw39sTJ0KmJJyaChqvqAJ8guulKlgucQJ0Et9ppZyet9iVwNKX/aW
9UlwGBMQdafQ36nd1QMEA8AbAw4D+hw/KQJBANJbHDUGQtk2hrSmZNoV5HXB9Uiq
7v4N71k5ER8XwgM5yVGs2tX8dMM3RhnBEtQXXs9LW1uJZSOQcv7JGXNnhN0CQBZe
nzrJAWxh3XtznHtBfsHWelyCYRIAj4rpCHCmaGUM6IjCVKFUawOYKp5mmAyObkUZ
f8ue87emJLEdynC1CLkCQHduNjP1hemAGWrd6v8BHhE3kKtcK6KHsPvJR5dOfzbd
HAqVePERhISfN6cwZt5p8B3/JUwSR8el66DF7Jm57BM=
-----END RSA PRIVATE KEY-----
`

const testPrivateKeyPKCS8 = `
-----BEGIN PRIVATE KEY-----
MIICdgIBADANBgkqhkiG9w0BAQEFAASCAmAwggJcAgEAAoGBAKNwapOQ6rQJHetP
HRlJBIh1OsOsUBiXb3rXXE3xpWAxAha0MH+UPRblOko+5T2JqIb+xKf9Vi3oTM3t
KvffaOPtzKXZauscjq6NGzA3LgeiMy6q19pvkUUOlGYK6+Xfl+B7Xw6+hBMkQuGE
nUS8nkpR5mK4ne7djIyfHFfMu4ptAgMBAAECgYA+s0PPtMq1osG9oi4xoxeAGikf
JB3eMUptP+2DYW7mRibc+ueYKhB9lhcUoKhlQUhL8bUUFVZYakP8xD21thmQqnC4
f63asad0ycteJMLb3r+z26LHuCyOdPg1pyLk3oQ32lVQHBCYathRMcVznxOG16VK
I8BFfstJTaJu0lK/wQJBANYFGusBiZsJQ3utrQMVPpKmloO2++4q1v6ZR4puDQHx
TjLjAIgrkYfwTJBLBRZxec0E7TmuVQ9uJ+wMu/+7zaUCQQDDf2xMnQqYknJoKGq+
oAnyC66UqWC5xAnQS32mlnJ632JXA0pf9pb1SXAYExB1p9Dfqd3VAwQDwBsDDgP6
HD8pAkEA0lscNQZC2TaGtKZk2hXkdcH1SKru/g3vWTkRHxfCAznJUaza1fx0wzdG
GcES1Bdez0tbW4llI5By/skZc2eE3QJAFl6fOskBbGHde3Oce0F+wdZ6XIJhEgCP
iukIcKZoZQzoiMJUoVRrA5gqnmaYDI5uRRl/y57zt6YksR3KcLUIuQJAd242M/WF
6YAZat3q/wEeETeQq1wrooew+8lHl05/Nt0cCpV48RGEhJ83pzBm3mnwHf8lTBJH
x6XroMXsmbnsEw==
-----END PRIVATE KEY-----
`

var (
	testAuthConfigGithubAppConf = &managementClient.GithubAppConfig{
		Name:                AuthConfigGithubAppName,
		Type:                managementClient.GithubAppConfigType,
		AccessMode:          "access",
		AllowedPrincipalIDs: []string{"allowed1", "allowed2"},
		Enabled:             true,
		ClientID:            "client_id",
		Hostname:            "hostname",
		TLS:                 true,
		AppID:               "12345",
		InstallationID:      "2345678",
		PrivateKey:          testPrivateKeyPKCS8,
	}
	testAuthConfigGithubAppInterface = map[string]interface{}{
		"name":                  AuthConfigGithubAppName,
		"type":                  managementClient.GithubAppConfigType,
		"access_mode":           "access",
		"allowed_principal_ids": []interface{}{"allowed1", "allowed2"},
		"enabled":               true,
		"client_id":             "client_id",
		"hostname":              "hostname",
		"tls":                   true,
		"app_id":                "12345",
		"installation_id":       "2345678",
		"private_key":           testPrivateKeyPKCS8,
	}
)

func TestFlattenAuthConfigGithubApp(t *testing.T) {
	cases := []struct {
		input          *managementClient.GithubAppConfig
		expectedOutput map[string]interface{}
	}{
		{
			testAuthConfigGithubAppConf,
			testAuthConfigGithubAppInterface,
		},
	}

	for _, tc := range cases {
		output := schema.TestResourceDataRaw(t, authConfigGithubAppFields(), map[string]interface{}{})
		err := flattenAuthConfigGithubApp(output, tc.input)
		if err != nil {
			assert.FailNow(t, "[ERROR] on flattener: %#v", err)
		}
		expectedOutput := map[string]interface{}{}
		for k := range tc.expectedOutput {
			expectedOutput[k] = output.Get(k)
		}
		assert.Equal(t, tc.expectedOutput, expectedOutput, "Unexpected output from flattener.")
	}
}

func TestExpandAuthConfigGithubApp(t *testing.T) {
	cases := []struct {
		input          map[string]interface{}
		expectedOutput *managementClient.GithubAppConfig
	}{
		{
			testAuthConfigGithubAppInterface,
			testAuthConfigGithubAppConf,
		},
	}
	for _, tc := range cases {
		inputResourceData := schema.TestResourceDataRaw(t, authConfigGithubAppFields(), tc.input)
		output, err := expandAuthConfigGithubApp(inputResourceData)
		if err != nil {
			assert.FailNow(t, "[ERROR] on expander: %#v", err)
		}
		assert.Equal(t, tc.expectedOutput, output, "Unexpected output from expander.")
	}
}
