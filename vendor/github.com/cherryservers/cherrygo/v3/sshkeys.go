package cherrygo

import "fmt"

const baseSSHPath = "/v1/ssh-keys"

// SSHKeysService is an interface for interfacing with the the SSH keys endpoints of the CherryServers API
// See: https://api.cherryservers.com/doc/#tag/SshKeys
type SSHKeysService interface {
	List(opts *GetOptions) ([]SSHKey, *Response, error)
	Get(sshKeyID int, opts *GetOptions) (SSHKey, *Response, error)
	Create(request *CreateSSHKey) (SSHKey, *Response, error)
	Delete(sshKeyID int) (SSHKey, *Response, error)
	Update(sshKeyID int, request *UpdateSSHKey) (SSHKey, *Response, error)
}

// SSHKeys fields for return values after creation
type SSHKey struct {
	ID          int    `json:"id,omitempty"`
	Label       string `json:"label,omitempty"`
	Key         string `json:"key,omitempty"`
	Fingerprint string `json:"fingerprint,omitempty"`
	User        User   `json:"user,omitempty"`
	Updated     string `json:"updated,omitempty"`
	Created     string `json:"created,omitempty"`
	Href        string `json:"href,omitempty"`
}

type SSHKeysClient struct {
	client *Client
}

// CreateSSHKey fields for adding new key with label and raw key
type CreateSSHKey struct {
	Label string `json:"label"`
	Key   string `json:"key"`
}

// UpdateSSHKey fields for label or key update
type UpdateSSHKey struct {
	Label *string `json:"label,omitempty"`
	Key   *string `json:"key,omitempty"`
}

// List all available ssh keys
func (s *SSHKeysClient) List(opts *GetOptions) ([]SSHKey, *Response, error) {
	var trans []SSHKey

	pathQuery := opts.WithQuery(baseSSHPath)
	resp, err := s.client.MakeRequest("GET", pathQuery, nil, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

func (s *SSHKeysClient) Get(sshKeyID int, opts *GetOptions) (SSHKey, *Response, error) {
	var trans SSHKey

	path := opts.WithQuery(fmt.Sprintf("%s/%d", baseSSHPath, sshKeyID))

	resp, err := s.client.MakeRequest("GET", path, nil, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

// Create adds new SSH key
func (s *SSHKeysClient) Create(request *CreateSSHKey) (SSHKey, *Response, error) {
	var trans SSHKey

	resp, err := s.client.MakeRequest("POST", baseSSHPath, request, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

// Delete removes desired SSH key by its ID
func (s *SSHKeysClient) Delete(sshKeyID int) (SSHKey, *Response, error) {
	var trans SSHKey

	path := fmt.Sprintf("%s/%d", baseSSHPath, sshKeyID)

	resp, err := s.client.MakeRequest("DELETE", path, nil, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

// Update function updates keys Label or key itself
func (s *SSHKeysClient) Update(sshKeyID int, request *UpdateSSHKey) (SSHKey, *Response, error) {
	var trans SSHKey

	path := fmt.Sprintf("%s/%d", baseSSHPath, sshKeyID)

	resp, err := s.client.MakeRequest("PUT", path, request, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}
