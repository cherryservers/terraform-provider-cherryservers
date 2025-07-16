package provider

import (
	"context"

	"github.com/cherryservers/cherrygo/v3"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// imageModel corresponds to the inner "image" object under softwares
type imageModel struct {
	Name types.String `tfsdk:"name"`
	Slug types.String `tfsdk:"slug"`
}

func imageAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
		"slug": types.StringType,
	}
}

// softwareModel wraps the imageModel
type softwareModel struct {
	Image types.Object `tfsdk:"image"`
}

func softwareAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"image": types.ObjectType{AttrTypes: imageAttributeTypes()},
	}
}

// cpusModel for the "cpus" block under specs
type cpusModel struct {
	Count     types.Int64   `tfsdk:"count"`
	Name      types.String  `tfsdk:"name"`
	Cores     types.Int64   `tfsdk:"cores"`
	Frequency types.Float64 `tfsdk:"frequency"`
	Unit      types.String  `tfsdk:"unit"`
}

func cpusAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"count":     types.Int64Type,
		"name":      types.StringType,
		"cores":     types.Int64Type,
		"frequency": types.Float64Type,
		"unit":      types.StringType,
	}
}

// memoryModel for the "memory" block under specs
type memoryModel struct {
	Count types.Int64  `tfsdk:"count"`
	Total types.Int64  `tfsdk:"total"`
	Unit  types.String `tfsdk:"unit"`
	Name  types.String `tfsdk:"name"`
}

func memoryAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"count": types.Int64Type,
		"total": types.Int64Type,
		"unit":  types.StringType,
		"name":  types.StringType,
	}
}

// storageModel for each element in the "storage" list under specs
type storageModel struct {
	Count types.Int64  `tfsdk:"count"`
	Name  types.String `tfsdk:"name"`
	Size  types.Int64  `tfsdk:"size"`
	Unit  types.String `tfsdk:"unit"`
}

func storageAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"count": types.Int64Type,
		"name":  types.StringType,
		"size":  types.Int64Type,
		"unit":  types.StringType,
	}
}

// nicsModel for the "nics" block under specs
type nicsModel struct {
	Name types.String `tfsdk:"name"`
}

func nicsAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
	}
}

// bandwidthModel for the "bandwidth" block under specs
type bandwidthModel struct {
	Name types.String `tfsdk:"name"`
}

func bandwidthAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
	}
}

// specsModel wraps all nested specs
type specsModel struct {
	CPUs      types.Object `tfsdk:"cpus"`
	Memory    types.Object `tfsdk:"memory"`
	Storage   types.List   `tfsdk:"storage"`
	NICs      types.Object `tfsdk:"nics"`
	Bandwidth types.Object `tfsdk:"bandwidth"`
}

func specsAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cpus":      types.ObjectType{AttrTypes: cpusAttributeTypes()},
		"memory":    types.ObjectType{AttrTypes: memoryAttributeTypes()},
		"storage":   types.ListType{ElemType: types.ObjectType{AttrTypes: storageAttributeTypes()}},
		"nics":      types.ObjectType{AttrTypes: nicsAttributeTypes()},
		"bandwidth": types.ObjectType{AttrTypes: bandwidthAttributeTypes()},
	}
}

// pricingModel for each element in the "pricing" list
type pricingModel struct {
	Unit     types.String  `tfsdk:"unit"`
	Price    types.Float64 `tfsdk:"price"`
	Currency types.String  `tfsdk:"currency"`
}

func pricingAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"unit":     types.StringType,
		"price":    types.Float64Type,
		"currency": types.StringType,
	}
}

// planRegionModel for each element in the "available_regions" list
type planRegionModel struct {
	ID         types.Int64  `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	RegionISO2 types.String `tfsdk:"region_iso_2"`
	StockQty   types.Int64  `tfsdk:"stock_qty"`
	SpotQty    types.Int64  `tfsdk:"spot_qty"`
	Slug       types.String `tfsdk:"slug"`
	BGP        types.Object `tfsdk:"bgp"`
	Location   types.String `tfsdk:"location"`
}

func planRegionAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":           types.Int64Type,
		"name":         types.StringType,
		"region_iso_2": types.StringType,
		"stock_qty":    types.Int64Type,
		"spot_qty":     types.Int64Type,
		"slug":         types.StringType,
		"bgp":          types.ObjectType{AttrTypes: bgpAttributeTypes()},
		"location":     types.StringType,
	}
}

// planModel is the top-level struct for each element in the API response array.
type planModel struct {
	ID               types.Int64  `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Slug             types.String `tfsdk:"slug"`
	Type             types.String `tfsdk:"type"`
	Softwares        types.List   `tfsdk:"softwares"`
	Specs            types.Object `tfsdk:"specs"`
	Pricing          types.List   `tfsdk:"pricing"`
	AvailableRegions types.List   `tfsdk:"available_regions"`
}

func planAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                types.Int64Type,
		"name":              types.StringType,
		"slug":              types.StringType,
		"type":              types.StringType,
		"softwares":         types.ListType{ElemType: types.ObjectType{AttrTypes: softwareAttributeTypes()}},
		"specs":             types.ObjectType{AttrTypes: specsAttributeTypes()},
		"pricing":           types.ListType{ElemType: types.ObjectType{AttrTypes: pricingAttributeTypes()}},
		"available_regions": types.ListType{ElemType: types.ObjectType{AttrTypes: planRegionAttributeTypes()}},
	}
}

func expandSoftwares(ctx context.Context, softwares []cherrygo.SoftwareImage) (types.List, diag.Diagnostics) {
	diagsCap := len(softwares)*2 + 1
	var diags diag.Diagnostics = make(diag.Diagnostics, 0, diagsCap)
	swObjs := make([]types.Object, len(softwares))

	for i, v := range softwares {
		imageModel := imageModel{Name: types.StringValue(v.Image.Name), Slug: types.StringValue(v.Image.Slug)}
		imageObj, d := types.ObjectValueFrom(ctx, imageAttributeTypes(), imageModel)
		diags.Append(d...)

		swModel := softwareModel{Image: imageObj}
		swObj, d := types.ObjectValueFrom(ctx, softwareAttributeTypes(), swModel)
		diags.Append(d...)

		swObjs[i] = swObj
	}

	swList, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: softwareAttributeTypes()}, swObjs)
	diags.Append(d...)
	return swList, diags
}

func expandPricings(ctx context.Context, pricings []cherrygo.Pricing) (types.List, diag.Diagnostics) {
	diagsCap := len(pricings) + 1
	var diags diag.Diagnostics = make(diag.Diagnostics, 0, diagsCap)
	objs := make([]types.Object, len(pricings))

	for i, v := range pricings {
		model := pricingModel{
			Unit:     types.StringValue(v.Unit),
			Price:    types.Float64Value(float64(v.Price)),
			Currency: types.StringValue(v.Currency),
		}

		obj, d := types.ObjectValueFrom(ctx, pricingAttributeTypes(), model)
		diags.Append(d...)
		objs[i] = obj
	}

	list, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: pricingAttributeTypes()}, objs)
	diags.Append(d...)
	return list, diags
}

func expandStorage(ctx context.Context, storages []cherrygo.Storage) (types.List, diag.Diagnostics) {
	diagsCap := len(storages) + 1
	var diags diag.Diagnostics = make(diag.Diagnostics, 0, diagsCap)
	objs := make([]types.Object, len(storages))

	for i, v := range storages {
		m := storageModel{
			Count: types.Int64Value(int64(v.Count)),
			Name:  types.StringValue(v.Name),
			Size:  types.Int64Value(int64(v.Size)),
			Unit:  types.StringValue(v.Unit),
		}

		obj, d := types.ObjectValueFrom(ctx, storageAttributeTypes(), m)
		diags.Append(d...)
		objs[i] = obj
	}

	list, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: storageAttributeTypes()}, objs)
	diags.Append(d...)
	return list, diags
}

func expandSpecs(ctx context.Context, specs cherrygo.Specs) (types.Object, diag.Diagnostics) {
	diagsCap := len(specs.Storage) + 5 // Number of object type elements in specs, plus length of 'storage'.
	var diags diag.Diagnostics = make(diag.Diagnostics, 0, diagsCap)

	cpu := cpusModel{
		Count:     types.Int64Value(int64(specs.Cpus.Count)),
		Name:      types.StringValue(specs.Cpus.Name),
		Cores:     types.Int64Value(int64(specs.Cpus.Cores)),
		Frequency: types.Float64Value(float64(specs.Cpus.Frequency)),
		Unit:      types.StringValue(specs.Cpus.Unit),
	}
	cpuObj, d := types.ObjectValueFrom(ctx, cpusAttributeTypes(), cpu)
	diags.Append(d...)

	memory := memoryModel{
		Count: types.Int64Value(int64(specs.Memory.Count)),
		Total: types.Int64Value(int64(specs.Memory.Total)),
		Unit:  types.StringValue(specs.Memory.Unit),
		Name:  types.StringValue(specs.Memory.Name),
	}
	memoryObj, d := types.ObjectValueFrom(ctx, memoryAttributeTypes(), memory)
	diags.Append(d...)

	storageList, d := expandStorage(ctx, specs.Storage)
	diags.Append(d...)

	nics := nicsModel{
		Name: types.StringValue(specs.Nics.Name),
	}
	nicsObj, d := types.ObjectValueFrom(ctx, nicsAttributeTypes(), nics)
	diags.Append(d...)

	bandwidth := bandwidthModel{
		Name: types.StringValue(specs.Bandwidth.Name),
	}
	bandwidthObj, d := types.ObjectValueFrom(ctx, bandwidthAttributeTypes(), bandwidth)
	diags.Append(d...)

	m := specsModel{CPUs: cpuObj, Memory: memoryObj, Storage: storageList, NICs: nicsObj, Bandwidth: bandwidthObj}
	obj, d := types.ObjectValueFrom(ctx, specsAttributeTypes(), m)

	return obj, d
}

func expandBGP(ctx context.Context, bgp cherrygo.RegionBGP) (types.Object, diag.Diagnostics) {
	diagsCap := 2
	var diags diag.Diagnostics = make(diag.Diagnostics, 0, diagsCap)

	hosts, d := types.ListValueFrom(ctx, types.StringType, bgp.Hosts)
	diags.Append(d...)

	m := bgpModel{
		Hosts: hosts,
		ASN:   types.Int64Value(int64(bgp.Asn)),
	}
	bgpObj, d := types.ObjectValueFrom(ctx, bgpAttributeTypes(), m)
	diags.Append(d...)
	return bgpObj, diags
}

func expandRegions(ctx context.Context, ar []cherrygo.AvailableRegions) (types.List, diag.Diagnostics) {
	diagsCap := len(ar)*3 + 1
	var diags diag.Diagnostics = make(diag.Diagnostics, 0, diagsCap)
	objs := make([]types.Object, len(ar))

	for i, v := range ar {
		bgpObj, d := expandBGP(ctx, v.BGP)
		diags.Append(d...)

		m := planRegionModel{
			ID:         types.Int64Value(int64(v.ID)),
			Name:       types.StringValue(v.Name),
			RegionISO2: types.StringValue(v.RegionIso2),
			StockQty:   types.Int64Value(int64(v.StockQty)),
			SpotQty:    types.Int64Value(int64(v.SpotQty)),
			Slug:       types.StringValue(v.Slug),
			BGP:        bgpObj,
			Location:   types.StringValue(v.Location),
		}

		obj, d := types.ObjectValueFrom(ctx, planRegionAttributeTypes(), m)
		diags.Append(d...)

		objs[i] = obj
	}

	list, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: planRegionAttributeTypes()}, objs)
	diags.Append(d...)
	return list, diags
}

func (m *planModel) populateState(ctx context.Context, plan cherrygo.Plan) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.Int64Value(int64(plan.ID))
	m.Name = types.StringValue(plan.Name)
	m.Slug = types.StringValue(plan.Slug)
	m.Type = types.StringValue(plan.Type)
	softwares, d := expandSoftwares(ctx, plan.Softwares)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}
	m.Softwares = softwares

	pricing, d := expandPricings(ctx, plan.Pricing)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}
	m.Pricing = pricing

	specs, d := expandSpecs(ctx, plan.Specs)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}
	m.Specs = specs

	ar, d := expandRegions(ctx, plan.AvailableRegions)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}
	m.AvailableRegions = ar

	return diags
}

var softwaresAttr = schema.ListNestedAttribute{
	Computed:    true,
	Description: "Plan OS images.",
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"image": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "OS image specification.",
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "Full image name.",
					},
					"slug": schema.StringAttribute{
						Computed:    true,
						Description: "Used as identifier for the image in requests.",
					},
				},
			},
		},
	},
}

var pricingAttr = schema.ListNestedAttribute{
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"unit": schema.StringAttribute{
				Computed:    true,
				Description: "Pricing period unit.",
			},
			"price": schema.Float64Attribute{
				Computed: true,
			},
			"currency": schema.StringAttribute{
				Computed:    true,
				Description: "Currency type.",
			},
		},
	},
	Computed:    true,
	Description: "Available pricing plans.",
}

var storagesAttr = schema.ListNestedAttribute{
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"count": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of storage devices.",
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"size": schema.Int64Attribute{
				Computed:    true,
				Description: "Storage capacity.",
			},
			"unit": schema.StringAttribute{
				Computed:    true,
				Description: "Storage capacity units.",
			},
		},
	},
	Computed:    true,
	Description: "Storage specification.",
}

var cpusAttr = schema.SingleNestedAttribute{
	Attributes: map[string]schema.Attribute{
		"count": schema.Int64Attribute{
			Computed:    true,
			Description: "Number of CPU devices.",
		},
		"name": schema.StringAttribute{
			Computed: true,
		},
		"cores": schema.Int64Attribute{
			Computed: true,
		},
		"frequency": schema.Float64Attribute{
			Computed: true,
		},
		"unit": schema.StringAttribute{
			Computed:    true,
			Description: "Frequency measurement unit.",
		},
	},
	Computed:    true,
	Description: "CPU specification.",
}

var memoryAttr = schema.SingleNestedAttribute{
	Attributes: map[string]schema.Attribute{
		"count": schema.Int64Attribute{
			Computed:    true,
			Description: "Number of memory devices.",
		},
		"name": schema.StringAttribute{
			Computed: true,
		},
		"total": schema.Int64Attribute{
			Computed:    true,
			Description: "Total memory capacity.",
		},
		"unit": schema.StringAttribute{
			Computed:    true,
			Description: "Memory capacity measurement unit.",
		},
	},
	Computed:    true,
	Description: "Memory specification.",
}

var nicsAttr = schema.SingleNestedAttribute{
	Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Computed: true,
		},
	},
	Computed:    true,
	Description: "NICS specification.",
}

var bandwidthAttr = schema.SingleNestedAttribute{
	Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Computed: true,
		},
	},
	Computed:    true,
	Description: "Bandwidth specification.",
}

var specsAttr = schema.SingleNestedAttribute{
	Attributes: map[string]schema.Attribute{
		"cpus":      cpusAttr,
		"memory":    memoryAttr,
		"storage":   storagesAttr,
		"nics":      nicsAttr,
		"bandwidth": bandwidthAttr,
	},
	Computed:    true,
	Description: "Server plan hardware specification.",
}

var bgpAttr = schema.SingleNestedAttribute{
	Attributes: map[string]schema.Attribute{
		"hosts": schema.ListAttribute{
			ElementType: types.StringType,
			Computed:    true,
			Description: "Host IP addresses.",
		},
		"asn": schema.Int64Attribute{
			Computed: true,
		},
	},
	Computed:    true,
	Description: "Region BGP specification.",
}

var regionsAttr = schema.ListNestedAttribute{
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"region_iso_2": schema.StringAttribute{
				Computed: true,
			},
			"stock_qty": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of instances in stock.",
			},
			"spot_qty": schema.Int64Attribute{
				Computed:    true,
				Description: "Number of spot instances in stock.",
			},
			"slug": schema.StringAttribute{
				Computed:    true,
				Description: "A more readable substitute for id.",
			},
			"bgp": bgpAttr,
			"location": schema.StringAttribute{
				Computed: true,
			},
		},
	},
	Computed:    true,
	Description: "Available region specification.",
}

func planAttr(readOnly bool) map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.Int64Attribute{
			Computed:    true,
			Optional:    !readOnly,
			Description: "Plan ID.",
		},
		"name": schema.StringAttribute{
			Computed:    true,
			Description: "Plan name.",
		},
		"slug": schema.StringAttribute{
			Computed:    true,
			Optional:    !readOnly,
			Description: "A more readable substitute for id.",
		},
		"type": schema.StringAttribute{
			Computed:    true,
			Description: "Machine type. Bare-metal, virtual, etc.",
		},
		"softwares":         softwaresAttr,
		"pricing":           pricingAttr,
		"specs":             specsAttr,
		"available_regions": regionsAttr,
	}
}
