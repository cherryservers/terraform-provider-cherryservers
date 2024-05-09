// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"github.com/cherryservers/cherrygo/v3"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure CherryServersProvider satisfies various provider interfaces.
var _ provider.Provider = &CherryServersProvider{}
var _ provider.ProviderWithFunctions = &CherryServersProvider{}

// CherryServersProvider defines the provider implementation.
type CherryServersProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// CherryServersProviderModel describes the provider data model.
type CherryServersProviderModel struct {
	APIKey types.String `tfsdk:"api_key"`
}

func (p *CherryServersProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cherryservers"
	resp.Version = p.version
}

func (p *CherryServersProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Description: "Cherry Servers [API Key](https://portal.cherryservers.com/settings/api-keys) that allows interactions with the API",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *CherryServersProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring CherryServers client")

	var data CherryServersProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }
	if data.APIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown CherryServers API Key",
			"The provider cannot create the CherryServers API client as there is an unknown configuration value for the CherryServers API Key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CHERRY_AUTH_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	apiKey := os.Getenv("CHERRY_AUTH_KEY")

	if !data.APIKey.IsNull() {
		apiKey = data.APIKey.ValueString()
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing CherryServers API Key",
			"The provider cannot create the CherryServers API client as there is a missing or empty value for the CherryServers API key. "+
				"Set the username value in the configuration or use the CHERRY_AUTH_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "cherryservers_api_key", apiKey)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "cherryservers_api_key")

	tflog.Debug(ctx, "Creating CherryServers client")

	// Example client configuration for data sources and resources
	userAgent := fmt.Sprintf("terraform-provider/cherryservers/%s terraform/%s", p.version, req.TerraformVersion)
	args := []cherrygo.ClientOpt{cherrygo.WithAuthToken(apiKey), cherrygo.WithUserAgent(userAgent)}
	client, err := cherrygo.NewClient(args...)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create CherryServers API Client",
			"An unexpected error occurred when creating the CherryServers API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"CherryServers Client Error: "+err.Error(),
		)
		return
	}
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Successfully created CherryServers client")
}

func (p *CherryServersProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewProjectResource,
		NewIpResource,
		NewServerResource,
	}
}

func (p *CherryServersProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewProjectDataSource,
	}
}

func (p *CherryServersProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewExampleFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CherryServersProvider{
			version: version,
		}
	}
}
