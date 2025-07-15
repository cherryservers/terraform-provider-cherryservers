package provider

import (
	"context"

	"github.com/cherryservers/cherrygo/v3"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource              = &cycleListDS{}
	_ datasource.DataSourceWithConfigure = &cycleListDS{}
)

func NewCycleListDS(configurator configurator) func() datasource.DataSource {
	return func() datasource.DataSource {
		return &cycleListDS{configurator: configurator}
	}
}

type cycleListDS struct {
	configurator
}

type cycleModel struct {
	ID   types.Int64  `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Slug types.String `tfsdk:"slug"`
}

func cycleAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.Int64Type,
		"name": types.StringType,
		"slug": types.StringType,
	}
}

type cycleListModel struct {
	Cycles types.List `tfsdk:"cycles"`
}

func (m *cycleListModel) populateState(ctx context.Context, cycles []cherrygo.ServerCycle) diag.Diagnostics {
	models := make([]cycleModel, len(cycles), cap(cycles))
	var diags diag.Diagnostics

	for i, v := range cycles {
		models[i] = cycleModel{
			ID:   types.Int64Value(int64(v.ID)),
			Name: types.StringValue(v.Name),
			Slug: types.StringValue(v.Slug),
		}
	}

	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: cycleAttributeTypes()}, models)
	if diags.HasError() {
		return diags
	}
	m.Cycles = list
	return diags
}

func (d *cycleListDS) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cycles"
}

func (d *cycleListDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Provides a CherryServers billing cycles data source. This can be used to read available billing cycle data.",

		Attributes: map[string]schema.Attribute{
			"cycles": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"slug": schema.StringAttribute{
							Computed:    true,
							Description: "Used when provisioning resources.",
						},
					},
				},
				Computed:    true,
				Description: "Available billing cycles.",
			},
		},
	}
}

func (d *cycleListDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state cycleListModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	cycles, _, err := d.configurator.Client().Servers.ListCycles(nil)
	if err != nil {
		resp.Diagnostics.AddError("cycle list failed", err.Error())
		return
	}

	resp.Diagnostics.Append(state.populateState(ctx, cycles)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
