package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/cherryservers/cherrygo/v4"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

const (
	apiKeyVar       = "CHERRY_API_KEY"
	oldApiKeyVar    = "CHERRY_AUTH_KEY"
	oldestApiKeyVar = "CHERRY_AUTH_TOKEN"
)

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
					"and will be removed in the next major version of the provider. " +
					"Cherry Servers [API Key](https://portal.cherryservers.com/settings/api-keys) " +
					"that allows interactions with the API.",
				Optional:  true,
				Sensitive: true,
				DeprecationMessage: "Use `api_key` instead, as this attribute is deprecated " +
					"and will be removed in the next major version of the provider.",
			},
			"api_key": schema.StringAttribute{
				Description: "Cherry Servers [API Key](https://portal.cherryservers.com/settings/api-keys) " +
					fmt.Sprintf("that allows interactions with the API. Can also be set with the %s ", apiKeyVar) +
					"environment variable.",
				Optional:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRoot("api_token")),
				},
			},
		},
	}
}

func apiKey(diags *diag.Diagnostics, cfg CherryServersProviderModel) string {
	if cfg.APIToken.IsUnknown() {
		diags.AddAttributeError(
			path.Root("api_token"),
			"Unknown CherryServers API token",
			"The provider cannot create the CherryServers API client as there "+
				"is an unknown configuration value for the CherryServers API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration,"+
				fmt.Sprintf(" or use the %s environment variable.", apiKeyVar),
		)
	}

	if cfg.APIKey.IsUnknown() {
		diags.AddAttributeError(
			path.Root("api_key"),
			"Unknown CherryServers API key",
			"The provider cannot create the CherryServers API client as there "+
				"is an unknown configuration value for the CherryServers API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration,"+
				fmt.Sprintf(" or use the %s environment variable.", apiKeyVar),
		)
	}

	if diags.HasError() {
		return ""
	}

	var source, key string

	// CHERRY_AUTH_TOKEN and CHERRY_AUTH_KEY are deprecated,
	// so CHERRY_API_KEY beats them.
	for _, envVar := range []string{oldestApiKeyVar, oldApiKeyVar, apiKeyVar} {
		if k := os.Getenv(envVar); k != "" {
			key = k
			source = envVar
		}
	}

	if !cfg.APIToken.IsNull() {
		key = cfg.APIToken.ValueString()
		source = "config"
	}
	if !cfg.APIKey.IsNull() {
		key = cfg.APIKey.ValueString()
		source = "config"
	}

	// Add a warning if deprecated environment variables are used.
	if source == oldApiKeyVar || source == oldestApiKeyVar {
		diags.AddWarning(fmt.Sprintf(
			"%s is deprecated", source,
		),
			fmt.Sprintf("%s is deprecated and will be removed in the next major ", source)+
				fmt.Sprintf("version of the provider, please use %s instead.", apiKeyVar))
	}

	if key == "" {
		diags.AddAttributeError(
			path.Root("api_key"),
			"Missing CherryServers API key",
			"The provider cannot create the CherryServers API client "+
				"as there is a missing or empty value for the CherryServers API key. "+
				"Set the API key value in the configuration or use the "+
				fmt.Sprintf("%s environment variable.", apiKeyVar),
		)
	}

	return key
}

func (p *CherryServersProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring CherryServers client")

	var cfg CherryServersProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &cfg)...)

	if resp.Diagnostics.HasError() {
		return
	}

	key := apiKey(&resp.Diagnostics, cfg)

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "cherryservers_api_key", key)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "cherryservers_api_key")

	tflog.Debug(ctx, "Creating CherryServers client")

	// Example client configuration for data sources and resources
	userAgent := fmt.Sprintf("terraform-provider/cherryservers/%s terraform/%s", p.version, req.TerraformVersion)
	args := []cherrygo.ClientOpt{cherrygo.WithAPIKey(key), cherrygo.WithUserAgent(userAgent)}
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
