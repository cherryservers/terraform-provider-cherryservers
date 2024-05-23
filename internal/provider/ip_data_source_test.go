package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIpDataSource_basic(t *testing.T) {
	teamID := os.Getenv("CHERRY_TEST_TEAM_ID")
	projectName := "terraform_test_project_" + acctest.RandString(5)
	resourceName := "cherryservers_ip.test_ip_ip"
	dataSourceName := "data.cherryservers_ip.test_ip_ip"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccIpDataSourceConfig(teamID, projectName),
				Check: resource.ComposeAggregateTestCheckFunc(
					//By ID checks.
					resource.TestCheckResourceAttrPair(dataSourceName, "project_id", resourceName, "project_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "region", resourceName, "region"),
					resource.TestCheckResourceAttrPair(dataSourceName, "target_id", resourceName, "target_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "target_hostname", resourceName, "target_hostname"),
					resource.TestCheckResourceAttrPair(dataSourceName, "route_ip_id", resourceName, "route_ip_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ddos_scrubbing", resourceName, "ddos_scrubbing"),
					resource.TestCheckResourceAttrPair(dataSourceName, "a_record_actual", resourceName, "a_record_actual"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ptr_record_actual", resourceName, "ptr_record_actual"),
					resource.TestCheckResourceAttrPair(dataSourceName, "address", resourceName, "address"),
					resource.TestCheckResourceAttrPair(dataSourceName, "address_family", resourceName, "address_family"),
					resource.TestCheckResourceAttrPair(dataSourceName, "cidr", resourceName, "cidr"),
					resource.TestCheckResourceAttrPair(dataSourceName, "gateway", resourceName, "gateway"),
					resource.TestCheckResourceAttrPair(dataSourceName, "type", resourceName, "type"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tags", resourceName, "tags"),
					//By address checks.
					resource.TestCheckResourceAttrPair("data.cherryservers_ip.test_ip_ip_by_address", "project_id", resourceName, "project_id"),
					resource.TestCheckResourceAttrPair("data.cherryservers_ip.test_ip_ip_by_address", "region", resourceName, "region"),
					resource.TestCheckResourceAttrPair("data.cherryservers_ip.test_ip_ip_by_address", "target_id", resourceName, "target_id"),
					resource.TestCheckResourceAttrPair("data.cherryservers_ip.test_ip_ip_by_address", "target_hostname", resourceName, "target_hostname"),
					resource.TestCheckResourceAttrPair("data.cherryservers_ip.test_ip_ip_by_address", "route_ip_id", resourceName, "route_ip_id"),
					resource.TestCheckResourceAttrPair("data.cherryservers_ip.test_ip_ip_by_address", "ddos_scrubbing", resourceName, "ddos_scrubbing"),
					resource.TestCheckResourceAttrPair("data.cherryservers_ip.test_ip_ip_by_address", "a_record_actual", resourceName, "a_record_actual"),
					resource.TestCheckResourceAttrPair("data.cherryservers_ip.test_ip_ip_by_address", "ptr_record_actual", resourceName, "ptr_record_actual"),
					resource.TestCheckResourceAttrPair("data.cherryservers_ip.test_ip_ip_by_address", "address", resourceName, "address"),
					resource.TestCheckResourceAttrPair("data.cherryservers_ip.test_ip_ip_by_address", "address_family", resourceName, "address_family"),
					resource.TestCheckResourceAttrPair("data.cherryservers_ip.test_ip_ip_by_address", "cidr", resourceName, "cidr"),
					resource.TestCheckResourceAttrPair("data.cherryservers_ip.test_ip_ip_by_address", "gateway", resourceName, "gateway"),
					resource.TestCheckResourceAttrPair("data.cherryservers_ip.test_ip_ip_by_address", "type", resourceName, "type"),
					resource.TestCheckResourceAttrPair("data.cherryservers_ip.test_ip_ip_by_address", "tags", resourceName, "tags"),
				),
			},
		},
	})
}

func testAccIpDataSourceConfig(teamID string, projectName string) string {
	return fmt.Sprintf(`
resource "cherryservers_project" "test_ip_project" {
  name = "%s"
  team_id = "%s"
}

resource "cherryservers_ip" "test_ip_ip" {
  region = "eu_nord_1"
  project_id = "${cherryservers_project.test_ip_project.id}"
}

data "cherryservers_ip" "test_ip_ip" {
  id = cherryservers_ip.test_ip_ip.id
}

data "cherryservers_ip" "test_ip_ip_by_address" {
  address = "${cherryservers_ip.test_ip_ip.address}"
  project_id = "${cherryservers_project.test_ip_project.id}"
}
`, projectName, teamID)
}
