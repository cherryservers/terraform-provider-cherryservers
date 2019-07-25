package cherryservers

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider init
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CHERRY_AUTH_TOKEN", nil),
				Description: "The token key for API operations.",
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
		Token: d.Get("token").(string),
	}
	return config.Client()
}
