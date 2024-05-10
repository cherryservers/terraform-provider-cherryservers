// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"github.com/cherryservers/cherrygo/v3"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strconv"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &projectResource{}
	_ resource.ResourceWithConfigure   = &projectResource{}
	_ resource.ResourceWithImportState = &projectResource{}
)

func NewProjectResource() resource.Resource {
	return &projectResource{}
}

// projectResource defines the resource implementation.
type projectResource struct {
	client *cherrygo.Client
}

// projectResourceModel describes the resource data model.
type projectResourceModel struct {
	Name   types.String `tfsdk:"name"`
	TeamId types.Int64  `tfsdk:"team_id"`
	Href   types.String `tfsdk:"href"`
	BGP    types.Object `tfsdk:"bgp"`
	Id     types.String `tfsdk:"id"`
}

func (d *projectResourceModel) populateState(project cherrygo.Project, ctx context.Context, diags diag.Diagnostics) {
	d.Id = types.StringValue(strconv.Itoa(project.ID))

	bgp := projectBGPModel{
		Enabled:  types.BoolValue(project.Bgp.Enabled),
		LocalASN: types.Int64Value(int64(project.Bgp.LocalASN)),
	}

	bgpTf, bgpDiags := types.ObjectValueFrom(ctx, bgp.AttributeTypes(), bgp)

	d.BGP = bgpTf
	diags.Append(bgpDiags...)

	d.Href = types.StringValue(project.Href)
	d.Name = types.StringValue(project.Name)

}

func (r *projectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_project"
}

func (r *projectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Provides a CherryServers project resource. This can be used to create, modify, and delete projects",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the project",
				Required:    true,
			},
			"team_id": schema.Int64Attribute{
				Description: "The ID of the team that owns the project",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"href": schema.StringAttribute{
				Description: "The hypertext reference attribute(href) of the project",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bgp": schema.SingleNestedAttribute{
				Description: "Project border gateway protocol(BGP) configuration.",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Required:    true,
						Description: "BGP is enabled for the project",
					},
					"local_asn": schema.Int64Attribute{
						Computed:    true,
						Description: "The local ASN of the project",
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
				},
				Computed: true,
				Optional: true,
				Default: objectdefault.StaticValue(
					types.ObjectValueMust(
						map[string]attr.Type{
							"enabled":   types.BoolType,
							"local_asn": types.Int64Type,
						},
						map[string]attr.Value{
							"enabled":   types.BoolValue(false),
							"local_asn": types.Int64Unknown(),
						})),
			},
			"id": schema.StringAttribute{
				Description: "Project identifier",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *projectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *projectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data projectResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var bgp projectBGPModel
	bgpDiags := data.BGP.As(ctx, &bgp, basetypes.ObjectAsOptions{})
	resp.Diagnostics.Append(bgpDiags...)

	teamId := data.TeamId.ValueInt64()
	request := &cherrygo.CreateProject{
		Name: data.Name.ValueString(),
		Bgp:  bgp.Enabled.ValueBool(),
	}

	project, _, err := r.client.Projects.Create(int(teamId), request)
	if err != nil {
		resp.Diagnostics.AddError(
			"unable to create a CherryServers project resource",
			err.Error(),
		)
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
	//     return
	// }

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.populateState(project, ctx, resp.Diagnostics)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	ctx = tflog.SetField(ctx, "project_id", project.ID)
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *projectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data projectResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	projectId, _ := strconv.Atoi(data.Id.ValueString())
	project, projectGetResp, err := r.client.Projects.Get(projectId, nil)
	if err != nil {
		if is404Error(projectGetResp) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"unable to read a CherryServers project resource",
			err.Error(),
		)
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read example, got error: %s", err))
	//     return
	// }

	data.populateState(project, ctx, resp.Diagnostics)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *projectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data projectResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var bgp projectBGPModel
	bgpDiags := data.BGP.As(ctx, &bgp, basetypes.ObjectAsOptions{})
	resp.Diagnostics.Append(bgpDiags...)

	name := data.Name.ValueString()
	bgpEnabled := bgp.Enabled.ValueBool()
	request := &cherrygo.UpdateProject{
		Name: &name,
		Bgp:  &bgpEnabled,
	}

	projectID, _ := strconv.Atoi(data.Id.ValueString())
	project, _, err := r.client.Projects.Update(projectID, request)
	if err != nil {
		resp.Diagnostics.AddError(
			"unable to update a CherryServers project resource",
			err.Error(),
		)
		return
	}

	data.populateState(project, ctx, resp.Diagnostics)

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	ctx = tflog.SetField(ctx, "project_id", project.ID)
	tflog.Trace(ctx, "updated a resource")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *projectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data projectResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	projectId, _ := strconv.Atoi(data.Id.ValueString())
	if _, err := r.client.Projects.Delete(projectId); err != nil {
		resp.Diagnostics.AddError(
			"unable to delete a CherryServers project resource",
			err.Error(),
		)
		return
	}

	ctx = tflog.SetField(ctx, "project_id", projectId)
	tflog.Trace(ctx, "deleted a resource")

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete example, got error: %s", err))
	//     return
	// }
}

func (r *projectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
