// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"github.com/cherryservers/cherrygo/v3"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ipDataSource{}
var _ datasource.DataSourceWithConfigure = &ipDataSource{}

func NewIpDataSource() datasource.DataSource {
	return &ipDataSource{}
}

// ipDataSource defines the data source implementation.
type ipDataSource struct {
	client *cherrygo.Client
}

func (d *ipDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ip"
}

func (d *ipDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Provides a CherryServers IP data source. This can be used to read IP addresses",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "IP identifier",
				Required:    true,
			},
			"project_id": schema.Int64Attribute{
				Description: "CherryServers project id, associated with the IP",
				Computed:    true,
			},
			"region": schema.StringAttribute{
				Description: "Slug of the region. Example: eu_nord_1 [See List Regions](https://api.cherryservers.com/doc/#tag/Regions/operation/get-regions)",
				Computed:    true,
			},
			"target_id": schema.StringAttribute{
				Description: "The ID of the server to which the IP is attached\n" +
					"Conflicts with target_hostname and route_ip_id",
				Computed: true,
			},
			"target_hostname": schema.StringAttribute{
				Description: "The hostname of the server to which the IP is attached\n" +
					"Conflicts with target_id and route_ip_id",
				Computed: true,
			},
			"route_ip_id": schema.StringAttribute{
				Description: "Subnet or primary-ip type IP ID to route the created IP to\n" +
					"Conflicts with target_hostname and target_id",
				Computed: true,
			},
			"ddos_scrubbing": schema.BoolAttribute{
				Description: "If true, DDOS scrubbing protection will be applied in real-time",
				Computed:    true,
			},
			"a_record": schema.StringAttribute{
				Description: "Relative DNS name for the IP address. Resulting FQDN will be '<relative-dns-name>.cloud.cherryservers.net' and must be globally unique",
				Computed:    true,
			},
			"a_record_actual": schema.StringAttribute{
				Description: "Relative DNS name for the IP address. Resulting FQDN will be '<relative-dns-name>.cloud.cherryservers.net' and must be globally unique\n" +
					"API return value",
				Computed: true,
			},
			"ptr_record": schema.StringAttribute{
				Computed:    true,
				Description: "Reverse DNS name for the IP address",
			},
			"ptr_record_actual": schema.StringAttribute{
				Description: "Reverse DNS name for the IP address, API return value",
				Computed:    true,
			},
			"address": schema.StringAttribute{
				Description: "The IP address in canonical format used in the reverse DNS record",
				Computed:    true,
			},
			"address_family": schema.Int64Attribute{
				Description: "IP address family IPv4 or IPv6",
				Computed:    true,
			},
			"cidr": schema.StringAttribute{
				Description: "The CIDR block of the IP",
				Computed:    true,
			},
			"gateway": schema.StringAttribute{
				Description: "The gateway IP address",
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of IP address",
				Computed:    true,
			},
			"tags": schema.MapAttribute{
				Description: "Key/value metadata for server tagging",
				ElementType: types.StringType,
				Computed:    true,
			},
		},
	}
}

func (d *ipDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ipDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ipResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ipID := state.Id.ValueString()

	ip, _, err := d.client.IPAddresses.Get(ipID, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error: unable to read a CherryServers IP data source",
			err.Error(),
		)
		return
	}

	state.ProjectId = types.Int64Value(int64(ip.Project.ID))
	state.populateState(ip, ctx, resp.Diagnostics)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
