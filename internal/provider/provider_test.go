package provider

import (
	"fmt"
	"github.com/cherryservers/cherrygo/v3"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

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
