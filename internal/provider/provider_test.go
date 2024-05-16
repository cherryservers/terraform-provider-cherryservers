package provider

import (
	"fmt"
	"github.com/cherryservers/cherrygo/v3"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// TODO
// Use a provider instead of a client, when that functionality is available.
var testCherryGoClient *cherrygo.Client

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
	if v := os.Getenv("CHERRY_AUTH_KEY"); v == "" {
		t.Fatal("CHERRY_AUTH_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("CHERRY_TEST_TEAM_ID"); v == "" {
		t.Fatal("CHERRY_TEST_TEAM_ID must be set for acceptance tests")
	}

	//TODO
	//Make version responsive (or better yet, use a provider instead of a client)
	userAgent := fmt.Sprintf("terraform-provider/cherryservers/%s terraform/%s", "test", "1.0.0")
	args := []cherrygo.ClientOpt{cherrygo.WithAuthToken(os.Getenv("CHERRY_AUTH_KEY")), cherrygo.WithUserAgent(userAgent)}
	client, err := cherrygo.NewClient(args...)
	if err != nil {
		t.Fatal("error: couldn't create client a cherry servers client for testing:", err)
	}
	testCherryGoClient = client

}
