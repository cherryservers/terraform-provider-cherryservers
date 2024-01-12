package cherrygo

import (
	"fmt"
)

const baseImagePath = "/v1/plans"

// ImagesService is an interface for interfacing with the the Images endpoints of the CherryServers API
// See: https://api.cherryservers.com/doc/#tag/Images
type ImagesService interface {
	List(plan string, opts *GetOptions) ([]Image, *Response, error)
}

type Image struct {
	ID      int       `json:"id,omitempty"`
	Name    string    `json:"name,omitempty"`
	Slug    string    `json:"slug,omitempty"`
	Pricing []Pricing `json:"pricing,omitempty"`
}

type ImagesClient struct {
	client *Client
}

// List func lists images
func (i *ImagesClient) List(plan string, opts *GetOptions) ([]Image, *Response, error) {
	path := opts.WithQuery(fmt.Sprintf("%s/%s/images", baseImagePath, plan))
	var trans []Image

	resp, err := i.client.MakeRequest("GET", path, nil, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}
