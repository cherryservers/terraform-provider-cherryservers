package provider

import (
	"context"
	"fmt"
	"github.com/cherryservers/cherrygo/v3"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strconv"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource                     = &serverDataSource{}
	_ datasource.DataSourceWithConfigure        = &serverDataSource{}
	_ datasource.DataSourceWithConfigValidators = &serverDataSource{}
)

func NewServerDataSource() datasource.DataSource {
	return &serverDataSource{}
}

// serverDataSource defines the data source implementation.
type serverDataSource struct {
	client *cherrygo.Client
}

type serverDataSourceModel struct {
	Plan                types.String   `tfsdk:"plan"`
	ProjectId           types.Int64    `tfsdk:"project_id"`
	Region              types.String   `tfsdk:"region"`
	Hostname            types.String   `tfsdk:"hostname"`
	Name                types.String   `tfsdk:"name"`
	Username            types.String   `tfsdk:"username"`
	Password            types.String   `tfsdk:"password"`
	BMC                 types.Object   `tfsdk:"bmc"`
	Image               types.String   `tfsdk:"image"`
	SSHKeyIds           types.Set      `tfsdk:"ssh_key_ids"`
	ExtraIPAddressesIds types.Set      `tfsdk:"extra_ip_addresses_ids"`
	Tags                types.Map      `tfsdk:"tags"`
	SpotInstance        types.Bool     `tfsdk:"spot_instance"`
	OSPartitionSize     types.Int64    `tfsdk:"os_partition_size"`
	PowerState          types.String   `tfsdk:"power_state"`
	State               types.String   `tfsdk:"state"`
	IpAddresses         types.Set      `tfsdk:"ip_addresses"`
	Id                  types.String   `tfsdk:"id"`
	Timeouts            timeouts.Value `tfsdk:"timeouts"`
}

func (d *serverDataSourceModel) populateModel(server cherrygo.Server, ctx context.Context, diags diag.Diagnostics, powerState string) {
	var resourceModel serverResourceModel
	resourceModel.populateModel(server, ctx, diags, powerState)

	d.Plan = resourceModel.Plan
	d.ProjectId = resourceModel.ProjectId
	d.Region = resourceModel.Region
	d.Hostname = resourceModel.Hostname
	d.Name = resourceModel.Name
	d.Username = resourceModel.Username
	d.Password = resourceModel.Password
	d.BMC = resourceModel.BMC
	d.Image = resourceModel.Image
	d.SSHKeyIds = resourceModel.SSHKeyIds
	//d.ExtraIPAddressesIds = resourceModel.ExtraIPAddressesIds
	d.Tags = resourceModel.Tags
	d.SpotInstance = resourceModel.SpotInstance
	d.OSPartitionSize = resourceModel.OSPartitionSize
	d.PowerState = resourceModel.PowerState
	d.State = resourceModel.State
	d.IpAddresses = resourceModel.IpAddresses
	d.Id = resourceModel.Id

}

func (d *serverDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server"
}

func (d *serverDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("hostname"), path.MatchRoot("id")),
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("project_id"), path.MatchRoot("id")),
		datasourcevalidator.RequiredTogether(path.MatchRoot("hostname"), path.MatchRoot("project_id")),
	}
}

func (d *serverDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a Cherry Servers server resource. This can be used to create, read, modify, and delete servers on your Cherry Servers account.",

		Attributes: map[string]schema.Attribute{
			"plan": schema.StringAttribute{
				Computed:    true,
				Description: "Slug of the plan. Example: e5_1620v4. [See List Plans](https://api.cherryservers.com/doc/#tag/Plans/operation/get-plans).",
			},
			"project_id": schema.Int64Attribute{
				Description: "CherryServers project id, associated with the server.",
				Computed:    true,
				Optional:    true,
			},
			"region": schema.StringAttribute{
				Description: "Slug of the region. Example: eu_nord_1 [See List Regions](https://api.cherryservers.com/doc/#tag/Regions/operation/get-regions).",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the server.",
				Computed:    true,
			},
			"hostname": schema.StringAttribute{
				Description: "Hostname of the server.",
				Computed:    true,
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "Server username credential.",
				Computed:    true,
			},
			"password": schema.StringAttribute{
				Description: "Server password credential.",
				Computed:    true,
				Sensitive:   true,
			},
			"bmc": schema.SingleNestedAttribute{
				Description: "Server BMC credentials.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"user": schema.StringAttribute{
						Computed: true,
					},
					"password": schema.StringAttribute{
						Computed:  true,
						Sensitive: true,
					},
				},
			},
			"image": schema.StringAttribute{
				Description: "Slug of the operating system. Example: ubuntu_22_04. [See List Images](https://api.cherryservers.com/doc/#tag/Images/operation/get-plan-images).",
				Computed:    true,
			},
			"ssh_key_ids": schema.SetAttribute{
				Description: "Set of the SSH key IDs allowed to SSH to the server.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"extra_ip_addresses_ids": schema.SetAttribute{
				Description: "Set of the IP address IDs to be embedded into the Server.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"tags": schema.MapAttribute{
				Description: "Key/value metadata for server tagging.",
				ElementType: types.StringType,
				Computed:    true,
			},
			"spot_instance": schema.BoolAttribute{
				Description: "If True, provisions the server as a spot instance.",
				Computed:    true,
			},
			"os_partition_size": schema.Int64Attribute{
				Description: "OS partition size in GB.",
				Computed:    true,
			},
			"power_state": schema.StringAttribute{
				Description: "The power state of the server, such as 'Powered off' or 'Powered on'.",
				Computed:    true,
			},
			"state": schema.StringAttribute{
				Description: "The state of the server, such as 'pending' or 'active'.",
				Computed:    true,
			},
			"ip_addresses": schema.SetNestedAttribute{
				Description: "IP addresses attached to the server.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "ID of the IP address.",
							Computed:    true,
						},
						"type": schema.StringAttribute{
							Description: "Type of the IP address.",
							Computed:    true,
						},
						"address": schema.StringAttribute{
							Description: "Address of the IP address.",
							Computed:    true,
						},
						"address_family": schema.Int64Attribute{
							Description: "Address family of the IP address.",
							Computed:    true,
						},
						"cidr": schema.StringAttribute{
							Description: "CIDR of the IP address.",
							Computed:    true,
						},
					},
				},
			},
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Server identifier.",
			},
			"timeouts": timeouts.Attributes(ctx),
		},
	}
}

func (d *serverDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *serverDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data serverDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var serverID int
	if data.Hostname.ValueString() != "" {
		var err error
		serverID, err = ServerHostnameToID(data.Hostname.ValueString(), int(data.ProjectId.ValueInt64()), d.client.Servers)
		if err != nil {
			resp.Diagnostics.AddError("couldn't find server ID from hostname", err.Error())
			return
		}
	} else {
		var err error
		serverID, err = strconv.Atoi(data.Id.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("invalid server ID", err.Error())
			return
		}
	}

	server, _, err := d.client.Servers.Get(serverID, nil)
	if err != nil {
		resp.Diagnostics.AddError("server not found", err.Error())
		return
	}

	powerState, _, err := d.client.Servers.PowerState(server.ID)
	if err != nil {
		resp.Diagnostics.AddError("unable to get CherryServers server power-state", err.Error())
		return
	}

	data.populateModel(server, ctx, resp.Diagnostics, powerState.Power)

	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
