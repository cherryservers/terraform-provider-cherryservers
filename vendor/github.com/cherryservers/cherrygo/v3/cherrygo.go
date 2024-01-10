package cherrygo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
)

const (
	libraryVersion     = "3.1.0"
	apiURL             = "https://api.cherryservers.com/v1/"
	cherryAuthTokenVar = "CHERRY_AUTH_TOKEN"
	mediaType          = "application/json"
	userAgent          = "cherry-agent-go/" + libraryVersion
	cherryDebugVar     = "CHERRY_DEBUG"
)

// Client returns struct for client
type Client struct {
	client *http.Client
	debug  bool

	BaseURL *url.URL

	UserAgent string
	AuthToken string

	Teams       TeamsService
	Plans       PlansService
	Images      ImagesService
	Projects    ProjectsService
	SSHKeys     SSHKeysService
	Servers     ServersService
	IPAddresses IpAddressesService
	Storages    StoragesService
	Regions     RegionsService
	Users       UsersService
	Backups     BackupsService
}

// Response is the http response from api calls
type Response struct {
	*http.Response
	Meta
}

type Meta struct {
	Total int
}

// MakeRequest makes request to API
func (c *Client) MakeRequest(method, path string, body, v interface{}) (*Response, error) {

	url, _ := url.Parse(path)

	u := c.BaseURL.ResolveReference(url)
	if c.debug {
		fmt.Printf("\nAPI Endpoint: %v\n", u)
	}

	buf := new(bytes.Buffer)
	if body != nil {
		coder := json.NewEncoder(buf)
		err := coder.Encode(body)
		if err != nil {
			log.Printf("Error while encoding body: %v -> %v", err, err.Error())
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Close = true

	bearer := "Bearer " + c.AuthToken
	req.Header.Add("Authorization", bearer)
	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)

	if c.debug {
		o, _ := httputil.DumpRequestOut(req, true)
		log.Printf("\n+++++++++++++REQUEST+++++++++++++\n%s\n+++++++++++++++++++++++++++++++++", string(o))
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	response := Response{Response: resp}
	response.populateTotal()

	if c.debug {
		o, _ := httputil.DumpResponse(response.Response, true)
		log.Printf("\n+++++++++++++RESPONSE+++++++++++++\n%s\n+++++++++++++++++++++++++++++++++", string(o))
	}

	if sc := response.StatusCode; sc >= 299 {
		type ErrorResponse struct {
			Response *http.Response
			Code     int    `json:"code"`
			Message  string `json:"message"`
		}

		var errorResponse ErrorResponse

		bod, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(bod, &errorResponse)
		if err != nil {
			return nil, err
		}

		err = fmt.Errorf("Error response from API: %v (error code: %v)", errorResponse.Message, errorResponse.Code)

		return &response, err
	}

	// Handling delete requests which EOF is not an error
	if method == "DELETE" && response.StatusCode == 204 {
		return &response, err
	}

	if v != nil {
		// if v implements the io.Writer interface, return the raw response
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {

			decoder := json.NewDecoder(resp.Body)
			err := decoder.Decode(&v)
			if err != nil {
				log.Printf("Error while decoding body: %v -> %v", err, err.Error())
				return &response, err
			}
		}
	}

	return &response, nil
}

type options struct {
	url       string
	client    *http.Client
	userAgent string
	authToken string
}

type ClientOpt func(*options) error

// NewClient initialization
func NewClient(opts ...ClientOpt) (*Client, error) {

	parsedOpts := &options{
		authToken: os.Getenv(cherryAuthTokenVar),
		client:    &http.Client{},
		url:       apiURL,
		userAgent: userAgent,
	}
	for _, opt := range opts {
		if err := opt(parsedOpts); err != nil {
			return nil, err
		}
	}
	if parsedOpts.authToken == "" {
		return nil, fmt.Errorf("auth token must be provided as parameter of environment variable %s", cherryAuthTokenVar)
	}

	url, err := url.Parse(parsedOpts.url)
	if err != nil {
		return nil, err
	}

	c := &Client{client: parsedOpts.client, AuthToken: parsedOpts.authToken, BaseURL: url, UserAgent: parsedOpts.userAgent}

	c.debug = os.Getenv(cherryDebugVar) != ""
	c.Teams = &TeamsClient{client: c}
	c.Plans = &PlansClient{client: c}
	c.Images = &ImagesClient{client: c}
	c.Projects = &ProjectsClient{client: c}
	c.SSHKeys = &SSHKeysClient{client: c}
	c.Servers = &ServersClient{client: c}
	c.IPAddresses = &IPsClient{client: c}
	c.Storages = &StoragesClient{client: c}
	c.Regions = &RegionsClient{client: c}
	c.Users = &UsersClient{client: c}
	c.Backups = &BackupsClient{client: c}

	return c, err
}

// ErrorResponse fields
type ErrorResponse struct {
	Response    *http.Response
	Errors      []string `json:"errors"`
	SingleError string   `json:"error"`
}

func checkResponseForErrors(r *http.Response) *ErrorResponse {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errR := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		json.Unmarshal(data, errR)
	}

	return errR

}

// WithUserAgent set user agent when making requests
func WithUserAgent(ua string) ClientOpt {
	return func(c *options) error {
		c.userAgent = fmt.Sprintf("%s %s", ua, userAgent)
		return nil
	}
}

// WithURL use url as endpoint for API requests
func WithURL(url string) ClientOpt {
	return func(c *options) error {
		c.url = url
		return nil
	}
}

// WithHTTPClient use client as the http.Client to make API requests
func WithHTTPClient(client *http.Client) ClientOpt {
	return func(c *options) error {
		c.client = client
		return nil
	}
}

// WithAuthToken use provided auth token to make requests, defaults to environment variable
// CHERRY_AUTH_TOKEN
func WithAuthToken(authToken string) ClientOpt {
	return func(c *options) error {
		c.authToken = authToken
		return nil
	}
}

func (r *Response) populateTotal() {
	// parse the headers and populate Meta.Total
	if total := r.Header.Get("X-Total-Count"); total != "" {
		r.Meta.Total, _ = strconv.Atoi(total)
	}
}
