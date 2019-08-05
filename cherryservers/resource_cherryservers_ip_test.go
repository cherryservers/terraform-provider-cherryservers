package cherryservers

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var floatingIP string
var projectID string

func TestAccCherryServersFloatingIPBasic(t *testing.T) {

	expectedURNRegEx, _ := regexp.Compile(`(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCherryServersFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCherryServersFloatingIPConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCherryServersFloatingIPExists("cherryservers_ip.foobar"),
					resource.TestCheckResourceAttr(
						"cherryservers_ip.foobar", "region", "EU-East-1"),
					resource.TestCheckResourceAttrSet(
						"cherryservers_ip.foobar", "project_id"),
					resource.TestMatchResourceAttr("cherryservers_ip.foobar", "address", expectedURNRegEx),
				),
			},
		},
	})
}
func testAccCheckCherryServersFloatingIPExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		client, _ := testAccProvider.Meta().(*Config).Client()
		projectID = rs.Primary.Attributes["project_id"]
		// Try to find the FloatingIP
		foundFloatingIP, _, err := client.client.IPAddresses.List(projectID)
		if err != nil {
			return err
		}
		if len(foundFloatingIP) == 1 {
			floatingIP = foundFloatingIP[0].Address
			return nil
		}
		return fmt.Errorf("IP not found")
	}
}
func testAccCheckCherryServersFloatingIPDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()
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
			results, _, err := client.client.IPAddresses.List(projectID)
			if len(results) == 0 {
				return nil
			}
			if err != nil {
				return err
			}
			if len(results) != 0 {
				return fmt.Errorf("Floating IP still exists")
			}
		}
	}
	return nil
}

func testAccCheckCherryServersFloatingIPConfigBasic() string {
	res := fmt.Sprintf(`
resource "cherryservers_project" "myproject" {
  team_id = "%s"
  name = "foobar-project-2"
}
resource "cherryservers_ip" "foobar" {
  project_id = "${cherryservers_project.myproject.id}"
  region = "EU-East-1"
}`, teamID)
	return res
}
