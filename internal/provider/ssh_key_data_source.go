package provider

import (
	"context"
	"fmt"
	"github.com/cherryservers/cherrygo/v3"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strconv"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &sshKeyDataSource{}
var _ datasource.DataSourceWithConfigure = &sshKeyDataSource{}

func NewSSHKeyDataSource() datasource.DataSource {
	return &sshKeyDataSource{}
}

// sshKeyDataSource defines the data source implementation.
type sshKeyDataSource struct {
	client *cherrygo.Client
}

func (d *sshKeyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_key"
}

func (d *sshKeyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Provides a CherryServers SSH Key data source. This can be used to read SSH Keys associated with your Cherry account",

		Attributes: map[string]schema.Attribute{
			"label": schema.StringAttribute{
				Computed:    true,
				Description: "Label of the SSH key",
			},
			"public_key": schema.StringAttribute{
				Computed:    true,
				Description: "Public SSH key",
			},
			"fingerprint": schema.StringAttribute{
				Computed:    true,
				Description: "Fingerprint of the SSH public key",
			},
			"created": schema.StringAttribute{
				Computed:    true,
				Description: "Date when this Key was created",
			},
			"updated": schema.StringAttribute{
				Computed:    true,
				Description: "Date when this Key was last modified",
			},
			"id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the SSH Key",
			},
		},
	}
}

func (d *sshKeyDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *sshKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state sshKeyResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	sshKeyID, err := strconv.Atoi(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("invalid SSH Key ID in state", err.Error())
		return
	}

	sshKey, _, err := d.client.SSHKeys.Get(sshKeyID, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"unable to read a CherryServers SSH Key data source",
			err.Error(),
		)
		return
	}

	state.populateModel(sshKey)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
