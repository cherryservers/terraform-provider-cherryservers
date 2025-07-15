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
	"terraform-provider-cherryservers/internal/provider/datasourcebase"
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
	APIToken types.String `tfsdk:"api_token"`
}

func (p *CherryServersProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cherryservers"
	resp.Version = p.version
}

func (p *CherryServersProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_token": schema.StringAttribute{
				Description: "Cherry Servers [API Key](https://portal.cherryservers.com/settings/api-keys) that allows interactions with the API.",
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
	if data.APIToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Unknown CherryServers API Token",
			"The provider cannot create the CherryServers API client as there is an unknown configuration value for the CherryServers API Token. "+
				"Either target apply the source of the value first, set the value statically in the configuration,"+
				" or use the CHERRY_AUTH_TOKEN or CHERRY_AUTH_KEY environment variables.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	apiToken := os.Getenv("CHERRY_AUTH_KEY")
	if apiToken == "" {
		apiToken = os.Getenv("CHERRY_AUTH_TOKEN")
	}

	if !data.APIToken.IsNull() {
		apiToken = data.APIToken.ValueString()
	}

	if apiToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Missing CherryServers API Token",
			"The provider cannot create the CherryServers API client as there is a missing or empty value for the CherryServers API token. "+
				"Set the API token value in the configuration or use the CHERRY_AUTH_TOKEN or CHERRY_AUTH_KEY environment variables. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "cherryservers_api_token", apiToken)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "cherryservers_api_token")

	tflog.Debug(ctx, "Creating CherryServers client")

	// Example client configuration for data sources and resources
	userAgent := fmt.Sprintf("terraform-provider/cherryservers/%s terraform/%s", p.version, req.TerraformVersion)
	args := []cherrygo.ClientOpt{cherrygo.WithAuthToken(apiToken), cherrygo.WithUserAgent(userAgent)}
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
		NewSSHKeyResource,
	}
}

func (p *CherryServersProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	cfg := &datasourcebase.Configurator{}
	return []func() datasource.DataSource{
		NewProjectDataSource,
		NewServerDataSource,
		NewIpDataSource,
		NewSSHKeyDataSource,
		NewRegionSingleDS(cfg),
		NewRegionListDS(cfg),
		NewPlanSingleDS(cfg),
		NewPlanListDS(cfg),
	}
}

func (p *CherryServersProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CherryServersProvider{
			version: version,
		}
	}
}
