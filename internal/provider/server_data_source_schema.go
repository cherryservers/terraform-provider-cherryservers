package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func serverDataSourceSchema(ctx context.Context) schema.Schema {
	return schema.Schema{
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
			"user_data_file": schema.StringAttribute{
				Description: "Base64 encoded User-Data blob. It should be either a bash or cloud-config script.",
				Computed:    true,
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
