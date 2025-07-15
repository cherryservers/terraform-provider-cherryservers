package provider

import (
	"context"

	"github.com/cherryservers/cherrygo/v3"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource                     = &planSingleDS{}
	_ datasource.DataSourceWithConfigure        = &planSingleDS{}
	_ datasource.DataSourceWithConfigValidators = &planSingleDS{}
)

func NewPlanSingleDS(configurator configurator) func() datasource.DataSource {
	return func() datasource.DataSource {
		return &planSingleDS{configurator: configurator}
	}
}

type planSingleDS struct {
	configurator
}

func (d *planSingleDS) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("slug"), path.MatchRoot("id")),
	}
}

func (d *planSingleDS) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_plan"
}

func (d *planSingleDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Provides a CherryServers plan data source. This can be used to read available plan data.",
		Attributes:  planAttr(false),
	}
}

func (d *planSingleDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state planModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	slug := state.Slug.ValueString()
	id := state.ID.ValueInt64()

	var plan cherrygo.Plan
	var err error = nil
	if slug != "" {
		plan, _, err = d.configurator.Client().Plans.GetBySlug(slug, nil)
	} else {
		plan, _, err = d.configurator.Client().Plans.GetByID(int(id), nil)
	}

	if err != nil {
		resp.Diagnostics.AddError("plan read failed", err.Error())
		return
	}

	resp.Diagnostics.Append(state.populateState(ctx, plan)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
