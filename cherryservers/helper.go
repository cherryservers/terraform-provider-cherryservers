package cherryservers

import (
	"github.com/cherryservers/cherrygo"

	"github.com/hashicorp/go-cleanhttp"
)

// Config for auth var
type Config struct {
	AuthToken string
}

// Client return client
func (c *Config) Client() *cherrygo.Client {
	client := cleanhttp.DefaultClient()

	return cherrygo.NewClientWithAuthVar(client, c.AuthToken)
}
