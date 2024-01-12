package cherrygo

import (
	"fmt"
)

const basePlanPath = "/v1/teams"

// PlansService is an interface for interfacing with the Plan endpoints of the CherryServers API
// See: https://api.cherryservers.com/doc/#tag/Plans
type PlansService interface {
	List(teamID int, opts *GetOptions) ([]Plan, *Response, error)
}

type Plan struct {
	ID               int                `json:"id,omitempty"`
	Name             string             `json:"name,omitempty"`
	Slug             string             `json:"slug,omitempty"`
	Custom           bool               `json:"custom,omitempty"`
	Type             string             `json:"type,omitempty"`
	Specs            Specs              `json:"specs,omitempty"`
	Pricing          []Pricing          `json:"pricing,omitempty"`
	AvailableRegions []AvailableRegions `json:"available_regions,omitempty"`
}

// Specs specifies fields for specs
type Specs struct {
	Cpus      Cpus      `json:"cpus,omitempty"`
	Memory    Memory    `json:"memory,omitempty"`
	Storage   []Storage `json:"storage,omitempty"`
	Raid      Raid      `json:"raid,omitempty"`
	Nics      Nics      `json:"nics,omitempty"`
	Bandwidth Bandwidth `json:"bandwidth,omitempty"`
}

// Cpus fields
type Cpus struct {
	Count     int     `json:"count,omitempty"`
	Name      string  `json:"name,omitempty"`
	Cores     int     `json:"cores,omitempty"`
	Frequency float32 `json:"frequency,omitempty"`
	Unit      string  `json:"unit,omitempty"`
}

// Memory fields
type Memory struct {
	Count int    `json:"count,omitempty"`
	Total int    `json:"total,omitempty"`
	Unit  string `json:"unit,omitempty"`
	Name  string `json:"name,omitempty"`
}

// Storage fields
type Storage struct {
	Count int     `json:"count,omitempty"`
	Name  string  `json:"name,omitempty"`
	Size  float32 `json:"size,omitempty"`
	Unit  string  `json:"unit,omitempty"`
}

// Raid fields
type Raid struct {
	Name string `json:"name,omitempty"`
}

// Nics fields
type Nics struct {
	Name string `json:"name,omitempty"`
}

// Bandwidth fields
type Bandwidth struct {
	Name string `json:"name,omitempty"`
}

type AvailableRegions struct {
	*Region
	StockQty int `json:"stock_qty,omitempty"`
	SpotQty  int `json:"spot_qty,omitempty"`
}

type PlansClient struct {
	client *Client
}

// List func lists plans
func (p *PlansClient) List(teamID int, opts *GetOptions) ([]Plan, *Response, error) {
	basePath := "/v1/plans"
	if teamID != 0 {
		basePath = fmt.Sprintf("%s/%d/plans", basePlanPath, teamID)
	}

	path := opts.WithQuery(basePath)
	var trans []Plan

	resp, err := p.client.MakeRequest("GET", path, nil, &trans)
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}

	return trans, resp, err
}
