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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
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
	Registry *c.ClientRegistry
	Version  string
}

// RancherClientResourceModel describes the resource data model.
type RancherClientResourceModel struct {
	Id             types.String `tfsdk:"id"`
	ApiURL         types.String `tfsdk:"api_url"`
	CACerts        types.String `tfsdk:"ca_certs"`
	IgnoreSystemCA types.Bool   `tfsdk:"ignore_system_ca"`
	Insecure       types.Bool   `tfsdk:"insecure"`
	MaxRedirects   types.String `tfsdk:"max_redirects"`
	ConnectTimeout types.String `tfsdk:"connect_timeout"`
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
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"insecure": schema.BoolAttribute{
				Description: "Allow insecure TLS connections. Defaults to false. " +
					"This can also be set using the RANCHER_INSECURE environment variable. " +
					"Environment variable takes precedence over this setting if both are set.",
				Default:  booldefault.StaticBool(false),
				Optional: true,
				Computed: true,
			},
			"max_redirects": schema.StringAttribute{
				Description: "Maximum number of redirects to follow when making requests to the Rancher API, defaults to 3." +
					"This can also be set using the RANCHER_MAX_REDIRECTS environment variable. " +
					"Environment variable takes precedence over this setting if both are set.",
				Default:  stringdefault.StaticString("3"),
				Optional: true,
				Computed: true,
			},
			"connect_timeout": schema.StringAttribute{
				Description: "Rancher connection timeout. Golang duration format, ex: '60s'. Defaults to '30s'. " +
					"This can also be set using the RANCHER_TIMEOUT environment variable. " +
					"Environment variable takes precedence over this setting if both are set.",
				Default:  stringdefault.StaticString("30s"),
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
	if req.ProviderData == nil {
		return // Prevent panic if the provider has not been configured.
	}

	registry, ok := req.ProviderData.(*c.ClientRegistry)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *c.ClientRegistry, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.Registry = registry
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
	var err error

	var plan RancherClientResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	getDefaultValues(&plan)

	var data RancherClientResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	getDefaultValues(&data)
	getEnvironmentValues(&data)
	err = validateValues(&data)
	if err != nil {
		resp.Diagnostics.AddError("Error getting values for create: ", err.Error())
		return
	}

	ID := data.Id.ValueString()
	pApiURL := data.ApiURL.ValueString()
	pCACerts := data.CACerts.ValueString()
	pIgnoreSystemCA := data.IgnoreSystemCA.ValueBool()
	pInsecure := data.Insecure.ValueBool()
	pMaxRedirects := data.MaxRedirects.ValueString()
	pTimeout := data.ConnectTimeout.ValueString()
	pAccessKey := data.AccessKey.ValueString()
	pSecretKey := data.SecretKey.ValueString()
	pTokenKey := data.TokenKey.ValueString()

	pMR, err := strconv.ParseInt(pMaxRedirects, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("[ERROR] Invalid MaxRedirects value: %v", pMaxRedirects), err.Error())
		return
	}

	_, found := r.Registry.Load(ID)
	if !found {
		tflog.Debug(ctx, "Using default http client.")
		newClient := c.NewHttpClient(ctx, pApiURL, pCACerts, pIgnoreSystemCA, pInsecure, pAccessKey, pSecretKey, pTokenKey, pMR, pTimeout)
		r.Registry.Store(ID, newClient)
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
	_, found := r.Registry.Load(ID)
	if !found {
		tflog.Debug(ctx, "Client not found in registry. Re-creating from state for ID: "+ID)

		getDefaultValues(&state)
		err := validateValues(&state)
		if err != nil {
			resp.Diagnostics.AddError("Error getting values for create: ", err.Error())
			return
		}

		sApiURL := state.ApiURL.ValueString()
		sCACerts := state.CACerts.ValueString()
		sIgnoreSystemCA := state.IgnoreSystemCA.ValueBool()
		sInsecure := state.Insecure.ValueBool()
		sMaxRedirects := state.MaxRedirects.ValueString()
		sTimeout := state.ConnectTimeout.ValueString()
		sAccessKey := state.AccessKey.ValueString()
		sSecretKey := state.SecretKey.ValueString()
		sTokenKey := state.TokenKey.ValueString()

		sMR, err := strconv.ParseInt(sMaxRedirects, 10, 64)
		if err != nil {
			resp.Diagnostics.AddError(fmt.Sprintf("[ERROR] Invalid MaxRedirects value: %v", sMaxRedirects), err.Error())
			return
		}

		// When re-hydrating, we assume the standard HttpClient. A testing client would be
		// already defined by the test logic, so would never hit this point.
		newClient := c.NewHttpClient(ctx, sApiURL, sCACerts, sIgnoreSystemCA, sInsecure, sAccessKey, sSecretKey, sTokenKey, sMR, sTimeout)
		if newClient == nil {
			resp.Diagnostics.AddError("Client Creation Error", "Failed to create new HTTP client from state.")
			return
		}
		r.Registry.Store(ID, newClient)
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

	getDefaultValues(&config)
	err := validateValues(&config)
	if err != nil {
		resp.Diagnostics.AddError("Error getting values for create: ", err.Error())
		return
	}

	cApiURL := config.ApiURL.ValueString()
	cCACerts := config.CACerts.ValueString()
	cIgnoreSystemCA := config.IgnoreSystemCA.ValueBool()
	cInsecure := config.Insecure.ValueBool()
	cMaxRedirects := config.MaxRedirects.ValueString()
	cTimeout := config.ConnectTimeout.ValueString()
	cAccessKey := config.AccessKey.ValueString()
	cSecretKey := config.SecretKey.ValueString()
	cTokenKey := config.TokenKey.ValueString()
	ID := config.Id.ValueString()

	cMR, err := strconv.ParseInt(cMaxRedirects, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("[ERROR] Invalid MaxRedirects value: %v", cMaxRedirects), err.Error())
		return
	}

	tflog.Debug(ctx, "Using default http client.")
	newClient := c.NewHttpClient(ctx, cApiURL, cCACerts, cIgnoreSystemCA, cInsecure, cAccessKey, cSecretKey, cTokenKey, cMR, cTimeout)
	r.Registry.Store(ID, newClient)

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

	getDefaultValues(&state)
	err := validateValues(&state)
	if err != nil {
		resp.Diagnostics.AddError("Error getting values for create: ", err.Error())
		return
	}

	ID := state.Id.ValueString()

	_, found := r.Registry.Load(ID)
	if found {
		r.Registry.Delete(ID)
	}

	tflog.Debug(ctx, fmt.Sprintf("Delete Response Object: %+v", *resp))
}

func (r *RancherClientResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getDefaultValues(d *RancherClientResourceModel) {

	if d.MaxRedirects.ValueString() == "" {
		d.MaxRedirects = types.StringValue("3")
	}

	if d.ConnectTimeout.ValueString() == "" {
		d.ConnectTimeout = types.StringValue("30s")
	}

	if d.TokenKey.ValueString() == "" && (d.AccessKey.ValueString() != "") && (d.SecretKey.ValueString() != "") {
		t := base64.StdEncoding.EncodeToString([]byte(d.AccessKey.ValueString() + ":" + d.SecretKey.ValueString()))
		d.TokenKey = types.StringValue(t)
	}
}

func validateValues(d *RancherClientResourceModel) error {
	if d.ApiURL.ValueString() == "" {
		return fmt.Errorf("[ERROR] No API URL detected, please either specify one in your config or in the RANCHER_API_URL environment variable")
	}
	err := isValidURL(d.ApiURL.ValueString(), d.Insecure.ValueBool())
	if err != nil {
		return fmt.Errorf("[ERROR] Invalid Rancher API URL, error: %v", err.Error())
	}

	// this is just validating that the string supplied is a valid duration
	_, err = time.ParseDuration(d.ConnectTimeout.ValueString())
	if err != nil {
		return fmt.Errorf("[ERROR] Timeout must be in golang duration format, have %v, error: %v", d.ConnectTimeout.ValueString(), err.Error())
	}

	return nil
}

func getEnvironmentValues(d *RancherClientResourceModel) {
	if os.Getenv("RANCHER_API_URL") != "" {
		d.ApiURL = types.StringValue(os.Getenv("RANCHER_API_URL"))
	}
	if os.Getenv("RANCHER_CA_CERTS") != "" {
		d.CACerts = types.StringValue(os.Getenv("RANCHER_CA_CERTS"))
	}
	if os.Getenv("RANCHER_IGNORE_SYSTEM_CA") == "true" {
		d.IgnoreSystemCA = types.BoolValue(true)
	}
	if os.Getenv("RANCHER_INSECURE") == "true" {
		d.Insecure = types.BoolValue(true)
	}
	if os.Getenv("RANCHER_MAX_REDIRECTS") != "" {
		d.MaxRedirects = types.StringValue(os.Getenv("RANCHER_MAX_REDIRECTS"))
	}
	if os.Getenv("RANCHER_TIMEOUT") != "" {
		d.ConnectTimeout = types.StringValue(os.Getenv("RANCHER_TIMEOUT"))
	}
	if os.Getenv("RANCHER_ACCESS_KEY") != "" {
		d.AccessKey = types.StringValue(os.Getenv("RANCHER_ACCESS_KEY"))
	}
	if os.Getenv("RANCHER_SECRET_KEY") != "" {
		d.SecretKey = types.StringValue(os.Getenv("RANCHER_SECRET_KEY"))
	}
	if os.Getenv("RANCHER_TOKEN_KEY") != "" {
		d.TokenKey = types.StringValue(os.Getenv("RANCHER_TOKEN_KEY"))
	}
	if d.TokenKey.ValueString() == "" && (d.AccessKey.ValueString() != "") && (d.SecretKey.ValueString() != "") {
		t := base64.StdEncoding.EncodeToString([]byte(d.AccessKey.ValueString() + ":" + d.SecretKey.ValueString()))
		d.TokenKey = types.StringValue(t)
	}
}

func isValidURL(u string, insecure bool) error {
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		return fmt.Errorf("invalid URL scheme: %s", u)
	}
	if strings.HasPrefix(u, "http://") && !insecure {
		return fmt.Errorf("insecure URL scheme 'http' not allowed without insecure flag set to true: %s", u)
	}
	if strings.Count(u, "/") > 2 {
		return fmt.Errorf("url path not allowed: %s", u)
	}
	_, err := url.ParseRequestURI(u)
	if err != nil {
		return fmt.Errorf("url parsing error: %s", err.Error())
	}
	return nil
}
