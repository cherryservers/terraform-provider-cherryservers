package provider

import (
	"fmt"
	"github.com/cherryservers/cherrygo/v3"
	"strings"
)

// is404Error returns true if err is an HTTP 404 error
func is404Error(httpResponse *cherrygo.Response) bool {
	return httpResponse.StatusCode == 404
}

func ServerHostnameToID(hostname string, projectID int, ServerService cherrygo.ServersService) (int, error) {
	serversList, err := serverList(projectID, ServerService)
	for _, s := range serversList {
		if strings.EqualFold(hostname, s.Hostname) {
			return s.ID, err
		}
	}

	return 0, fmt.Errorf("Could not find server with `%s` hostname", hostname)
}

func serverList(projectID int, ServerService cherrygo.ServersService) ([]cherrygo.Server, error) {
	getOptions := cherrygo.GetOptions{
		Fields: []string{"id", "name", "hostname"},
	}
	serverList, _, err := ServerService.List(projectID, &getOptions)

	return serverList, err
}
