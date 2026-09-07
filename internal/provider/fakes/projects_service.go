package fakes

import (
	"context"
	"github.com/cherryservers/cherrygo/v4"
)

var _ cherrygo.ProjectsService = (*ProjectsService)(nil)

type ProjectsService struct {
	GetFunc func(ctx context.Context, id int, opts *cherrygo.GetOptions) (cherrygo.Project, *cherrygo.Response, error)
}

// Create implements [cherrygo.ProjectsService].
func (p *ProjectsService) Create(ctx context.Context, teamID int, request *cherrygo.CreateProject) (cherrygo.Project, *cherrygo.Response, error) {
	panic("unimplemented")
}

// Delete implements [cherrygo.ProjectsService].
func (p *ProjectsService) Delete(ctx context.Context, projectID int) (*cherrygo.Response, error) {
	panic("unimplemented")
}

// Get implements [cherrygo.ProjectsService].
func (p *ProjectsService) Get(ctx context.Context, projectID int, opts *cherrygo.GetOptions) (cherrygo.Project, *cherrygo.Response, error) {
	return p.GetFunc(ctx, projectID, opts)
}

// List implements [cherrygo.ProjectsService].
func (p *ProjectsService) List(ctx context.Context, teamID int, opts *cherrygo.GetOptions) ([]cherrygo.Project, *cherrygo.Response, error) {
	panic("unimplemented")
}

// ListSSHKeys implements [cherrygo.ProjectsService].
func (p *ProjectsService) ListSSHKeys(ctx context.Context, projectID int, opts *cherrygo.GetOptions) ([]cherrygo.SSHKey, *cherrygo.Response, error) {
	panic("unimplemented")
}

// Update implements [cherrygo.ProjectsService].
func (p *ProjectsService) Update(ctx context.Context, projectID int, request *cherrygo.UpdateProject) (cherrygo.Project, *cherrygo.Response, error) {
	panic("unimplemented")
}
