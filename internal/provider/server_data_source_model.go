package provider

import (
	"context"
	"github.com/cherryservers/cherrygo/v3"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type serverDataSourceModel struct {
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
	UserDataFile        types.String   `tfsdk:"user_data_file"`
	Tags                types.Map      `tfsdk:"tags"`
	SpotInstance        types.Bool     `tfsdk:"spot_instance"`
	OSPartitionSize     types.Int64    `tfsdk:"os_partition_size"`
	PowerState          types.String   `tfsdk:"power_state"`
	State               types.String   `tfsdk:"state"`
	IpAddresses         types.Set      `tfsdk:"ip_addresses"`
	Id                  types.String   `tfsdk:"id"`
	Timeouts            timeouts.Value `tfsdk:"timeouts"`
}

func (d *serverDataSourceModel) populateModel(server cherrygo.Server, ctx context.Context, diags diag.Diagnostics, powerState string) {
	var resourceModel serverResourceModel
	resourceModel.populateModel(server, ctx, diags, powerState)

	d.Plan = resourceModel.Plan
	d.ProjectId = resourceModel.ProjectId
	d.Region = resourceModel.Region
	d.Hostname = resourceModel.Hostname
	d.Name = resourceModel.Name
	d.Username = resourceModel.Username
	d.Password = resourceModel.Password
	d.BMC = resourceModel.BMC
	d.Image = resourceModel.Image
	d.SSHKeyIds = resourceModel.SSHKeyIds
	//d.ExtraIPAddressesIds = resourceModel.ExtraIPAddressesIds
	d.UserDataFile = resourceModel.UserDataFile
	d.Tags = resourceModel.Tags
	d.SpotInstance = resourceModel.SpotInstance
	d.OSPartitionSize = resourceModel.OSPartitionSize
	d.PowerState = resourceModel.PowerState
	d.State = resourceModel.State
	d.IpAddresses = resourceModel.IpAddresses
	d.Id = resourceModel.Id

}
