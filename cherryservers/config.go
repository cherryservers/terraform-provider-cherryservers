package cherryservers

import (
	"github.com/cherryservers/cherrygo"
)

type Config struct {
	Token string
}
type CombinedConfig struct {
	client *cherrygo.Client
	token  string
}

func (c *CombinedConfig) Client() *cherrygo.Client { return c.client }

func (c *Config) Client() (*CombinedConfig, error) {
	cherryClient, err := cherrygo.NewClient()
	if err != nil {
		return nil, err
	}
	return &CombinedConfig{
		client: cherryClient,
	}, nil
}
