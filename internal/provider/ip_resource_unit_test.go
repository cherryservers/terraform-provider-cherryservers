package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

// newIPModel runs f over a model with zero valued fields set to null,
// as required for a valid model, and returns the model. Nil f is ignored.
func newIPModel(f func(m *ipResourceModel)) ipResourceModel {
	m := ipResourceModel{
		Id:                 types.StringNull(),
		ProjectId:          types.Int64Null(),
		Region:             types.StringNull(),
		TargetId:           types.StringNull(),
		TargetHostname:     types.StringNull(),
		TargetIPID:         types.StringNull(),
		ARecord:            types.StringNull(),
		ARecordEffective:   types.StringNull(),
		PTRRecord:          types.StringNull(),
		PTRRecordEffective: types.StringNull(),
		Address:            types.StringNull(),
		AddressFamily:      types.Int64Null(),
		CIDR:               types.StringNull(),
		Gateway:            types.StringNull(),
		Type:               types.StringNull(),
		Tags:               types.MapNull(types.StringType),
	}

	if f != nil {
		f(&m)
	}
	return m
}
