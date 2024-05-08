// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"github.com/cherryservers/cherrygo/v3"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strconv"
)

// Ensure provider defined types fully satisfy framework interfaces.

var (
	_ resource.Resource                = &ipResource{}
	_ resource.ResourceWithConfigure   = &ipResource{}
	_ resource.ResourceWithImportState = &ipResource{}
)

func NewIpResource() resource.Resource {
	return &ipResource{}
}

// ipResource defines the resource implementation.
type ipResource struct {
	client *cherrygo.Client
}

// ipResourceModel describes the resource data model.
type ipResourceModel struct {
	Id             types.String `tfsdk:"id"`
	ProjectId      types.Int64  `tfsdk:"project_id"`
	Region         types.String `tfsdk:"region"`
	TargetId       types.String `tfsdk:"target_id"`
	TargetHostname types.String `tfsdk:"target_hostname"`
	DDOSScrubbing  types.Bool   `tfsdk:"ddos_scrubbing"`
	ARecord        types.String `tfsdk:"a_record"`
	PTRRecord      types.String `tfsdk:"ptr_record"`
	Address        types.String `tfsdk:"address"`
	AddressFamily  types.Int64  `tfsdk:"address_family"`
	CIDR           types.String `tfsdk:"cidr"`
	Gateway        types.String `tfsdk:"gateway"`
	Type           types.String `tfsdk:"type"`
	Tags           types.Map    `tfsdk:"tags"`
}

func (d *ipResourceModel) populateState(ip cherrygo.IPAddress, ctx context.Context, diags diag.Diagnostics) {
	d.Id = types.StringValue(ip.ID)
	d.ProjectId = types.Int64Value(int64(ip.Project.ID))
	d.Region = types.StringValue(ip.Region.Slug)
	d.TargetId = types.StringValue(strconv.Itoa(ip.TargetedTo.ID))
	d.TargetHostname = types.StringValue(ip.TargetedTo.Hostname)
	d.DDOSScrubbing = types.BoolValue(ip.DDoSScrubbing)
	d.ARecord = types.StringValue(ip.ARecord)
	d.PTRRecord = types.StringValue(ip.PtrRecord)
	d.Address = types.StringValue(ip.Address)
	d.AddressFamily = types.Int64Value(int64(ip.AddressFamily))
	d.CIDR = types.StringValue(ip.Cidr)
	d.Gateway = types.StringValue(ip.Gateway)
	d.Type = types.StringValue(ip.Type)

	tags, mapDiag := types.MapValueFrom(ctx, types.StringType, ip.Tags)
	d.Tags = tags

	diags.Append(mapDiag...)
}

func (r *ipResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ip"
}

func (r *ipResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Provides a CherryServers IP resource. This can be used to create, modify, and delete IP addresses",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "IP identifier",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"project_id": schema.Int64Attribute{
				Description: "CherryServers project id, associated with the IP",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"region": schema.StringAttribute{
				Description: "Slug of the region. Example: eu_nord_1 [See List Regions](https://api.cherryservers.com/doc/#tag/Regions/operation/get-regions)",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"target_id": schema.StringAttribute{
				Description: "The ID of the server to which the IP is attached\n" +
					"Conflicts with target_hostname",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("target_hostname"),
					}...),
				},
				Default: stringdefault.StaticString("0"),
			},
			"target_hostname": schema.StringAttribute{
				Description: "The hostname of the server to which the IP is attached\n" +
					"Conflicts with target_id",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(""),
			},
			"ddos_scrubbing": schema.BoolAttribute{
				Description: "If true, DDOS scrubbing protection will be applied in real-time",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"a_record": schema.StringAttribute{
				Description: "Relative DNS name for the IP address. Resulting FQDN will be '<relative-dns-name>.cloud.cherryservers.net' and must be globally unique",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"ptr_record": schema.StringAttribute{
				Optional:    true,
				Description: "Reverse DNS name for the IP address",
				Computed:    true,
				Default:     stringdefault.StaticString(""),
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
				Optional:    true,
				ElementType: types.StringType,
				Default:     mapdefault.StaticValue(types.MapValueMust(types.StringType, map[string]attr.Value{})),
				Computed:    true,
			},
		},
	}
}

func (r *ipResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
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

func (r *ipResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ipResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	request := &cherrygo.CreateIPAddress{
		Region:        data.Region.ValueString(),
		DDoSScrubbing: data.DDOSScrubbing.ValueBool(),
		PtrRecord:     data.PTRRecord.ValueString(),
		ARecord:       data.ARecord.ValueString(),
	}

	tagsMap := make(map[string]string, len(data.Tags.Elements()))
	diags := data.Tags.ElementsAs(ctx, &tagsMap, false)
	resp.Diagnostics.Append(diags...)

	request.Tags = &tagsMap

	target, err := data.getTargetId(r)
	if err != nil {
		resp.Diagnostics.AddError("invalid target server ID or hostname", err.Error())
		return
	}
	request.TargetedTo = target

	ip, _, err := r.client.IPAddresses.Create(int(data.ProjectId.ValueInt64()), request)
	if err != nil {
		resp.Diagnostics.AddError("unable to create a CherryServers IP resource", err.Error())
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
	//     return
	// }
	data.populateState(ip, ctx, resp.Diagnostics)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.SetField(ctx, "ip_id", data.Id)
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ipResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ipResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ip, ipGetResp, err := r.client.IPAddresses.Get(data.Id.ValueString(), nil)
	if err != nil {
		if is404Error(ipGetResp) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"unable to read a CherryServers ip resource",
			err.Error(),
		)
		return
	}

	data.populateState(ip, ctx, resp.Diagnostics)

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ipResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ipResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	request := cherrygo.UpdateIPAddress{
		PtrRecord: data.PTRRecord.ValueString(),
		ARecord:   data.ARecord.ValueString(),
	}

	tagsMap := make(map[string]string, len(data.Tags.Elements()))
	diags := data.Tags.ElementsAs(ctx, &tagsMap, false)
	resp.Diagnostics.Append(diags...)

	request.Tags = &tagsMap

	target, err := data.getTargetId(r)
	if err != nil {
		resp.Diagnostics.AddError("invalid target server ID or hostname", err.Error())
		return
	}
	request.TargetedTo = target

	ipID := data.Id.ValueString()
	ip, _, err := r.client.IPAddresses.Update(ipID, &request)
	if err != nil {
		resp.Diagnostics.AddError(
			"unable to update a CherryServers ip resource",
			err.Error(),
		)
		return
	}

	data.populateState(ip, ctx, resp.Diagnostics)

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }
	ctx = tflog.SetField(ctx, "ip_id", data.Id)
	tflog.Trace(ctx, "updated a resource")

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ipResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ipResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.TargetId.ValueString() != "0" {
		if _, err := r.client.IPAddresses.Unassign(data.Id.ValueString()); err != nil {
			resp.Diagnostics.AddError(
				"unable to unassign a CherryServers ip resource from target, before deleting",
				err.Error(),
			)
			return
		}
	}

	IpID := data.Id.ValueString()
	if _, err := r.client.IPAddresses.Remove(IpID); err != nil {
		resp.Diagnostics.AddError(
			"unable to delete a CherryServers ip resource",
			err.Error(),
		)
		return
	}

	ctx = tflog.SetField(ctx, "ip_id", data.Id)
	tflog.Trace(ctx, "deleted a resource")

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r *ipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (d *ipResourceModel) getTargetId(r *ipResource) (string, error) {
	if !(d.TargetId.ValueString() != "0") {
		return d.TargetId.ValueString(), nil
	} else if !d.TargetHostname.IsNull() {
		srvID, err := ServerHostnameToID(d.TargetHostname.ValueString(), int(d.ProjectId.ValueInt64()), r.client.Servers)
		return strconv.Itoa(srvID), err
	}
	return "", nil
}
