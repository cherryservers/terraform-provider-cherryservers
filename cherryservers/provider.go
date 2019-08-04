package cherryservers

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider init
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CHERRY_AUTH_TOKEN", nil),
				Description: "The API token",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"cherryservers_server":  resourceServer(),
			"cherryservers_ssh":     resourceSSHKey(),
			"cherryservers_ip":      resourceIP(),
			"cherryservers_project": resourceProject(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		AuthToken: d.Get("auth_token").(string),
	}
	return config.Client(), nil
}
