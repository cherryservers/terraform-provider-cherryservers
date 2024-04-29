package provider

import (
	"fmt"
	"github.com/cherryservers/cherrygo/v3"
	"runtime/debug"
)

const terraformSDKPath string = "github.com/hashicorp/terraform-plugin-framework"

// Config for auth variable
type Config struct {
	AuthKey string
}

// Client wrap cherrygo client
type Client struct {
	client *cherrygo.Client
}

func (c *Client) cherrygoClient() *cherrygo.Client {
	return c.client
}

// terraformFrameworkVersion looks up the module version of the Terraform Framework for use
// in the User Agent client string
func terraformFrameworkVersion() string {
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
func (c *Config) Client(version string) (*Client, error) {
	userAgent := fmt.Sprintf("terraform-provider/cherryservers/%s terraform/%s", version, terraformFrameworkVersion())
	args := []cherrygo.ClientOpt{cherrygo.WithAuthToken(c.AuthKey), cherrygo.WithUserAgent(userAgent)}
	cherryClient, _ := cherrygo.NewClient(args...)

	return &Client{client: cherryClient}, nil
}
