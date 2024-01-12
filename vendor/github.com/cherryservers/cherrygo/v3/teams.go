package cherrygo

import "fmt"

const teamsPath = "/v1/teams"

// TeamsService is an interface for interfacing with the Teams endpoints of the CherryServers API
// See: https://api.cherryservers.com/doc/#tag/Teams
type TeamsService interface {
	List(opts *GetOptions) ([]Team, *Response, error)
	Get(teamID int, opts *GetOptions) (Team, *Response, error)
	Create(request *CreateTeam) (Team, *Response, error)
	Update(teamID int, request *UpdateTeam) (Team, *Response, error)
	Delete(teamID int) (*Response, error)
}

type Team struct {
	ID          int          `json:"id,omitempty"`
	Name        string       `json:"name,omitempty"`
	Credit      Credit       `json:"credit,omitempty"`
	Billing     Billing      `json:"billing,omitempty"`
	Projects    []Project    `json:"projects,omitempty"`
	Memberships []Membership `json:"memberships,omitempty"`
	Href        string       `json:"href,omitempty"`
}

type Credit struct {
	Account   CreditDetails `json:"account,omitempty"`
	Promo     CreditDetails `json:"promo,omitempty"`
	Resources Resources     `json:"resources,omitempty"`
}

type CreditDetails struct {
	Remaining float32 `json:"remaining,omitempty"`
	Usage     float32 `json:"usage,omitempty"`
	Currency  string  `json:"currency,omitempty"`
}

type Resources struct {
	Pricing   Pricing       `json:"pricing,omitempty"`
	Remaining RemainingTime `json:"remaining,omitempty"`
}

type RemainingTime struct {
	Time int    `json:"time,omitempty"`
	Unit string `json:"unit,omitempty"`
}

type Pricing struct {
	Price    float32 `json:"price,omitempty"`
	Taxed    bool    `json:"taxed,omitempty"`
	Currency string  `json:"currency,omitempty"`
	Unit     string  `json:"unit,omitempty"`
}

type Billing struct {
	Type        string `json:"type,omitempty"`
	CompanyName string `json:"company_name,omitempty"`
	CompanyCode string `json:"company_code,omitempty"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Address1    string `json:"address_1,omitempty"`
	Address2    string `json:"address_2,omitempty"`
	CountryIso2 string `json:"country_iso_2,omitempty"`
	City        string `json:"city,omitempty"`
	Vat         Vat    `json:"vat,omitempty"`
	Currency    string `json:"currency,omitempty"`
}

type Vat struct {
	Amount int    `json:"amount"`
	Number string `json:"number,omitempty"`
	Valid  bool   `json:"valid"`
}

type TeamsClient struct {
	client *Client
}

type CreateTeam struct {
	Name     string `json:"name,omitempty"`
	Type     string `json:"type,omitempty"`
	Currency string `json:"currency,omitempty"`
}

type UpdateTeam struct {
	Name     *string `json:"name,omitempty"`
	Type     *string `json:"type,omitempty"`
	Currency *string `json:"currency,omitempty"`
}

// List func lists teams
func (t *TeamsClient) List(opts *GetOptions) ([]Team, *Response, error) {
	var trans []Team

	pathQuery := opts.WithQuery(teamsPath)
	resp, err := t.client.MakeRequest("GET", pathQuery, nil, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

func (p *TeamsClient) Get(teamID int, opts *GetOptions) (Team, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("%s/%d", teamsPath, teamID))

	var trans Team

	resp, err := p.client.MakeRequest("GET", path, nil, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

func (p *TeamsClient) Create(request *CreateTeam) (Team, *Response, error) {
	path := fmt.Sprintf("%s", teamsPath)

	var trans Team

	resp, err := p.client.MakeRequest("POST", path, request, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

func (p *TeamsClient) Update(teamID int, request *UpdateTeam) (Team, *Response, error) {
	path := fmt.Sprintf(fmt.Sprintf("%s/%d", teamsPath, teamID))

	var trans Team

	resp, err := p.client.MakeRequest("PUT", path, request, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

func (p *TeamsClient) Delete(teamID int) (*Response, error) {
	path := fmt.Sprintf("%s/%d", teamsPath, teamID)

	resp, err := p.client.MakeRequest("DELETE", path, nil, nil)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return resp, err
}
