package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"os"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServerResource_basic(t *testing.T) {
	serverResourceName := "terraform_test_server_" + acctest.RandString(5)
	projectName := "terraform_test_project_" + acctest.RandString(5)
	testPlan := "cloud_vps_1"
	testRegion := "eu_nord_1"
	teamID := os.Getenv("CHERRY_TEST_TEAM_ID")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCherryServersServerDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccServerResourceConfigOnlyReq(projectName, testPlan, testRegion, serverResourceName, teamID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCherryServersServerExists("cherryservers_server."+serverResourceName),
					resource.TestCheckResourceAttr("cherryservers_server."+serverResourceName, "bmc.password", ""),
					resource.TestCheckResourceAttr("cherryservers_server."+serverResourceName, "bmc.user", ""),
					resource.TestMatchResourceAttr("cherryservers_server."+serverResourceName, "hostname", regexp.MustCompile("[a-z]+-[a-z]+")),
					resource.TestMatchResourceAttr("cherryservers_server."+serverResourceName, "id", regexp.MustCompile("[0-9]+")),
					resource.TestMatchResourceAttr("cherryservers_server."+serverResourceName, "ip_addresses.0.address", regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`)),
					resource.TestMatchResourceAttr("cherryservers_server."+serverResourceName, "ip_addresses.1.address", regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`)),
					resource.TestCheckResourceAttrSet("cherryservers_server."+serverResourceName, "name"),
					resource.TestCheckResourceAttrSet("cherryservers_server."+serverResourceName, "password"),
					resource.TestCheckResourceAttr("cherryservers_server."+serverResourceName, "power_state", "on"),
					resource.TestMatchResourceAttr("cherryservers_server."+serverResourceName, "project_id", regexp.MustCompile(`[0-9]+`)),
					resource.TestCheckResourceAttr("cherryservers_server."+serverResourceName, "spot_instance", "false"),
					resource.TestCheckResourceAttr("cherryservers_server."+serverResourceName, "state", "active"),
					resource.TestCheckResourceAttr("cherryservers_server."+serverResourceName, "username", "root"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "cherryservers_server." + serverResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccServerResourceConfigUpdate(projectName, testPlan, testRegion, serverResourceName, teamID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cherryservers_server."+serverResourceName, "name", "update"),
					resource.TestCheckResourceAttr("cherryservers_server."+serverResourceName, "hostname", "server-update-test"),
					resource.TestCheckResourceAttr("cherryservers_server."+serverResourceName, "tags.env", "test"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCheckCherryServersServerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("server ID is not set")
		}
		client := testCherryGoClient
		serverID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("unable to convert Server ID")
		}

		// Try to get the server id
		_, _, err = client.Servers.Get(serverID, nil)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckCherryServersServerDestroy(s *terraform.State) error {
	client := testCherryGoClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cherryservers_server" {
			continue
		}

		serverID, converr := strconv.Atoi(rs.Primary.ID)
		if converr != nil {
			return fmt.Errorf("unable to convert Server ID")
		}

		server, resp, err := client.Servers.Get(serverID, nil)

		if err != nil {
			if is404Error(resp) {
				continue
			}

			return fmt.Errorf("server listing error: %#v", err)
		}

		if server.State != "terminating" {
			return fmt.Errorf("server state is not terminating: %s", server.State)
		}
	}
	return nil
}

func testAccServerResourceConfigOnlyReq(projectName string, plan string, region string, serverResourceName string, teamID string) string {
	return fmt.Sprintf(`
resource "cherryservers_project" "test_server_project" {
  name = "%s"
  team_id = "%s"
}

resource "cherryservers_server" "%s" {
  region = "%s"
  plan = "%s"
  project_id = "${cherryservers_project.test_server_project.id}"
}
`, projectName, teamID, serverResourceName, region, plan)
}

func testAccServerResourceConfigUpdate(projectName string, plan string, region string, serverResourceName string, teamID string) string {
	return fmt.Sprintf(`
resource "cherryservers_project" "test_server_project" {
  name = "%s"
  team_id = "%s"
}

resource "cherryservers_server" "%s" {
  region = "%s"
  plan = "%s"
  project_id = "${cherryservers_project.test_server_project.id}"
  name = "update"
  hostname = "server-update-test"
  tags = {
    env = "test"
  }
}
`, projectName, teamID, serverResourceName, region, plan)
}
