package cherrygo

import "fmt"

const baseUserPath = "/v1/users"

// UsersService is an interface for interfacing with the the User endpoints of the CherryServers API
// See: https://api.cherryservers.com/doc/#tag/Users
type UsersService interface {
	CurrentUser(opts *GetOptions) (User, *Response, error)
	Get(userID int, opts *GetOptions) (User, *Response, error)
}

type User struct {
	ID                    int    `json:"id,omitempty"`
	FirstName             string `json:"first_name,omitempty"`
	LastName              string `json:"last_name,omitempty"`
	Email                 string `json:"email,omitempty"`
	EmailVerified         bool   `json:"email_verified,omitempty"`
	Phone                 string `json:"phone,omitempty"`
	SecurityPhone         string `json:"security_phone,omitempty"`
	SecurityPhoneVerified bool   `json:"security_phone_verified,omitempty"`
	Href                  string `json:"href,omitempty"`
}

type UsersClient struct {
	client *Client
}

func (s *UsersClient) CurrentUser(opts *GetOptions) (User, *Response, error) {
	var trans User

	path := opts.WithQuery(fmt.Sprintf("/v1/user"))

	resp, err := s.client.MakeRequest("GET", path, nil, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

func (s *UsersClient) Get(userID int, opts *GetOptions) (User, *Response, error) {
	var trans User

	path := opts.WithQuery(fmt.Sprintf("%s/%d", baseUserPath, userID))

	resp, err := s.client.MakeRequest("GET", path, nil, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}
