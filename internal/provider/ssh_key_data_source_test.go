package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSSHKeyDataSource_basic(t *testing.T) {
	label := "terraform_test_ssh_" + acctest.RandString(5)
	publicKey, _, err := acctest.RandSSHKeyPair("cherryservers@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}
	resourceName := "cherryservers_ssh_key.test_ssh_key_ssh_key"
	dataSourceName := "data.cherryservers_ssh_key.test_ssh_key_ssh_key"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccSSHKeyDataSourceConfig(label, publicKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "label", dataSourceName, "label"),
					resource.TestCheckResourceAttrPair(resourceName, "created", dataSourceName, "created"),
					resource.TestCheckResourceAttrPair(resourceName, "fingerprint", dataSourceName, "fingerprint"),
					resource.TestCheckResourceAttrPair(resourceName, "id", dataSourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "label", dataSourceName, "label"),
					resource.TestCheckResourceAttrPair(resourceName, "public_key", dataSourceName, "public_key"),
					resource.TestCheckResourceAttrPair(resourceName, "updated", dataSourceName, "updated"),
				),
			},
		},
	})
}

func testAccSSHKeyDataSourceConfig(label string, publicKey string) string {
	return fmt.Sprintf(`
resource "cherryservers_ssh_key" "test_ssh_key_ssh_key" {
  label = "%s"
  public_key = "%s"
}

data "cherryservers_ssh_key" "test_ssh_key_ssh_key" {
  id = "${cherryservers_ssh_key.test_ssh_key_ssh_key.id}"
}
`, label, publicKey)
}
