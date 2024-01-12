package cherrygo

import (
	"fmt"
)

const baseProjectPath = "/v1/projects"

// ProjectsService is an interface for interfacing with the Projects endpoints of the CherryServers API
// See: https://api.cherryservers.com/doc/#tag/Projects
type ProjectsService interface {
	List(teamID int, opts *GetOptions) ([]Project, *Response, error)
	Get(projectID int, opts *GetOptions) (Project, *Response, error)
	Create(teamID int, request *CreateProject) (Project, *Response, error)
	Update(projectID int, request *UpdateProject) (Project, *Response, error)
	ListSSHKeys(projectID int, opts *GetOptions) ([]SSHKey, *Response, error)
	Delete(projectID int) (*Response, error)
}

type Project struct {
	ID   int        `json:"id,omitempty"`
	Name string     `json:"name,omitempty"`
	Bgp  ProjectBGP `json:"bgp,omitempty"`
	Href string     `json:"href,omitempty"`
}

// CreateProject fields for adding new project with specified name
type CreateProject struct {
	Name string `json:"name,omitempty"`
	Bgp  bool   `json:"bgp,omitempty"`
}

// UpdateProject fields for updating a project with specified name
type UpdateProject struct {
	Name *string `json:"name,omitempty"`
	Bgp  *bool   `json:"bgp,omitempty"`
}

type ProjectsClient struct {
	client *Client
}

// List func lists projects
func (p *ProjectsClient) List(teamID int, opts *GetOptions) ([]Project, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("/v1/teams/%d/projects", teamID))

	var trans []Project

	resp, err := p.client.MakeRequest("GET", path, nil, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

func (p *ProjectsClient) Get(projectID int, opts *GetOptions) (Project, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("%s/%d", baseProjectPath, projectID))

	var trans Project

	resp, err := p.client.MakeRequest("GET", path, nil, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

// Create func will create new Project for specified team
func (p *ProjectsClient) Create(teamID int, request *CreateProject) (Project, *Response, error) {
	var trans Project

	path := fmt.Sprintf("/v1/teams/%d/projects", teamID)

	resp, err := p.client.MakeRequest("POST", path, request, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

// Update func will update a project
func (p *ProjectsClient) Update(projectID int, request *UpdateProject) (Project, *Response, error) {
	var trans Project

	path := fmt.Sprintf("%s/%d", baseProjectPath, projectID)

	resp, err := p.client.MakeRequest("PUT", path, request, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

// Delete func will delete a project
func (p *ProjectsClient) Delete(projectID int) (*Response, error) {
	path := fmt.Sprintf("%s/%d", baseProjectPath, projectID)

	resp, err := p.client.MakeRequest("DELETE", path, nil, nil)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return resp, err
}

func (p *ProjectsClient) ListSSHKeys(projectID int, opts *GetOptions) ([]SSHKey, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("/v1/projects/%d/ssh-keys", projectID))

	var trans []SSHKey

	resp, err := p.client.MakeRequest("GET", path, nil, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}
