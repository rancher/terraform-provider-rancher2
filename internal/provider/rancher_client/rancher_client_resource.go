package rancher_client

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	c "github.com/rancher/terraform-provider-rancher2/internal/provider/client"
)

// The "var _" is a special Go construct that results in an unusable variable.
// The purpose of these lines is to make sure our LocalFileResource correctly implements the "resource.Resource“ interface.
// These will fail at compilation time if the implementation is not satisfied.
var _ resource.Resource = &RancherClientResource{}
var _ resource.ResourceWithImportState = &RancherClientResource{}

func NewRancherClientResource() resource.Resource {
	return &RancherClientResource{}
}

type RancherClientResource struct {
	registry *c.ClientRegistry
}

// RancherClientResourceModel describes the resource data model.
type RancherClientResourceModel struct {
	Id             types.String `tfsdk:"id"`
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

func (r *RancherClientResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_client" // rancher_client
}

func (r *RancherClientResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Rancher Client resource. \n" +
			"This resource manages the provider's client used to talk to Rancher. " +
			"Since this is a resource it prevents the provider from accessing Rancher until apply time, respecting the dependency chain." +
			"This also acts as a way to separate out access between resources; " +
			"you can implement least privilege access to resources by tailoring the client to your resources." +
			"Every other resource has a required argument 'rancher_client_id' which refers to this resource's id attribute." +
			"When set up properly, in subsequent applies, only the client which has experienced a resource change will be updated, " +
			"meaning the user making the changes will only have access appropriate to the change.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique identifier for the resource.",
				Required:            true,
			},
			"api_url": schema.StringAttribute{
				Description: "The URL to the rancher API. Example: https://rancher.example.com. " +
					"This can also be set using the RANCHER_API_URL environment variable. " +
					"Environment variable takes precedence over this setting if both are set.",
				Optional: true,
				Computed: true,
			},
			"ca_certs": schema.StringAttribute{
				Description: "CA certificates used to sign rancher server TLS certificates. " +
					"This can also be set using the RANCHER_CA_CERTS environment variable. " +
					"Environment variable takes precedence over this setting if both are set.",
				Optional: true,
				Computed: true,
			},
			"ignore_system_ca": schema.BoolAttribute{
				Description: "Ignore system CA certificates when validating TLS connections to Rancher. Defaults to false. " +
					"This can also be set using the RANCHER_IGNORE_SYSTEM_CA environment variable. " +
					"Environment variable takes precedence over this setting if both are set.",
				Optional: true,
				Computed: true,
			},
			"insecure": schema.BoolAttribute{
				Description: "Allow insecure TLS connections. Defaults to false. " +
					"This can also be set using the RANCHER_INSECURE environment variable. " +
					"Environment variable takes precedence over this setting if both are set.",
				Optional: true,
				Computed: true,
			},
			"max_redirects": schema.Int64Attribute{
				Description: "Maximum number of redirects to follow when making requests to the Rancher API. Defaults to 0 (no redirects)." +
					"This can also be set using the RANCHER_MAX_REDIRECTS environment variable. " +
					"Environment variable takes precedence over this setting if both are set.",
				Optional: true,
				Computed: true,
			},
			"timeout": schema.StringAttribute{
				Description: "Rancher connection timeout. Golang duration format, ex: '60s'. Defaults to '30s'. " +
					"This can also be set using the RANCHER_TIMEOUT environment variable. " +
					"Environment variable takes precedence over this setting if both are set.",
				Optional: true,
				Computed: true,
			},
			"access_key": schema.StringAttribute{
				Description: "API Key used to authenticate with the rancher server. One of access_key and secret_key or token_key must be provided." +
					"This can also be set using the RANCHER_ACCESS_KEY environment variable. " +
					"Environment variable takes precedence over this setting if both are set.",
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"secret_key": schema.StringAttribute{
				Description: "API secret used to authenticate with the rancher server. One of access_key and secret_key or token_key must be provided." +
					"This can also be set using the RANCHER_SECRET_KEY environment variable. " +
					"Environment variable takes precedence over this setting if both are set.",
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
			"token_key": schema.StringAttribute{
				Description: "API token used to authenticate with the rancher server. One of access_key and secret_key or token_key must be provided. " +
					"This can also be set using the RANCHER_TOKEN_KEY environment variable. " +
					"Environment variable takes precedence over this setting if both are set.",
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func (r *RancherClientResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	r.registry = req.ProviderData.(*c.ClientRegistry)
}

// - Create generates reality and state to match plan.
// - Read updates state to match reality.
// - Update updates reality and state to match plan/config (don't compare old state, just override).
// - Destroy destroys reality (state is handled automatically).

// - The Config Phase triggers when the provider loads, it doesn't respect the dependency chain.
// - The Plan Phase detects changes in the config and state, after Read.
// - The Plan Phase respects the dependency chain, but only for reading, it happens before any resource has a chance to Update.
// - The Apply Phase triggers create and update functions, resources respect dependency chain.
// - The Destroy Phase triggers the destroy function, resources destroy in reverse dependency chain order.

func (r *RancherClientResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Create Request Object: %+v", req))

	var plan RancherClientResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := getValues(&plan)
	if err != nil {
		resp.Diagnostics.AddError("Error getting values for create: ", err.Error())
		return
	}
	ID := plan.Id.ValueString()
	pApiURL := plan.ApiURL.ValueString()
	pCACerts := plan.CACerts.ValueString()
	pIgnoreSystemCA := plan.IgnoreSystemCA.ValueBool()
	pInsecure := plan.Insecure.ValueBool()
	pMaxRedirects := plan.MaxRedirects.ValueInt64()
	pTimeout := plan.Timeout.ValueString()
	pAccessKey := plan.AccessKey.ValueString()
	pSecretKey := plan.SecretKey.ValueString()
	pTokenKey := plan.TokenKey.ValueString()

	_, found := r.registry.Load(ID)
	if !found {
		tflog.Debug(ctx, "Using default http client.")
		newClient := c.NewHttpClient(ctx, pApiURL, pCACerts, pIgnoreSystemCA, pInsecure, pAccessKey, pSecretKey, pTokenKey, pMaxRedirects, pTimeout)
		r.registry.Store(ID, newClient)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Debug(ctx, fmt.Sprintf("Create Response Object: %+v", *resp))
}

// Read runs at refresh time which happens before all other functions and every time another function would be called.
func (r *RancherClientResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Read Request Object: %+v", req))

	var state RancherClientResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ID := state.Id.ValueString()

	// On a new terraform run, the registry is empty.
	// We check if the client is in the registry. If not, we re-create it from state and add it.
	// This ensures other resources can find the client and that we don't plan a "destroy and re-create" on every run.
	_, found := r.registry.Load(ID)
	if !found {
		tflog.Debug(ctx, "Client not found in registry. Re-creating from state for ID: "+ID)

		// The getValues function handles environment variables and defaults, which we need.
		err := getValues(&state)
		if err != nil {
			resp.Diagnostics.AddError("Error getting values for read: ", err.Error())
			return
		}

		sApiURL := state.ApiURL.ValueString()
		sCACerts := state.CACerts.ValueString()
		sIgnoreSystemCA := state.IgnoreSystemCA.ValueBool()
		sInsecure := state.Insecure.ValueBool()
		sMaxRedirects := state.MaxRedirects.ValueInt64()
		sTimeout := state.Timeout.ValueString()
		sAccessKey := state.AccessKey.ValueString()
		sSecretKey := state.SecretKey.ValueString()
		sTokenKey := state.TokenKey.ValueString()

		// When re-hydrating, we assume the standard HttpClient. A testing client would be
		// already defined by the test logic, so would never hit this point.
		newClient := c.NewHttpClient(ctx, sApiURL, sCACerts, sIgnoreSystemCA, sInsecure, sAccessKey, sSecretKey, sTokenKey, sMaxRedirects, sTimeout)
		if newClient == nil {
			resp.Diagnostics.AddError("Client Creation Error", "Failed to create new HTTP client from state.")
			return
		}
		r.registry.Store(ID, newClient)
	} else {
		tflog.Debug(ctx, "Client found in registry for ID: "+ID)
	}

	// The state from the request is the source of truth, so we just set it back to confirm the resource exists.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("Read Response Object: %+v", *resp))
}

func (r *RancherClientResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Update Request Object: %+v", req))

	var config RancherClientResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := getValues(&config)
	if err != nil {
		resp.Diagnostics.AddError("Error getting values for update: ", err.Error())
		return
	}
	cApiURL := config.ApiURL.ValueString()
	cCACerts := config.CACerts.ValueString()
	cIgnoreSystemCA := config.IgnoreSystemCA.ValueBool()
	cInsecure := config.Insecure.ValueBool()
	cMaxRedirects := config.MaxRedirects.ValueInt64()
	cTimeout := config.Timeout.ValueString()
	cAccessKey := config.AccessKey.ValueString()
	cSecretKey := config.SecretKey.ValueString()
	cTokenKey := config.TokenKey.ValueString()
	ID := config.Id.ValueString()

	tflog.Debug(ctx, "Using default http client.")
	newClient := c.NewHttpClient(ctx, cApiURL, cCACerts, cIgnoreSystemCA, cInsecure, cAccessKey, cSecretKey, cTokenKey, cMaxRedirects, cTimeout)
	r.registry.Store(ID, newClient)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
	tflog.Debug(ctx, fmt.Sprintf("Update Response Object: %+v", *resp))
}

func (r *RancherClientResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Delete Request Object: %+v", req))

	var state RancherClientResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := getValues(&state)
	if err != nil {
		resp.Diagnostics.AddError("Error getting values for delete: ", err.Error())
		return
	}

	ID := state.Id.ValueString()

	_, found := r.registry.Load(ID)
	if found {
		r.registry.Delete(ID)
	}

	tflog.Debug(ctx, fmt.Sprintf("Delete Response Object: %+v", *resp))
}

func (r *RancherClientResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getValues(d *RancherClientResourceModel) error {
	// Environment variables override config settings
	var err error

	var ApiURL string
	if !d.ApiURL.IsNull() {
		ApiURL = d.ApiURL.ValueString()
	}
	if os.Getenv("RANCHER_API_URL") != "" {
		ApiURL = os.Getenv("RANCHER_API_URL")
	}

	var CACerts string
	if !d.CACerts.IsNull() {
		CACerts = d.CACerts.ValueString()
	}
	if os.Getenv("RANCHER_CA_CERTS") != "" {
		CACerts = os.Getenv("RANCHER_CA_CERTS")
	}

	IgnoreSystemCA := false
	if !d.IgnoreSystemCA.IsNull() {
		IgnoreSystemCA = d.IgnoreSystemCA.ValueBool()
	}
	if os.Getenv("RANCHER_IGNORE_SYSTEM_CA") == "true" {
		IgnoreSystemCA = true
	}

	Insecure := false
	if !d.Insecure.IsNull() {
		Insecure = d.Insecure.ValueBool()
	}
	if os.Getenv("RANCHER_INSECURE") == "true" {
		Insecure = true
	}

	MaxRedirects := 0
	if !d.MaxRedirects.IsNull() {
		MaxRedirects = int(d.MaxRedirects.ValueInt64())
	}
	if os.Getenv("RANCHER_MAX_REDIRECTS") != "" {
		MaxRedirects, err = strconv.Atoi(os.Getenv("RANCHER_MAX_REDIRECTS"))
		if err != nil {
			return fmt.Errorf("[ERROR] Invalid RANCHER_MAX_REDIRECTS value: %v", err)
		}
	}

	var Timeout string
	if !d.Timeout.IsNull() {
		Timeout = d.Timeout.ValueString()
	}
	if os.Getenv("RANCHER_TIMEOUT") != "" {
		Timeout = os.Getenv("RANCHER_TIMEOUT")
	}

	var AccessKey string
	if !d.AccessKey.IsNull() {
		AccessKey = d.AccessKey.ValueString()
	}
	if os.Getenv("RANCHER_ACCESS_KEY") != "" {
		AccessKey = os.Getenv("RANCHER_ACCESS_KEY")
	}

	var SecretKey string
	if !d.SecretKey.IsNull() {
		SecretKey = d.SecretKey.ValueString()
	}
	if os.Getenv("RANCHER_SECRET_KEY") != "" {
		SecretKey = os.Getenv("RANCHER_SECRET_KEY")
	}

	var TokenKey string
	if !d.TokenKey.IsNull() {
		TokenKey = d.TokenKey.ValueString()
	}
	if os.Getenv("RANCHER_TOKEN_KEY") != "" {
		TokenKey = os.Getenv("RANCHER_TOKEN_KEY")
	}

	// Validate Data Here //
	if ApiURL == "" {
		return fmt.Errorf("[ERROR] No API URL detected, please either specify one in your config or in the RANCHER_API_URL environment variable.")
	}
	err = isValidURL(ApiURL, Insecure)
	if err != nil {
		return fmt.Errorf("[ERROR] Invalid Rancher API URL, error: %v", err.Error())
	}

	// this is just validating that the string supplied is a valid duration
	_, err = time.ParseDuration(Timeout)
	if err != nil {
		return fmt.Errorf("[ERROR] Timeout must be in golang duration format, error: %v", err.Error())
	}

	if TokenKey == "" && (AccessKey != "") && (SecretKey != "") {
		TokenKey = base64.StdEncoding.EncodeToString([]byte(AccessKey + ":" + SecretKey))
	}

	d.ApiURL = types.StringValue(ApiURL)
	d.CACerts = types.StringValue(CACerts)
	d.IgnoreSystemCA = types.BoolValue(IgnoreSystemCA)
	d.Insecure = types.BoolValue(Insecure)
	d.MaxRedirects = types.Int64Value(int64(MaxRedirects))
	d.Timeout = types.StringValue(Timeout)
	d.AccessKey = types.StringValue(AccessKey)
	d.SecretKey = types.StringValue(SecretKey)
	d.TokenKey = types.StringValue(TokenKey)

	return nil
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
