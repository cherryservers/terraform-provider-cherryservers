package provider

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/cherryservers/cherrygo/v3"
	"strings"
)

// is404Error returns true if err is an HTTP 404 error.
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

	return 0, fmt.Errorf("could not find server with `%s` hostname", hostname)
}

func serverList(projectID int, ServerService cherrygo.ServersService) ([]cherrygo.Server, error) {
	getOptions := cherrygo.GetOptions{
		Fields: []string{"id", "name", "hostname"},
	}
	srvList, _, err := ServerService.List(projectID, &getOptions)

	return srvList, err
}

func IsBase64(s string) error {
	_, err := base64.StdEncoding.DecodeString(s)
	return err
}

func NormalizeServerImage(server *cherrygo.Server, client *cherrygo.Client) error {
	images, _, err := client.Images.List(server.Plan.Slug, nil)
	if err != nil {
		return err
	}

	for _, image := range images {
		if image.Name == server.Image {
			server.Image = image.Slug
			return nil
		}
	}

	return errors.New("could not find image slug for image with name `" + server.Image + "`")
}
