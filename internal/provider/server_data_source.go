package provider

import (
	"context"
	"fmt"
	"github.com/cherryservers/cherrygo/v3"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
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
	resp.Schema = serverDataSourceSchema(ctx)
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
