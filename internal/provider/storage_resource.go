package provider

import (
	"context"
	"fmt"

	"github.com/cherryservers/cherrygo/v3"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types satisfy framework interfaces.
var (
	_ resource.Resource                = &storageResource{}
	_ resource.ResourceWithConfigure   = &storageResource{}
	_ resource.ResourceWithImportState = &storageResource{}
)

// NewStorageResource is a helper function to simplify the provider implementation.
func NewStorageResource() resource.Resource {
	return &storageResource{}
}

// storageResource is the resource implementation.
type storageResource struct {
	client *cherrygo.Client
}

// storageResourceModel describes the resource data model.
type storageResourceModel struct {
	// Identifiers
	Id types.Int64 `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	ProjectId types.Int64 `tfsdk:"project_id"`

	// Required on creation
	Region types.String `tfsdk:"region"`
	Size types.Int64 `tfsdk:"size"`

	// Optional on creation
	Description types.String `tfsdk:"description"`
	AttachedTo types.Int64 `tfsdk:"attached_to"` // Server ID

	// Computed (iSCSI details)
	VlanId types.String `tfsdk:"vlan_id"`
	VlanIp types.String `tfsdk:"vlan_ip"`
	Initiator types.String `tfsdk:"initiator"`
	DiscoveryIp types.String `tfsdk:"discovery_ip"`

	// Computed (metadata)
	AllowEditSize types.Bool `tfsdk:"allow_edit_size"`
	Unit types.String `tfsdk:"unit"`
}

// populateState fills the model with API response data.
func (d *storageResourceModel) populateState(storage cherrygo.BlockStorage, ctx context.Context, diags diag.Diagnostics) {
	d.Id = types.Int64Value(int64(storage.ID))
	d.Name = types.StringValue(storage.Name)
	d.Size = types.Int64Value(int64(storage.Size))
	d.Description = types.StringValue(storage.Description)
	d.VlanId = types.StringValue(storage.VlanID)
	d.VlanIp = types.StringValue(storage.VlanIP)
	d.Initiator = types.StringValue(storage.Initiator)
	d.DiscoveryIp = types.StringValue(storage.DiscoveryIP)
	d.AllowEditSize = types.BoolValue(storage.AllowEditSize)
	d.Unit = types.StringValue(storage.Unit)
	d.Region = types.StringValue(storage.Region.Slug)

	if storage.AttachedTo.ID != 0 {
		d.AttachedTo = types.Int64Value(int64(storage.AttachedTo.ID))
	} else {
		d.AttachedTo = types.Int64Null()
	}
}

// Metadata returns the resource type name.
func (r *storageResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storage"
}

// Schema defines the schema for the resource.
func (r *storageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a CherryServers Storage (EBS-like) resource. This can be used to create, modify, and delete elastic block storage volumes.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Storage volume identifier.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the storage volume.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.Int64Attribute{
				Description: "CherryServers project ID associated with the storage.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"region": schema.StringAttribute{
				Description: "Slug of the region. Example: LT-Siauliai. Cannot be changed after creation.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"size": schema.Int64Attribute{
				Description: "Storage size in gigabytes. Can only be increased. Resize creates new storage ID.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Optional description for the storage volume.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"attached_to": schema.Int64Attribute{
				Description: "Server ID that this storage is attached to. Null if unattached.",
				Optional:    true,
				Computed:    true,
			},
			"vlan_id": schema.StringAttribute{
				Description: "iSCSI VLAN ID for connecting to the storage.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vlan_ip": schema.StringAttribute{
				Description: "iSCSI VLAN IP address.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"initiator": schema.StringAttribute{
				Description: "iSCSI initiator name.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"discovery_ip": schema.StringAttribute{
				Description: "iSCSI discovery IP address.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_edit_size": schema.BoolAttribute{
				Description: "Whether this storage volume can be resized.",
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
				},
			},
			"unit": schema.StringAttribute{
				Description: "Unit type for storage (typically 'GB').",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *storageResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cherrygo.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *cherrygo.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *storageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data storageResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create API request
	createReq := &cherrygo.CreateStorage{
		ProjectID:   int(data.ProjectId.ValueInt64()),
		Size:        int(data.Size.ValueInt64()),
		Region:      data.Region.ValueString(),
		Description: data.Description.ValueString(),
	}

	// Create storage
	storage, _, err := r.client.Storages.Create(createReq)
	if err != nil {
		resp.Diagnostics.AddError("unable to create storage", err.Error())
		return
	}

	tflog.Trace(ctx, "created storage", map[string]any{"storage_id": storage.ID})

	// Attach to server if specified (only if attached_to is not null)
	if !data.AttachedTo.IsNull() && !data.AttachedTo.IsUnknown() {
		attachReq := &cherrygo.AttachTo{
			StorageID: storage.ID,
			AttachTo:  int(data.AttachedTo.ValueInt64()),
		}
		_, _, err := r.client.Storages.Attach(attachReq)
		if err != nil {
			resp.Diagnostics.AddError("unable to attach storage to server", err.Error())
			return
		}
		tflog.Trace(ctx, "attached storage", map[string]any{"storage_id": storage.ID, "server_id": data.AttachedTo.ValueInt64()})
	}

	// Refresh to get all computed fields
	storage, _, err = r.client.Storages.Get(storage.ID, nil)
	if err != nil {
		resp.Diagnostics.AddError("unable to read created storage", err.Error())
		return
	}

	// Save data into Terraform state
	data.populateState(storage, ctx, resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Read refreshes the Terraform state with the latest data.
func (r *storageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data storageResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get storage from API
	storage, storageResp, err := r.client.Storages.Get(int(data.Id.ValueInt64()), nil)
	if err != nil {
		if storageResp != nil && storageResp.StatusCode == 404 {
			tflog.Warn(ctx, "storage not found, removing from state", map[string]any{"storage_id": data.Id.ValueInt64()})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("unable to read storage", err.Error())
		return
	}

	// Save updated data into Terraform state
	data.populateState(storage, ctx, resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Update updates the resource and sets the updated Terraform state.
func (r *storageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data storageResourceModel
	var stateData storageResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	storageId := int(data.Id.ValueInt64())

	// Handle description update
	if !data.Description.Equal(stateData.Description) {
		updateReq := &cherrygo.UpdateStorage{
			StorageID:   storageId,
			Description: data.Description.ValueString(),
		}
		_, _, err := r.client.Storages.Update(updateReq)
		if err != nil {
			resp.Diagnostics.AddError("unable to update storage description", err.Error())
			return
		}
	}

	// Handle attach/detach
	stateAttached := !stateData.AttachedTo.IsNull() && !stateData.AttachedTo.IsUnknown()
	dataAttached := !data.AttachedTo.IsNull() && !data.AttachedTo.IsUnknown()

	if dataAttached && !stateAttached {
		// Attach to server
		attachReq := &cherrygo.AttachTo{
			StorageID: storageId,
			AttachTo:  int(data.AttachedTo.ValueInt64()),
		}
		_, _, err := r.client.Storages.Attach(attachReq)
		if err != nil {
			resp.Diagnostics.AddError("unable to attach storage", err.Error())
			return
		}
	} else if !dataAttached && stateAttached {
		// Detach from server
		_, err := r.client.Storages.Detach(storageId)
		if err != nil {
			resp.Diagnostics.AddError("unable to detach storage", err.Error())
			return
		}
	} else if dataAttached && stateAttached && data.AttachedTo.ValueInt64() != stateData.AttachedTo.ValueInt64() {
		// Change attachment to different server
		_, err := r.client.Storages.Detach(storageId)
		if err != nil {
			resp.Diagnostics.AddError("unable to detach storage", err.Error())
			return
		}
		attachReq := &cherrygo.AttachTo{
			StorageID: storageId,
			AttachTo:  int(data.AttachedTo.ValueInt64()),
		}
		_, _, err = r.client.Storages.Attach(attachReq)
		if err != nil {
			resp.Diagnostics.AddError("unable to attach storage to new server", err.Error())
			return
		}
	}

	// Refresh state
	storage, _, err := r.client.Storages.Get(storageId, nil)
	if err != nil {
		resp.Diagnostics.AddError("unable to read updated storage", err.Error())
		return
	}

	data.populateState(storage, ctx, resp.Diagnostics)
	tflog.Trace(ctx, "updated storage", map[string]any{"storage_id": data.Id.ValueInt64()})
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *storageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data storageResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	storageId := int(data.Id.ValueInt64())

	// Detach from server if attached
	if !data.AttachedTo.IsNull() && !data.AttachedTo.IsUnknown() {
		_, err := r.client.Storages.Detach(storageId)
		if err != nil {
			resp.Diagnostics.AddError("unable to detach storage before deletion", err.Error())
			return
		}
	}

	// Delete the storage
	_, err := r.client.Storages.Delete(storageId)
	if err != nil {
		resp.Diagnostics.AddError("unable to delete storage", err.Error())
		return
	}

	tflog.Trace(ctx, "deleted storage", map[string]any{"storage_id": data.Id.ValueInt64()})
}

// ImportState imports the resource into Terraform state.
func (r *storageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Use the ID passed in as the storage_id
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
