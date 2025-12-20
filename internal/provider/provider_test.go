package provider

import (
	"fmt"
	"github.com/cherryservers/cherrygo/v3"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// Test name prefix for resources
const testProjectNamePrefix = "terraform-test-"

var testCherryGoClient *cherrygo.Client

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"cherryservers": providerserver.NewProtocol6WithError(New("test")()),
}

// sharedClient creates and returns a shared client for acceptance testing.
// It reads the API token from CHERRY_AUTH_TOKEN or CHERRY_AUTH_KEY environment variables.
func sharedClient() (interface{}, error) {
	apiToken := os.Getenv("CHERRY_AUTH_TOKEN")
	if apiToken == "" {
		apiToken = os.Getenv("CHERRY_AUTH_KEY")
	}

	if apiToken == "" {
		return nil, fmt.Errorf("CHERRY_AUTH_TOKEN or CHERRY_AUTH_KEY environment variable not set")
	}

	client, err := cherrygo.NewClient(cherrygo.WithAuthToken(apiToken))
	if err != nil {
		return nil, fmt.Errorf("failed to create cherrygo client: %w", err)
	}

	return client, nil
}

func testAccPreCheck(t *testing.T) {
	// You can add code here to run prior to any test case execution, for example assertions
	// about the appropriate environment variables being set are common to see in a pre-check
	// function.

	client, err := sharedClient()
	if err != nil {
		t.Fatal(err)
	}

	var ok bool
	testCherryGoClient, ok = client.(*cherrygo.Client)
	if !ok {
		errStr := fmt.Sprintf("expected cherrygo.Client, got %T", client)
		t.Fatal(errStr)
	}

}
