package provider

import (
	"context"
	"fmt"
	"github.com/cherryservers/cherrygo/v3"
	"github.com/hashicorp/terraform-plugin-framework-validators/datasourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strconv"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ datasource.DataSource                     = &sshKeyDataSource{}
	_ datasource.DataSourceWithConfigure        = &sshKeyDataSource{}
	_ datasource.DataSourceWithConfigValidators = &sshKeyDataSource{}
)

func NewSSHKeyDataSource() datasource.DataSource {
	return &sshKeyDataSource{}
}

// sshKeyDataSource defines the data source implementation.
type sshKeyDataSource struct {
	client *cherrygo.Client
}

// sshKeyDataSourceModel describes the resource data model.
type sshKeyDataSourceModel struct {
	Name        types.String `tfsdk:"name"`
	PublicKey   types.String `tfsdk:"public_key"`
	Fingerprint types.String `tfsdk:"fingerprint"`
	Created     types.String `tfsdk:"created"`
	Updated     types.String `tfsdk:"updated"`
	ID          types.String `tfsdk:"id"`
}

func (d *sshKeyDataSourceModel) populateModel(sshKey cherrygo.SSHKey) {
	d.Name = types.StringValue(sshKey.Label)
	d.PublicKey = types.StringValue(sshKey.Key)
	d.Fingerprint = types.StringValue(sshKey.Fingerprint)
	d.Created = types.StringValue(sshKey.Created)
	d.Updated = types.StringValue(sshKey.Updated)
	d.ID = types.StringValue(strconv.Itoa(sshKey.ID))
}

func (d *sshKeyDataSource) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	return []datasource.ConfigValidator{
		datasourcevalidator.ExactlyOneOf(path.MatchRoot("name"), path.MatchRoot("id")),
	}
}

func (d *sshKeyDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_key"
}

func (d *sshKeyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Provides a CherryServers SSH Key data source. This can be used to read SSH Keys associated with your Cherry account.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "Label of the SSH key.",
			},
			"public_key": schema.StringAttribute{
				Computed:    true,
				Description: "Public SSH key.",
			},
			"fingerprint": schema.StringAttribute{
				Computed:    true,
				Description: "Fingerprint of the SSH public key.",
			},
			"created": schema.StringAttribute{
				Computed:    true,
				Description: "Date when this Key was created.",
			},
			"updated": schema.StringAttribute{
				Computed:    true,
				Description: "Date when this Key was last modified.",
			},
			"id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "ID of the SSH Key.",
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
	var state sshKeyDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var sshKeyID int
	//Get SSH key ID by label and project ID.
	if !state.Name.IsNull() {
		getOptions := cherrygo.GetOptions{}
		getOptions.Fields = []string{"ssh_key", "email"}

		sshKeys, _, err := d.client.SSHKeys.List(&getOptions)
		if err != nil {
			resp.Diagnostics.AddError("couldn't read project SSHKeys", err.Error())
			return
		}
		for _, sshKey := range sshKeys {
			if sshKey.Label == state.Name.ValueString() {
				if sshKeyID != 0 {
					resp.Diagnostics.AddError("multiple SSH keys with the same name", "multiple SSH keys with the same name")
					return
				}
				sshKeyID = sshKey.ID
			}
		}
		//Get SSH key ID straight from schema.
	} else {
		var err error
		sshKeyID, err = strconv.Atoi(state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("invalid SSH Key ID in state", err.Error())
			return
		}
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
