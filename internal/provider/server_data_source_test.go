package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServerDataSource_basic(t *testing.T) {
	teamID := os.Getenv("CHERRY_TEST_TEAM_ID")
	projectName := "terraform_test_project_" + acctest.RandString(5)
	resourceName := "cherryservers_server.test_server_server"
	dataSourceName := "data.cherryservers_server.test_server_server"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccServerDataSourceConfig(teamID, projectName),
				Check: resource.ComposeAggregateTestCheckFunc(
					//Check by ID.
					resource.TestCheckResourceAttrPair(dataSourceName, "plan", resourceName, "plan"),
					resource.TestCheckResourceAttrPair(dataSourceName, "project_id", resourceName, "project_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "region", resourceName, "region"),
					resource.TestCheckResourceAttrPair(dataSourceName, "hostname", resourceName, "hostname"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "username", resourceName, "username"),
					resource.TestCheckResourceAttrPair(dataSourceName, "password", resourceName, "password"),
					resource.TestCheckResourceAttrPair(dataSourceName, "bmc", resourceName, "bmc"),
					resource.TestCheckResourceAttrPair(dataSourceName, "image", resourceName, "image"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ssh_key_ids", resourceName, "ssh_key_ids"),
					resource.TestCheckResourceAttrPair(dataSourceName, "extra_ip_addresses_ids", resourceName, "extra_ip_addresses_ids"),
					resource.TestCheckResourceAttrPair(dataSourceName, "user_data_file", resourceName, "user_data_file"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tags", resourceName, "tags"),
					resource.TestCheckResourceAttrPair(dataSourceName, "spot_instance", resourceName, "spot_instance"),
					resource.TestCheckResourceAttrPair(dataSourceName, "os_partition_size", resourceName, "os_partition_size"),
					resource.TestCheckResourceAttrPair(dataSourceName, "power_state", resourceName, "power_state"),
					resource.TestCheckResourceAttrPair(dataSourceName, "state", resourceName, "state"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ip_addresses", resourceName, "ip_addresses"),
					//Check by hostname.
					resource.TestCheckResourceAttrPair("data.cherryservers_server.test_server_server_by_hostname", "plan", resourceName, "plan"),
					resource.TestCheckResourceAttrPair("data.cherryservers_server.test_server_server_by_hostname", "project_id", resourceName, "project_id"),
					resource.TestCheckResourceAttrPair("data.cherryservers_server.test_server_server_by_hostname", "region", resourceName, "region"),
					resource.TestCheckResourceAttrPair("data.cherryservers_server.test_server_server_by_hostname", "hostname", resourceName, "hostname"),
					resource.TestCheckResourceAttrPair("data.cherryservers_server.test_server_server_by_hostname", "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair("data.cherryservers_server.test_server_server_by_hostname", "username", resourceName, "username"),
					resource.TestCheckResourceAttrPair("data.cherryservers_server.test_server_server_by_hostname", "password", resourceName, "password"),
					resource.TestCheckResourceAttrPair("data.cherryservers_server.test_server_server_by_hostname", "bmc", resourceName, "bmc"),
					resource.TestCheckResourceAttrPair("data.cherryservers_server.test_server_server_by_hostname", "image", resourceName, "image"),
					resource.TestCheckResourceAttrPair("data.cherryservers_server.test_server_server_by_hostname", "ssh_key_ids", resourceName, "ssh_key_ids"),
					resource.TestCheckResourceAttrPair("data.cherryservers_server.test_server_server_by_hostname", "extra_ip_addresses_ids", resourceName, "extra_ip_addresses_ids"),
					resource.TestCheckResourceAttrPair("data.cherryservers_server.test_server_server_by_hostname", "user_data_file", resourceName, "user_data_file"),
					resource.TestCheckResourceAttrPair("data.cherryservers_server.test_server_server_by_hostname", "tags", resourceName, "tags"),
					resource.TestCheckResourceAttrPair("data.cherryservers_server.test_server_server_by_hostname", "spot_instance", resourceName, "spot_instance"),
					resource.TestCheckResourceAttrPair("data.cherryservers_server.test_server_server_by_hostname", "os_partition_size", resourceName, "os_partition_size"),
					resource.TestCheckResourceAttrPair("data.cherryservers_server.test_server_server_by_hostname", "power_state", resourceName, "power_state"),
					resource.TestCheckResourceAttrPair("data.cherryservers_server.test_server_server_by_hostname", "state", resourceName, "state"),
					resource.TestCheckResourceAttrPair("data.cherryservers_server.test_server_server_by_hostname", "ip_addresses", resourceName, "ip_addresses"),
				),
			},
		},
	})
}

func testAccServerDataSourceConfig(teamID string, projectName string) string {
	return fmt.Sprintf(`
resource "cherryservers_project" "test_server_project" {
  name = "%s"
  team_id = "%s"
}

resource "cherryservers_server" "test_server_server" {
  plan = "cloud_vps_1"
  region = "eu_nord_1"
  project_id = "${cherryservers_project.test_server_project.id}"
}

data "cherryservers_server" "test_server_server" {
  id = cherryservers_server.test_server_server.id
}

data "cherryservers_server" "test_server_server_by_hostname" {
  hostname = "${cherryservers_server.test_server_server.hostname}"
  project_id = "${cherryservers_server.test_server_server.project_id}"
}
`, projectName, teamID)
}
