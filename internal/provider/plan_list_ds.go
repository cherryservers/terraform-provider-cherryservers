package provider

import (
	"context"

	"github.com/cherryservers/cherrygo/v3"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &planListDS{}
	_ datasource.DataSourceWithConfigure = &planListDS{}
)

func NewPlanListDS(configurator configurator) func() datasource.DataSource {
	return func() datasource.DataSource {
		return &planListDS{configurator: configurator}
	}
}

type planListDS struct {
	configurator
}

type planListModel struct {
	Plans types.List `tfsdk:"plans"`
}

func (m *planListModel) populateState(ctx context.Context, plans []cherrygo.Plan) diag.Diagnostics {
	planModels := make([]planModel, len(plans), cap(plans))
	var diags diag.Diagnostics = make(diag.Diagnostics, 0, cap(plans))

	for i, v := range plans {
		diags.Append(planModels[i].populateState(ctx, v)...)
	}

	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: planAttributeTypes()}, planModels)
	if diags.HasError() {
		return diags
	}
	m.Plans = list
	return diags
}

func (d *planListDS) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_plans"
}

func (d *planListDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Provides a CherryServers plans data source. This can be used to read available plan data.",

		Attributes: map[string]schema.Attribute{
			"plans": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: planAttr(true),
				},
				Computed:    true,
				Description: "Available server plans.",
			},
		},
	}
}

func (d *planListDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state planListModel
	// cherrygo API needs a team id to list plans,
	// but will ignore 0 and call the proper team-less endpoint.
	const null_team_id = 0

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	plans, _, err := d.configurator.Client().Plans.List(null_team_id, nil)
	if err != nil {
		resp.Diagnostics.AddError("plan list failed", err.Error())
		return
	}

	resp.Diagnostics.Append(state.populateState(ctx, plans)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
