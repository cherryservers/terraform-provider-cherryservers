package provider

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/cherryservers/cherrygo/v4"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerReinstallModifierReturnsErrorsWhenNotAllowed(t *testing.T) {
	cases := []struct {
		name        string
		state, plan serverResourceModel
	}{
		{
			name: "known values > known values",
			state: newServerModel(func(m *serverResourceModel) {
				m.Image = types.StringValue("test")
				m.OSPartitionSize = types.Int64Value(1)
				m.SSHKeyIds = types.SetValueMust(types.StringType, []attr.Value{types.StringValue("1")})
				m.UserData = types.StringValue("test")
				m.IPXE = types.StringValue("test")
				m.PersistIPXE = types.BoolValue(false)
			}),
			plan: newServerModel(func(m *serverResourceModel) {
				m.Image = types.StringValue("test-new")
				m.OSPartitionSize = types.Int64Value(2)
				m.SSHKeyIds = types.SetValueMust(types.StringType, []attr.Value{types.StringValue("2")})
				m.UserData = types.StringValue("test-new")
				m.IPXE = types.StringValue("test-new")
				m.PersistIPXE = types.BoolValue(true)
			}),
		},
		{
			name:  "null values > unknown values",
			state: newServerModel(nil),
			plan: newServerModel(func(m *serverResourceModel) {
				m.Image = types.StringUnknown()
				m.OSPartitionSize = types.Int64Unknown()
				m.SSHKeyIds = types.SetUnknown(types.StringType)
				m.UserData = types.StringUnknown()
				m.IPXE = types.StringUnknown()
				m.PersistIPXE = types.BoolUnknown()
			}),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := t.Context()

			r := serverResource{client: &cherrygo.Client{}}
			s := resourceSchema(t, &r)

			req := newModifyReqFromModels(t, newServerModel(nil), tc.state, tc.plan, s)
			resp := newModifyResp(s)

			r.ModifyPlan(ctx, req, &resp)
			d := resp.Diagnostics

			attrs := []string{"image", "os_partition_size", "ssh_key_ids", "user_data", "ipxe", "persist_ipxe"}
			for _, attr := range attrs {
				assert.Contains(
					t, d, diag.NewAttributeErrorDiagnostic(
						path.Root(attr),
						"Re-installation not allowed.",
						fmt.Sprintf("Updating `%s` requires `allow_reinstall` to be enabled.", attr),
					),
				)
			}

			assert.Equal(t, 6, d.ErrorsCount())
		})
	}
}

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
		s := resourceSchema(t, &r)

		req := resource.ModifyPlanRequest{
			Plan:   planFromModelMust(t, tc.plan, s),
			Config: configFromModelMust(t, tc.plan, s),
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
			s := resourceSchema(t, &r)

			req := newModifyReqFromModels(
				t, newServerModel(nil), newServerModel(nil), tc.plan, s,
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
			s := resourceSchema(t, &r)

			req := newModifyReqFromModels(
				t, newServerModel(nil), newServerModel(nil), tc.plan, s,
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
	s := resourceSchema(t, &r)

	state := newServerModel(func(m *serverResourceModel) {
		m.Image = types.StringValue("not-ipxe")
	})
	plan := newServerModel(func(m *serverResourceModel) {
		m.AllowReinstall = types.BoolValue(true)
		m.IPXE = types.StringValue("test")
		m.SSHKeyIds = types.SetValueMust(types.StringType, []attr.Value{types.StringValue("1")})
	})
	req := newModifyReqFromModels(t, newServerModel(nil), state, plan, s)
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
	s := resourceSchema(t, &r)

	state := newServerModel(func(m *serverResourceModel) {
		m.Image = types.StringValue("not-ipxe")
	})
	plan := newServerModel(func(m *serverResourceModel) {
		m.AllowReinstall = types.BoolValue(true)
		m.Image = types.StringValue("another-not-ipxe")
		m.SSHKeyIds = types.SetValueMust(types.StringType, []attr.Value{types.StringValue("1")})
	})
	req := newModifyReqFromModels(t, newServerModel(nil), state, plan, s)
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
		s := resourceSchema(t, &r)

		state := newServerModel(func(m *serverResourceModel) {
			m.Image = types.StringValue("custom_ipxe_install")
		})
		plan := newServerModel(func(m *serverResourceModel) {
			m.AllowReinstall = types.BoolValue(true)
			m.Plan = types.StringValue("test-plan")
		})
		req := newModifyReqFromModels(t, newServerModel(nil), state, plan, s)
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
		s := resourceSchema(t, &r)

		state := newServerModel(func(m *serverResourceModel) {
			m.Image = types.StringValue("custom_ipxe_install")
		})
		plan := newServerModel(func(m *serverResourceModel) {
			m.AllowReinstall = types.BoolValue(true)
			m.Plan = types.StringValue("test-plan")
		})
		req := newModifyReqFromModels(t, newServerModel(nil), state, plan, s)
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
	s := resourceSchema(t, &r)

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
	req := newModifyReqFromModels(t, config, state, plan, s)
	resp := newModifyResp(s)

	r.ModifyPlan(t.Context(), req, &resp)
	require.Empty(t, resp.Diagnostics)

	var gotPlan serverResourceModel
	require.Empty(t, resp.Plan.Get(t.Context(), &gotPlan))
	assert.Equal(t, "test-img", gotPlan.Image.ValueString())
}

func TestServerModifyPlanReturnsErrorWithInvalidRequests(t *testing.T) {
	r := serverResource{client: &cherrygo.Client{}}
	s := resourceSchema(t, &r)
	okVal := planFromModelMust(t, newServerModel(nil), s).Raw

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
	s := resourceSchema(t, &r)

	var req resource.ModifyPlanRequest
	resp := newModifyResp(s)

	r.ModifyPlan(t.Context(), req, &resp)
	require.Empty(t, resp.Diagnostics)
	assert.True(t, resp.Plan.Raw.IsNull())
}

func TestServerModifyPlanSetsIPXEImageWhenIPXEIsPlanned(t *testing.T) {
	r := serverResource{client: &cherrygo.Client{}}
	s := resourceSchema(t, &r)
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
			state: stateFromModelMust(t, newServerModel(nil), s),
		},
	}

	for _, tc := range cases {
		req := resource.ModifyPlanRequest{
			State:  tc.state,
			Config: configFromModelMust(t, newServerModel(nil), s),
			Plan:   planFromModelMust(t, planModel, s),
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
	s := resourceSchema(t, &r)
	planModel := newServerModel(func(m *serverResourceModel) {
		m.Hostname = types.StringValue("test")
	})
	stateModel := newServerModel(func(m *serverResourceModel) {
		m.Hostname = types.StringValue("old")
	})
	req := newModifyReqFromModels(t, planModel, stateModel, planModel, s)
	resp := newModifyResp(s)

	r.ModifyPlan(t.Context(), req, &resp)
	require.Empty(t, resp.Diagnostics)

	var gotPlan serverResourceModel
	require.Empty(t, resp.Plan.Get(t.Context(), &gotPlan))

	assert.Equal(t, planModel, gotPlan)
}

func configFromModelMust(t *testing.T, m any, s schema.Schema) tfsdk.Config {
	t.Helper()

	// tfsdk.Config doesn't have Set, due to its immutable nature,
	// so we need this workaround.
	plan := tfsdk.Plan{
		Schema: s,
	}
	diags := plan.Set(t.Context(), m)
	require.Empty(t, diags)

	return tfsdk.Config{Raw: plan.Raw, Schema: s}
}

func planFromModelMust(t *testing.T, m any, s schema.Schema) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{
		Schema: s,
	}
	diags := plan.Set(t.Context(), m)
	require.Empty(t, diags)

	return plan
}

func stateFromModelMust(t *testing.T, m any, s schema.Schema) tfsdk.State {
	t.Helper()

	state := tfsdk.State{
		Schema: s,
	}
	diags := state.Set(t.Context(), m)
	require.Empty(t, diags)

	return state
}

func resourceSchema(t *testing.T, r resource.Resource) schema.Schema {
	t.Helper()

	var resp resource.SchemaResponse
	r.Schema(t.Context(), resource.SchemaRequest{}, &resp)
	return resp.Schema
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

func newModifyReqFromModels[T serverResourceModel](t *testing.T, config, state, plan T, s schema.Schema) resource.ModifyPlanRequest {
	return resource.ModifyPlanRequest{
		Config: configFromModelMust(t, config, s),
		Plan:   planFromModelMust(t, plan, s),
		State:  stateFromModelMust(t, state, s),
	}
}

func newModifyResp(s schema.Schema) resource.ModifyPlanResponse {
	return resource.ModifyPlanResponse{
		Plan: tfsdk.Plan{
			Schema: s,
		},
	}
}
