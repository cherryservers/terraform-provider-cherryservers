package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource                     = &regionSingleDS{}
	_ datasource.DataSourceWithConfigure        = &regionSingleDS{}
	_ datasource.DataSourceWithConfigValidators = &regionSingleDS{}
)

func NewRegionSingleDS(configurator configurator) func() datasource.DataSource {
	return func() datasource.DataSource {
		return &regionSingleDS{configurator: configurator}
	}
}

type regionSingleDS struct {
	configurator
}

func (d *regionSingleDS) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("slug"), path.MatchRoot("id")),
	}
}

func (d *regionSingleDS) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_region"
}

func (d *regionSingleDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Provides a CherryServers region data source. This can be used to read available region data.",
		Attributes:  regionSchema(false),
	}
}

func (d *regionSingleDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state regionModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := state.Identifier()
	if err != nil {
		resp.Diagnostics.AddError("invalid region identifier", err.Error())
		return
	}

	region, _, err := d.Client().Regions.Get(id, nil)
	if err != nil {
		resp.Diagnostics.AddError("region read failed", err.Error())
		return
	}

	resp.Diagnostics.Append(state.populateState(ctx, region)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
