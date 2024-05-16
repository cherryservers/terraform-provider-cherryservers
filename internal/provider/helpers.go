package provider

import (
	"encoding/base64"
	"fmt"
	"github.com/cherryservers/cherrygo/v3"
	"strings"
)

// is404Error returns true if err is an HTTP 404 error.
func is404Error(httpResponse *cherrygo.Response) bool {
	return httpResponse.StatusCode == 404
}
func is403Error(httpResponse *cherrygo.Response) bool {
	return httpResponse.StatusCode == 403
}

func ServerHostnameToID(hostname string, projectID int, ServerService cherrygo.ServersService) (int, error) {
	serversList, err := serverList(projectID, ServerService)
	for _, s := range serversList {
		if strings.EqualFold(hostname, s.Hostname) {
			return s.ID, err
		}
	}

	return 0, fmt.Errorf("could not find server with `%s` hostname", hostname)
}

func serverList(projectID int, ServerService cherrygo.ServersService) ([]cherrygo.Server, error) {
	getOptions := cherrygo.GetOptions{
		Fields: []string{"id", "name", "hostname"},
	}
	srvList, _, err := ServerService.List(projectID, &getOptions)

	return srvList, err
}

func IsBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}
