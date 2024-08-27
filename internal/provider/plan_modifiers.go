package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ planmodifier.String = useStateIfNoConfigurationChangesAttributePlanModifier{}

// UseStateIfNoConfigurationChangesAttributePlanModifier matches a computed attribute to its state if there are no
// configuration changes in `attributePaths`. `attributePaths` should consist of required, optional or optional&computed attributes.
func UseStateIfNoConfigurationChangesAttributePlanModifier(attributePaths []string) planmodifier.String {
	return &useStateIfNoConfigurationChangesAttributePlanModifier{
		attributePaths: attributePaths,
	}
}

type useStateIfNoConfigurationChangesAttributePlanModifier struct {
	attributePaths []string
}

func (d useStateIfNoConfigurationChangesAttributePlanModifier) Description(ctx context.Context) string {
	return "Matches attribute plan to state if the practitioner has not updated the configuration."
}

func (d useStateIfNoConfigurationChangesAttributePlanModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d useStateIfNoConfigurationChangesAttributePlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	//Ignore create or destroy cases.
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	//Ignore cases where the attribute has been configured.
	if !req.ConfigValue.IsNull() {
		return
	}

	for _, attributePath := range d.attributePaths {
		var attributeConfig types.String
		diags := req.Config.GetAttribute(ctx, path.Root(attributePath), &attributeConfig)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		var attributeState types.String
		diags = req.State.GetAttribute(ctx, path.Root(attributePath), &attributeState)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		if !attributeConfig.IsNull() && !attributeState.Equal(attributeConfig) {
			resp.PlanValue = types.StringUnknown()
			return
		}

	}

	resp.PlanValue = req.StateValue
}

var _ planmodifier.String = warnIfChangedModifier{}
var _ planmodifier.Set = warnIfChangedModifier{}
var _ planmodifier.Int64 = warnIfChangedModifier{}

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

	resp.Diagnostics.AddWarning(d.warningSummary, d.warningDetail)
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

	resp.Diagnostics.AddWarning(d.warningSummary, d.warningDetail)
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

	resp.Diagnostics.AddWarning(d.warningSummary, d.warningDetail)
}
