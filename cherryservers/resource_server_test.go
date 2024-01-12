package cherryservers

import (
	"fmt"
	"strconv"
	"testing"

	acctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCherryServersServer_Basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCherryServersServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCherryServersServerConfigBasic(teamID, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCherryServersServerExists("cherryservers_server.foobar"),
					resource.TestCheckResourceAttr("cherryservers_server.foobar", "plan", "cloud_vps_1"),
					resource.TestCheckResourceAttr("cherryservers_server.foobar", "hostname", fmt.Sprintf("terraform-test-%d", rInt)),
					resource.TestCheckResourceAttrSet("cherryservers_server.foobar", "name"),
					resource.TestCheckResourceAttrSet("cherryservers_server.foobar", "power_state"),
					resource.TestCheckResourceAttrSet("cherryservers_server.foobar", "state"),
					//resource.TestCheckResourceAttrSet("cherryservers_server.foobar", "tags.#"),
					//resource.TestCheckResourceAttrSet("cherryservers_server.foobar", "ip_addresses.#"),
				),
			},
		},
	})
}
func testAccCheckCherryServersServerExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource cherryservers_server.foobar not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Server ID is set")
		}
		serverID, _ := strconv.Atoi(rs.Primary.ID)

		client := testAccProvider.Meta().(*Client).cherrygoClient()
		// Try to get the Server
		_, _, err := client.Servers.Get(serverID, nil)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckCherryServersServerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client).cherrygoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cherryservers_server" {
			continue
		}

		serverID, converr := strconv.Atoi(rs.Primary.ID)
		if converr != nil {
			return fmt.Errorf("Unable to convert Server ID")
		}

		// Try to find the key
		_, resp, err := client.Servers.Get(serverID, nil)
		if err != nil {
			if is404Error(resp) {
				return nil
			}

			return fmt.Errorf("Error getting Server ID (%s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckCherryServersServerConfigBasic(teamID string, rInt int) string {
	return fmt.Sprintf(`
resource "cherryservers_project" "server_project" {
  team_id = %s
  name = "terraform-test-project-server"
}

resource "cherryservers_server" "foobar" {
  hostname      = "terraform-test-%d"
  plan = "cloud_vps_1"
  project_id = "${cherryservers_project.server_project.id}"
  image     = "ubuntu_22_04"
  region    = "eu_nord_1"
  tags = {
    Name        = "VPS"
    Environment = "Test"
  }
}`, teamID, rInt)
}
