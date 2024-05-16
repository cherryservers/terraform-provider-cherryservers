package provider

import (
	"context"
	"github.com/cherryservers/cherrygo/v3"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
)

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

	//d.Image = types.StringValue(server.Image)

	//var sshKeyIds, ipIds []string
	var sshKeyIds []string
	for _, sshKey := range server.SSHKeys {
		sshKeyID := strconv.Itoa(sshKey.ID)
		sshKeyIds = append(sshKeyIds, sshKeyID)
	}
	sshKeyIdsTf, sshDiags := types.SetValueFrom(ctx, types.StringType, sshKeyIds)
	d.SSHKeyIds = sshKeyIdsTf
	diags.Append(sshDiags...)

	var ips []attr.Value
	for _, ip := range server.IPAddresses {

		// ExtraIPAddresses shouldn't have unmodifiable (primary and private type) IPs
		/*if ip.Type == "subnet" || ip.Type == "floating-ip" {
			ipIds = append(ipIds, ip.ID)
		}*/

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

	/*ipIdsTf, ipIdDiags := types.SetValueFrom(ctx, types.StringType, ipIds)
	d.ExtraIPAddressesIds = ipIdsTf
	diags.Append(ipIdDiags...)*/

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
