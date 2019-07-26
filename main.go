package main

import (
	"terraform-provider-cherryservers/cherryservers"

	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return cherryservers.Provider()
		},
	})
}
