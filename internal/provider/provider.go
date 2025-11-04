package provider

import (
	"context"
  "encoding/base64"
  "fmt"
  "net/url"
  "os"
  "strconv"
  "strings"
  "time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
  "github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
  "github.com/hashicorp/terraform-plugin-framework/resource"
  "github.com/hashicorp/terraform-plugin-framework/types"
  "github.com/hashicorp/terraform-plugin-log/tflog"
  client "github.com/rancher/terraform-provider-rancher2/internal/provider/rancher_client"
)

// The `var _` is a special Go construct that results in an unusable variable.
// The purpose of these lines is to make sure our class implements the provider.Provider interface.
// These will fail at compilation time if the implementation is not satisfied.
var _ provider.Provider = &RancherProvider{}

type RancherProvider struct {
	version string
}

type RancherProviderModel struct{
  ApiURL         types.String `tfsdk:"api_url"`
  CACerts        types.String `tfsdk:"ca_certs"`
  IgnoreSystemCA types.Bool   `tfsdk:"ignore_system_ca"`
  Insecure       types.Bool   `tfsdk:"insecure"`
  MaxRedirects   types.Int64  `tfsdk:"max_redirects"`
  Timeout        types.String `tfsdk:"timeout"`
  AccessKey      types.String `tfsdk:"access_key"`
  SecretKey      types.String `tfsdk:"secret_key"`
  TokenKey       types.String `tfsdk:"token_key"`
}

func (p *RancherProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "rancher"
	resp.Version = p.version
}

func (p *RancherProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
      "api_url": schema.StringAttribute{
        Description: "The URL to the rancher API. Example: https://rancher.my-domain.com. "+
          "This can also be set using the RANCHER_API_URL environment variable.",
        Optional:    true,
      },
      "ca_certs": schema.StringAttribute{
        Description: "CA certificates used to sign rancher server TLS certificates. "+
          "This can also be set using the RANCHER_CA_CERTS environment variable.",
        Optional:    true,
      },
      "ignore_system_ca": schema.BoolAttribute{
        Description: "Ignore system CA certificates when validating TLS connections to Rancher. Defaults to false. "+
          "This can also be set using the RANCHER_IGNORE_SYSTEM_CA environment variable.",
        Optional:    true,
      },
      "insecure": schema.BoolAttribute{
        Description: "Allow insecure TLS connections. Defaults to false. "+
          "This can also be set using the RANCHER_INSECURE environment variable.",
        Optional:    true,
      },
      "max_redirects": schema.Int64Attribute{
        Description: "Maximum number of redirects to follow when making requests to the Rancher API. Defaults to 0 (no redirects)."+
          "This can also be set using the RANCHER_MAX_REDIRECTS environment variable.",
        Optional:    true,
      },
      "timeout": schema.StringAttribute{
        Description: "Rancher connection timeout. Golang duration format, ex: '60s'. Defaults to '30s'. "+
          "This can also be set using the RANCHER_TIMEOUT environment variable.",
        Optional:    true,
      },
      "access_key": schema.StringAttribute{
        Description: "API Key used to authenticate with the rancher server. One of access_key and secret_key or token_key must be provided."+
          "This can also be set using the RANCHER_ACCESS_KEY environment variable.",
        Optional:    true,
        Sensitive:   true,
      },
      "secret_key": schema.StringAttribute{
        Description: "API secret used to authenticate with the rancher server. One of access_key and secret_key or token_key must be provided."+
          "This can also be set using the RANCHER_SECRET_KEY environment variable.",
        Optional:    true,
        Sensitive:   true,
      },
      "token_key": schema.StringAttribute{
        Description: "API token used to authenticate with the rancher server. One of access_key and secret_key or token_key must be provided. "+
          "This can also be set using the RANCHER_TOKEN_KEY environment variable.",
        Optional:    true,
        Sensitive:   true,
      },
    },
	}
}

func (p *RancherProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Request Object: %#v", req))
	var err error

  var config RancherProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

  // Connection settings
  ApiURL := os.Getenv("RANCHER_API_URL")
  if !config.ApiURL.IsNull() {
    ApiURL = config.ApiURL.ValueString()
  }
  CACerts := os.Getenv("RANCHER_CA_CERTS")
  if !config.CACerts.IsNull() {
    CACerts = config.CACerts.ValueString()
  }
  IgnoreSystemCA := false
  if os.Getenv("RANCHER_IGNORE_SYSTEM_CA") == "true" {
    IgnoreSystemCA = true
  }
  if !config.IgnoreSystemCA.IsNull() {
    IgnoreSystemCA = config.IgnoreSystemCA.ValueBool()
  }
  Insecure := false
  if os.Getenv("RANCHER_INSECURE") == "true" {
    Insecure = true
  }
  if !config.Insecure.IsNull() {
    Insecure = config.Insecure.ValueBool()
  }
  MaxRedirects := 0
  if os.Getenv("RANCHER_MAX_REDIRECTS") != "" {
    MaxRedirects, err = strconv.Atoi(os.Getenv("RANCHER_MAX_REDIRECTS"))
    if err != nil {
      resp.Diagnostics.AddError("[ERROR] Invalid RANCHER_MAX_REDIRECTS value", err.Error())
      return
    }
  }
  if !config.MaxRedirects.IsNull() {
    MaxRedirects = int(config.MaxRedirects.ValueInt64())
  }

  // Auth settings
  AccessKey := os.Getenv("RANCHER_ACCESS_KEY")
  if !config.AccessKey.IsNull() {
    AccessKey = config.AccessKey.ValueString()
  }
  SecretKey := os.Getenv("RANCHER_SECRET_KEY")
  if !config.SecretKey.IsNull() {
    SecretKey = config.SecretKey.ValueString()
  }
  TokenKey  := os.Getenv("RANCHER_TOKEN_KEY")
  if !config.TokenKey.IsNull() {
    TokenKey = config.TokenKey.ValueString()
  }
  
  // Validate settings below here //

  to := config.Timeout.ValueString()
  Timeout, err := time.ParseDuration(to)
	if err != nil {
    resp.Diagnostics.AddError("[ERROR] Timeout must be in golang duration format, error: %v", err.Error())
    return
	}
  
	if TokenKey == "" && (AccessKey != "") && (SecretKey != "") {
    TokenKey = base64.StdEncoding.EncodeToString([]byte(AccessKey + ":" + SecretKey))
	}

  if ApiURL == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_url"),
			"Missing Rancher API URL",
			"The provider cannot create the Rancher API client as there is no value found for the Rancher API URL. "+
			"Either set the value statically in the configuration, or use the RANCHER_API_URL environment variable.",
		)
	}

  err = isValidURL(ApiURL, Insecure)
  if err != nil {
    resp.Diagnostics.AddAttributeError(
      path.Root("api_url"),
      "Invalid Rancher API URL",
      err.Error(),
    )
  }

  if TokenKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token_key"),
			"Missing Rancher API Token Key",
			"The provider cannot create the Rancher API client as there is no value found for the Rancher Token Key. "+
			"Either set the value statically in the configuration, or use the RANCHER_TOKEN_KEY environment variable. "+
      "Alternatively, you can provide both Access Key and Secret Key to generate the Token Key automatically. "+
      "You can also set the RANCHER_ACCESS_KEY and RANCHER_SECRET_KEY environment variables. ",
		)
	}

  // Check if there were any errors added to diagnostics
	if resp.Diagnostics.HasError() {
		return
	}
  var rancherClient client.RancherClient
  if p.version == "test" {
    tflog.Debug(ctx, "Rancher Provider configured (test version), creating memory client.")
    rancherClient = client.NewRancherMemoryClient(ApiURL, CACerts, IgnoreSystemCA, Insecure, TokenKey, MaxRedirects, Timeout)
  } else {
    rancherClient = client.NewRancherHttpClient(ApiURL, CACerts, IgnoreSystemCA, Insecure, TokenKey, MaxRedirects, Timeout)
  }
	resp.DataSourceData = rancherClient
	resp.ResourceData   = rancherClient

  tflog.Debug(ctx, fmt.Sprintf("Response Object: %#v", resp))
}

func (p *RancherProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
	}
}

func (p *RancherProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &RancherProvider{
			version: version,
		}
	}
}

func isValidURL(u string, insecure bool) error {
  if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
    return fmt.Errorf("Invalid URL scheme: %s", u)
  }
  if strings.HasPrefix(u, "http://") && !insecure {
    return fmt.Errorf("Insecure URL scheme 'http' not allowed without insecure flag set to true: %s", u)
  }
  if strings.Count(u, "/") > 2 {
    return fmt.Errorf("URL path not allowed: %s", u)
  }
  _, err := url.ParseRequestURI(u)
  if err != nil {
    return fmt.Errorf("URL parsing error: %s", err.Error())
  }
  return nil
}
