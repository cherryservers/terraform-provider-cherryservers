package datasourcebase

import (
	"context"
	"fmt"

	"github.com/cherryservers/cherrygo/v3"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

type Configurator struct {
	client *cherrygo.Client
}

func (c *Configurator) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*cherrygo.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *cherrygo.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	c.client = client
}

func (c *Configurator) Client() *cherrygo.Client {
	return c.client
}
