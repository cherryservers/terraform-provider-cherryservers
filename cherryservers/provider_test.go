package cherryservers

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider
var testAccProviderFactories map[string]func() (*schema.Provider, error)
var teamID = "48787"

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"cherryservers": testAccProvider,
	}
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"cherryservers": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("CHERRY_AUTH_TOKEN"); v == "" {
		t.Fatal("CHERRY_AUTH_TOKEN must be set for acceptance tests")
	}
}
