package fakes

import (
	"context"

	"github.com/cherryservers/cherrygo/v4"
)

var _ cherrygo.IPAddressesService = (*IPAddressesService)(nil)

type IPAddressesService struct {
	GetFunc func(ctx context.Context, id string, opts *cherrygo.GetOptions) (cherrygo.IPAddress, *cherrygo.Response, error)
}

// Assign implements [cherrygo.IPAddressesService].
func (i *IPAddressesService) Assign(ctx context.Context, ipID string, request *cherrygo.AssignIPAddress) (cherrygo.IPAddress, *cherrygo.Response, error) {
	panic("unimplemented")
}

// Create implements [cherrygo.IPAddressesService].
func (i *IPAddressesService) Create(ctx context.Context, projectID int, request *cherrygo.CreateIPAddress) (cherrygo.IPAddress, *cherrygo.Response, error) {
	panic("unimplemented")
}

// Get implements [cherrygo.IPAddressesService].
func (i *IPAddressesService) Get(ctx context.Context, ipID string, opts *cherrygo.GetOptions) (cherrygo.IPAddress, *cherrygo.Response, error) {
	return i.GetFunc(ctx, ipID, opts)
}

// List implements [cherrygo.IPAddressesService].
func (i *IPAddressesService) List(ctx context.Context, projectID int, opts *cherrygo.GetOptions) ([]cherrygo.IPAddress, *cherrygo.Response, error) {
	panic("unimplemented")
}

// Remove implements [cherrygo.IPAddressesService].
func (i *IPAddressesService) Remove(ctx context.Context, ipID string) (*cherrygo.Response, error) {
	panic("unimplemented")
}

// Unassign implements [cherrygo.IPAddressesService].
func (i *IPAddressesService) Unassign(ctx context.Context, ipID string) (*cherrygo.Response, error) {
	panic("unimplemented")
}

// Update implements [cherrygo.IPAddressesService].
func (i *IPAddressesService) Update(ctx context.Context, ipID string, request *cherrygo.UpdateIPAddress) (cherrygo.IPAddress, *cherrygo.Response, error) {
	panic("unimplemented")
}
