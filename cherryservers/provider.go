package cherryservers

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CHERRY_AUTH_TOKEN", nil),
				Description: "Cherry Servers [API Token](https://portal.cherryservers.com/settings/api-keys) that allows interactions with the API",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"cherryservers_server":  dataSourceCherryServersServer(),
			"cherryservers_ip":      dataSourceCherryServersIP(),
			"cherryservers_project": dataSourceCherryServersProject(),
			"cherryservers_ssh_key": dataSourceCherryServersSSHKey(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"cherryservers_server":  resourceCherryServersServer(),
			"cherryservers_ip":      resourceCherryServersIP(),
			"cherryservers_project": resourceCherryServersProject(),
			"cherryservers_ssh_key": resourceCherryServersSSHKey(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{}
	if apiKey, exist := d.GetOk("api_token"); exist {
		config.AuthToken = apiKey.(string)
	}

	return config.Client()
}
