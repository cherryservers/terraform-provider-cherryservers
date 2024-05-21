package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ planmodifier.String = unknownDependingOnUpdateAttributePlanModifier{}

// UnknownDependingOnUpdateAttributePlanModifier marks attribute as unknown if any of the `attributePaths` are being updated.
func UnknownDependingOnUpdateAttributePlanModifier(attributePaths []string) planmodifier.String {
	return &unknownDependingOnUpdateAttributePlanModifier{
		attributePaths: attributePaths,
	}
}

type unknownDependingOnUpdateAttributePlanModifier struct {
	attributePaths []string
}

func (d unknownDependingOnUpdateAttributePlanModifier) Description(ctx context.Context) string {
	return "Marks attribute as unknown if any of the attributePaths are being updated"
}

func (d unknownDependingOnUpdateAttributePlanModifier) MarkdownDescription(ctx context.Context) string {
	return d.Description(ctx)
}

func (d unknownDependingOnUpdateAttributePlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	//Ignore create or destroy cases.
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	if !req.PlanValue.IsUnknown() {
		return
	}

	for _, attributePath := range d.attributePaths {
		var attributePlan types.String
		diags := req.Plan.GetAttribute(ctx, path.Root(attributePath), &attributePlan)
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

		if (attributePlan.IsNull() && !attributeState.IsNull()) || (!attributePlan.IsNull() && attributeState.IsNull()) ||
			attributePlan.IsUnknown() || attributePlan.ValueString() != attributeState.ValueString() {
			resp.PlanValue = types.StringUnknown()
			return
		}
	}

	resp.PlanValue = req.StateValue
}
