package provider

import (
	"context"
	"fmt"
	"github.com/cherryservers/cherrygo/v3"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &projectDataSource{}
var _ datasource.DataSourceWithConfigure = &projectDataSource{}

func NewProjectDataSource() datasource.DataSource {
	return &projectDataSource{}
}

// projectDataSource defines the data source implementation.
type projectDataSource struct {
	client *cherrygo.Client
}

// projectDataSourceModel describes the data source data model.
type projectDataSourceModel struct {
	Name types.String     `tfsdk:"name"`
	Href types.String     `tfsdk:"href"`
	BGP  *projectBGPModel `tfsdk:"bgp"`
	Id   types.Int64      `tfsdk:"id"`
}

type projectBGPModel struct {
	Enabled  types.Bool  `tfsdk:"enabled"`
	LocalASN types.Int64 `tfsdk:"local_asn"`
}

func (m projectBGPModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled":   types.BoolType,
		"local_asn": types.Int64Type,
	}
}

func (d *projectDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (d *projectDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Provides a CherryServers Project data source. This can be used to read project data from CherryServers.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the project",
				Computed:    true,
			},
			"href": schema.StringAttribute{
				Description: "The hypertext reference attribute(href) of the project",
				Computed:    true,
			},
			"bgp": schema.SingleNestedAttribute{
				Description: "Project border gateway protocol(BGP) configuration.",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Computed:    true,
						Description: "BGP is enabled for the project",
					},
					"local_asn": schema.Int64Attribute{
						Computed:    true,
						Description: "The local ASN of the project",
					},
				},
				Computed: true,
			},
			"id": schema.Int64Attribute{
				Description: "Project identifier",
				Required:    true,
			},
		},
	}
}

func (d *projectDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cherrygo.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *cherrygo.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *projectDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state projectDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	projectID := state.Id.ValueInt64()

	if projectID == 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("id"),
			"Project ID must be set",
			"The provider cannot create the project data source as there is a missing or empty value for project ID. ")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	project, _, err := d.client.Projects.Get(int(projectID), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error: unable to read a CherryServers project data source",
			err.Error(),
		)
		return
	}

	state.Id = types.Int64Value(int64(project.ID))
	state.Href = types.StringValue(project.Href)
	state.Name = types.StringValue(project.Name)
	state.BGP = &projectBGPModel{
		Enabled:  types.BoolValue(project.Bgp.Enabled),
		LocalASN: types.Int64Value(int64(project.Bgp.LocalASN)),
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
