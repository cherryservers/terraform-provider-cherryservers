package cherrygo

import "fmt"

const baseBackupPath = "/v1/backup-storages"

type BackupsService interface {
	ListPlans(opts *GetOptions) ([]BackupStoragePlan, *Response, error)
	ListBackups(projectID int, opts *GetOptions) ([]BackupStorage, *Response, error)
	Get(backupID int, opts *GetOptions) (BackupStorage, *Response, error)
	Create(request *CreateBackup) (BackupStorage, *Response, error)
	Update(request *UpdateBackupStorage) (BackupStorage, *Response, error)
	UpdateBackupMethod(request *UpdateBackupMethod) ([]BackupMethod, *Response, error)
	Delete(backupID int) (*Response, error)
}

type BackupsClient struct {
	client *Client
}

type BackupStoragePlan struct {
	ID            int       `json:"id,omitempty"`
	Name          string    `json:"name,omitempty"`
	Slug          string    `json:"slug,omitempty"`
	SizeGigabytes int       `json:"size_gigabytes,omitempty"`
	Pricing       []Pricing `json:"pricing,omitempty"`
	Regions       []Region  `json:"regions,omitempty"`
	Href          string    `json:"href,omitempty"`
}

type BackupStorage struct {
	ID                   int            `json:"id,omitempty"`
	Status               string         `json:"status,omitempty"`
	State                string         `json:"state,omitempty"`
	PrivateIP            string         `json:"private_ip,omitempty"`
	PublicIP             string         `json:"public_ip,omitempty"`
	SizeGigabytes        int            `json:"size_gigabytes,omitempty"`
	UsedGigabytes        int            `json:"used_gigabytes,omitempty"`
	AttachedTo           AttachedTo     `json:"attached_to,omitempty"`
	Methods              []BackupMethod `json:"methods,omitempty"`
	AvailableIPAddresses []IPAddress    `json:"available_addresses,omitempty"`
	Rules                []Rule         `json:"rules,omitempty"`
	Plan                 Plan           `json:"plan,omitempty"`
	Pricing              Pricing        `json:"pricing,omitempty"`
	Region               Region         `json:"region,omitempty"`
	Href                 string         `json:"href,omitempty"`
}

type BackupMethod struct {
	Name       string   `json:"name,omitempty"`
	Username   string   `json:"username,omitempty"`
	Password   string   `json:"password,omitempty"`
	Port       int      `json:"port,omitempty"`
	Host       string   `json:"host,omitempty"`
	SSHKey     string   `json:"ssh_key,omitempty"`
	WhiteList  []string `json:"whitelist,omitempty"`
	Enabled    bool     `json:"enabled,omitempty"`
	Processing bool     `json:"processing,omitempty"`
}

type Rule struct {
	IPAddress      IPAddress      `json:"ip,omitempty"`
	EnabledMethods EnabledMethods `json:"methods,omitempty"`
}

type EnabledMethods struct {
	BORG bool `json:"borg,omitempty"`
	FTP  bool `json:"ftp,omitempty"`
	NFS  bool `json:"nfs,omitempty"`
	SMB  bool `json:"smb,omitempty"`
}

type CreateBackup struct {
	ServerID       int    `json:"server_id,omitempty"`
	BackupPlanSlug string `json:"slug"`
	RegionSlug     string `json:"region"`
	SSHKey         string `json:"ssh_key,omitempty"`
}

type UpdateBackupStorage struct {
	BackupStorageID int    `json:"id"`
	BackupPlanSlug  string `json:"slug,omitempty"`
	Password        string `json:"password,omitempty"`
	SSHKey          string `json:"ssh_key,omitempty"`
}

type UpdateBackupMethod struct {
	BackupStorageID  int      `json:"id"`
	BackupMethodName string   `json:"name"`
	Enabled          bool     `json:"enabled"`
	Whitelist        []string `json:"whitelist"`
}

func (s *BackupsClient) ListPlans(opts *GetOptions) ([]BackupStoragePlan, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("/v1/backup-storage-plans"))

	var trans []BackupStoragePlan
	resp, err := s.client.MakeRequest("GET", path, nil, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

func (s *BackupsClient) ListBackups(projectID int, opts *GetOptions) ([]BackupStorage, *Response, error) {
	var trans []BackupStorage

	path := opts.WithQuery(fmt.Sprintf("/v1/projects/%d/backup-storages", projectID))
	resp, err := s.client.MakeRequest("GET", path, nil, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

func (s *BackupsClient) Get(backupID int, opts *GetOptions) (BackupStorage, *Response, error) {
	var trans BackupStorage

	path := opts.WithQuery(fmt.Sprintf("%s/%d", baseBackupPath, backupID))
	resp, err := s.client.MakeRequest("GET", path, nil, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

func (s *BackupsClient) Create(request *CreateBackup) (BackupStorage, *Response, error) {
	var trans BackupStorage

	path := fmt.Sprintf("/v1/servers/%d/backup-storages", request.ServerID)
	resp, err := s.client.MakeRequest("POST", path, request, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

func (s *BackupsClient) Update(request *UpdateBackupStorage) (BackupStorage, *Response, error) {
	var trans BackupStorage

	path := fmt.Sprintf("%s/%d", baseBackupPath, request.BackupStorageID)

	resp, err := s.client.MakeRequest("PUT", path, request, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

func (s *BackupsClient) UpdateBackupMethod(request *UpdateBackupMethod) ([]BackupMethod, *Response, error) {
	var trans []BackupMethod

	path := fmt.Sprintf("%s/%d/methods/%s", baseBackupPath, request.BackupStorageID, request.BackupMethodName)
	resp, err := s.client.MakeRequest("PATCH", path, request, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

func (s *BackupsClient) Delete(backupID int) (*Response, error) {
	path := fmt.Sprintf("%s/%d", baseBackupPath, backupID)
	resp, err := s.client.MakeRequest("DELETE", path, nil, nil)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return resp, err
}
