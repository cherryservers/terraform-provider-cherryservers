package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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
