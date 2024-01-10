package cherryservers

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCherryServersFloatingIPBasic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCherryServersFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCherryServersFloatingIPConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCherryServersFloatingIPExists("cherryservers_ip.simple_ip"),
					resource.TestCheckResourceAttr("cherryservers_ip.simple_ip", "region", "eu_nord_1"),
					resource.TestCheckResourceAttrSet("cherryservers_ip.simple_ip", "project_id"),
					resource.TestCheckResourceAttrSet("cherryservers_ip.simple_ip", "address"),
					resource.TestCheckResourceAttrSet("cherryservers_ip.simple_ip", "cidr"),
					resource.TestCheckResourceAttrSet("cherryservers_ip.simple_ip", "address_family"),
					resource.TestCheckResourceAttrSet("cherryservers_ip.simple_ip", "ddos_scrubbing"),
					resource.TestCheckResourceAttrSet("cherryservers_ip.simple_ip", "type"),
					//resource.TestCheckResourceAttrSet("cherryservers_ip.simple_ip", "tags.1"),
					//resource.TestCheckResourceAttrSet("cherryservers_ip.simple_ip", "tags.2"),
				),
			},
		},
	})
}

func TestAccCherryServersAttachFloatingIPToServerHostname(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCherryServersFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCherryServersAttachFloatingIPToServerHostnameConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCherryServersFloatingIPExists("cherryservers_ip.attached_ip"),
					resource.TestCheckResourceAttr("cherryservers_ip.attached_ip", "target_hostname", "terraform-test-ip"),
				),
			},
		},
	})
}

func testAccCheckCherryServersFloatingIPExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		client := testAccProvider.Meta().(*Client).cherrygoClient()
		ipID := rs.Primary.ID

		// Try to get the FloatingIP
		_, _, err := client.IPAddresses.Get(ipID, nil)
		if err != nil {
			return err
		}

		return nil
	}
}
func testAccCheckCherryServersFloatingIPDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client).cherrygoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cherryservers_ip" {
			continue
		}

		ipID := rs.Primary.ID
		if ipID == "" {
			return fmt.Errorf("IP address ID is not set")
		}

		// Try to find the key
		_, resp, err := client.IPAddresses.Get(ipID, nil)
		if err != nil {
			if is404Error(resp) {
				return nil
			}

			return fmt.Errorf("Error getting IP address (%s)", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckCherryServersFloatingIPConfigBasic() string {
	return fmt.Sprintf(`
resource "cherryservers_project" "ip_project" {
  team_id = "%s"
  name = "terraform-test-project-ip"
}
resource "cherryservers_ip" "simple_ip" {
  project_id = "${cherryservers_project.ip_project.id}"
  region = "eu_nord_1"
  tags = {
    Name        = "IP"
    Environment = "Test"
  }
}`, teamID)
}

func testAccCherryServersAttachFloatingIPToServerHostnameConfig() string {
	return fmt.Sprintf(`
resource "cherryservers_project" "ip_project2" {
  team_id = "%s"
  name = "terraform-test-project-ip2"
}
resource "cherryservers_server" "ip-server" {
  plan = "cloud_vps_1"
  hostname = "terraform-test-ip"
  project_id = "${cherryservers_project.ip_project2.id}"
  region    = "eu_nord_1"
}
resource "cherryservers_ip" "attached_ip" {
  project_id = "${cherryservers_project.ip_project2.id}"
  target_hostname = "${cherryservers_server.ip-server.hostname}"
  region = "eu_nord_1"
}`, teamID)
}
