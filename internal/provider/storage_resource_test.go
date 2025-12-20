package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccStorageResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccStorageResourceConfig("test"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("cherryservers_storage.test", "id"),
					resource.TestCheckResourceAttr("cherryservers_storage.test", "size", "10"),
					resource.TestCheckResourceAttr("cherryservers_storage.test", "region", "LT-Siauliai"),
					resource.TestCheckResourceAttr("cherryservers_storage.test", "description", "Test storage"),
					resource.TestCheckResourceAttrSet("cherryservers_storage.test", "vlan_id"),
					resource.TestCheckResourceAttrSet("cherryservers_storage.test", "vlan_ip"),
					resource.TestCheckResourceAttrSet("cherryservers_storage.test", "initiator"),
					resource.TestCheckResourceAttrSet("cherryservers_storage.test", "discovery_ip"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "cherryservers_storage.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update description
			{
				Config: testAccStorageResourceConfig("updated"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cherryservers_storage.test", "description", "Updated storage"),
				),
			},
			// Delete testing automatically occurs after the last TestStep
		},
	})
}

func testAccStorageResourceConfig(name string) string {
	return fmt.Sprintf(`
data "cherryservers_project" "main" {
}

resource "cherryservers_storage" "test" {
  project_id  = data.cherryservers_project.main.id
  region      = "LT-Siauliai"
  size        = 10
  description = "%s storage"
}
`, name)
}
