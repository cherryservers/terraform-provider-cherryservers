package provider

import (
	"context"
	"github.com/cherryservers/cherrygo/v3"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strconv"
	"strings"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &sshKeyResource{}
	_ resource.ResourceWithConfigure   = &sshKeyResource{}
	_ resource.ResourceWithImportState = &sshKeyResource{}
)

func NewSSHKeyResource() resource.Resource {
	return &sshKeyResource{}
}

// sshKeyResource defines the resource implementation.
type sshKeyResource struct {
	client *cherrygo.Client
}

// sshKeyResourceModel describes the resource data model.
type sshKeyResourceModel struct {
	Label       types.String `tfsdk:"label"`
	PublicKey   types.String `tfsdk:"public_key"`
	Fingerprint types.String `tfsdk:"fingerprint"`
	Created     types.String `tfsdk:"created"`
	Updated     types.String `tfsdk:"updated"`
	ID          types.String `tfsdk:"id"`
}

func (d *sshKeyResourceModel) populateModel(sshKey cherrygo.SSHKey) {
	d.Label = types.StringValue(sshKey.Label)
	//d.PublicKey = types.StringValue(sshKey.Key)
	d.Fingerprint = types.StringValue(sshKey.Fingerprint)
	d.Created = types.StringValue(sshKey.Created)
	d.Updated = types.StringValue(sshKey.Updated)
	d.ID = types.StringValue(strconv.Itoa(sshKey.ID))

	if types.StringValue(strings.TrimSpace(d.PublicKey.ValueString())) != types.StringValue(sshKey.Key) {
		d.PublicKey = types.StringValue(sshKey.Key)
	}
}

func (r *sshKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_key"
}

func (r *sshKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Provides a CherryServers SSH Key resource. This can be used to create, and delete SSH Keys associated with your Cherry account.",

		Attributes: map[string]schema.Attribute{
			"label": schema.StringAttribute{
				Required:    true,
				Description: "Label of the SSH key.",
			},
			"public_key": schema.StringAttribute{
				Required:    true,
				Description: "Public SSH key.",
			},
			"fingerprint": schema.StringAttribute{
				Computed:    true,
				Description: "Fingerprint of the SSH public key.",
				PlanModifiers: []planmodifier.String{
					UseStateIfNoConfigurationChanges(path.Expressions{
						path.MatchRoot("public_key"),
					}...),
				},
			},
			"created": schema.StringAttribute{
				Computed:    true,
				Description: "Date when this Key was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated": schema.StringAttribute{
				Computed:    true,
				Description: "Date when this Key was last modified.",
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the SSH Key.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *sshKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	r.client = DefaultClientConfigure(req, resp)
}

func (r *sshKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data sshKeyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	request := &cherrygo.CreateSSHKey{
		Label: data.Label.ValueString(),
		Key:   strings.TrimSpace(data.PublicKey.ValueString()),
	}

	sshKey, _, err := r.client.SSHKeys.Create(request)
	if err != nil {
		resp.Diagnostics.AddError("Error creating SSH key", err.Error())
		return
	}

	sshKey, _, err = r.client.SSHKeys.Get(sshKey.ID, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error reading SSH key", err.Error())
		return
	}

	data.populateModel(sshKey)

	// Write logs using the tflog package
	tflog.SetField(ctx, "ssh_key_id", data.ID)
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sshKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data sshKeyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("invalid ssh_key ID in state", err.Error())
		return
	}

	sshKey, sshResp, err := r.client.SSHKeys.Get(id, nil)
	if err != nil {
		if is404Error(sshResp) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"unable to read a CherryServers ssh_key resource",
			err.Error(),
		)
		return
	}

	data.populateModel(sshKey)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sshKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data sshKeyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("invalid ssh_key ID in state", err.Error())
		return
	}

	label := data.Label.ValueString()
	publicKey := strings.TrimSpace(data.PublicKey.ValueString())

	request := cherrygo.UpdateSSHKey{
		Label: &label,
		Key:   &publicKey,
	}

	_, _, err = r.client.SSHKeys.Update(id, &request)
	if err != nil {
		resp.Diagnostics.AddError("error updating SSH key", err.Error())
		return
	}

	sshKey, _, err := r.client.SSHKeys.Get(id, nil)
	if err != nil {
		resp.Diagnostics.AddError("error reading SSH key", err.Error())
		return
	}

	data.populateModel(sshKey)

	ctx = tflog.SetField(ctx, "ssh_key_id", data.ID)
	tflog.Trace(ctx, "updated a resource")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *sshKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data sshKeyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("invalid ssh_key ID in state", err.Error())
		return
	}

	if _, _, err = r.client.SSHKeys.Delete(id); err != nil {
		resp.Diagnostics.AddError("error deleting SSH key", err.Error())
		return
	}

	ctx = tflog.SetField(ctx, "ssh_key_id", data.ID)
	tflog.Trace(ctx, "deleted a resource")

}

func (r *sshKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
