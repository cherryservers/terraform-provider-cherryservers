package cherryservers

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCherryServersDataSourceServer_Basic(t *testing.T) {
	rInt := acctest.RandInt()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckCherryServersServerDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCherryServersDataSourceServerConfigBasic(teamID, rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cherryservers_server.by_id", "hostname"),
					resource.TestCheckResourceAttrSet("data.cherryservers_server.by_id", "plan"),

					resource.TestCheckResourceAttrSet("data.cherryservers_server.by_hostname", "hostname"),
					resource.TestCheckResourceAttrSet("data.cherryservers_server.by_hostname", "plan"),
				),
			},
		},
	})
}

func testAccCheckCherryServersDataSourceServerConfigBasic(teamID string, rInt int) string {
	return fmt.Sprintf(`
resource "cherryservers_project" "server_project" {
	team_id = %s
  	name = "terraform-test-project-data-server"
}

resource "cherryservers_server" "data_server" {
	hostname = "terraform-test-%d"
	plan = "cloud_vps_1"
	project_id = "${cherryservers_project.server_project.id}"
	image     = "ubuntu_22_04"
	region    = "eu_nord_1"
	tags = {
		Name        = "VPS"
		Environment = "Test"
	}
}

data "cherryservers_server" "by_id" {
	server_id = "${cherryservers_server.data_server.id}"
}

data "cherryservers_server" "by_hostname" {
	hostname = "${cherryservers_server.data_server.hostname}"
	project_id = "${cherryservers_project.server_project.id}"
}`, teamID, rInt)
}
