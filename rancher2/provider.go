package rancher2

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const (
	providerDefaultEmptyString = "nil"
)

var (
	descriptions                        map[string]string
	rancher2ClusterRKEK8SDefaultVersion string
	rancher2ClusterRKEK8SVersions       []string
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
			"api_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_URL", providerDefaultEmptyString),
				Description: descriptions["api_url"],
			},
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_ACCESS_KEY", providerDefaultEmptyString),
				Description: descriptions["access_key"],
			},
			"bootstrap": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_BOOTSTRAP", false),
				Description: descriptions["bootstrap"],
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_SECRET_KEY", providerDefaultEmptyString),
				Description: descriptions["secret_key"],
			},
			"token_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_TOKEN_KEY", providerDefaultEmptyString),
				Description: descriptions["token_key"],
			},
			"ca_certs": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_CA_CERTS", ""),
				Description: descriptions["ca_certs"],
			},
			"insecure": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("RANCHER_INSECURE", false),
				Description: descriptions["insecure"],
			},
			"retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Deprecated:  "Use timeout instead",
				Description: descriptions["retries"],
			},
			"timeout": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     rancher2DefaultTimeout,
				Description: descriptions["timeout"],
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v, ok := val.(string)
					if !ok || len(v) == 0 {
						return
					}
					_, err := time.ParseDuration(v)
					if err != nil {
						errs = append(errs, fmt.Errorf("%q must be in golang duration format, error: %v", key, err))
					}
					return
				},
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"rancher2_app":                           resourceRancher2App(),
			"rancher2_app_v2":                        resourceRancher2AppV2(),
			"rancher2_auth_config_activedirectory":   resourceRancher2AuthConfigActiveDirectory(),
			"rancher2_auth_config_adfs":              resourceRancher2AuthConfigADFS(),
			"rancher2_auth_config_azuread":           resourceRancher2AuthConfigAzureAD(),
			"rancher2_auth_config_freeipa":           resourceRancher2AuthConfigFreeIpa(),
			"rancher2_auth_config_github":            resourceRancher2AuthConfigGithub(),
			"rancher2_auth_config_keycloak":          resourceRancher2AuthConfigKeyCloak(),
			"rancher2_auth_config_okta":              resourceRancher2AuthConfigOKTA(),
			"rancher2_auth_config_openldap":          resourceRancher2AuthConfigOpenLdap(),
			"rancher2_auth_config_ping":              resourceRancher2AuthConfigPing(),
			"rancher2_bootstrap":                     resourceRancher2Bootstrap(),
			"rancher2_catalog":                       resourceRancher2Catalog(),
			"rancher2_catalog_v2":                    resourceRancher2CatalogV2(),
			"rancher2_certificate":                   resourceRancher2Certificate(),
			"rancher2_cloud_credential":              resourceRancher2CloudCredential(),
			"rancher2_cluster":                       resourceRancher2Cluster(),
			"rancher2_cluster_v2":                    resourceRancher2ClusterV2(),
			"rancher2_cluster_alert_group":           resourceRancher2ClusterAlertGroup(),
			"rancher2_cluster_alert_rule":            resourceRancher2ClusterAlertRule(),
			"rancher2_cluster_driver":                resourceRancher2ClusterDriver(),
			"rancher2_cluster_role_template_binding": resourceRancher2ClusterRoleTemplateBinding(),
			"rancher2_cluster_sync":                  resourceRancher2ClusterSync(),
			"rancher2_cluster_template":              resourceRancher2ClusterTemplate(),
			"rancher2_config_map_v2":                 resourceRancher2ConfigMapV2(),
			"rancher2_etcd_backup":                   resourceRancher2EtcdBackup(),
			"rancher2_feature":                       resourceRancher2Feature(),
			"rancher2_global_dns":                    resourceRancher2GlobalDNS(),
			"rancher2_global_dns_provider":           resourceRancher2GlobalDNSProvider(),
			"rancher2_global_role":                   resourceRancher2GlobalRole(),
			"rancher2_global_role_binding":           resourceRancher2GlobalRoleBinding(),
			"rancher2_machine_config_v2":             resourceRancher2MachineConfigV2(),
			"rancher2_multi_cluster_app":             resourceRancher2MultiClusterApp(),
			"rancher2_namespace":                     resourceRancher2Namespace(),
			"rancher2_node_driver":                   resourceRancher2NodeDriver(),
			"rancher2_node_pool":                     resourceRancher2NodePool(),
			"rancher2_node_template":                 resourceRancher2NodeTemplate(),
			"rancher2_notifier":                      resourceRancher2Notifier(),
			"rancher2_pod_security_policy_template":  resourceRancher2PodSecurityPolicyTemplate(),
			"rancher2_project":                       resourceRancher2Project(),
			"rancher2_project_alert_group":           resourceRancher2ProjectAlertGroup(),
			"rancher2_project_alert_rule":            resourceRancher2ProjectAlertRule(),
			"rancher2_project_role_template_binding": resourceRancher2ProjectRoleTemplateBinding(),
			"rancher2_registry":                      resourceRancher2Registry(),
			"rancher2_role_template":                 resourceRancher2RoleTemplate(),
			"rancher2_secret":                        resourceRancher2Secret(),
			"rancher2_secret_v2":                     resourceRancher2SecretV2(),
			"rancher2_setting":                       resourceRancher2Setting(),
			"rancher2_storage_class_v2":              resourceRancher2StorageClassV2(),
			"rancher2_token":                         resourceRancher2Token(),
			"rancher2_user":                          resourceRancher2User(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"rancher2_app":                           dataSourceRancher2App(),
			"rancher2_catalog":                       dataSourceRancher2Catalog(),
			"rancher2_catalog_v2":                    dataSourceRancher2CatalogV2(),
			"rancher2_certificate":                   dataSourceRancher2Certificate(),
			"rancher2_cloud_credential":              dataSourceRancher2CloudCredential(),
			"rancher2_cluster":                       dataSourceRancher2Cluster(),
			"rancher2_cluster_v2":                    dataSourceRancher2ClusterV2(),
			"rancher2_cluster_alert_group":           dataSourceRancher2ClusterAlertGroup(),
			"rancher2_cluster_alert_rule":            dataSourceRancher2ClusterAlertRule(),
			"rancher2_cluster_driver":                dataSourceRancher2ClusterDriver(),
			"rancher2_cluster_role_template_binding": dataSourceRancher2ClusterRoleTemplateBinding(),
			"rancher2_cluster_template":              dataSourceRancher2ClusterTemplate(),
			"rancher2_config_map_v2":                 dataSourceRancher2ConfigMapV2(),
			"rancher2_etcd_backup":                   dataSourceRancher2EtcdBackup(),
			"rancher2_global_dns_provider":           dataSourceRancher2GlobalDNSProvider(),
			"rancher2_global_role":                   dataSourceRancher2GlobalRole(),
			"rancher2_global_role_binding":           dataSourceRancher2GlobalRoleBinding(),
			"rancher2_multi_cluster_app":             dataSourceRancher2MultiClusterApp(),
			"rancher2_namespace":                     dataSourceRancher2Namespace(),
			"rancher2_node_driver":                   dataSourceRancher2NodeDriver(),
			"rancher2_node_pool":                     dataSourceRancher2NodePool(),
			"rancher2_node_template":                 dataSourceRancher2NodeTemplate(),
			"rancher2_notifier":                      dataSourceRancher2Notifier(),
			"rancher2_pod_security_policy_template":  dataSourceRancher2PodSecurityPolicyTemplate(),
			"rancher2_principal":                     dataSourceRancher2Principal(),
			"rancher2_project":                       dataSourceRancher2Project(),
			"rancher2_project_alert_group":           dataSourceRancher2ProjectAlertGroup(),
			"rancher2_project_alert_rule":            dataSourceRancher2ProjectAlertRule(),
			"rancher2_project_role_template_binding": dataSourceRancher2ProjectRoleTemplateBinding(),
			"rancher2_registry":                      dataSourceRancher2Registry(),
			"rancher2_role_template":                 dataSourceRancher2RoleTemplate(),
			"rancher2_secret":                        dataSourceRancher2Secret(),
			"rancher2_secret_v2":                     dataSourceRancher2SecretV2(),
			"rancher2_setting":                       dataSourceRancher2Setting(),
			"rancher2_storage_class_v2":              dataSourceRancher2StorageClassV2(),
			"rancher2_user":                          dataSourceRancher2User(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func init() {
	descriptions = map[string]string{
		"access_key": "API Key used to authenticate with the rancher server",
		"secret_key": "API secret used to authenticate with the rancher server",
		"token_key":  "API token used to authenticate with the rancher server",
		"ca_certs":   "CA certificates used to sign rancher server tls certificates. Mandatory if self signed tls and insecure option false",
		"insecure":   "Allow insecure connections to Rancher. Mandatory if self signed tls and not ca_certs provided",
		"api_url":    "The URL to the rancher API",
		"bootstrap":  "Bootstrap rancher server",
		"retries":    "Rancher connection retries",
		"timeout":    "Rancher connection timeout (retry every 5s). Golang duration format, ex: \"60s\"",
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

	timeout, err := time.ParseDuration(d.Get("timeout").(string))
	if err != nil {
		return nil, fmt.Errorf("[ERROR] timeout must be in golang duration format, error: %v", err)
	}

	// Set tokenKey based on accessKey and secretKey if needed
	if tokenKey == providerDefaultEmptyString && accessKey != providerDefaultEmptyString && secretKey != providerDefaultEmptyString {
		tokenKey = accessKey + ":" + secretKey
	}

	config := &Config{
		URL:       apiURL,
		TokenKey:  tokenKey,
		CACerts:   caCerts,
		Insecure:  insecure,
		Bootstrap: bootstrap,
		Timeout:   timeout,
	}

	return providerValidateConfig(config)
}

func providerValidateConfig(config *Config) (*Config, error) {
	err := config.NormalizeURL()
	if err != nil {
		return &Config{}, fmt.Errorf("[ERROR] %v", err)
	}
	if config.Bootstrap {
		// If bootstrap tokenkey accesskey nor secretkey can be provided
		if config.TokenKey != providerDefaultEmptyString {
			return &Config{}, fmt.Errorf("[ERROR] Bootsrap mode activated. Token_key or access_key and secret_key can not be provided")
		}
	} else {
		// Else token or access key and secret key should be provided
		if config.TokenKey == providerDefaultEmptyString {
			return &Config{}, fmt.Errorf("[ERROR] No token_key nor access_key and secret_key are provided")
		}
	}

	// Setting default timeout if not specified
	if config.Timeout.Seconds() == 0 {
		timeout, err := time.ParseDuration(rancher2DefaultTimeout)
		if err != nil {
			return nil, fmt.Errorf("[ERROR] default timeout is not in golang duration format, error: %v", err)
		}
		config.Timeout = timeout
	}

	return config, nil
}
