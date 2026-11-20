package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStorageResource(t *testing.T) {
	storageResourceName := "terraform_test_storage_" + acctest.RandString(5)
	projectName := testProjectNamePrefix + acctest.RandString(5)
	teamID := os.Getenv("CHERRY_TEST_TEAM_ID")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccStorageResourceConfig(projectName, storageResourceName, teamID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("cherryservers_storage."+storageResourceName, "id"),
					resource.TestCheckResourceAttr("cherryservers_storage."+storageResourceName, "size", "10"),
					resource.TestCheckResourceAttr("cherryservers_storage."+storageResourceName, "region", "LT-Siauliai"),
					resource.TestCheckResourceAttr("cherryservers_storage."+storageResourceName, "description", "Test storage"),
					// iSCSI details are only populated after attachment, so skip checking them here
				),
			},
			// ImportState testing - skip project_id verification since API doesn't return it
			{
				ResourceName:            "cherryservers_storage." + storageResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"project_id"},
			},
			// Update description
			{
				Config: testAccStorageResourceConfigUpdate(projectName, storageResourceName, teamID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cherryservers_storage."+storageResourceName, "description", "Updated storage"),
				),
			},
			// Delete testing automatically occurs after the last TestStep
		},
	})
}

func testAccStorageResourceConfig(projectName string, storageName string, teamID string) string {
	return fmt.Sprintf(`
resource "cherryservers_project" "test_storage_project" {
  name = "%s"
  team_id = "%s"
}

resource "cherryservers_storage" "%s" {
  project_id  = cherryservers_project.test_storage_project.id
  region      = "LT-Siauliai"
  size        = 10
  description = "Test storage"
}
`, projectName, teamID, storageName)
}

func testAccStorageResourceConfigUpdate(projectName string, storageName string, teamID string) string {
	return fmt.Sprintf(`
resource "cherryservers_project" "test_storage_project" {
  name = "%s"
  team_id = "%s"
}

resource "cherryservers_storage" "%s" {
  project_id  = cherryservers_project.test_storage_project.id
  region      = "LT-Siauliai"
  size        = 10
  description = "Updated storage"
}
`, projectName, teamID, storageName)
}
