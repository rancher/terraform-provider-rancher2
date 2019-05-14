package rancher2

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// CLIConfig used to store data from file.
type CLIConfig struct {
	AdminPass string `json:"adminpass"`
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	TokenKey  string `json:"tokenKey"`
	CACerts   string `json:"caCerts"`
	Insecure  bool   `json:"insecure,omitempty"`
	URL       string `json:"url"`
	Project   string `json:"project"`
	Path      string `json:"path,omitempty"`
}

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_URL", ""),
				Description: descriptions["api_url"],
			},
			"access_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_ACCESS_KEY", ""),
				Description: descriptions["access_key"],
			},
			"bootstrap": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_BOOTSTRAP", false),
				Description: descriptions["bootstrap"],
			},
			"secret_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_SECRET_KEY", ""),
				Description: descriptions["secret_key"],
			},
			"token_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_TOKEN_KEY", ""),
				Description: descriptions["token_key"],
			},
			"ca_certs": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_CA_CERTS", ""),
				Description: descriptions["ca_certs"],
			},
			"insecure": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_INSECURE", false),
				Description: descriptions["insecure"],
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"rancher2_auth_config_activedirectory":   resourceRancher2AuthConfigActiveDirectory(),
			"rancher2_auth_config_adfs":              resourceRancher2AuthConfigADFS(),
			"rancher2_auth_config_azuread":           resourceRancher2AuthConfigAzureAD(),
			"rancher2_auth_config_freeipa":           resourceRancher2AuthConfigFreeIpa(),
			"rancher2_auth_config_github":            resourceRancher2AuthConfigGithub(),
			"rancher2_auth_config_openldap":          resourceRancher2AuthConfigOpenLdap(),
			"rancher2_auth_config_ping":              resourceRancher2AuthConfigPing(),
			"rancher2_bootstrap":                     resourceRancher2Bootstrap(),
			"rancher2_catalog":                       resourceRancher2Catalog(),
			"rancher2_cloud_credential":              resourceRancher2CloudCredential(),
			"rancher2_cluster":                       resourceRancher2Cluster(),
			"rancher2_cluster_driver":                resourceRancher2ClusterDriver(),
			"rancher2_cluster_logging":               resourceRancher2ClusterLogging(),
			"rancher2_cluster_role_template_binding": resourceRancher2ClusterRoleTemplateBinding(),
			"rancher2_etcd_backup":                   resourceRancher2EtcdBackup(),
			"rancher2_node_driver":                   resourceRancher2NodeDriver(),
			"rancher2_node_pool":                     resourceRancher2NodePool(),
			"rancher2_node_template":                 resourceRancher2NodeTemplate(),
			"rancher2_project":                       resourceRancher2Project(),
			"rancher2_project_logging":               resourceRancher2ProjectLogging(),
			"rancher2_project_role_template_binding": resourceRancher2ProjectRoleTemplateBinding(),
			"rancher2_namespace":                     resourceRancher2Namespace(),
			"rancher2_setting":                       resourceRancher2Setting(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"rancher2_setting": dataSourceRancher2Setting(),
		},

		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"access_key": "API Key used to authenticate with the rancher server",

		"secret_key": "API secret used to authenticate with the rancher server",

		"token_key": "API token used to authenticate with the rancher server",

		"ca_certs": "CA certificates used to sign rancher server tls certificates. Mandatory if self signed tls and insecure option false",

		"insecure": "Allow insecure connections to Rancher. Mandatory if self signed tls and not ca_certs provided",

		"api_url": "The URL to the rancher API",

		"bootstrap": "Bootstrap rancher server",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	apiURL := d.Get("api_url").(string)
	accessKey := d.Get("access_key").(string)
	secretKey := d.Get("secret_key").(string)
	tokenKey := d.Get("token_key").(string)
	caCerts := d.Get("ca_certs").(string)
	insecure := d.Get("insecure").(bool)
	bootstrap := d.Get("bootstrap").(bool)

	if apiURL == "" {
		return &Config{}, fmt.Errorf("[ERROR] No api_url provided")
	}

	config := &Config{
		URL:       NormalizeURL(apiURL),
		AccessKey: accessKey,
		SecretKey: secretKey,
		TokenKey:  tokenKey,
		CACerts:   caCerts,
		Insecure:  insecure,
		Bootstrap: bootstrap,
	}

	// If bootstrap tokenkey accesskey nor secretkey can be provided
	if bootstrap {
		if config.TokenKey != "" || config.AccessKey != "" || config.SecretKey != "" {
			return &Config{}, fmt.Errorf("[ERROR] Bootsrap mode activated. Token_key or access_key and secret_key can not be provided")
		}
	} else {
		// Else token or access key and secret key should be provided
		if config.TokenKey == "" && (config.AccessKey == "" || config.SecretKey == "") {
			return &Config{}, fmt.Errorf("[ERROR] No token_key nor access_key and secret_key are provided")
		}

		_, err := config.ManagementClient()
		if err != nil {
			return &Config{}, err
		}
	}

	return config, nil
}
