package main

import (
	"github.com/cherryservers/terraform-provider-cherryservers/cherryservers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: cherryservers.Provider,
	})
}
