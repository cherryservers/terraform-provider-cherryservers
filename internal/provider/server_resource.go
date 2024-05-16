package provider

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/cenkalti/backoff/v4"
	"github.com/cherryservers/cherrygo/v3"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"os"
	"strconv"
	"time"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &serverResource{}
	_ resource.ResourceWithConfigure   = &serverResource{}
	_ resource.ResourceWithImportState = &serverResource{}
)

func NewServerResource() resource.Resource {
	return &serverResource{}
}

// serverResource defines the resource implementation.
type serverResource struct {
	client *cherrygo.Client
}

func (r *serverResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server"
}

func (r *serverResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = serverResourceSchema(ctx)
}

func (r *serverResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *serverResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data serverResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	request := &cherrygo.CreateServer{
		ProjectID:    int(data.ProjectId.ValueInt64()),
		Plan:         data.Plan.ValueString(),
		Region:       data.Region.ValueString(),
		Image:        data.Image.ValueString(),
		Hostname:     data.Hostname.ValueString(),
		SpotInstance: data.SpotInstance.ValueBool(),
	}

	if !data.SSHKeyIds.IsUnknown() {
		sshIds := make([]string, len(data.SSHKeyIds.Elements()))
		diags := data.SSHKeyIds.ElementsAs(ctx, &sshIds, false)
		resp.Diagnostics.Append(diags...)

		request.SSHKeys = sshIds
	}

	if !data.ExtraIPAddressesIds.IsUnknown() {
		ipsIds := make([]string, len(data.ExtraIPAddressesIds.Elements()))
		diags := data.ExtraIPAddressesIds.ElementsAs(ctx, &ipsIds, false)
		resp.Diagnostics.Append(diags...)

		request.IPAddresses = ipsIds

	}

	tagsMap := make(map[string]string, len(data.Tags.Elements()))
	diags := data.Tags.ElementsAs(ctx, &tagsMap, false)
	resp.Diagnostics.Append(diags...)

	request.Tags = &tagsMap

	if !data.UserDataFile.IsNull() {
		userdataRaw, err := os.ReadFile(data.UserDataFile.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("unable to read user data file", err.Error())
			return
		}
		userData := base64.StdEncoding.EncodeToString(userdataRaw)
		request.UserData = userData
	}

	if !data.OSPartitionSize.IsNull() {
		request.OSPartitionSize = int(data.OSPartitionSize.ValueInt64())
	}

	server, _, err := r.client.Servers.Create(request)
	if err != nil {
		resp.Diagnostics.AddError(
			"unable to create a CherryServers server resource",
			err.Error(),
		)
		return
	}

	createTimeout, diags := data.Timeouts.Create(ctx, 60*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = backoff.Retry(
		func() error {
			stateOption := cherrygo.GetOptions{Fields: []string{"state"}}
			s, _, e := r.client.Servers.Get(server.ID, &stateOption)
			if e != nil {
				return backoff.Permanent(e)
			}

			if s.State == "pending" || s.State == "provisioning" {
				return errors.New("server is in inactive state")
			}

			if s.State == "active" {
				return nil
			}

			return backoff.Permanent(errors.New("failed to deploy server"))

		}, backoff.NewExponentialBackOff(
			backoff.WithMaxElapsedTime(createTimeout),
			backoff.WithInitialInterval(time.Second*10)))
	if err != nil {
		resp.Diagnostics.AddError("unable to deploy CherryServers server", err.Error())
		return
	}

	powerState, _, err := r.client.Servers.PowerState(server.ID)
	if err != nil {
		resp.Diagnostics.AddError("unable to get CherryServers server power-state", err.Error())
		return
	}

	//Workaround for not being able to set BGP and Name on "Request a server" request in API
	// TODO: add BGP
	updateRequest := cherrygo.UpdateServer{
		Name: data.Name.ValueString(),
	}

	server, _, err = r.client.Servers.Update(server.ID, &updateRequest)
	if err != nil {
		resp.Diagnostics.AddError("unable to update a CherryServers server resource with name/bgp after it's creation", err.Error())
		return
	}

	server, _, err = r.client.Servers.Get(server.ID, nil)
	if err != nil {
		resp.Diagnostics.AddError("unable to read a CherryServers server resource", err.Error())
		return
	}

	data.populateModel(server, ctx, resp.Diagnostics, powerState.Power)

	// Write logs using the tflog package
	tflog.SetField(ctx, "server_id", data.Id)
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *serverResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data serverResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	serverID, err := strconv.Atoi(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("invalid server ID in state", err.Error())
		return
	}
	server, serverGetResp, err := r.client.Servers.Get(serverID, nil)
	if err != nil {
		if is404Error(serverGetResp) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"unable to read a CherryServers server resource",
			err.Error(),
		)
		return
	}

	if server.State == "terminating" {
		resp.State.RemoveResource(ctx)
		return
	}

	powerState, _, err := r.client.Servers.PowerState(server.ID)
	if err != nil {
		resp.Diagnostics.AddError("unable to get CherryServers server power-state", err.Error())
		return
	}

	data.populateModel(server, ctx, resp.Diagnostics, powerState.Power)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *serverResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state serverResourceModel

	// Read Terraform plan and state data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	serverID, _ := strconv.Atoi(plan.Id.ValueString())

	/*requestReinstall := cherrygo.ReinstallServerFields{}
	reinstallNeeded := false
	if !plan.Image.Equal(state.Image) {
		requestReinstall.Image = plan.Image.ValueString()
		reinstallNeeded = true
	}

	if !plan.SSHKeyIds.Equal(state.SSHKeyIds) {
		sshIds := make([]string, len(plan.SSHKeyIds.Elements()))
		diags := plan.SSHKeyIds.ElementsAs(ctx, &sshIds, false)
		resp.Diagnostics.Append(diags...)

		requestReinstall.SSHKeys = sshIds
		reinstallNeeded = true
	}

	if !plan.OSPartitionSize.Equal(state.OSPartitionSize) {
		requestReinstall.OSPartitionSize = int(plan.OSPartitionSize.ValueInt64())
		reinstallNeeded = true
	}

	if !plan.UserData.Equal(state.UserData) {
		if !IsBase64(plan.UserData.ValueString()) {
			resp.Diagnostics.AddError("invalid UserData", "error reinstalling server, user_data property must be base64 encoded value")
			return
		}
		reinstallNeeded = true
	}

	if reinstallNeeded {
		_, _, err := r.client.Servers.Reinstall(serverID, &requestReinstall)
		if err != nil {
			resp.Diagnostics.AddError("unable to reinstall a CherryServers server resource", err.Error())
		}
		return
	}*/

	/*if !plan.ExtraIPAddressesIds.Equal(state.ExtraIPAddressesIds) {
		for _, ip := range plan.ExtraIPAddressesIds.Elements() {
			if !slices.Contains(state.ExtraIPAddressesIds.Elements(), ip) {
				ipRequest := cherrygo.UpdateIPAddress{
					TargetedTo: plan.Id.ValueString(),
				}
				ipTf, err := ip.ToTerraformValue(ctx)
				if err != nil {
					resp.Diagnostics.AddError("invalid IP value in plan", err.Error())
					return
				}
				if ipTf.IsKnown() {
					var ipStr string
					_ = ipTf.As(&ipStr)
					_, _, err = r.client.IPAddresses.Update(ipStr, &ipRequest)
					if err != nil {
						resp.Diagnostics.AddError("unable to update IP address in CherryServers server update operation", err.Error())
					}
				}
			}
		}
	}*/

	requestUpdate := cherrygo.UpdateServer{
		Hostname: plan.Hostname.ValueString(),
		Name:     plan.Name.ValueString(),
	}

	tagsMap := make(map[string]string, len(plan.Tags.Elements()))
	diags := plan.Tags.ElementsAs(ctx, &tagsMap, false)
	resp.Diagnostics.Append(diags...)

	requestUpdate.Tags = &tagsMap

	server, _, err := r.client.Servers.Update(serverID, &requestUpdate)
	if err != nil {
		resp.Diagnostics.AddError(
			"unable to update a CherryServers server resource",
			err.Error(),
		)
		return
	}

	server, _, err = r.client.Servers.Get(serverID, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"unable to update a CherryServers server resource",
			err.Error(),
		)
		return
	}

	powerState, _, err := r.client.Servers.PowerState(server.ID)
	if err != nil {
		resp.Diagnostics.AddError("unable to get CherryServers server power-state", err.Error())
		return
	}

	plan.populateModel(server, ctx, resp.Diagnostics, powerState.Power)

	ctx = tflog.SetField(ctx, "server_id", plan.Id)
	tflog.Trace(ctx, "updated a resource")

	// Save updated plan into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *serverResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data serverResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	serverID, _ := strconv.Atoi(data.Id.ValueString())

	if _, _, err := r.client.Servers.Delete(serverID); err != nil {
		resp.Diagnostics.AddError(
			"unable to delete a CherryServers server resource",
			err.Error(),
		)
		return
	}

	ctx = tflog.SetField(ctx, "server_id", data.Id)
	tflog.Trace(ctx, "deleted a resource")

}

func (r *serverResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
