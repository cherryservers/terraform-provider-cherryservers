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
