package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/cherryservers/cherrygo/v4"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ planmodifier.String = useStateIfNoConfigurationChangesModifier{}

// UseStateIfNoConfigurationChanges assigns a computed attribute its previous state if there are no
// user configuration changes in `attributePaths`. Known limitation:
// ignores null attribute configurations. For example, if an attribute configuration has changed from
// a previously known value to null, it will be treated as unchanged. This is because there is no
// way to differentiate between attributes that are intentionally null and those that are not configured.
func UseStateIfNoConfigurationChanges(expressions ...path.Expression) planmodifier.String {
	return &useStateIfNoConfigurationChangesModifier{
		expressions: expressions,
	}
}

type useStateIfNoConfigurationChangesModifier struct {
	expressions path.Expressions
}

func (d useStateIfNoConfigurationChangesModifier) Description(ctx context.Context) string {
	return "Assigns its previous state to an attribute, if the practitioner has not updated the configuration."
}

func (d useStateIfNoConfigurationChangesModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d useStateIfNoConfigurationChangesModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Ignore create or destroy cases.
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	// Ignore cases where the attribute has a known planned value.
	if !req.PlanValue.IsUnknown() {
		return
	}

	// Ignore cases where the attribute is already unknown in the configuration.
	if req.ConfigValue.IsUnknown() {
		return
	}

	expressions := req.PathExpression.MergeExpressions(d.expressions...)

	for _, expression := range expressions {
		// Find paths matching the expression in the configuration data.
		matchedPaths, diags := req.Config.PathMatches(ctx, expression)

		resp.Diagnostics.Append(diags...)

		// Collect all errors
		if diags.HasError() {
			continue
		}

		for _, matchedPath := range matchedPaths {
			// Fetch the generic attr.Value at the given path. This ensures any
			// potential parent value of a different type, which can be a null
			// or unknown value, can be safely checked without raising a type
			// conversion error.
			var matchedPathConfigValue attr.Value
			var matchedPathStateValue attr.Value

			diags = req.Config.GetAttribute(ctx, matchedPath, &matchedPathConfigValue)
			resp.Diagnostics.Append(diags...)
			diags = req.State.GetAttribute(ctx, matchedPath, &matchedPathStateValue)
			resp.Diagnostics.Append(diags...)

			// Collect all errors
			if diags.HasError() {
				continue
			}

			if !matchedPathConfigValue.IsNull() && !matchedPathStateValue.Equal(matchedPathConfigValue) {
				return
			}
		}

	}

	resp.PlanValue = req.StateValue
}

var (
	_ planmodifier.String = warnIfChangedModifier{}
	_ planmodifier.Set    = warnIfChangedModifier{}
	_ planmodifier.Int64  = warnIfChangedModifier{}
	_ planmodifier.Bool   = warnIfChangedModifier{}
)

// WarnIfChangedString returns a plan modifier that displays a warning if an attribute will be changed on update.
func WarnIfChangedString(warningSummary, warningDetail string) planmodifier.String {
	return warnIfChangedModifier{
		warningSummary: warningSummary,
		warningDetail:  warningDetail,
	}
}

// WarnIfChangedSet returns a plan modifier that displays a warning if an attribute will be changed on update.
func WarnIfChangedSet(warningSummary, warningDetail string) planmodifier.Set {
	return warnIfChangedModifier{
		warningSummary: warningSummary,
		warningDetail:  warningDetail,
	}
}

// WarnIfChangedInt64 returns a plan modifier that displays a warning if an attribute will be changed on update.
func WarnIfChangedInt64(warningSummary, warningDetail string) planmodifier.Int64 {
	return warnIfChangedModifier{
		warningSummary: warningSummary,
		warningDetail:  warningDetail,
	}
}

func WarnIfChangedBool(warningSummary, warningDetail string) planmodifier.Bool {
	return warnIfChangedModifier{
		warningSummary: warningSummary,
		warningDetail:  warningDetail,
	}
}

type warnIfChangedModifier struct {
	warningSummary string
	warningDetail  string
}

func (d warnIfChangedModifier) Description(ctx context.Context) string {
	return "Display a warning, if the attribute will be changed on update."
}

func (d warnIfChangedModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d warnIfChangedModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Ignore create or destroy cases.
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	// Ignore if attribute has not changed.
	if req.PlanValue.Equal(req.StateValue) {
		return
	}

	resp.Diagnostics.AddAttributeWarning(req.Path, d.warningSummary, d.warningDetail)
}

func (d warnIfChangedModifier) PlanModifySet(ctx context.Context, req planmodifier.SetRequest, resp *planmodifier.SetResponse) {
	// Ignore create or destroy cases.
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	// Ignore if attribute has not changed.
	if req.PlanValue.Equal(req.StateValue) {
		return
	}

	resp.Diagnostics.AddAttributeWarning(req.Path, d.warningSummary, d.warningDetail)
}

func (d warnIfChangedModifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	// Ignore create or destroy cases.
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	// Ignore if attribute has not changed.
	if req.PlanValue.Equal(req.StateValue) {
		return
	}

	resp.Diagnostics.AddAttributeWarning(req.Path, d.warningSummary, d.warningDetail)
}

func (d warnIfChangedModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	// Ignore create or destroy cases.
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	// Ignore if attribute has not changed.
	if req.PlanValue.Equal(req.StateValue) {
		return
	}

	resp.Diagnostics.AddAttributeWarning(req.Path, d.warningSummary, d.warningDetail)
}

type imgLister interface {
	List(ctx context.Context, serverPlan string, opts *cherrygo.GetOptions) ([]cherrygo.Image, *cherrygo.Response, error)
}

type serverReinstallModifier struct {
	plan          *serverResourceModel
	state, config serverResourceModel
	lister        imgLister
}

// modify modifies the plan for server re-installation.
// Returns attribute errors if `allow_reinstall` is not enabled.
func (m *serverReinstallModifier) modify(ctx context.Context) diag.Diagnostics {
	var d diag.Diagnostics

	if !m.plan.AllowReinstall.ValueBool() {
		m.appendNotAllowedErrors(&d)
		return d
	}

	// Power state can be unpredictable when re-installing.
	m.plan.PowerState = types.StringUnknown()

	// Private IP may change on reinstall.
	d.Append(m.modifyPrivateIP(ctx)...)
	if d.HasError() {
		return d
	}

	// Ensure we don't pass SSH keys, when reinstalling non-iPXE -> iPXE,
	// since SSH keys use state, if unconfigured.
	if !m.plan.IPXE.IsNull() && m.state.Image.ValueString() != ipxeImage {
		m.plan.SSHKeyIds = types.SetValueMust(types.StringType, []attr.Value{})
	}

	// If we need to reinstall iPXE -> non-iPXE, we may need
	// to find a default image, if the user didn't configure one.
	if m.plan.IPXE.IsNull() && m.state.Image.ValueString() == ipxeImage {
		if m.config.Image.IsNull() {
			d.Append(m.modifyImage(ctx)...)
		}
	}

	return d
}

func (m *serverReinstallModifier) modifyPrivateIP(ctx context.Context) diag.Diagnostics {
	var d diag.Diagnostics

	ips := make([]ipAddressFlatResourceModel, 0, len(m.plan.IpAddresses.Elements()))
	diags := m.plan.IpAddresses.ElementsAs(ctx, &ips, false)
	d.Append(diags...)
	if d.HasError() {
		return d
	}

	ipsAttrs := make([]attr.Value, 0, len(ips))

	for i := range ips {
		if ips[i].Type.ValueString() == "private-ip" {
			ips[i].Address = types.StringUnknown()
			ips[i].Id = types.StringUnknown()
			ips[i].CIDR = types.StringUnknown()
		}

		ipAttr, ipDiags := types.ObjectValueFrom(ctx, ips[i].AttributeTypes(), ips[i])
		d.Append(ipDiags...)
		if d.HasError() {
			return d
		}

		ipsAttrs = append(ipsAttrs, ipAttr)
	}

	ipsTf, ipsDiags := types.SetValue(types.ObjectType{AttrTypes: ipAddressFlatResourceModel{}.AttributeTypes()}, ipsAttrs)
	d.Append(ipsDiags...)
	if d.HasError() {
		return d
	}

	m.plan.IpAddresses = ipsTf
	return d
}

func (m *serverReinstallModifier) modifyImage(ctx context.Context) diag.Diagnostics {
	var d diag.Diagnostics

	img, err := m.defaultImage(ctx, m.plan.Plan.ValueString())
	if err != nil {
		d.AddError("No Default Plan Image",
			fmt.Sprintf("Failed to get a default image for plan: %s.", err.Error()))
		return d
	}
	m.plan.Image = types.StringValue(img)

	return d
}

// defaultImage gets a default OS image for the given server plan.
// Specifically, it tries to find the latest Ubuntu version,
// with a fallback to a random image.
func (m *serverReinstallModifier) defaultImage(ctx context.Context, plan string) (string, error) {
	images, _, err := m.lister.List(ctx, plan, nil)
	if err != nil {
		return "", err
	}

	var newImage string
	for _, image := range images {
		// Try to pick the latest ubuntu version.
		if strings.HasPrefix(image.Slug, "ubuntu") && image.Slug > newImage {
			newImage = image.Slug
		}
	}

	// If for some reason we couldn't find an image, try to fall back to the first
	// one in the slice.
	if newImage == "" {
		if len(images) == 0 {
			return "", fmt.Errorf("no images found for plan %q", plan)
		}
		newImage = images[0].Slug
	}

	return newImage, nil
}

func (m *serverReinstallModifier) appendNotAllowedErrors(d *diag.Diagnostics) {
	var errs diag.Diagnostics

	const (
		summary = "Re-installation not allowed."
		detail  = "Updating `%s` requires `allow_reinstall` to be enabled."
	)
	if !m.plan.Image.Equal(m.state.Image) {
		d.AddAttributeError(path.Root("image"), summary, fmt.Sprintf(detail, "image"))
	}
	if !m.plan.OSPartitionSize.Equal(m.state.OSPartitionSize) {
		d.AddAttributeError(path.Root("os_partition_size"), summary, fmt.Sprintf(detail, "os_partition_size"))
	}
	if !m.plan.SSHKeyIds.Equal(m.state.SSHKeyIds) {
		d.AddAttributeError(path.Root("ssh_key_ids"), summary, fmt.Sprintf(detail, "ssh_key_ids"))
	}
	if !m.plan.UserData.Equal(m.state.UserData) {
		d.AddAttributeError(path.Root("user_data"), summary, fmt.Sprintf(detail, "user_data"))
	}
	if !m.plan.IPXE.Equal(m.state.IPXE) {
		d.AddAttributeError(path.Root("ipxe"), summary, fmt.Sprintf(detail, "ipxe"))
	}
	if !m.plan.PersistIPXE.Equal(m.state.PersistIPXE) {
		d.AddAttributeError(path.Root("persist_ipxe"), summary, fmt.Sprintf(detail, "persist_ipxe"))
	}

	d.Append(errs...)
}
