package region

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/cherryservers/cherrygo/v3"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource                     = &regionDataSource{}
	_ datasource.DataSourceWithConfigure        = &regionDataSource{}
	_ datasource.DataSourceWithConfigValidators = &regionDataSource{}
)

func NewRegionDataSource() datasource.DataSource {
	return &regionDataSource{}
}

type regionDataSource struct {
	client *cherrygo.Client
}

type bgpDataSourceModel struct {
	Hosts types.List  `tfsdk:"hosts"`
	ASN   types.Int64 `tfsdk:"asn"`
}

func (b bgpDataSourceModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"hosts": types.ListType{ElemType: types.StringType},
		"asn":   types.Int64Type,
	}
}

type RegionDataSourceModel struct {
	Name       types.String `tfsdk:"name"`
	ID         types.Int64  `tfsdk:"id"`
	Slug       types.String `tfsdk:"slug"`
	RegionISO2 types.String `tfsdk:"region_iso_2"`
	BGP        types.Object `tfsdk:"bgp"`
}

func (r *RegionDataSourceModel) populateState(ctx context.Context, region cherrygo.Region) diag.Diagnostics {
	r.Name = types.StringValue(region.Name)
	r.ID = types.Int64Value(int64(region.ID))
	r.Slug = types.StringValue((region.Slug))
	r.RegionISO2 = types.StringValue((region.RegionIso2))

	hosts, diags := types.ListValueFrom(ctx, types.StringType, region.BGP.Hosts)
	if diags.HasError() {
		return diags
	}

	bgp := bgpDataSourceModel{
		Hosts: hosts,
		ASN:   types.Int64Value(int64(region.BGP.Asn)),
	}
	bgpObject, diags := types.ObjectValueFrom(ctx, bgp.AttributeTypes(), bgp)
	r.BGP = bgpObject

	return diags
}

// Region can be identified by ID or slug.
func (r *RegionDataSourceModel) Identifier() (string, error) {
	slug := r.Slug.ValueString()
	id := r.ID.ValueInt64()

	if slug == "" && id == 0 {
		return "", errors.New("unidentifiable region, no slug or ID set")
	}
	if slug != "" && id != 0 {
		return "", errors.New("unidentifiable region, both slug and ID set")
	}
	if slug != "" {
		return slug, nil
	}
	return strconv.Itoa(int(id)), nil
}

func (d *regionDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("slug"), path.MatchRoot("id")),
	}
}

func (d *regionDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_region"
}

func (d *regionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Provides a CherryServers region data source. This can be used to read available region data.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Region name.",
			},
			"id": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Region ID.",
			},
			"slug": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Region slug.",
			},
			"region_iso_2": schema.StringAttribute{
				Computed:    true,
				Description: "Region ISO 3166-1 alpha-2 code.",
			},
			"bgp": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"hosts": schema.ListAttribute{
						ElementType: types.StringType,
						Computed:    true,
						Description: "BGP host addresses.",
					},
					"asn": schema.Int64Attribute{
						Computed:    true,
						Description: "BGP ASN.",
					},
				},
				Computed:    true,
				Description: "Region BGP data.",
			},
		},
	}
}

func (d *regionDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *regionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state RegionDataSourceModel

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

	region, _, err := d.client.Regions.Get(id, nil)
	if err != nil {
		resp.Diagnostics.AddError("region read failed", err.Error())
		return
	}

	resp.Diagnostics.Append(state.populateState(ctx, region)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
