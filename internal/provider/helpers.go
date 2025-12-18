package provider

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/cherryservers/cherrygo/v3"
)

// is404Error returns true if err is an HTTP 404 error.
func is404Error(httpResponse *cherrygo.Response) bool {
	return httpResponse.StatusCode == 404
}

func serverHostnameToID(hostname string, projectID int, ServerService cherrygo.ServersService) (int, error) {
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

func isBase64(s string) error {
	_, err := base64.StdEncoding.DecodeString(s)
	return err
}

// normalizeServerImage is used to transform the server image field into the same type of slug
// that is used in the schema.
func normalizeServerImage(server *cherrygo.Server, client *cherrygo.Client) error {
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

// generatePassword is used to generate a password that matches CherryServers password constraints.
func generatePassword() (string, error) {
	const (
		lowercase = "abcdefghijklmnopqrstuvwxyz"
		uppercase = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		digits    = "0123456789"
		all       = lowercase + uppercase + digits
		length    = 20
	)
	password := make([]byte, length)

	var charset string
	for i := range length {
		switch i {
		case 0:
			// Ensure there is at least one lower-case alphabetical character.
			charset = lowercase
		case 1:
			// Ensure there is at least one upper-case alphabetical
			// character that is not first.
			charset = uppercase
		case 2:
			// Ensure there is at least one digit that is not last.
			charset = digits
		default:
			charset = all
		}
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[idx.Int64()]
	}
	return string(password), nil
}
