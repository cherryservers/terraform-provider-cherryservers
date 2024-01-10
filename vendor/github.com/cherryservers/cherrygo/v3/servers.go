package cherrygo

import (
	"fmt"
)

const baseServerPath = "/v1/servers"
const endServersPath = "servers"

// ServersService is an interface for interfacing with the the Server endpoints of the CherryServers API
// See: https://api.cherryservers.com/doc/#tag/Servers
type ServersService interface {
	List(projectID int, opts *GetOptions) ([]Server, *Response, error)
	Get(serverID int, opts *GetOptions) (Server, *Response, error)
	PowerOff(serverID int) (Server, *Response, error)
	PowerOn(serverID int) (Server, *Response, error)
	Create(request *CreateServer) (Server, *Response, error)
	Delete(serverID int) (Server, *Response, error)
	PowerState(serverID int) (PowerState, *Response, error)
	Reboot(serverID int) (Server, *Response, error)
	Update(serverID int, request *UpdateServer) (Server, *Response, error)
	Reinstall(serverID int, fields *ReinstallServerFields) (Server, *Response, error)
	ListSSHKeys(serverID int, opts *GetOptions) ([]SSHKey, *Response, error)
}

// Server response object
type Server struct {
	ID               int               `json:"id,omitempty"`
	Name             string            `json:"name,omitempty"`
	Href             string            `json:"href,omitempty"`
	Hostname         string            `json:"hostname,omitempty"`
	Image            string            `json:"image,omitempty"`
	SpotInstance     bool              `json:"spot_instance"`
	BGP              ServerBGP         `json:"bgp,omitempty"`
	Project          Project           `json:"project,omitempty"`
	Region           Region            `json:"region,omitempty"`
	State            string            `json:"state,omitempty"`
	Plan             Plan              `json:"plan,omitempty"`
	AvailableRegions AvailableRegions  `json:"availableregions,omitempty"`
	Pricing          Pricing           `json:"pricing,omitempty"`
	IPAddresses      []IPAddress       `json:"ip_addresses,omitempty"`
	SSHKeys          []SSHKey          `json:"ssh_keys,omitempty"`
	Tags             map[string]string `json:"tags,omitempty"`
	Storage          BlockStorage      `json:"storage,omitempty"`
	Backup           BackupStorage     `json:"backup_storage,omitempty"`
	Created          string            `json:"created_at,omitempty"`
	TerminationDate  string            `json:"termination_date,omitempty"`
}

type ReinstallServer struct {
	ServerAction
	*ReinstallServerFields
}

type ReinstallServerFields struct {
	Image           string   `json:"image"`
	Hostname        string   `json:"hostname,omitempty"`
	Password        string   `json:"password"`
	SSHKeys         []string `json:"ssh_key,omitempty"`
	OSPartitionSize int      `json:"os_partition_size,omitempty"`
}

// ServerAction fields for performed action on server
type ServerAction struct {
	Type string `json:"type"`
}

// PowerState fields
type PowerState struct {
	Power string `json:"power"`
}

// CreateServer fields for ordering new server
type CreateServer struct {
	ProjectID       int                `json:"project_id"`
	Plan            string             `json:"plan"`
	Hostname        string             `json:"hostname,omitempty"`
	Image           string             `json:"image,omitempty"`
	Region          string             `json:"region"`
	SSHKeys         []string           `json:"ssh_keys,omitempty"`
	IPAddresses     []string           `json:"ip_addresses,omitempty"`
	UserData        string             `json:"user_data,omitempty"`
	Tags            *map[string]string `json:"tags,omitempty"`
	SpotInstance    bool               `json:"spot_market"`
	OSPartitionSize int                `json:"os_partition_size,omitempty"`
}

// UpdateServer fields for updating a server with specified tags
type UpdateServer struct {
	Name     string             `json:"name,omitempty"`
	Hostname string             `json:"hostname,omitempty"`
	Tags     *map[string]string `json:"tags,omitempty"`
	Bgp      bool               `json:"bgp"`
}

type ServersClient struct {
	client *Client
}

// List func lists teams
func (s *ServersClient) List(projectID int, opts *GetOptions) ([]Server, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("/v1/projects/%d/servers", projectID))

	var trans []Server
	resp, err := s.client.MakeRequest("GET", path, nil, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

func (s *ServersClient) Get(serverID int, opts *GetOptions) (Server, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("%s/%d", baseServerPath, serverID))

	var trans Server

	resp, err := s.client.MakeRequest("GET", path, nil, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

func (s *ServersClient) action(serverID int, serverAction ServerAction) (Server, *Response, error) {
	var trans Server

	path := fmt.Sprintf("%s/%d/actions", baseServerPath, serverID)
	resp, err := s.client.MakeRequest("POST", path, serverAction, &trans)

	return trans, resp, err
}

// PowerOff function turns server off
func (s *ServersClient) PowerOff(serverID int) (Server, *Response, error) {
	action := ServerAction{
		Type: "power_off",
	}

	return s.action(serverID, action)
}

// PowerOn function turns server on
func (s *ServersClient) PowerOn(serverID int) (Server, *Response, error) {
	action := ServerAction{
		Type: "power_on",
	}

	return s.action(serverID, action)
}

// Reboot function restarts desired server
func (s *ServersClient) Reboot(serverID int) (Server, *Response, error) {
	action := ServerAction{
		Type: "reboot",
	}

	return s.action(serverID, action)
}

func (s *ServersClient) Reinstall(serverID int, fields *ReinstallServerFields) (Server, *Response, error) {
	var trans Server

	request := &ReinstallServer{ServerAction{Type: "reinstall"}, fields}
	path := fmt.Sprintf("%s/%d/actions", baseServerPath, serverID)
	resp, err := s.client.MakeRequest("POST", path, request, &trans)

	return trans, resp, err
}

func (s *ServersClient) PowerState(serverID int) (PowerState, *Response, error) {
	path := fmt.Sprintf("%s/%d?fields=power", baseServerPath, serverID)

	var trans PowerState

	resp, err := s.client.MakeRequest("GET", path, nil, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

func (s *ServersClient) Create(request *CreateServer) (Server, *Response, error) {
	var trans Server

	path := fmt.Sprintf("/v1/projects/%d/servers", request.ProjectID)
	resp, err := s.client.MakeRequest("POST", path, request, &trans)

	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

func (s *ServersClient) Update(serverID int, request *UpdateServer) (Server, *Response, error) {
	var trans Server

	path := fmt.Sprintf("%s/%d", baseServerPath, serverID)
	resp, err := s.client.MakeRequest("PUT", path, request, &trans)

	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

func (s *ServersClient) Delete(serverID int) (Server, *Response, error) {
	var trans Server

	path := fmt.Sprintf("%s/%d", baseServerPath, serverID)
	resp, err := s.client.MakeRequest("DELETE", path, nil, &trans)

	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

func (s *ServersClient) ListSSHKeys(serverID int, opts *GetOptions) ([]SSHKey, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("%s/%d/ssh-keys", baseServerPath, serverID))

	var trans []SSHKey
	resp, err := s.client.MakeRequest("GET", path, nil, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}
