package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"os"
	"regexp"
	"testing"
)

func TestAccIPResource_basic(t *testing.T) {
	teamId := os.Getenv("CHERRY_TEST_TEAM_ID")
	projectName := testProjectNamePrefix + acctest.RandString(5)
	aRecord := generateAlphaString(8)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccIPResourceBasicConfig(projectName, teamId, "LT-Siauliai"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCherryServersIPExists("cherryservers_ip.test_ip_ip"),
					resource.TestCheckResourceAttrSet("cherryservers_ip.test_ip_ip", "id"),
					resource.TestCheckResourceAttr("cherryservers_ip.test_ip_ip", "target_id", "0"),
					resource.TestCheckResourceAttr("cherryservers_ip.test_ip_ip", "target_hostname", ""),
					resource.TestCheckResourceAttr("cherryservers_ip.test_ip_ip", "target_ip_id", ""),
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
				Config: testAccIPResourceBasicUpdateConfig(projectName, teamId, "LT-Siauliai", aRecord),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cherryservers_ip.test_ip_ip", "a_record_effective", aRecord+".cloud.cherryservers.net."),
					resource.TestCheckResourceAttr("cherryservers_ip.test_ip_ip", "ptr_record_effective", "test."),
					resource.TestCheckResourceAttr("cherryservers_ip.test_ip_ip", "tags.env", "test"),
					resource.TestMatchResourceAttr("cherryservers_ip.test_ip_ip", "target_id", regexp.MustCompile(`[0-9]+`)),
					resource.TestMatchResourceAttr("cherryservers_ip.test_ip_ip", "target_hostname", regexp.MustCompile("[a-z]+-[a-z]+")),
					resource.TestCheckResourceAttrSet("cherryservers_ip.test_ip_ip", "target_ip_id"),
				),
			},
			// Delete testing automatically occurs in TestCase

		},
	})
}

func TestAccIPResource_fullConfig(t *testing.T) {
	teamId := os.Getenv("CHERRY_TEST_TEAM_ID")
	aRecord := generateAlphaString(8)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPResourceFullConfig(teamId, aRecord),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCherryServersIPExists("cherryservers_ip.test_ip_ip"),
					resource.TestCheckResourceAttr("cherryservers_ip.test_ip_ip", "a_record_effective", aRecord+".cloud.cherryservers.net."),
					resource.TestCheckResourceAttr("cherryservers_ip.test_ip_ip", "ptr_record_effective", "test."),
					resource.TestCheckResourceAttrSet("cherryservers_ip.test_ip_ip", "target_ip_id"),
					resource.TestMatchResourceAttr("cherryservers_ip.test_ip_ip", "target_id", regexp.MustCompile(`[0-9]+`)),
					resource.TestMatchResourceAttr("cherryservers_ip.test_ip_ip", "address", regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`)),
					resource.TestMatchResourceAttr("cherryservers_ip.test_ip_ip", "address_family", regexp.MustCompile(`^[0-9]`)),
					resource.TestCheckResourceAttrSet("cherryservers_ip.test_ip_ip", "cidr"),
					resource.TestCheckResourceAttr("cherryservers_ip.test_ip_ip", "gateway", ""),
					resource.TestCheckResourceAttr("cherryservers_ip.test_ip_ip", "type", "floating-ip"),
					resource.TestCheckResourceAttrSet("cherryservers_ip.test_ip_ip", "id"),
				),
			},
		},
	})
}

func testAccIPResourceBasicConfig(projectName string, teamId string, region string) string {
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

func testAccIPResourceBasicUpdateConfig(projectName string, teamID string, region string, aRecord string) string {
	return fmt.Sprintf(`
resource "cherryservers_project" "test_ip_project" {
  name = "%s"
  team_id = "%s"
}

resource "cherryservers_server" "test_ip_server" {
  plan = "B1-1-1gb-20s-shared"
  region = "LT-Siauliai"
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

func testAccIPResourceFullConfig(teamId string, aRecord string) string {
	return fmt.Sprintf(`
resource "cherryservers_project" "test_ip_project" {
  name = "%s"
  team_id = "%s"
}

resource "cherryservers_server" "test_ip_server" {
  plan = "B1-1-1gb-20s-shared"
  region = "LT-Siauliai"
  project_id = "${cherryservers_project.test_ip_project.id}"
}

resource "cherryservers_ip" "test_ip_ip" {
  project_id = "${cherryservers_project.test_ip_project.id}"
  region = "LT-Siauliai"
  target_hostname = "${cherryservers_server.test_ip_server.hostname}"
  a_record = "%s"
  ptr_record = "test"
  tags = {
    env = "test"
  }
}
`, testProjectNamePrefix+acctest.RandString(5), teamId, aRecord)
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
