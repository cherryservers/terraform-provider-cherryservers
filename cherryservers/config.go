package cherryservers

import (
	"fmt"
	"runtime/debug"

	"github.com/cherryservers/cherrygo/v3"
)

const terraformSDKPath string = "github.com/hashicorp/terraform-plugin-sdk/v2"

// Config for auth variable
type Config struct {
	AuthToken string
}

// Client wrap cherrygo client
type Client struct {
	client *cherrygo.Client
}

func (c *Client) cherrygoClient() *cherrygo.Client {
	return c.client
}

// terraformSDKVersion looks up the module version of the Terraform SDK for use
// in the User Agent client string
func terraformSDKVersion() string {
	i, ok := debug.ReadBuildInfo()
	if !ok {
		return "0.0.0"
	}

	for _, module := range i.Deps {
		if module.Path == terraformSDKPath {
			return module.Version
		}
	}

	return "0.0.0"
}

// Client initialize cherrygo client
func (c *Config) Client() (*Client, error) {
	userAgent := fmt.Sprintf("terraform-provider/cherryservers/%s terraform/%s", version, terraformSDKVersion())
	args := []cherrygo.ClientOpt{cherrygo.WithAuthToken(c.AuthToken), cherrygo.WithUserAgent(userAgent)}
	cherryClient, _ := cherrygo.NewClient(args...)

	return &Client{client: cherryClient}, nil
}
