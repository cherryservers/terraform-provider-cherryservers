package provider

import (
	"context"
	"fmt"

	"github.com/cherryservers/cherrygo/v3"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &StorageListDataSource{}
	_ datasource.DataSourceWithConfigure = &StorageListDataSource{}
)

// NewStorageListDataSource is a helper function to simplify the provider implementation.
func NewStorageListDataSource() datasource.DataSource {
	return &StorageListDataSource{}
}

// StorageListDataSource is the data source implementation.
type StorageListDataSource struct {
	client *cherrygo.Client
}

// StorageListDataSourceModel describes the data source data model.
type StorageListDataSourceModel struct {
	ProjectId types.Int64 `tfsdk:"project_id"`
	Storages []StorageListItem `tfsdk:"storages"`
}

type StorageListItem struct {
	Id types.Int64 `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Region types.String `tfsdk:"region"`
	Size types.Int64 `tfsdk:"size"`
	Description types.String `tfsdk:"description"`
	VlanId types.String `tfsdk:"vlan_id"`
	VlanIp types.String `tfsdk:"vlan_ip"`
	Initiator types.String `tfsdk:"initiator"`
	DiscoveryIp types.String `tfsdk:"discovery_ip"`
	AllowEditSize types.Bool `tfsdk:"allow_edit_size"`
	Unit types.String `tfsdk:"unit"`
	AttachedTo types.Int64 `tfsdk:"attached_to"`
}

// Metadata returns the data source type name.
func (d *StorageListDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storages"
}

// Schema defines the schema for the data source.
func (d *StorageListDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List all CherryServers storage volumes in a project.",
		Attributes: map[string]schema.Attribute{
			"project_id": schema.Int64Attribute{
				Description: "Project ID to list storage volumes for.",
				Required:    true,
			},
			"storages": schema.ListNestedAttribute{
				Description: "List of storage volumes.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Description: "Storage volume ID.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Storage volume name.",
							Computed:    true,
						},
						"region": schema.StringAttribute{
							Description: "Region slug.",
							Computed:    true,
						},
						"size": schema.Int64Attribute{
							Description: "Storage size in gigabytes.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Storage description.",
							Computed:    true,
						},
						"vlan_id": schema.StringAttribute{
							Description: "iSCSI VLAN ID.",
							Computed:    true,
						},
						"vlan_ip": schema.StringAttribute{
							Description: "iSCSI VLAN IP.",
							Computed:    true,
						},
						"initiator": schema.StringAttribute{
							Description: "iSCSI initiator name.",
							Computed:    true,
						},
						"discovery_ip": schema.StringAttribute{
							Description: "iSCSI discovery IP.",
							Computed:    true,
						},
						"allow_edit_size": schema.BoolAttribute{
							Description: "Whether storage can be resized.",
							Computed:    true,
						},
						"unit": schema.StringAttribute{
							Description: "Storage unit.",
							Computed:    true,
						},
						"attached_to": schema.Int64Attribute{
							Description: "Server ID if attached.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *StorageListDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// Read refreshes the Terraform state with the latest data.
func (d *StorageListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config StorageListDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// List storages from API
	storages, _, err := d.client.Storages.List(int(config.ProjectId.ValueInt64()), nil)
	if err != nil {
		resp.Diagnostics.AddError("unable to list storages", err.Error())
		return
	}

	// Map response to model
	var items []StorageListItem
	for _, storage := range storages {
		item := StorageListItem{
			Id:            types.Int64Value(int64(storage.ID)),
			Name:          types.StringValue(storage.Name),
			Region:        types.StringValue(storage.Region.Slug),
			Size:          types.Int64Value(int64(storage.Size)),
			Description:  types.StringValue(storage.Description),
			VlanId:        types.StringValue(storage.VlanID),
			VlanIp:        types.StringValue(storage.VlanIP),
			Initiator:     types.StringValue(storage.Initiator),
			DiscoveryIp:   types.StringValue(storage.DiscoveryIP),
			AllowEditSize: types.BoolValue(storage.AllowEditSize),
			Unit:          types.StringValue(storage.Unit),
		}

		if storage.AttachedTo.ID != 0 {
			item.AttachedTo = types.Int64Value(int64(storage.AttachedTo.ID))
		} else {
			item.AttachedTo = types.Int64Null()
		}

		items = append(items, item)
	}

	config.Storages = items

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
