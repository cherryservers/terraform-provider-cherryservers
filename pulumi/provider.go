// Package pulumi exposes provider internals for pulumi bridging.
package pulumi

import (
	"github.com/cherryservers/terraform-provider-cherryservers/internal/provider"
	tfprovider "github.com/hashicorp/terraform-plugin-framework/provider"
)

// New exports the provider factory for the pulumi bridge.
func New(version string) func() tfprovider.Provider {
	return provider.New(version)
}
