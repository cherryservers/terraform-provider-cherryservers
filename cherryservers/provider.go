package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

// Provider init
func Provider() *schema.Provider {
	return &schema.Provider{
    Schema: map[string]*schema.Schema{
			"token": {
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"CHERRY_AUTH_TOKEN"
				}, nil),
				Description: "The token key for API operations.",
    },
		ResourcesMap: map[string]*schema.Resource{
			"cherryservers_server":  resourceServer(),
			"cherryservers_ssh":     resourceSSHKey(),
			"cherryservers_ip":      resourceIP(),
			"cherryservers_project": resourceProject(),
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Token:       d.Get("token").(string)
	}

	return config.Client()
}

