package provider

import (
	"context"
	"net"
	"regexp"
	"testing"

	"github.com/cherryservers/cherrygo/v4"
	"github.com/cherryservers/terraform-provider-cherryservers/internal/provider/fakes"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func configFromResourceModelMust(t *testing.T, m any, r resource.Resource) tfsdk.Config {
	t.Helper()

	// tfsdk.Config doesn't have Set, due to its immutable nature,
	// so we need this workaround.
	plan := tfsdk.Plan{
		Schema: schemaFromResourceMust(t, r),
	}
	diags := plan.Set(t.Context(), m)
	require.Empty(t, diags)

	return tfsdk.Config(plan)
}

func planFromResourceModelMust(t *testing.T, m any, r resource.Resource) tfsdk.Plan {
	t.Helper()

	plan := tfsdk.Plan{
		Schema: schemaFromResourceMust(t, r),
	}
	diags := plan.Set(t.Context(), m)
	require.Empty(t, diags)

	return plan
}

func stateFromResourceModelMust(t *testing.T, m any, r resource.Resource) tfsdk.State {
	t.Helper()

	state := tfsdk.State{
		Schema: schemaFromResourceMust(t, r),
	}

	diags := state.Set(t.Context(), &m)
	require.Empty(t, diags)

	return state
}

func schemaFromResourceMust(t *testing.T, r resource.Resource) schema.Schema {
	t.Helper()

	var (
		req  resource.SchemaRequest
		resp resource.SchemaResponse
	)
	r.Schema(t.Context(), req, &resp)
	require.Empty(t, resp.Diagnostics)

	return resp.Schema
}

func serverGetNetError(_ context.Context, id int, _ *cherrygo.GetOptions) (cherrygo.Server, *cherrygo.Response, error) {
	return cherrygo.Server{}, nil, net.UnknownNetworkError("test")
}

func ipGetNetError(_ context.Context, id string, _ *cherrygo.GetOptions) (cherrygo.IPAddress, *cherrygo.Response, error) {
	return cherrygo.IPAddress{}, nil, net.UnknownNetworkError("test")
}

func projectGetNetError(_ context.Context, id int, _ *cherrygo.GetOptions) (cherrygo.Project, *cherrygo.Response, error) {
	return cherrygo.Project{}, nil, net.UnknownNetworkError("test")
}

func sshkeyGetNetError(_ context.Context, id int, _ *cherrygo.GetOptions) (cherrygo.SSHKey, *cherrygo.Response, error) {
	return cherrygo.SSHKey{}, nil, net.UnknownNetworkError("test")
}

func TestResourceReadsHandleNetErrors(t *testing.T) {
	cases := []struct {
		name     string
		resource resource.Resource
		state    tfsdk.State
	}{
		{
			name: "server",
			resource: &serverResource{client: &cherrygo.Client{
				Servers: &fakes.ServersService{GetFunc: serverGetNetError},
			}},
			state: stateFromResourceModelMust(
				t, newServerModel(func(m *serverResourceModel) { m.Id = types.StringValue("1") }), new(serverResource),
			),
		},
		{
			name: "ip",
			resource: &ipResource{client: &cherrygo.Client{
				IPAddresses: &fakes.IPAddressesService{GetFunc: ipGetNetError},
			}},
			state: stateFromResourceModelMust(
				t, newIPModel(func(m *ipResourceModel) { m.Id = types.StringValue("1") }), new(ipResource),
			),
		},
		{
			name: "project",
			resource: &projectResource{client: &cherrygo.Client{
				Projects: &fakes.ProjectsService{GetFunc: projectGetNetError},
			}},
			state: stateFromResourceModelMust(
				t, newProjectModel(func(m *projectResourceModel) { m.Id = types.StringValue("1") }), new(projectResource),
			),
		},
		{
			name: "ssh-key",
			resource: &sshKeyResource{client: &cherrygo.Client{
				SSHKeys: &fakes.SSHKeysService{GetFunc: sshkeyGetNetError},
			}},
			state: stateFromResourceModelMust(
				t, newSSHKeyModel(func(m *sshKeyResourceModel) { m.ID = types.StringValue("1") }), new(sshKeyResource),
			),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := resource.ReadRequest{State: tc.state}
			resp := resource.ReadResponse{State: req.State}

			tc.resource.Read(t.Context(), req, &resp)

			require.Len(t, resp.Diagnostics, 1)
			assert.Regexp(t, regexp.MustCompile(`unable to read`), resp.Diagnostics[0].Summary())
			assert.Equal(t, req.State, resp.State)
		})
	}
}
