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
func (c *Config) Client(terraformVersion string) (*cherrygo.Client, error) {
	userAgent := fmt.Sprintf("terraform-provider/cherryservers/%s terraform/%s", terraformVersion, terraformFrameworkVersion())
	args := []cherrygo.ClientOpt{cherrygo.WithAuthToken(c.AuthKey), cherrygo.WithUserAgent(userAgent)}
	cherryClient, _ := cherrygo.NewClient(args...)

	return cherryClient, nil
}
