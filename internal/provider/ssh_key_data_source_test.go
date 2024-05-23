package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"os"
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

func TestAccSSHKeyDataSource_byLabel(t *testing.T) {
	teamId := os.Getenv("CHERRY_TEST_TEAM_ID")
	label := "terraform_test_ssh_" + acctest.RandString(5)
	publicKey, _, err := acctest.RandSSHKeyPair("cherryservers@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}
	projectName := "terraform_test_project_" + acctest.RandString(5)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccSSHKeyDataSourceByLabelConfig(projectName, teamId, label, publicKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair("cherryservers_ssh_key.test_ssh_key_ssh_key_by_label", "label", "data.cherryservers_ssh_key.test_ssh_key_ssh_key_by_label", "label"),
					resource.TestCheckResourceAttrPair("cherryservers_ssh_key.test_ssh_key_ssh_key_by_label", "created", "data.cherryservers_ssh_key.test_ssh_key_ssh_key_by_label", "created"),
					resource.TestCheckResourceAttrPair("cherryservers_ssh_key.test_ssh_key_ssh_key_by_label", "fingerprint", "data.cherryservers_ssh_key.test_ssh_key_ssh_key_by_label", "fingerprint"),
					resource.TestCheckResourceAttrPair("cherryservers_ssh_key.test_ssh_key_ssh_key_by_label", "id", "data.cherryservers_ssh_key.test_ssh_key_ssh_key_by_label", "id"),
					resource.TestCheckResourceAttrPair("cherryservers_ssh_key.test_ssh_key_ssh_key_by_label", "label", "data.cherryservers_ssh_key.test_ssh_key_ssh_key_by_label", "label"),
					resource.TestCheckResourceAttrPair("cherryservers_ssh_key.test_ssh_key_ssh_key_by_label", "public_key", "data.cherryservers_ssh_key.test_ssh_key_ssh_key_by_label", "public_key"),
					resource.TestCheckResourceAttrPair("cherryservers_ssh_key.test_ssh_key_ssh_key_by_label", "updated", "data.cherryservers_ssh_key.test_ssh_key_ssh_key_by_label", "updated"),
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

func testAccSSHKeyDataSourceByLabelConfig(projectName string, teamID string, label string, publicKey string) string {
	return fmt.Sprintf(`
resource "cherryservers_project" "test_ssh_key_project" {
  name = "%s"
  team_id      = "%s"
}

resource "cherryservers_ssh_key" "test_ssh_key_ssh_key_by_label" {
  label = "%s"
  public_key = "%s"
}

data "cherryservers_ssh_key" "test_ssh_key_ssh_key_by_label" {
  label = "${cherryservers_ssh_key.test_ssh_key_ssh_key_by_label.label}"
  project_id = "${cherryservers_project.test_ssh_key_project.id}"
}
`, projectName, teamID, label, publicKey)
}
