package cherryservers

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
  "github.com/cherryservers/cherrygo"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)


func TestAccCherryServersServer_Basic(t *testing.T) {
	var server cherrygo.Servers
  var teamID = "35587"
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCherryServersServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCherryServersServerConfig_basic(teamID,rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCherryServersServerExists("cherryservers_server.foobar", &server),
					testAccCheckCherryServersServerAttributes(&server),
					resource.TestCheckResourceAttr(
						"cherryservers_server.foobar", "name", fmt.Sprintf("foo-%d", rInt)),
					resource.TestCheckResourceAttr(
						"cherryservers_server.foobar", "plans", "512mb"),
					resource.TestCheckResourceAttr(
						"cherryservers_server.foobar", "price_hourly", "0.00744"),
					resource.TestCheckResourceAttr(
						"cherryservers_server.foobar", "price_monthly", "5"),
					resource.TestCheckResourceAttr(
						"cherryservers_server.foobar", "image", "centos-7-x64"),
					resource.TestCheckResourceAttr(
						"cherryservers_server.foobar", "region", "EU-East-1"),
					resource.TestCheckResourceAttr(
						"cherryservers_server.foobar", "ipv4_address_private", ""),
				),
			},
		},
	})
}
func testAccCheckCherryServersServerExists(n string, server *cherrygo.Servers) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
    project_id := rs.Primary.Attributes["project_id"]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Server ID is set")
		}

		client := testAccProvider.Meta().(*CombinedConfig).Client()

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		// Try to find the Droplet
		retrieveDroplet, _, err := client.Servers.List(project_id)

		if err != nil {
			return err
		}

		/*if strconv.Itoa(retrieveDroplet.ID) != rs.Primary.ID {
			return fmt.Errorf("Droplet not found")
		}
    */

		return nil
	}
}


func testAccCheckCherryServersServerDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).Client()

	for _, rs := range s.RootModule().Resources {
    project_id := rs.Primary.Attributes["project_id"]
		if rs.Type != "cherryservers_server" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		// Try to find the Server
		_, _, err = client.Servers.List(project_id)

		// Wait

		if err != nil && !strings.Contains(err.Error(), "404") {
			return fmt.Errorf(
				"Error waiting for server (%s) to be destroyed: %s",
				rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckCherryServersServerAttributes(server *cherrygo.Servers) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if server.Image != "centos-7-x64" {
			return fmt.Errorf("Bad image_slug: %s", server.Image)
		}

		if server.Plans.Name != "512mb" {
			return fmt.Errorf("Bad size_slug: %s", server.Plans)
		}

		if server.Pricing.Price != 0.00744 {
			return fmt.Errorf("Bad price_hourly: %v", server.Pricing.Price)
		}

		if server.Region.Name != "EU-East-1" {
			return fmt.Errorf("Bad region_slug: %s", server.Region.Name)
		}

		return nil
	}
}

func testAccCheckCherryServersServerConfig_basic(teamID string,rInt int) string {
	return fmt.Sprintf(`
resource "cherryservers_project" "myproject" {
  team_id = "%v"
  name = "foobar-project"
}
resource "cherryservers_server" "foobar" {
  name      = "foo-%d"
  plans = "512mb"
  image     = "ubuntu-1804"
  region    = "EU-East-1"
}`, teamID,rInt)
}
