package provider

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	c "github.com/rancher/terraform-provider-rancher2/internal/provider/client"
	pp "github.com/rancher/terraform-provider-rancher2/internal/provider/pretty_print"
	"github.com/rancher/terraform-provider-rancher2/internal/provider/rancher2_dev"
	"github.com/rancher/terraform-provider-rancher2/internal/provider/rancher2_login"
	"github.com/rancher/terraform-provider-rancher2/internal/provider/validators"
)

// The `var _` is a special Go construct that results in an unusable variable.
// The purpose of these lines is to make sure our class implements the provider.Provider interface.
// These will fail at compilation time if the implementation is not satisfied.
var _ provider.Provider = &RancherProvider{}

type RancherProvider struct {
	version string
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &RancherProvider{
			version: version,
		}
	}
}

type RancherProviderModel struct {
	ApiUrl         types.String `tfsdk:"api_url"`
	CaCert         types.String `tfsdk:"ca_cert"`
	Insecure       types.Bool   `tfsdk:"insecure"`
	IgnoreSystemCa types.Bool   `tfsdk:"ignore_system_ca"`
	Timeout        types.Int64  `tfsdk:"timeout"`
	MaxRedirects   types.Int64  `tfsdk:"max_redirects"`
	Token          types.String `tfsdk:"token"`
}

// ToPlan returns a tfsdk.Plan with the data from the model.
func (m *RancherProviderModel) ToConfig(ctx context.Context) (tfsdk.Config, diag.Diagnostics) {
	p := RancherProvider{}
	s := &provider.SchemaResponse{}
	p.Schema(ctx, provider.SchemaRequest{}, s)

	var raw attr.Value
	diags := tfsdk.ValueFrom(ctx, m, s.Schema.Type(), &raw)
	if diags.HasError() {
		return tfsdk.Config{}, diags
	}

	tfRaw, err := raw.ToTerraformValue(ctx)
	if err != nil {
		dgs := diag.Diagnostics{}
		dgs.AddError("Error converting to terraform value", err.Error())
		return tfsdk.Config{}, dgs
	}

	return tfsdk.Config{
		Raw:    tfRaw,
		Schema: s.Schema,
	}, nil
}

func (p *RancherProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "rancher2"
	resp.Version = p.version
}

func (p *RancherProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_url": schema.StringAttribute{
				Description: "The URL of the rancher API without paths included, e.g. https://rancher.example.com.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^https?:\/\/`),
						"must be a valid URL starting with http:// or https://",
					),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}`),
						"must be a valid domain",
					),
				},
			},
			"ca_cert": schema.StringAttribute{
				Description: "The CA certificate to use for the Rancher API.",
				Optional:    true,
				Validators: []validator.String{
					validators.IsCertificate(),
				},
			},
			"insecure": schema.BoolAttribute{
				Description: "Whether to allow insecure connections to the Rancher API.",
				Optional:    true,
			},
			"ignore_system_ca": schema.BoolAttribute{
				Description: "Whether to ignore the system's CA certificates.",
				Optional:    true,
			},
			"timeout": schema.Int64Attribute{
				Description: "The connection timeout to use for the Rancher API.",
				Optional:    true,
			},
			"max_redirects": schema.Int64Attribute{
				Description: "The maximum number of redirects to follow.",
				Optional:    true,
			},
			"token": schema.StringAttribute{
				Description: "The token to use for authentication, can be set with RANCHER_TOKEN.",
				Optional:    true,
			},
		},
	}
}

func (p *RancherProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Provider Config: %s", pp.PrettyPrint(req.Config.Raw)))

	var config RancherProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	if config.ApiUrl.IsNull() || config.ApiUrl.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing Rancher API URL",
			"The API URL must be set at the provider level for resources to work properly.",
		)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Rancher API URL: %s", config.ApiUrl.ValueString()))

	if config.Insecure.IsNull() {
		config.Insecure = types.BoolValue(false)
	}

	if config.IgnoreSystemCa.IsNull() {
		config.IgnoreSystemCa = types.BoolValue(false)
	}

	if config.Timeout.IsNull() {
		config.Timeout = types.Int64Value(30)
	}

	if config.MaxRedirects.IsNull() {
		config.MaxRedirects = types.Int64Value(10)
	}

	_, err := url.ParseRequestURI(config.ApiUrl.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Rancher API URL",
			fmt.Sprintf("The provider received an invalid Rancher API URL: %s", err.Error()),
		)
		return
	}

	token := config.Token.ValueString()
	if config.Token.IsNull() {
		token = os.Getenv("RANCHER_TOKEN")
	}

	var client c.Client
	if p.version == "test" {
		client = c.NewTestClient(
			ctx,
			config.ApiUrl.ValueString(),
			config.CaCert.ValueString(),
			config.Insecure.ValueBool(),
			config.IgnoreSystemCa.ValueBool(),
			time.Duration(config.Timeout.ValueInt64())*time.Second,
			config.MaxRedirects.ValueInt64(),
			token,
		)
	} else {
		client = c.NewHttpClient(
			ctx,
			config.ApiUrl.ValueString(),
			config.CaCert.ValueString(),
			config.Insecure.ValueBool(),
			config.IgnoreSystemCa.ValueBool(),
			time.Duration(config.Timeout.ValueInt64())*time.Second,
			config.MaxRedirects.ValueInt64(),
			token,
		)
	}

	// The client variable is an interface that holds a pointer to a struct.
	resp.ResourceData = client
	resp.DataSourceData = client

	tflog.Debug(ctx, fmt.Sprintf("Provider Configure Response: %+v", pp.PrettyPrint(resp)))
}

func (p *RancherProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		rancher2_dev.NewRancherDevResource,
		rancher2_login.NewRancherLoginResource,
	}
}

func (p *RancherProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
