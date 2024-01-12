package cherrygo

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type GetOptions struct {
	Fields []string `url:"fields,omitempty,comma"`
	Limit  int      `url:"limit,omitempty"`
	Offset int      `url:"offset,omitempty"`
	Type   []string `url:"type,ommitempty"`
	Status []string `url:"status,ommitempty"`
	// QueryParams for API URL, used for arbitrary filters
	QueryParams map[string]string `url:"-"`
}

func (g *GetOptions) WithQuery(apiPath string) string {
	params := g.Encode()
	if params != "" {
		return fmt.Sprintf("%s?%s", apiPath, params)
	}
	return apiPath
}

func (g *GetOptions) Encode() string {
	if g == nil {
		return ""
	}
	v := url.Values{}
	if g.Fields != nil && len(g.Fields) > 0 {
		v.Add("fields", strings.Join(g.Fields, ","))
	}

	if g.QueryParams != nil {
		for k, val := range g.QueryParams {
			v.Add(k, val)
		}
	}

	if g.Type != nil && len(g.Type) > 0 {
		for _, el := range g.Type {
			v.Add("type[]", el)
		}
	}

	if g.Status != nil && len(g.Status) > 0 {
		for _, el := range g.Status {
			v.Add("status[]", el)
		}
	}

	if g.Limit != 0 {
		v.Add("limit", strconv.Itoa(g.Limit))
	}

	if g.Offset != 0 {
		v.Add("offset", strconv.Itoa(g.Offset))
	}

	return v.Encode()
}
