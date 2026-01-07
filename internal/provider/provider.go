package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	c "github.com/rancher/terraform-provider-rancher2/internal/provider/client"
	"github.com/rancher/terraform-provider-rancher2/internal/provider/rancher_client"
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

// This provider has no configuration.
type RancherProviderModel struct{}

func (p *RancherProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "rancher"
	resp.Version = p.version
}

// This provider has no configuration arguments.
func (p *RancherProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
}

func (p *RancherProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Provider Configure Request Object: %#v", req))

	var config RancherProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	registry := c.NewRegistry()
	resp.ResourceData = registry
	resp.DataSourceData = registry

	tflog.Debug(ctx, fmt.Sprintf("Provider Configure Response Object: %#v", resp))
}

func (p *RancherProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		rancher_client.NewRancherClientResource,
	}
}

func (p *RancherProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
