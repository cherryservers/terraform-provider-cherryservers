package provider

//All are WIP!

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ planmodifier.String = warnServerReinstallNeededModifier{}
	_ planmodifier.Int64  = warnServerReinstallNeededModifier{}
	_ planmodifier.Set    = warnServerReinstallNeededModifier{}
)

type warnServerReinstallNeededModifier struct {
}

const reinstallWarningDetail string = `When updating "image", "ssh_key_ids", "os_partition_size" or "user_data" values, the server OS has to be reinstalled.`
const reinstallWarningSummary string = `Warning: Server reinstall required.`

func (m warnServerReinstallNeededModifier) Description(_ context.Context) string {
	return "Diagnostics warning that a server reinstall will be needed"
}

func (m warnServerReinstallNeededModifier) MarkdownDescription(_ context.Context) string {
	return "Diagnostics warning that a server reinstall will be needed"
}

func (m warnServerReinstallNeededModifier) PlanModifyString(_ context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Not applicable on resource creation and destruction
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	if req.PlanValue.IsUnknown() {
		return
	}

	if req.StateValue.Equal(req.PlanValue) {
		return
	}

	resp.Diagnostics.AddAttributeWarning(req.Path, reinstallWarningSummary,
		reinstallWarningDetail)
}

func WarnServerReinstallNeededString() planmodifier.String {
	return warnServerReinstallNeededModifier{}
}

func (m warnServerReinstallNeededModifier) PlanModifyInt64(_ context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	// Not applicable on resource creation and destruction
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	if req.PlanValue.IsUnknown() {
		return
	}

	if req.StateValue.Equal(req.PlanValue) {
		return
	}

	resp.Diagnostics.AddAttributeWarning(req.Path, reinstallWarningSummary,
		reinstallWarningDetail)
}

func WarnServerReinstallNeededInt64() planmodifier.Int64 {
	return warnServerReinstallNeededModifier{}
}

func (m warnServerReinstallNeededModifier) PlanModifySet(_ context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	// Not applicable on resource creation and destruction
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	if req.PlanValue.IsUnknown() {
		return
	}

	if req.StateValue.Equal(req.PlanValue) {
		return
	}

	resp.Diagnostics.AddAttributeWarning(req.Path, reinstallWarningSummary,
		reinstallWarningDetail)
}

func WarnServerReinstallNeededSet() planmodifier.Set {
	return warnServerReinstallNeededModifier{}
}

type stateUnknownIfReinstallingModifier struct {
}

func (m stateUnknownIfReinstallingModifier) Description(_ context.Context) string {
	return "Use state for unknown if condition is met."
}

func (m stateUnknownIfReinstallingModifier) MarkdownDescription(_ context.Context) string {
	return "Use state for unknown if condition is met."
}

func (m stateUnknownIfReinstallingModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Not applicable on resource creation and destruction
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {

		return
	}

	var plan serverResourceModel
	var state serverResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.Image.Equal(state.Image) || !plan.OSPartitionSize.Equal(state.OSPartitionSize) ||
		!plan.SSHKeyIds.Equal(state.SSHKeyIds) ||
		!plan.UserDataFile.Equal(state.UserDataFile) {
		resp.PlanValue = types.StringUnknown()
	} else {
		resp.PlanValue = req.StateValue
	}
}

func StateUnknownIfReinstallingString() planmodifier.String {
	return stateUnknownIfReinstallingModifier{}
}

type ipAddressDependantOnExtraModifier struct {
}

func (m ipAddressDependantOnExtraModifier) Description(_ context.Context) string {
	return "Use state for unknown if condition is met."
}

func (m ipAddressDependantOnExtraModifier) MarkdownDescription(_ context.Context) string {
	return "Use state for unknown if condition is met."
}

func (m ipAddressDependantOnExtraModifier) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	// Not applicable on resource creation and destruction
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {

		return
	}

	var plan serverResourceModel
	var state serverResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !req.StateValue.Equal(req.PlanValue) {
		resp.PlanValue = types.SetUnknown(types.StringType)
	} else {
		resp.PlanValue = types.SetUnknown(types.StringType)
	}
}

func IpAddressDependantOnExtraModifier() planmodifier.Set {
	return ipAddressDependantOnExtraModifier{}
}
