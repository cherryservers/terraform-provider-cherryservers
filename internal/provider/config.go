package provider

import (
	"fmt"
	"github.com/cherryservers/cherrygo/v3"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func DefaultClientConfigure(req resource.ConfigureRequest, resp *resource.ConfigureResponse) *cherrygo.Client {
	client, ok := req.ProviderData.(*cherrygo.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *cherrygo.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
	}

	return client
}
