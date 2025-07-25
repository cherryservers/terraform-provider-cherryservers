package provider

import (
	"context"
	"errors"
	"strconv"

	"github.com/cherryservers/cherrygo/v3"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type bgpModel struct {
	Hosts types.List  `tfsdk:"hosts"`
	ASN   types.Int64 `tfsdk:"asn"`
}

func bgpAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"hosts": types.ListType{ElemType: types.StringType},
		"asn":   types.Int64Type,
	}
}

type regionModel struct {
	Name       types.String `tfsdk:"name"`
	ID         types.Int64  `tfsdk:"id"`
	Slug       types.String `tfsdk:"slug"`
	RegionISO2 types.String `tfsdk:"region_iso_2"`
	BGP        types.Object `tfsdk:"bgp"`
}

func regionAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":         types.StringType,
		"id":           types.Int64Type,
		"slug":         types.StringType,
		"region_iso_2": types.StringType,
		"bgp":          types.ObjectType{AttrTypes: bgpAttributeTypes()},
	}
}

func (m *regionModel) populateState(ctx context.Context, region cherrygo.Region) diag.Diagnostics {
	m.Name = types.StringValue(region.Name)
	m.ID = types.Int64Value(int64(region.ID))
	m.Slug = types.StringValue((region.Slug))
	m.RegionISO2 = types.StringValue((region.RegionIso2))

	hosts, diags := types.ListValueFrom(ctx, types.StringType, region.BGP.Hosts)
	if diags.HasError() {
		return diags
	}

	bgp := bgpModel{
		Hosts: hosts,
		ASN:   types.Int64Value(int64(region.BGP.Asn)),
	}
	bgpObject, diags := types.ObjectValueFrom(ctx, bgpAttributeTypes(), bgp)
	m.BGP = bgpObject

	return diags
}

// Region can be identified by ID or slug.
func (m *regionModel) Identifier() (string, error) {
	slug := m.Slug.ValueString()
	id := m.ID.ValueInt64()

	if slug == "" && id == 0 {
		return "", errors.New("unidentifiable region, no slug or ID set")
	}
	if slug != "" && id != 0 {
		return "", errors.New("unidentifiable region, both slug and ID set")
	}
	if slug != "" {
		return slug, nil
	}
	return strconv.Itoa(int(id)), nil
}

func regionSchema(readOnly bool) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Computed:    true,
			Description: "Region name.",
		},
		"id": schema.Int64Attribute{
			Optional:    !readOnly,
			Computed:    true,
			Description: "Region ID.",
		},
		"slug": schema.StringAttribute{
			Optional:    !readOnly,
			Computed:    true,
			Description: "Region slug.",
		},
		"region_iso_2": schema.StringAttribute{
			Computed:    true,
			Description: "Region ISO 3166-1 alpha-2 code.",
		},
		"bgp": schema.SingleNestedAttribute{
			Attributes: map[string]schema.Attribute{
				"hosts": schema.ListAttribute{
					ElementType: types.StringType,
					Computed:    true,
					Description: "BGP host addresses.",
				},
				"asn": schema.Int64Attribute{
					Computed:    true,
					Description: "BGP ASN.",
				},
			},
			Computed:    true,
			Description: "Region BGP data.",
		},
	}
}
