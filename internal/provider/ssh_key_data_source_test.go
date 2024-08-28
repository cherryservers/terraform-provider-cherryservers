package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSSHKeyDataSource_basic(t *testing.T) {
	name := "terraform_test_ssh_" + acctest.RandString(5)
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
				Config: testAccSSHKeyDataSourceConfig(name, publicKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "created", dataSourceName, "created"),
					resource.TestCheckResourceAttrPair(resourceName, "fingerprint", dataSourceName, "fingerprint"),
					resource.TestCheckResourceAttrPair(resourceName, "id", dataSourceName, "id"),
					resource.TestCheckResourceAttrPair(resourceName, "name", dataSourceName, "name"),
					resource.TestCheckResourceAttrPair(resourceName, "public_key", dataSourceName, "public_key"),
					resource.TestCheckResourceAttrPair(resourceName, "updated", dataSourceName, "updated"),
				),
			},
		},
	})
}

func TestAccSSHKeyDataSource_byName(t *testing.T) {
	name := "terraform_test_ssh_" + acctest.RandString(5)
	publicKey, _, err := acctest.RandSSHKeyPair("cherryservers@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccSSHKeyDataSourceByNameConfig(name, publicKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("cherryservers_ssh_key.test_ssh_key_ssh_key_by_name", "name", "data.cherryservers_ssh_key.test_ssh_key_ssh_key_by_name", "name"),
					resource.TestCheckResourceAttrPair("cherryservers_ssh_key.test_ssh_key_ssh_key_by_name", "created", "data.cherryservers_ssh_key.test_ssh_key_ssh_key_by_name", "created"),
					resource.TestCheckResourceAttrPair("cherryservers_ssh_key.test_ssh_key_ssh_key_by_name", "fingerprint", "data.cherryservers_ssh_key.test_ssh_key_ssh_key_by_name", "fingerprint"),
					resource.TestCheckResourceAttrPair("cherryservers_ssh_key.test_ssh_key_ssh_key_by_name", "id", "data.cherryservers_ssh_key.test_ssh_key_ssh_key_by_name", "id"),
					resource.TestCheckResourceAttrPair("cherryservers_ssh_key.test_ssh_key_ssh_key_by_name", "name", "data.cherryservers_ssh_key.test_ssh_key_ssh_key_by_name", "name"),
					resource.TestCheckResourceAttrPair("cherryservers_ssh_key.test_ssh_key_ssh_key_by_name", "public_key", "data.cherryservers_ssh_key.test_ssh_key_ssh_key_by_name", "public_key"),
					resource.TestCheckResourceAttrPair("cherryservers_ssh_key.test_ssh_key_ssh_key_by_name", "updated", "data.cherryservers_ssh_key.test_ssh_key_ssh_key_by_name", "updated"),
				),
			},
		},
	})
}

func testAccSSHKeyDataSourceConfig(name string, publicKey string) string {
	return fmt.Sprintf(`
resource "cherryservers_ssh_key" "test_ssh_key_ssh_key" {
  name = "%s"
  public_key = "%s"
}

data "cherryservers_ssh_key" "test_ssh_key_ssh_key" {
  id = "${cherryservers_ssh_key.test_ssh_key_ssh_key.id}"
}
`, name, publicKey)
}

func testAccSSHKeyDataSourceByNameConfig(name string, publicKey string) string {
	return fmt.Sprintf(`
resource "cherryservers_ssh_key" "test_ssh_key_ssh_key_by_name" {
  name = "%s"
  public_key = "%s"
}

data "cherryservers_ssh_key" "test_ssh_key_ssh_key_by_name" {
  name = "${cherryservers_ssh_key.test_ssh_key_ssh_key_by_name.name}"
}
`, name, publicKey)
}
