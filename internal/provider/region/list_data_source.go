package region

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
	_ datasource.DataSource              = &listDataSource{}
	_ datasource.DataSourceWithConfigure = &listDataSource{}
)

func NewListDataSource(configurator configurator) func() datasource.DataSource {
	return func() datasource.DataSource {
		return &listDataSource{configurator: configurator}
	}
}

type listDataSource struct {
	configurator
}

type listModel struct {
	Regions types.List `tfsdk:"regions"`
}

func (m *listModel) populateState(ctx context.Context, regions []cherrygo.Region) diag.Diagnostics {
	regionModels := make([]model, len(regions), cap(regions))
	var diags diag.Diagnostics

	for i, v := range regions {
		diags.Append(regionModels[i].populateState(ctx, v)...)
	}

	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: attributeTypes()}, regionModels)
	if diags.HasError() {
		return diags
	}
	m.Regions = list
	return diags
}

func (d *listDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_regions"
}

func (d *listDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Provides a CherryServers regions data source. This can be used to read available region data.",

		Attributes: map[string]schema.Attribute{
			"regions": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: Schema(true),
				},
				Computed:    true,
				Description: "Available regions.",
			},
		},
	}
}

func (d *listDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state listModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	regions, _, err := d.configurator.Client().Regions.List(nil)
	if err != nil {
		resp.Diagnostics.AddError("regions list failed", err.Error())
		return
	}

	resp.Diagnostics.Append(state.populateState(ctx, regions)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
