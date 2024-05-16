package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPResource_basic(t *testing.T) {
	teamId := os.Getenv("CHERRY_TEST_TEAM_ID")
	projectName := "terraform_test_project_" + acctest.RandString(5)
	aRecord := acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCherryServersIPDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIPResourceConfig(projectName, teamId, "eu_nord_1"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCherryServersIPExists("cherryservers_ip.test_ip_ip"),
					resource.TestCheckResourceAttrSet("cherryservers_ip.test_ip_ip", "id"),
					resource.TestCheckResourceAttr("cherryservers_ip.test_ip_ip", "target_id", "0"),
					resource.TestCheckResourceAttr("cherryservers_ip.test_ip_ip", "target_hostname", ""),
					resource.TestCheckResourceAttr("cherryservers_ip.test_ip_ip", "route_ip_id", ""),
					resource.TestCheckResourceAttr("cherryservers_ip.test_ip_ip", "ddos_scrubbing", "false"),
					resource.TestMatchResourceAttr("cherryservers_ip.test_ip_ip", "address", regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`)),
					resource.TestMatchResourceAttr("cherryservers_ip.test_ip_ip", "address_family", regexp.MustCompile(`^[0-9]`)),
					resource.TestCheckResourceAttrSet("cherryservers_ip.test_ip_ip", "cidr"),
					resource.TestCheckResourceAttr("cherryservers_ip.test_ip_ip", "gateway", ""),
					resource.TestCheckResourceAttr("cherryservers_ip.test_ip_ip", "type", "floating-ip"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "cherryservers_ip.test_ip_ip",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccIPResourceUpdateConfig(projectName, teamId, "eu_nord_1", aRecord),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cherryservers_ip.test_ip_ip", "a_record_actual", aRecord+".cloud.cherryservers.net."),
					resource.TestCheckResourceAttr("cherryservers_ip.test_ip_ip", "ptr_record_actual", "test."),
					resource.TestCheckResourceAttr("cherryservers_ip.test_ip_ip", "tags.env", "test"),
					resource.TestMatchResourceAttr("cherryservers_ip.test_ip_ip", "target_id", regexp.MustCompile(`[0-9]+`)),
					resource.TestMatchResourceAttr("cherryservers_ip.test_ip_ip", "target_hostname", regexp.MustCompile("[a-z]+-[a-z]+")),
					resource.TestCheckResourceAttrSet("cherryservers_ip.test_ip_ip", "route_ip_id"),
				),
			},
			// Delete testing automatically occurs in TestCase

		},
	})
}

func testAccIPResourceConfig(projectName string, teamId string, region string) string {
	return fmt.Sprintf(`
resource "cherryservers_project" "test_ip_project" {
  name = "%s"
  team_id = "%s"
}

resource "cherryservers_ip" "test_ip_ip" {
  project_id = "${cherryservers_project.test_ip_project.id}"
  region = "%s"
}
`, projectName, teamId, region)
}

func testAccIPResourceUpdateConfig(projectName string, teamID string, region string, aRecord string) string {
	return fmt.Sprintf(`
resource "cherryservers_project" "test_ip_project" {
  name = "%s"
  team_id = "%s"
}

resource "cherryservers_server" "test_ip_server" {
  plan = "cloud_vps_1"
  region = "eu_nord_1"
  project_id = "${cherryservers_project.test_ip_project.id}"
}

resource "cherryservers_ip" "test_ip_ip" {
  project_id = "${cherryservers_project.test_ip_project.id}"
  region = "%s"
  target_id = "${cherryservers_server.test_ip_server.id}"
  a_record = "%s"
  ptr_record = "test"
  tags = {
    env = "test"
  }
}
`, projectName, teamID, region, aRecord)
}

func testAccCheckCherryServersIPExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("IP ID is not set")
		}
		client := testCherryGoClient

		// Try to get the IP
		_, _, err := client.IPAddresses.Get(rs.Primary.ID, nil)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckCherryServersIPDestroy(s *terraform.State) error {
	client := testCherryGoClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cherryservers_ip" {
			continue
		}

		// There is a delay with IP destruction
		time.Sleep(3 * time.Second)

		// Try to get the project
		_, resp, err := client.IPAddresses.Get(rs.Primary.ID, nil)

		if err != nil {
			if is404Error(resp) {
				continue
			}
			//API returns access denied instead of resource missing
			if is403Error(resp) {
				continue
			}

			return fmt.Errorf("IP listing error: %#v", err)
		}

		return fmt.Errorf("IP still exists" + " ID:" + rs.Primary.ID)
	}

	return nil
}
