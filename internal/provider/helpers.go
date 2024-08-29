package provider

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/cherryservers/cherrygo/v3"
	"math/rand"
	"strings"
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
func generatePassword() string {
	const (
		lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
		uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		digits           = "0123456789"
		all              = lowercaseLetters + uppercaseLetters + digits
	)

	password := make([]byte, 16)
	password[0] = lowercaseLetters[rand.Intn(len(lowercaseLetters))]
	password[1] = uppercaseLetters[rand.Intn(len(uppercaseLetters))]
	password[2] = digits[rand.Intn(len(digits))]
	for i := 3; i < 16; i++ {
		password[i] = all[rand.Intn(len(all))]
	}

	return string(password)
}
