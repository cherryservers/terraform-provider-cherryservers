package cherrygo

import "fmt"

const baseRegionPath = "/v1/regions"

// RegionsService is an interface for interfacing with the the Images endpoints of the CherryServers API
// See: https://api.cherryservers.com/doc/#tag/Regions
type RegionsService interface {
	List(opts *GetOptions) ([]Region, *Response, error)
	Get(region string, opts *GetOptions) (Region, *Response, error)
}

// Region fields
type Region struct {
	ID         int       `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	Slug       string    `json:"slug,omitempty"`
	RegionIso2 string    `json:"region_iso_2,omitempty"`
	BGP        RegionBGP `json:"bgp,omitempty"`
	Location   string    `json:"location,omitempty"`
	Href       string    `json:"href,omitempty"`
}

type RegionsClient struct {
	client *Client
}

func (i *RegionsClient) List(opts *GetOptions) ([]Region, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("%s", baseRegionPath))
	var trans []Region

	resp, err := i.client.MakeRequest("GET", path, nil, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}

func (i *RegionsClient) Get(region string, opts *GetOptions) (Region, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("%s/%s", baseRegionPath, region))
	var trans Region

	resp, err := i.client.MakeRequest("GET", path, nil, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}
