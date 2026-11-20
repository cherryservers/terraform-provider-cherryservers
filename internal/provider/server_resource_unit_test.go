package provider

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/cherryservers/cherrygo/v4"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newReinstallAttributeCases() []struct {
	name string
	plan serverResourceModel
} {
	return []struct {
		name string
		plan serverResourceModel
	}{
		{
			name: "planned image",
			plan: newServerModel(func(m *serverResourceModel) {
				m.Image = types.StringValue("test")
			}),
		},
		{
			name: "planned os_partition_size",
			plan: newServerModel(func(m *serverResourceModel) {
				m.OSPartitionSize = types.Int64Value(1)
			}),
		},
		{
			name: "planned ssh_key_ids",
			plan: newServerModel(func(m *serverResourceModel) {
				m.SSHKeyIds = types.SetValueMust(types.StringType, []attr.Value{types.StringValue("1")})
			}),
		},
		{
			name: "planned user_data",
			plan: newServerModel(func(m *serverResourceModel) {
				m.UserData = types.StringValue("test")
			}),
		},
		{
			name: "planned ipxe",
			plan: newServerModel(func(m *serverResourceModel) {
				m.IPXE = types.StringValue("test")
			}),
		},
		{
			name: "planned persist_ipxe",
			plan: newServerModel(func(m *serverResourceModel) {
				m.PersistIPXE = types.BoolValue(true)
			}),
		},
	}
}

func TestReinstallModifierDoesNotRequireAllowReinstallForCreation(t *testing.T) {
	for _, tc := range newReinstallAttributeCases() {
		r := serverResource{client: &cherrygo.Client{}}
		s := schemaFromResourceMust(t, &r)

		req := resource.ModifyPlanRequest{
			Plan:   planFromResourceModelMust(t, tc.plan, &r),
			Config: configFromResourceModelMust(t, tc.plan, &r),
		}
		resp := newModifyResp(s)

		t.Run(tc.name, func(t *testing.T) {
			r.ModifyPlan(t.Context(), req, &resp)
			require.Empty(t, resp.Diagnostics)
		})
	}
}

func TestServerReinstallModifierSetsUnknownPowerState(t *testing.T) {
	for _, tc := range newReinstallAttributeCases() {
		tc.plan.AllowReinstall = types.BoolValue(true)
		t.Run(tc.name, func(t *testing.T) {
			r := serverResource{client: &cherrygo.Client{}}
			s := schemaFromResourceMust(t, &r)

			req := newModifyReqFromModels(
				t, newServerModel(nil), newServerModel(nil), tc.plan, &r,
			)
			resp := newModifyResp(s)

			r.ModifyPlan(t.Context(), req, &resp)
			require.Empty(t, resp.Diagnostics)

			var gotPlan serverResourceModel
			require.Empty(t, resp.Plan.Get(t.Context(), &gotPlan))
			assert.Equal(t, types.StringUnknown(), gotPlan.PowerState)
		})
	}
}

func TestServerReinstallModifierSetsUnknownPrivateIP(t *testing.T) {
	ipModels := []ipAddressFlatResourceModel{
		{
			Id:            types.StringValue("uuid1"),
			Type:          types.StringValue("primary-ip"),
			Address:       types.StringValue("84.32.109.107"),
			AddressFamily: types.Int64Value(4),
			CIDR:          types.StringValue("84.32.109.0/25"),
		},
		{
			Id:            types.StringValue("uuid2"),
			Type:          types.StringValue("floating-ip"),
			Address:       types.StringValue("84.32.109.108"),
			AddressFamily: types.Int64Value(4),
			CIDR:          types.StringValue("84.32.109.0/25"),
		},
		{
			Id:            types.StringValue("uuid3"),
			Type:          types.StringValue("private-ip"),
			Address:       types.StringValue("10.190.213.5"),
			AddressFamily: types.Int64Value(4),
			CIDR:          types.StringValue("10.190.213.0/24"),
			VLAN:          types.Int64Value(1235),
		},
		{
			Id:            types.StringValue("uuid4"),
			Type:          types.StringValue("private-ip"),
			Address:       types.StringValue("10.190.213.6"),
			AddressFamily: types.Int64Value(4),
			CIDR:          types.StringValue("10.190.213.0/24"),
			VLAN:          types.Int64Value(1236),
		},
	}

	ipObjects := make([]attr.Value, len(ipModels))
	for i := range ipModels {
		ip, diags := types.ObjectValueFrom(t.Context(), ipModels[i].AttributeTypes(), ipModels[i])
		require.Empty(t, diags)
		ipObjects[i] = ip
	}

	ips := types.SetValueMust(types.ObjectType{AttrTypes: ipModels[0].AttributeTypes()}, ipObjects)

	for _, tc := range newReinstallAttributeCases() {
		tc.plan.AllowReinstall = types.BoolValue(true)
		t.Run(tc.name, func(t *testing.T) {
			tc.plan.IpAddresses = ips

			r := serverResource{client: &cherrygo.Client{}}
			s := schemaFromResourceMust(t, &r)

			req := newModifyReqFromModels(
				t, newServerModel(nil), newServerModel(nil), tc.plan, &r,
			)
			resp := newModifyResp(s)

			r.ModifyPlan(t.Context(), req, &resp)

			require.Empty(t, resp.Diagnostics)

			plannedIPs := make([]ipAddressFlatResourceModel, len(ipModels))
			var gotPlan serverResourceModel
			require.Empty(t, resp.Plan.Get(t.Context(), &gotPlan))
			require.Empty(t, gotPlan.IpAddresses.ElementsAs(t.Context(), &plannedIPs, false))

			assert.Contains(t, plannedIPs, ipModels[0])
			assert.Contains(t, plannedIPs, ipModels[1])
			assert.Contains(t, plannedIPs, ipAddressFlatResourceModel{
				Id:            types.StringUnknown(),
				Type:          types.StringValue("private-ip"),
				Address:       types.StringUnknown(),
				AddressFamily: types.Int64Value(4),
				CIDR:          types.StringUnknown(),
				VLAN:          types.Int64Value(1235),
			})
			assert.Contains(t, plannedIPs, ipAddressFlatResourceModel{
				Id:            types.StringUnknown(),
				Type:          types.StringValue("private-ip"),
				Address:       types.StringUnknown(),
				AddressFamily: types.Int64Value(4),
				CIDR:          types.StringUnknown(),
				VLAN:          types.Int64Value(1236),
			})
			assert.Len(t, plannedIPs, len(ipModels))
		})
	}
}

func TestServerReinstallModifierSetsSSHToEmptyWhenInstallingIPXE(t *testing.T) {
	r := serverResource{client: &cherrygo.Client{}}
	s := schemaFromResourceMust(t, &r)

	state := newServerModel(func(m *serverResourceModel) {
		m.Image = types.StringValue("not-ipxe")
	})
	plan := newServerModel(func(m *serverResourceModel) {
		m.AllowReinstall = types.BoolValue(true)
		m.IPXE = types.StringValue("test")
		m.SSHKeyIds = types.SetValueMust(types.StringType, []attr.Value{types.StringValue("1")})
	})
	req := newModifyReqFromModels(t, newServerModel(nil), state, plan, &r)
	resp := newModifyResp(s)

	r.ModifyPlan(t.Context(), req, &resp)

	require.Empty(t, resp.Diagnostics)

	var gotPlan serverResourceModel
	require.Empty(t, resp.Plan.Get(t.Context(), &gotPlan))
	want := types.SetValueMust(types.StringType, []attr.Value{})

	assert.Equal(t, want, gotPlan.SSHKeyIds)
}

func TestServerReinstallModifierRetainsSSHWhenNonIPXE(t *testing.T) {
	r := serverResource{client: &cherrygo.Client{}}
	s := schemaFromResourceMust(t, &r)

	state := newServerModel(func(m *serverResourceModel) {
		m.Image = types.StringValue("not-ipxe")
	})
	plan := newServerModel(func(m *serverResourceModel) {
		m.AllowReinstall = types.BoolValue(true)
		m.Image = types.StringValue("another-not-ipxe")
		m.SSHKeyIds = types.SetValueMust(types.StringType, []attr.Value{types.StringValue("1")})
	})
	req := newModifyReqFromModels(t, newServerModel(nil), state, plan, &r)
	resp := newModifyResp(s)

	r.ModifyPlan(t.Context(), req, &resp)

	require.Empty(t, resp.Diagnostics)

	var gotPlan serverResourceModel
	require.Empty(t, resp.Plan.Get(t.Context(), &gotPlan))
	want := types.SetValueMust(types.StringType, []attr.Value{types.StringValue("1")})

	assert.Equal(t, want, gotPlan.SSHKeyIds)
}

type fakeImgLister struct {
	sendImages []cherrygo.Image
	sendErr    error
	gotPlan    string
}

func (l *fakeImgLister) List(_ context.Context, plan string, _ *cherrygo.GetOptions) ([]cherrygo.Image, *cherrygo.Response, error) {
	l.gotPlan = plan
	return l.sendImages, nil, l.sendErr
}

func TestServerReinstallModifierSetsDefaultImageWhenInstallingNonIPXE(t *testing.T) {
	cases := []struct {
		name      string
		lister    fakeImgLister
		wantImage string
	}{
		{
			name: "gets latest ubuntu if available",
			lister: fakeImgLister{
				sendImages: []cherrygo.Image{
					{Slug: "some-img"},
					{Slug: "ubuntu_24_04_64bit"},
					{Slug: "ubuntu_26_04_64bit"},
					{Slug: "another-img"},
				},
				sendErr: nil,
			},
			wantImage: "ubuntu_26_04_64bit",
		},
		{
			name: "gets first if no ubuntu",
			lister: fakeImgLister{
				sendImages: []cherrygo.Image{
					{Slug: "some-img"},
					{Slug: "another-img"},
				},
				sendErr: nil,
			},
			wantImage: "some-img",
		},
	}

	for _, tc := range cases {
		r := serverResource{client: &cherrygo.Client{Images: &tc.lister}}
		s := schemaFromResourceMust(t, &r)

		state := newServerModel(func(m *serverResourceModel) {
			m.Image = types.StringValue("custom_ipxe_install")
		})
		plan := newServerModel(func(m *serverResourceModel) {
			m.AllowReinstall = types.BoolValue(true)
			m.Plan = types.StringValue("test-plan")
		})
		req := newModifyReqFromModels(t, newServerModel(nil), state, plan, &r)
		resp := newModifyResp(s)

		t.Run(tc.name, func(t *testing.T) {
			r.ModifyPlan(t.Context(), req, &resp)
			require.Empty(t, resp.Diagnostics)

			var gotPlan serverResourceModel
			require.Empty(t, resp.Plan.Get(t.Context(), &gotPlan))

			assert.Equal(t, tc.wantImage, gotPlan.Image.ValueString())
			assert.Equal(t, "test-plan", tc.lister.gotPlan)
		})
	}
}

func TestServerReinstallModifierReturnsErrorWhenDefaultImageFails(t *testing.T) {
	cases := []struct {
		name   string
		lister fakeImgLister
	}{
		{
			name: "when image list error",
			lister: fakeImgLister{
				sendErr: errors.New("test-error"),
			},
		},
		{
			name: "when no images",
			lister: fakeImgLister{
				sendImages: []cherrygo.Image{},
			},
		},
	}

	for _, tc := range cases {
		r := serverResource{client: &cherrygo.Client{Images: &tc.lister}}
		s := schemaFromResourceMust(t, &r)

		state := newServerModel(func(m *serverResourceModel) {
			m.Image = types.StringValue("custom_ipxe_install")
		})
		plan := newServerModel(func(m *serverResourceModel) {
			m.AllowReinstall = types.BoolValue(true)
			m.Plan = types.StringValue("test-plan")
		})
		req := newModifyReqFromModels(t, newServerModel(nil), state, plan, &r)
		resp := newModifyResp(s)

		t.Run(tc.name, func(t *testing.T) {
			r.ModifyPlan(t.Context(), req, &resp)
			diags := resp.Diagnostics

			require.Len(t, diags, 1)
			assert.Regexp(
				t,
				regexp.MustCompile("Failed to get a default image for plan"),
				diags[0].Detail(),
			)

			assert.Equal(t, "test-plan", tc.lister.gotPlan)
		})
	}
}

func TestServerReinstallModifierDoesNotOverrideImageWhenInstallingNonIPXE(t *testing.T) {
	r := serverResource{client: &cherrygo.Client{}}
	s := schemaFromResourceMust(t, &r)

	state := newServerModel(func(m *serverResourceModel) {
		m.Image = types.StringValue("custom_ipxe_install")
	})
	config := newServerModel(func(m *serverResourceModel) {
		m.Image = types.StringValue("test-img")
	})
	plan := newServerModel(func(m *serverResourceModel) {
		m.AllowReinstall = types.BoolValue(true)
		m.Plan = types.StringValue("test-plan")
		m.Image = types.StringValue("test-img")
	})
	req := newModifyReqFromModels(t, config, state, plan, &r)
	resp := newModifyResp(s)

	r.ModifyPlan(t.Context(), req, &resp)
	require.Empty(t, resp.Diagnostics)

	var gotPlan serverResourceModel
	require.Empty(t, resp.Plan.Get(t.Context(), &gotPlan))
	assert.Equal(t, "test-img", gotPlan.Image.ValueString())
}

func TestServerModifyPlanReturnsErrorWithInvalidRequests(t *testing.T) {
	r := serverResource{client: &cherrygo.Client{}}
	s := schemaFromResourceMust(t, &r)
	okVal := planFromResourceModelMust(t, newServerModel(nil), &r).Raw

	cases := []struct {
		name                         string
		configRaw, planRaw, stateRaw tftypes.Value
	}{
		{
			name:      "bad config",
			configRaw: tftypes.NewValue(tftypes.Number, 1),
			planRaw:   okVal,
			stateRaw:  okVal,
		},
		{
			name:      "bad plan",
			configRaw: okVal,
			planRaw:   tftypes.NewValue(tftypes.Number, 1),
			stateRaw:  okVal,
		},
		{
			name:      "bad state",
			configRaw: okVal,
			planRaw:   okVal,
			stateRaw:  tftypes.NewValue(tftypes.Number, 1),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := resource.ModifyPlanRequest{
				Plan:   tfsdk.Plan{Schema: s, Raw: tc.planRaw},
				Config: tfsdk.Config{Schema: s, Raw: tc.configRaw},
				State:  tfsdk.State{Schema: s, Raw: tc.stateRaw},
			}

			resp := newModifyResp(s)

			r.ModifyPlan(t.Context(), req, &resp)
			diags := resp.Diagnostics
			require.Len(t, diags, 1)
			assert.Regexp(
				t, regexp.MustCompile(`trying to convert`), diags[0].Detail(),
			)
		})
	}
}

func TestServerModifyPlanSkipsDeleteRequests(t *testing.T) {
	r := serverResource{client: &cherrygo.Client{}}
	s := schemaFromResourceMust(t, &r)

	var req resource.ModifyPlanRequest
	resp := newModifyResp(s)

	r.ModifyPlan(t.Context(), req, &resp)
	require.Empty(t, resp.Diagnostics)
	assert.True(t, resp.Plan.Raw.IsNull())
}

func TestServerModifyPlanSetsIPXEImageWhenIPXEIsPlanned(t *testing.T) {
	r := serverResource{client: &cherrygo.Client{}}
	s := schemaFromResourceMust(t, &r)
	planModel := newServerModel(func(m *serverResourceModel) {
		m.AllowReinstall = types.BoolValue(true)
		m.IPXE = types.StringValue("test")
	})

	cases := []struct {
		name  string
		state tfsdk.State
	}{
		{
			name:  "creation request",
			state: tfsdk.State{Schema: s},
		},
		{
			name:  "update request",
			state: stateFromResourceModelMust(t, newServerModel(nil), &r),
		},
	}

	for _, tc := range cases {
		req := resource.ModifyPlanRequest{
			State:  tc.state,
			Config: configFromResourceModelMust(t, newServerModel(nil), &r),
			Plan:   planFromResourceModelMust(t, planModel, &r),
		}
		resp := newModifyResp(s)

		t.Run(tc.name, func(t *testing.T) {
			r.ModifyPlan(t.Context(), req, &resp)
			require.Empty(t, resp.Diagnostics)

			var gotPlan serverResourceModel
			require.Empty(t, resp.Plan.Get(t.Context(), &gotPlan))

			assert.Equal(t, "custom_ipxe_install", gotPlan.Image.ValueString())
		})
	}
}

func TestServerModifyPlanDoesNothingWhenUpdatingWithoutReinstall(t *testing.T) {
	r := serverResource{client: &cherrygo.Client{}}
	s := schemaFromResourceMust(t, &r)
	planModel := newServerModel(func(m *serverResourceModel) {
		m.Hostname = types.StringValue("test")
	})
	stateModel := newServerModel(func(m *serverResourceModel) {
		m.Hostname = types.StringValue("old")
	})
	req := newModifyReqFromModels(t, planModel, stateModel, planModel, &r)
	resp := newModifyResp(s)

	r.ModifyPlan(t.Context(), req, &resp)
	require.Empty(t, resp.Diagnostics)

	var gotPlan serverResourceModel
	require.Empty(t, resp.Plan.Get(t.Context(), &gotPlan))

	assert.Equal(t, planModel, gotPlan)
}

// newServerModel runs f over a model with zero valued fields set to null,
// as required for a valid model, and returns the model. Nil f is ignored.
func newServerModel(f func(m *serverResourceModel)) serverResourceModel {
	ipsType := ipAddressFlatResourceModel{}.AttributeTypes()
	pricingType := serverPricingModel{}.AttributeTypes()
	timeoutsType := map[string]attr.Type{
		"create": types.StringType,
		"update": types.StringType,
	}

	m := serverResourceModel{
		Plan:                types.StringNull(),
		ProjectId:           types.Int64Null(),
		Region:              types.StringNull(),
		Hostname:            types.StringNull(),
		Image:               types.StringNull(),
		SSHKeyIds:           types.SetNull(types.StringType),
		ExtraIPAddressesIds: types.SetNull(types.StringType),
		ConfigureIPv6:       types.BoolNull(),
		IPAddressesIds:      types.SetNull(types.StringType),
		UserData:            types.StringNull(),
		IPXE:                types.StringNull(),
		PersistIPXE:         types.BoolNull(),
		Tags:                types.MapNull(types.StringType),
		SpotInstance:        types.BoolNull(),
		OSPartitionSize:     types.Int64Null(),
		PowerState:          types.StringNull(),
		State:               types.StringNull(),
		IpAddresses:         types.SetNull(types.ObjectType{AttrTypes: ipsType}),
		Id:                  types.StringNull(),
		Timeouts:            timeouts.Value{Object: types.ObjectNull(timeoutsType)},
		AllowReinstall:      types.BoolNull(),
		Cycle:               types.StringNull(),
		DiscountCode:        types.StringNull(),
		Pricing:             types.ObjectNull(pricingType),
	}

	if f != nil {
		f(&m)
	}
	return m
}

func newModifyReqFromModels[T serverResourceModel](t *testing.T, config, state, plan T, r resource.Resource) resource.ModifyPlanRequest {
	return resource.ModifyPlanRequest{
		Config: configFromResourceModelMust(t, config, r),
		Plan:   planFromResourceModelMust(t, plan, r),
		State:  stateFromResourceModelMust(t, state, r),
	}
}

func newModifyResp(s schema.Schema) resource.ModifyPlanResponse {
	return resource.ModifyPlanResponse{
		Plan: tfsdk.Plan{
			Schema: s,
		},
	}
}
