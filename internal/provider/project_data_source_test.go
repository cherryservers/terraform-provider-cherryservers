package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProjectDataSource(t *testing.T) {
	teamId := os.Getenv("CHERRY_TEST_TEAM_ID")
	name := testProjectNamePrefix + acctest.RandString(5)
	resourceName := "cherryservers_project.test_data"
	datasourceName := "data.cherryservers_project.test_data"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccProjectDataSourceConfig(name, teamId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(datasourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(datasourceName, "bgp.enabled", resourceName, "bgp.enabled"),
					resource.TestCheckResourceAttrPair(datasourceName, "bgp.local_asn", resourceName, "bgp.local_asn"),
				),
			},
		},
	})
}

func testAccProjectDataSourceConfig(name string, teamId string) string {
	return fmt.Sprintf(`
resource "cherryservers_project" "test_data" {
  name = "%s"
  team_id = "%s"
}
data "cherryservers_project" "test_data" {
  id = cherryservers_project.test_data.id
}
`, name, teamId)
}
