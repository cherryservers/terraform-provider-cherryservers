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
	Id              types.String `tfsdk:"id"`
	ProjectId       types.Int64  `tfsdk:"project_id"`
	Region          types.String `tfsdk:"region"`
	TargetId        types.String `tfsdk:"target_id"`
	TargetHostname  types.String `tfsdk:"target_hostname"`
	RouteIPID       types.String `tfsdk:"route_ip_id"`
	DDOSScrubbing   types.Bool   `tfsdk:"ddos_scrubbing"`
	ARecord         types.String `tfsdk:"a_record"`
	ARecordActual   types.String `tfsdk:"a_record_actual"`
	PTRRecord       types.String `tfsdk:"ptr_record"`
	PTRRecordActual types.String `tfsdk:"ptr_record_actual"`
	Address         types.String `tfsdk:"address"`
	AddressFamily   types.Int64  `tfsdk:"address_family"`
	CIDR            types.String `tfsdk:"cidr"`
	Gateway         types.String `tfsdk:"gateway"`
	Type            types.String `tfsdk:"type"`
	Tags            types.Map    `tfsdk:"tags"`
}

func (d *ipResourceModel) populateState(ip cherrygo.IPAddress, ctx context.Context, diags diag.Diagnostics) {
	d.Id = types.StringValue(ip.ID)
	d.ProjectId = types.Int64Value(int64(ip.Project.ID))
	d.Region = types.StringValue(ip.Region.Slug)
	d.TargetId = types.StringValue(strconv.Itoa(ip.TargetedTo.ID))
	d.TargetHostname = types.StringValue(ip.TargetedTo.Hostname)
	d.RouteIPID = types.StringValue(ip.RoutedTo.ID)
	d.DDOSScrubbing = types.BoolValue(ip.DDoSScrubbing)
	d.ARecordActual = types.StringValue(ip.ARecord)
	d.PTRRecordActual = types.StringValue(ip.PtrRecord)
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
				Description: "The ID of the server to which the IP is attached.\n" +
					"Conflicts with target_hostname and route_ip_id",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("target_hostname"),
						path.MatchRoot("route_ip_id"),
					}...),
				},
			},
			"target_hostname": schema.StringAttribute{
				Description: "The hostname of the server to which the IP is attached.\n" +
					"Conflicts with target_id and route_ip_id",
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("route_ip_id"),
					}...),
				},
			},
			"route_ip_id": schema.StringAttribute{
				Description: "Subnet or primary-ip type IP ID to route the created IP to.\n" +
					"Conflicts with target_hostname and target_id",
				Optional: true,
				Computed: true,
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
			},
			"a_record_actual": schema.StringAttribute{
				Description: "Relative DNS name for the IP address. Resulting FQDN will be '<relative-dns-name>.cloud.cherryservers.net' and must be globally unique.\n" +
					"API return value",
				Computed: true,
			},
			"ptr_record": schema.StringAttribute{
				Optional:    true,
				Description: "Reverse DNS name for the IP address",
			},
			"ptr_record_actual": schema.StringAttribute{
				Description: "Reverse DNS name for the IP address. API return value",
				Computed:    true,
			},
			"address": schema.StringAttribute{
				Description: "The IP address in canonical format used in the reverse DNS record",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"address_family": schema.Int64Attribute{
				Description: "IP address family IPv4 or IPv6",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"cidr": schema.StringAttribute{
				Description: "The CIDR block of the IP",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"gateway": schema.StringAttribute{
				Description: "The gateway IP address",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type": schema.StringAttribute{
				Description: "The type of IP address",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
		RoutedTo:      data.RouteIPID.ValueString(),
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
		return
	}

	ip, _, err = r.client.IPAddresses.Get(ip.ID, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"unable to read a CherryServers ip resource",
			err.Error(),
		)
		return
	}

	data.populateState(ip, ctx, resp.Diagnostics)

	// Write logs using the tflog package
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
		RoutedTo:  data.RouteIPID.ValueString(),
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
	_, _, err = r.client.IPAddresses.Update(ipID, &request)
	if err != nil {
		resp.Diagnostics.AddError(
			"unable to update a CherryServers ip resource",
			err.Error(),
		)
		return
	}

	ip, _, err := r.client.IPAddresses.Get(ipID, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"unable to read a CherryServers ip resource",
			err.Error(),
		)
		return
	}

	data.populateState(ip, ctx, resp.Diagnostics)

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

}

func (r *ipResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (d *ipResourceModel) getTargetId(r *ipResource) (string, error) {
	if d.TargetId.ValueString() != "0" && d.TargetId.ValueString() != "" {
		return d.TargetId.ValueString(), nil
	} else if d.TargetHostname.ValueString() != "" {
		srvID, err := ServerHostnameToID(d.TargetHostname.ValueString(), int(d.ProjectId.ValueInt64()), r.client.Servers)
		return strconv.Itoa(srvID), err
	}
	return "0", nil
}
