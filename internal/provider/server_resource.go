package provider

import (
	"context"
	"errors"
	"github.com/cenkalti/backoff/v4"
	"github.com/cherryservers/cherrygo/v3"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

// serverResourceModel describes the resource data model.
type serverResourceModel struct {
	Plan                types.String   `tfsdk:"plan"`
	ProjectId           types.Int64    `tfsdk:"project_id"`
	Region              types.String   `tfsdk:"region"`
	Hostname            types.String   `tfsdk:"hostname"`
	Name                types.String   `tfsdk:"name"`
	Username            types.String   `tfsdk:"username"`
	Password            types.String   `tfsdk:"password"`
	BMC                 types.Object   `tfsdk:"bmc"`
	Image               types.String   `tfsdk:"image"`
	SSHKeyIds           types.Set      `tfsdk:"ssh_key_ids"`
	ExtraIPAddressesIds types.Set      `tfsdk:"extra_ip_addresses_ids"`
	UserData            types.String   `tfsdk:"user_data"`
	Tags                types.Map      `tfsdk:"tags"`
	SpotInstance        types.Bool     `tfsdk:"spot_instance"`
	OSPartitionSize     types.Int64    `tfsdk:"os_partition_size"`
	PowerState          types.String   `tfsdk:"power_state"`
	State               types.String   `tfsdk:"state"`
	IpAddresses         types.Set      `tfsdk:"ip_addresses"`
	Id                  types.String   `tfsdk:"id"`
	Timeouts            timeouts.Value `tfsdk:"timeouts"`
	AllowReinstall      types.Bool     `tfsdk:"allow_reinstall"`
}

func (d *serverResourceModel) populateModel(server cherrygo.Server, ctx context.Context, diags diag.Diagnostics, powerState string) {
	d.Plan = types.StringValue(server.Plan.Slug)
	d.ProjectId = types.Int64Value(int64(server.Project.ID))
	d.Region = types.StringValue(server.Region.Slug)
	d.Hostname = types.StringValue(server.Hostname)
	d.Name = types.StringValue(server.Name)
	d.Username = types.StringValue(server.Username)
	d.Password = types.StringValue(server.Password)

	bmcModel := bmcResourceModel{
		User:     types.StringValue(server.BMC.User),
		Password: types.StringValue(server.BMC.Password),
	}
	bmcTf, bmcDiags := types.ObjectValueFrom(ctx, bmcModel.AttributeTypes(), bmcModel)
	diags.Append(bmcDiags...)

	d.BMC = bmcTf

	d.Image = types.StringValue(server.Image)

	sshKeyIds := make([]string, 0, len(server.SSHKeys))
	for _, sshKey := range server.SSHKeys {
		sshKeyID := strconv.Itoa(sshKey.ID)
		sshKeyIds = append(sshKeyIds, sshKeyID)
	}
	sshKeyIdsTf, sshDiags := types.SetValueFrom(ctx, types.StringType, sshKeyIds)
	d.SSHKeyIds = sshKeyIdsTf
	diags.Append(sshDiags...)

	ips := make([]attr.Value, 0, len(server.IPAddresses))

	for _, ip := range server.IPAddresses {
		ipModel := ipAddressFlatResourceModel{
			Id:            types.StringValue(ip.ID),
			Type:          types.StringValue(ip.Type),
			Address:       types.StringValue(ip.Address),
			AddressFamily: types.Int64Value(int64(ip.AddressFamily)),
			CIDR:          types.StringValue(ip.Cidr),
		}

		ipTf, ipDiags := types.ObjectValueFrom(ctx, ipModel.AttributeTypes(), ipModel)
		diags.Append(ipDiags...)

		ips = append(ips, ipTf)
	}

	ipsTf, ipsDiags := types.SetValue(types.ObjectType{AttrTypes: ipAddressFlatResourceModel{}.AttributeTypes()}, ips)
	diags.Append(ipsDiags...)
	d.IpAddresses = ipsTf

	tags, tagsDiags := types.MapValueFrom(ctx, types.StringType, server.Tags)
	d.Tags = tags
	diags.Append(tagsDiags...)

	d.SpotInstance = types.BoolValue(server.SpotInstance)
	d.PowerState = types.StringValue(powerState)
	d.State = types.StringValue(server.State)
	d.Id = types.StringValue(strconv.Itoa(server.ID))
}

type ipAddressFlatResourceModel struct {
	Id            types.String `tfsdk:"id"`
	Type          types.String `tfsdk:"type"`
	Address       types.String `tfsdk:"address"`
	AddressFamily types.Int64  `tfsdk:"address_family"`
	CIDR          types.String `tfsdk:"cidr"`
}

func (m ipAddressFlatResourceModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":             types.StringType,
		"type":           types.StringType,
		"address":        types.StringType,
		"address_family": types.Int64Type,
		"cidr":           types.StringType,
	}
}

type bmcResourceModel struct {
	User     types.String `tfsdk:"user"`
	Password types.String `tfsdk:"password"`
}

func (m bmcResourceModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"user":     types.StringType,
		"password": types.StringType,
	}
}

func (r *serverResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server"
}

func (r *serverResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		Description: "Provides a Cherry Servers server resource. This can be used to create, read, modify, and delete servers on your Cherry Servers account.",

		Attributes: map[string]schema.Attribute{
			"plan": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "Slug of the plan. Example: e5_1620v4. [See List Plans](https://api.cherryservers.com/doc/#tag/Plans/operation/get-plans).",
			},
			"project_id": schema.Int64Attribute{
				Description: "CherryServers project id, associated with the server.",
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"region": schema.StringAttribute{
				Description: "Slug of the region. Example: eu_nord_1 [See List Regions](https://api.cherryservers.com/doc/#tag/Regions/operation/get-regions).",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"hostname": schema.StringAttribute{
				Description: "Hostname of the server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"username": schema.StringAttribute{
				Description: "Server username credential.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"password": schema.StringAttribute{
				Description: "Server password credential.",
				Computed:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bmc": schema.SingleNestedAttribute{
				Description: "Server BMC credentials.",
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"user": schema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"password": schema.StringAttribute{
						Computed:  true,
						Sensitive: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"image": schema.StringAttribute{
				Description: "Slug of the server operating system. " +
					"Updating this attribute requires a server re-install.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					WarnIfChangedString("Server re-install required.",
						"You are updating attributes that require a server re-install."+
							" This will wipe all of your data and may take awhile."),
				},
			},
			"ssh_key_ids": schema.SetAttribute{
				Description: "Set of the SSH key IDs allowed to SSH to the server. " +
					"Updating this attribute requires a server re-install.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
					WarnIfChangedSet("Server re-install required.",
						"You are updating attributes that require a server re-install."+
							" This will wipe all of your data and may take awhile."),
				},
			},
			"extra_ip_addresses_ids": schema.SetAttribute{
				Description: "Set of the IP address IDs to be embedded into the server.",
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"user_data": schema.StringAttribute{
				Description: "Base64 encoded user-data blob. It should be a bash or cloud-config script. " +
					"Updating this attribute requires a server re-install.",
				Optional: true,
				PlanModifiers: []planmodifier.String{
					WarnIfChangedString("Server re-install required.",
						"You are updating attributes that require a server re-install."+
							" This will wipe all of your data and may take awhile."),
				},
			},
			"tags": schema.MapAttribute{
				Description: "Key/value metadata for server tagging.",
				Optional:    true,
				ElementType: types.StringType,
				Default:     mapdefault.StaticValue(types.MapValueMust(types.StringType, map[string]attr.Value{})),
				Computed:    true,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"spot_instance": schema.BoolAttribute{
				Description: "If True, provisions the server as a spot instance.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"os_partition_size": schema.Int64Attribute{
				Description: "OS partition size in GB. " +
					"Updating this attribute requires a server re-install.",
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					WarnIfChangedInt64("Server re-install required.",
						"You are updating attributes that require a server re-install."+
							" This will wipe all of your data and may take awhile."),
				},
			},
			"power_state": schema.StringAttribute{
				Description: "The power state of the server, such as 'Powered off' or 'Powered on'.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": schema.StringAttribute{
				Description: "The state of the server, such as 'pending' or 'active'.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ip_addresses": schema.SetNestedAttribute{
				Description: "IP addresses attached to the server.",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "ID of the IP address.",
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"type": schema.StringAttribute{
							Description: "Type of the IP address.",
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"address": schema.StringAttribute{
							Description: "Address of the IP address.",
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"address_family": schema.Int64Attribute{
							Description: "Address family of the IP address.",
							Computed:    true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"cidr": schema.StringAttribute{
							Description: "CIDR of the IP address.",
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Server identifier.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_reinstall": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Description: "Allow server re-installation when updating `image`, `ssh_key_ids`, `os_partition_size` or `user_data`. " +
					"WARNING: The reinstall will be triggered even if Terraform reports an in-place update.",
				Default: booldefault.StaticBool(false),
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
			}),
		},
	}
}

func (r *serverResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	r.client = DefaultClientConfigure(req, resp)
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
		sshIds := make([]string, 0, len(data.SSHKeyIds.Elements()))
		diags := data.SSHKeyIds.ElementsAs(ctx, &sshIds, false)
		resp.Diagnostics.Append(diags...)

		request.SSHKeys = sshIds
	}

	if !data.ExtraIPAddressesIds.IsUnknown() {
		ipsIds := make([]string, 0, len(data.ExtraIPAddressesIds.Elements()))
		diags := data.ExtraIPAddressesIds.ElementsAs(ctx, &ipsIds, false)
		resp.Diagnostics.Append(diags...)

		request.IPAddresses = ipsIds

	}

	tagsMap := make(map[string]string, len(data.Tags.Elements()))
	diags := data.Tags.ElementsAs(ctx, &tagsMap, false)
	resp.Diagnostics.Append(diags...)

	request.Tags = &tagsMap

	if !data.UserData.IsNull() {
		userData := data.UserData.ValueString()
		if err := isBase64(userData); err == nil {
			request.UserData = userData
		} else {
			resp.Diagnostics.AddError("unable to read user data", err.Error())
			return
		}
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

	if err = normalizeServerImage(&server, r.client); err != nil {
		resp.Diagnostics.AddError("Unable to normalize CherryServers server image", err.Error())
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

	if err = normalizeServerImage(&server, r.client); err != nil {
		resp.Diagnostics.AddError("Unable to normalize CherryServers server image", err.Error())
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

	if !plan.Image.Equal(state.Image) ||
		!plan.SSHKeyIds.Equal(state.SSHKeyIds) ||
		!plan.OSPartitionSize.Equal(state.OSPartitionSize) ||
		!plan.UserData.Equal(state.UserData) {
		if !plan.AllowReinstall.ValueBool() {
			resp.Diagnostics.AddError("allow_reinstall attribute not set",
				"updating image, ssh_key_ids, os_partition_size or user_data, requires setting allow_reinstall to true")
			return
		}

		r.reinstall(ctx, plan, resp)

	}

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

	if err = normalizeServerImage(&server, r.client); err != nil {
		resp.Diagnostics.AddError("Unable to normalize CherryServers server image", err.Error())
	}

	plan.populateModel(server, ctx, resp.Diagnostics, powerState.Power)

	ctx = tflog.SetField(ctx, "server_id", plan.Id)
	tflog.Trace(ctx, "updated a resource")

	// Save updated plan into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *serverResource) reinstall(ctx context.Context, plan serverResourceModel, resp *resource.UpdateResponse) {
	password := generatePassword()
	serverID, _ := strconv.Atoi(plan.Id.ValueString())

	sshIds := make([]string, 0, len(plan.SSHKeyIds.Elements()))
	diags := plan.SSHKeyIds.ElementsAs(ctx, &sshIds, false)
	resp.Diagnostics.Append(diags...)

	requestReinstall := &cherrygo.ReinstallServerFields{
		Image:           plan.Image.ValueString(),
		Hostname:        plan.Hostname.ValueString(),
		Password:        password,
		SSHKeys:         sshIds,
		OSPartitionSize: int(plan.OSPartitionSize.ValueInt64()),
	}

	if !plan.UserData.IsNull() {
		userData := plan.UserData.ValueString()
		if err := isBase64(userData); err == nil {
			requestReinstall.UserData = userData
		} else {
			resp.Diagnostics.AddError("unable to read user data", err.Error())
			return
		}
	}

	server, _, err := r.client.Servers.Reinstall(serverID, requestReinstall)
	if err != nil {
		resp.Diagnostics.AddError(
			"unable to create a CherryServers server resource",
			err.Error(),
		)
		return
	}

	updateTimeout, diags := plan.Timeouts.Update(ctx, 60*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = backoff.Retry(
		func() error {
			statusOption := cherrygo.GetOptions{Fields: []string{"status"}}
			s, _, e := r.client.Servers.Get(server.ID, &statusOption)
			if e != nil {
				return backoff.Permanent(e)
			}

			if s.Status == "deploying" {
				return errors.New("server is in inactive state")
			}

			if s.Status == "deployed" {
				return nil
			}

			return backoff.Permanent(errors.New("server is in unknown status"))

		}, backoff.NewExponentialBackOff(
			backoff.WithMaxElapsedTime(updateTimeout),
			backoff.WithInitialInterval(time.Second*10)))
	if err != nil {
		resp.Diagnostics.AddError("unable to reinstall CherryServers server", err.Error())
		return
	}
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
