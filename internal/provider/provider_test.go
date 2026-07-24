package provider

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/cherryservers/cherrygo/v3"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	testProjectNamePrefix = "terraform_test_project_"
	defaultTestImage      = "ubuntu_26_04_64bit"
)

var (
	testCherryGoClient *cherrygo.Client
	testTeam           int
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"cherryservers": providerserver.NewProtocol6WithError(New("test")()),
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.
}

func setTestTeam() error {
	const teamIDVar = "CHERRY_TEST_TEAM_ID"

	team, ok := os.LookupEnv(teamIDVar)
	if !ok {
		return fmt.Errorf("%s must be set for acceptance tests", teamIDVar)
	}
	id, err := strconv.Atoi(team)
	if err != nil {
		return fmt.Errorf("%s must be an integer: %s", teamIDVar, err.Error())
	}

	testTeam = id
	return nil
}

func setupClient() (*cherrygo.Client, error) {
	apiKey := os.Getenv("CHERRY_AUTH_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("CHERRY_AUTH_TOKEN")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("CHERRY_AUTH_KEY or CHERRY_AUTH_TOKEN must be set for acceptance tests")
	}

	userAgent := "terraform-provider/cherryservers/test terraform/dev"
	args := []cherrygo.ClientOpt{cherrygo.WithAuthToken(apiKey), cherrygo.WithUserAgent(userAgent)}
	client, err := cherrygo.NewClient(args...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func TestMain(m *testing.M) {
	var err error
	testCherryGoClient, err = setupClient()
	if err != nil {
		log.Fatalf("failed to initialize api client: %s", err.Error())
	}

	err = setTestTeam()
	if err != nil {
		log.Fatalf("failed to get test team: %s", err.Error())
	}

	resource.TestMain(m)
}
