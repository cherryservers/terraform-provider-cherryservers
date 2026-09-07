package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
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
