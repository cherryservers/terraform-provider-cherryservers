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
	_ datasource.DataSource              = &StorageDataSource{}
	_ datasource.DataSourceWithConfigure = &StorageDataSource{}
)

// NewStorageDataSource is a helper function to simplify the provider implementation.
func NewStorageDataSource() datasource.DataSource {
	return &StorageDataSource{}
}

// StorageDataSource is the data source implementation.
type StorageDataSource struct {
	client *cherrygo.Client
}

// StorageDataSourceModel describes the data source data model.
type StorageDataSourceModel struct {
	Id types.Int64 `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	ProjectId types.Int64 `tfsdk:"project_id"`
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
	CreatedAt types.String `tfsdk:"created_at"`
}

// Metadata returns the data source type name.
func (d *StorageDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storage"
}

// Schema defines the schema for the data source.
func (d *StorageDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Fetch information about a specific CherryServers storage volume.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Storage volume ID to lookup.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the storage volume.",
				Computed:    true,
			},
			"project_id": schema.Int64Attribute{
				Description: "Project ID the storage belongs to.",
				Computed:    true,
			},
			"region": schema.StringAttribute{
				Description: "Region slug where the storage is located.",
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
				Description: "iSCSI VLAN IP address.",
				Computed:    true,
			},
			"initiator": schema.StringAttribute{
				Description: "iSCSI initiator name.",
				Computed:    true,
			},
			"discovery_ip": schema.StringAttribute{
				Description: "iSCSI discovery IP address.",
				Computed:    true,
			},
			"allow_edit_size": schema.BoolAttribute{
				Description: "Whether the storage can be resized.",
				Computed:    true,
			},
			"unit": schema.StringAttribute{
				Description: "Storage unit (typically 'GB').",
				Computed:    true,
			},
			"attached_to": schema.Int64Attribute{
				Description: "Server ID if storage is attached.",
				Computed:    true,
			},
			"created_at": schema.StringAttribute{
				Description: "Creation timestamp.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *StorageDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *StorageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config StorageDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get storage from API
	storage, _, err := d.client.Storages.Get(int(config.Id.ValueInt64()), nil)
	if err != nil {
		resp.Diagnostics.AddError("unable to read storage", err.Error())
		return
	}

	// Map response to model
	config.Name = types.StringValue(storage.Name)
	config.Size = types.Int64Value(int64(storage.Size))
	config.Description = types.StringValue(storage.Description)
	config.VlanId = types.StringValue(storage.VlanID)
	config.VlanIp = types.StringValue(storage.VlanIP)
	config.Initiator = types.StringValue(storage.Initiator)
	config.DiscoveryIp = types.StringValue(storage.DiscoveryIP)
	config.AllowEditSize = types.BoolValue(storage.AllowEditSize)
	config.Unit = types.StringValue(storage.Unit)
	config.Region = types.StringValue(storage.Region.Slug)
	config.CreatedAt = types.StringValue(storage.CreatedAt)

	if storage.AttachedTo.ID != 0 {
		config.AttachedTo = types.Int64Value(int64(storage.AttachedTo.ID))
	} else {
		config.AttachedTo = types.Int64Null()
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
