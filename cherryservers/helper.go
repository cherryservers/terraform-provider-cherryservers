package cherryservers

import (
	"github.com/cherryservers/cherrygo"

	"github.com/hashicorp/go-cleanhttp"
)

// Config for auth variable
type Config struct {
	AuthToken string
}

// Client return client
func (c *CombinedConfig) Client() *cherrygo.Client {
	return c.client
}

// CombinedConfig including client and auth_token
type CombinedConfig struct {
	client    *cherrygo.Client
	AuthToken string
}

// Client returns CombinedConfig
func (c *Config) Client() (*CombinedConfig, error) {
	client := cleanhttp.DefaultClient()
	cherryClient := cherrygo.NewClientWithAuthVar(client, c.AuthToken)
	return &CombinedConfig{
		client:    cherryClient,
		AuthToken: c.AuthToken}, nil
}
