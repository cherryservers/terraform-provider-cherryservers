package fakes

import (
	"context"

	"github.com/cherryservers/cherrygo/v4"
)

var _ cherrygo.ServersService = (*ServersService)(nil)

type ServersService struct {
	GetFunc func(ctx context.Context, serverID int, opts *cherrygo.GetOptions) (cherrygo.Server, *cherrygo.Response, error)
}

// AllowBMCAccess implements [cherrygo.ServersService].
func (f *ServersService) AllowBMCAccess(ctx context.Context, serverID int, ip4 string) (cherrygo.Server, *cherrygo.Response, error) {
	panic("unimplemented")
}

// Create implements [cherrygo.ServersService].
func (f *ServersService) Create(ctx context.Context, request *cherrygo.CreateServer) (cherrygo.Server, *cherrygo.Response, error) {
	panic("unimplemented")
}

// Delete implements [cherrygo.ServersService].
func (f *ServersService) Delete(ctx context.Context, serverID int) (*cherrygo.Response, error) {
	panic("unimplemented")
}

// EnterRescueMode implements [cherrygo.ServersService].
func (f *ServersService) EnterRescueMode(ctx context.Context, serverID int, fields *cherrygo.RescueServerFields) (cherrygo.Server, *cherrygo.Response, error) {
	panic("unimplemented")
}

// ExitRescueMode implements [cherrygo.ServersService].
func (f *ServersService) ExitRescueMode(ctx context.Context, serverID int) (cherrygo.Server, *cherrygo.Response, error) {
	panic("unimplemented")
}

// Get implements [cherrygo.ServersService].
func (f *ServersService) Get(ctx context.Context, serverID int, opts *cherrygo.GetOptions) (cherrygo.Server, *cherrygo.Response, error) {
	return f.GetFunc(ctx, serverID, opts)
}

// List implements [cherrygo.ServersService].
func (f *ServersService) List(ctx context.Context, projectID int, opts *cherrygo.GetOptions) ([]cherrygo.Server, *cherrygo.Response, error) {
	panic("unimplemented")
}

// ListCycles implements [cherrygo.ServersService].
func (f *ServersService) ListCycles(ctx context.Context, opts *cherrygo.GetOptions) ([]cherrygo.ServerCycle, *cherrygo.Response, error) {
	panic("unimplemented")
}

// ListSSHKeys implements [cherrygo.ServersService].
func (f *ServersService) ListSSHKeys(ctx context.Context, serverID int, opts *cherrygo.GetOptions) ([]cherrygo.SSHKey, *cherrygo.Response, error) {
	panic("unimplemented")
}

// PowerOff implements [cherrygo.ServersService].
func (f *ServersService) PowerOff(ctx context.Context, serverID int) (cherrygo.Server, *cherrygo.Response, error) {
	panic("unimplemented")
}

// PowerOn implements [cherrygo.ServersService].
func (f *ServersService) PowerOn(ctx context.Context, serverID int) (cherrygo.Server, *cherrygo.Response, error) {
	panic("unimplemented")
}

// PowerState implements [cherrygo.ServersService].
func (f *ServersService) PowerState(ctx context.Context, serverID int) (cherrygo.PowerState, *cherrygo.Response, error) {
	panic("unimplemented")
}

// Reboot implements [cherrygo.ServersService].
func (f *ServersService) Reboot(ctx context.Context, serverID int) (cherrygo.Server, *cherrygo.Response, error) {
	panic("unimplemented")
}

// Reinstall implements [cherrygo.ServersService].
func (f *ServersService) Reinstall(ctx context.Context, serverID int, fields *cherrygo.ReinstallServerFields) (cherrygo.Server, *cherrygo.Response, error) {
	panic("unimplemented")
}

// ResetBMCPassword implements [cherrygo.ServersService].
func (f *ServersService) ResetBMCPassword(ctx context.Context, serverID int) (cherrygo.Server, *cherrygo.Response, error) {
	panic("unimplemented")
}

// Update implements [cherrygo.ServersService].
func (f *ServersService) Update(ctx context.Context, serverID int, request *cherrygo.UpdateServer) (cherrygo.Server, *cherrygo.Response, error) {
	panic("unimplemented")
}

// Upgrade implements [cherrygo.ServersService].
func (f *ServersService) Upgrade(ctx context.Context, serverID int, plan string) (cherrygo.Server, *cherrygo.Response, error) {
	panic("unimplemented")
}

// WaitForStatus implements [cherrygo.ServersService].
func (f *ServersService) WaitForStatus(ctx context.Context, serverID int, status cherrygo.ServerStatus) (cherrygo.Server, *cherrygo.Response, error) {
	panic("unimplemented")
}
