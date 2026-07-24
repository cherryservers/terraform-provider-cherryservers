package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/cherryservers/cherrygo/v3"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-cherryservers/internal/provider/datasourcebase"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const apiKeyVar = "CHERRY_API_KEY"

// Ensure CherryServersProvider satisfies various provider interfaces.
var (
	_ provider.Provider              = &CherryServersProvider{}
	_ provider.ProviderWithFunctions = &CherryServersProvider{}
)

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
	APIKey   types.String `tfsdk:"api_key"`
}

func (p *CherryServersProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "cherryservers"
	resp.Version = p.version
}

func (p *CherryServersProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_token": schema.StringAttribute{
				Description: "**Deprecated**: use `api_key` instead, as this attribute is deprecated " +
					"and will removed in the next major version of the provider. " +
					"Cherry Servers [API Key](https://portal.cherryservers.com/settings/api-keys)" +
					"that allows interactions with the API.",
				Optional:  true,
				Sensitive: true,
				DeprecationMessage: "Use `api_key` instead, as this attribute is deprecated " +
					"and will removed in the next major version of the provider.",
			},
			"api_key": schema.StringAttribute{
				Description: "Cherry Servers [API Key](https://portal.cherryservers.com/settings/api-keys)" +
					fmt.Sprintf("that allows interactions with the API. Can also be set with the %s ", apiKeyVar) +
					"environment variable.",
				Optional:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{path.MatchRoot("api_token")}...),
				},
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

	if data.APIToken.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_token"),
			"Unknown CherryServers API Token",
			"The provider cannot create the CherryServers API client as there "+
				"is an unknown configuration value for the CherryServers API Token. "+
				"Either target apply the source of the value first, set the value statically in the configuration,"+
				" or use the CHERRY_AUTH_TOKEN or CHERRY_AUTH_KEY environment variables.",
		)
	}

	if data.APIKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown CherryServers API key",
			"The provider cannot create the CherryServers API client as there "+
				"is an unknown configuration value for the CherryServers API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration,"+
				fmt.Sprintf(" or use the %s environment variable.", apiKeyVar),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// CHERRY_AUTH_TOKEN and CHERRY_AUTH_KEY are deprecated,
	// so CHERRY_API_KEY beats them.
	apiKey := os.Getenv("CHERRY_AUTH_KEY")
	source := "CHERRY_AUTH_KEY"
	if apiKey == "" {
		apiKey = os.Getenv("CHERRY_AUTH_TOKEN")
		source = "CHERRY_AUTH_TOKEN"
	}
	if k, ok := os.LookupEnv(apiKeyVar); ok {
		apiKey = k
		source = apiKeyVar
	}

	if !data.APIToken.IsNull() {
		apiKey = data.APIToken.ValueString()
		source = "config"
	}
	if !data.APIKey.IsNull() {
		apiKey = data.APIKey.ValueString()
		source = "config"
	}

	if apiKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing CherryServers API key",
			"The provider cannot create the CherryServers API client "+
				"as there is a missing or empty value for the CherryServers API key. "+
				"Set the API key value in the configuration or use the "+
				fmt.Sprintf("%s environment variables. ", apiKeyVar),
		)
	}

	// Add a warning if deprecated environment variables are used.
	if source == "CHERRY_AUTH_KEY" || source == "CHERRY_AUTH_TOKEN" {
		resp.Diagnostics.AddWarning(fmt.Sprintf(
			"%s is deprecated", source),
			fmt.Sprintf("%s is deprecated and will be removed in the next major ", source)+
				fmt.Sprintf("version of the provider, please use %s instead.", apiKeyVar))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "cherryservers_api_token", apiKey)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "cherryservers_api_token")

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
		NewCycleListDS(cfg),
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
