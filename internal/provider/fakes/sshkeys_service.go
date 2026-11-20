package fakes

import (
	"context"

	"github.com/cherryservers/cherrygo/v4"
)

var _ cherrygo.SSHKeysService = (*SSHKeysService)(nil)

type SSHKeysService struct {
	GetFunc func(ctx context.Context, id int, opts *cherrygo.GetOptions) (cherrygo.SSHKey, *cherrygo.Response, error)
}

// Create implements [cherrygo.SSHKeysService].
func (s *SSHKeysService) Create(ctx context.Context, request *cherrygo.CreateSSHKey) (cherrygo.SSHKey, *cherrygo.Response, error) {
	panic("unimplemented")
}

// Delete implements [cherrygo.SSHKeysService].
func (s *SSHKeysService) Delete(ctx context.Context, sshKeyID int) (*cherrygo.Response, error) {
	panic("unimplemented")
}

// Get implements [cherrygo.SSHKeysService].
func (s *SSHKeysService) Get(ctx context.Context, sshKeyID int, opts *cherrygo.GetOptions) (cherrygo.SSHKey, *cherrygo.Response, error) {
	return s.GetFunc(ctx, sshKeyID, opts)
}

// List implements [cherrygo.SSHKeysService].
func (s *SSHKeysService) List(ctx context.Context, opts *cherrygo.GetOptions) ([]cherrygo.SSHKey, *cherrygo.Response, error) {
	panic("unimplemented")
}

// Update implements [cherrygo.SSHKeysService].
func (s *SSHKeysService) Update(ctx context.Context, sshKeyID int, request *cherrygo.UpdateSSHKey) (cherrygo.SSHKey, *cherrygo.Response, error) {
	panic("unimplemented")
}
