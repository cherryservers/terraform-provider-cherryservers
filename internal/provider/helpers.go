package provider

import "github.com/cherryservers/cherrygo/v3"

// is404Error returns true if err is an HTTP 404 error
func is404Error(httpResponse *cherrygo.Response) bool {
	return httpResponse.StatusCode == 404
}
