package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

// Provider init
func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"cherryservers_server": resourceServer(),
			"cherryservers_ssh":    resourceSSHKey(),
			"cherryservers_ip":     resourceIP(),
		},
	}
}
