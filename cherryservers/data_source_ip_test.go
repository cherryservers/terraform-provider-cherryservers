package cherryservers

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCherryServersDataSourceIp_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckCherryServersFloatingIPDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCherryServersDataSourceIpConfigBasic(teamID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cherryservers_ip.by_id", "id"),
					resource.TestCheckResourceAttrSet("data.cherryservers_ip.by_id", "region"),
					resource.TestCheckResourceAttrSet("data.cherryservers_ip.by_id", "address"),
					resource.TestCheckResourceAttrSet("data.cherryservers_ip.by_id", "type"),

					resource.TestCheckResourceAttrSet("data.cherryservers_ip.by_address", "id"),
					resource.TestCheckResourceAttrSet("data.cherryservers_ip.by_address", "region"),
					resource.TestCheckResourceAttrSet("data.cherryservers_ip.by_address", "address"),
					resource.TestCheckResourceAttrSet("data.cherryservers_ip.by_address", "type"),
				),
			},
		},
	})
}

func testAccCheckCherryServersDataSourceIpConfigBasic(teamID string) string {
	return fmt.Sprintf(`
resource "cherryservers_project" "ip_project" {
	team_id = %s
  	name = "terraform-test-data-project-ip"
}

resource "cherryservers_ip" "foobar" {
  project_id = "${cherryservers_project.ip_project.id}"
  region = "eu_nord_1"
}

data "cherryservers_ip" "by_id" {
	ip_id = "${cherryservers_ip.foobar.id}"
}

data "cherryservers_ip" "by_address" {
	ip_address = "${cherryservers_ip.foobar.address}"
	project_id = "${cherryservers_project.ip_project.id}"
}`, teamID)
}
