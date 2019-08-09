package cherryservers

import (
	"fmt"
	"strconv"

	//	"strconv"

	"testing"

	"github.com/cherryservers/cherrygo"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var serverID string

func TestAccCherryServersServer_Basic(t *testing.T) {
	var server cherrygo.Servers
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCherryServersServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCherryServersServerConfigBasic(teamID, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCherryServersServerExists("cherryservers_server.foobar", &server),
					testAccCheckCherryServersServerAttributes(&server),
				),
			},
		},
	})
}
func testAccCheckCherryServersServerExists(n string, server *cherrygo.Servers) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// TODO: check this shit
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Server ID is set")
		}
		serverID = rs.Primary.ID
		rs, ok = s.RootModule().Resources["cherryservers_project.myproject"]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Project ID is set")
		}
		projectID = rs.Primary.ID

		client, _ := testAccProvider.Meta().(*Config).Client()
		// Try to find the Server
		servers, _, err := client.client.Servers.List(projectID)
		if err != nil {
			return err
		}
		if len(servers) < 1 {
			return fmt.Errorf("Servers not found")
		}
		for _, mserver := range servers {
			if strconv.Itoa(mserver.ID) == serverID {
				*server = mserver

				return nil
			}
		}
		return fmt.Errorf("Servers not found in root module")
	}
}

func testAccCheckCherryServersServerDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
	rs, ok := s.RootModule().Resources["cherryservers_server.foobar"]
	if !ok {
		return fmt.Errorf("Not found: cherryservers_server.foobar")
	}
	if rs.Primary.ID == "" {
		return fmt.Errorf("No Server ID is set")
	}
	serverID = rs.Primary.ID
	rs, ok = s.RootModule().Resources["cherryservers_project.myproject"]
	if !ok {
		return fmt.Errorf("Not found: cherryservers_project.myproject")
	}
	if rs.Primary.ID == "" {
		return fmt.Errorf("No Project ID is set")
	}
	projectID = rs.Primary.ID
	teamIDstr, converr2 := strconv.Atoi(teamID)

	projects, _, err := client.client.Projects.List(teamIDstr)
	if converr2 != nil {
		return fmt.Errorf("Unable to convert Team ID")
	}
	if err != nil {
		return err
	}
	for _, rs := range projects {
		if strconv.Itoa(rs.ID) == projectID {
			results, _, err := client.client.Servers.List(projectID)
			if len(results) == 0 {
				return nil
			}
			if err != nil {
				return err
			}
			if len(results) != 0 {
				return fmt.Errorf("Server still exists")
			}
		}
	}

	return nil
}

func testAccCheckCherryServersServerAttributes(server *cherrygo.Servers) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if server.Image != "Ubuntu 18.04 64bit" {
			return fmt.Errorf("Bad image_slug: %s", server.Image)
		}

		if server.Plans.ID != 86 { // ID?
			return fmt.Errorf("Bad size_slug: %#v", server.Plans.ID)
		}

		if server.Pricing.Price != 0.2 {
			return fmt.Errorf("Bad price_hourly: %v", server.Pricing.Price)
		}

		if server.Region.Name != "EU-East-1" {
			return fmt.Errorf("Bad region_slug: %s", server.Region.Name)
		}

		return nil
	}
}

func testAccCheckCherryServersServerConfigBasic(teamID string, rInt int) string {
	return fmt.Sprintf(`
resource "cherryservers_project" "myproject" {
  team_id = "%s"
  name = "foobar-project-server-%d"
}

resource "cherryservers_server" "foobar" {
  hostname      = "foo-%d"
  plan_id = "86"
  project_id = "${cherryservers_project.myproject.id}"
  image     = "Ubuntu 18.04 64bit"
  region    = "EU-East-1"
}`, teamID, rInt, rInt)
}
