package main

import (
	"terraform-provider-cherryservers/cherryservers"

	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
	"github.com/mm0/terraform-provider-cherryservers/cherryservers"
	//"./cherryservers"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return cherryservers.Provider()
		},
	})
}
