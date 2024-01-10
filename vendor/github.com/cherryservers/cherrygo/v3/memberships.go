package cherrygo

type Membership struct {
	ID       int       `json:"id,omitempty"`
	Roles    []string  `json:"roles,omitempty"`
	User     User      `json:"user,omitempty"`
	Projects []Project `json:"projects,omitempty"`
	Href     string    `json:"href,omitempty"`
}
