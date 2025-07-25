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
	_ datasource.DataSource              = &regionListDS{}
	_ datasource.DataSourceWithConfigure = &regionListDS{}
)

func NewRegionListDS(configurator configurator) func() datasource.DataSource {
	return func() datasource.DataSource {
		return &regionListDS{configurator: configurator}
	}
}

type regionListDS struct {
	configurator
}

type regionListModel struct {
	Regions types.List `tfsdk:"regions"`
}

func (m *regionListModel) populateState(ctx context.Context, regions []cherrygo.Region) diag.Diagnostics {
	regionModels := make([]regionModel, len(regions), cap(regions))
	var diags diag.Diagnostics

	for i, v := range regions {
		diags.Append(regionModels[i].populateState(ctx, v)...)
	}

	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: regionAttributeTypes()}, regionModels)
	if diags.HasError() {
		return diags
	}
	m.Regions = list
	return diags
}

func (d *regionListDS) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_regions"
}

func (d *regionListDS) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Provides a CherryServers regions data source. This can be used to read available region data.",

		Attributes: map[string]schema.Attribute{
			"regions": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: regionSchema(true),
				},
				Computed:    true,
				Description: "Available regions.",
			},
		},
	}
}

func (d *regionListDS) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state regionListModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	regions, _, err := d.Client().Regions.List(nil)
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
