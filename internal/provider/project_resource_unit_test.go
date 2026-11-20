package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

// newProjectModel runs f over a model with zero valued fields set to null,
// as required for a valid model, and returns the model. Nil f is ignored.
func newProjectModel(f func(m *projectResourceModel)) projectResourceModel {
	m := projectResourceModel{
		Name:   types.StringNull(),
		TeamId: types.Int64Null(),
		BGP:    types.ObjectNull(projectBGPModel{}.AttributeTypes()),
		Id:     types.StringNull(),
	}

	if f != nil {
		f(&m)
	}
	return m
}
