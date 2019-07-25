package cherryservers
import (
	"fmt"
  //"github.com/cherryservers/cherrygo"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)
func TestAccCherryServersProject_Basic(t *testing.T){
	projectName := fmt.Sprintf(testAccCheckCherryServersProjectConfig_basic, "test_project_"+acctest.RandString(5))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCherryServersProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: projectName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCherryServersProjectExists("cherryservers_project.foobar"),
					testAccCheckCherryServersProjectAttributes(projectName),
					resource.TestCheckResourceAttr(
						"cherryservers_project.foobar", "name", projectName),
					resource.TestCheckResourceAttr(
						"cherryservers_project.foobar", "href", "192.168.0.10"),
				),
			},
		},
	})
}

func testAccCheckCherryServersProjectDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).Client()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cherryservers_project" {
			continue
		}

		// Try to find the domain
		_, _, err := client.Projects.List(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Domain still exists")
		}
	}

	return nil
}
func testAccCheckCherryServersProjectAttributes(projectName string, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if projectName != name {
			return fmt.Errorf("Bad name: %s", projectName)
		}

		return nil
	}
}
func testAccCheckCherryServersProjectExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*CombinedConfig).Client()

		foundDomain, _, err := client.Projects.List(rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundDomain.Name != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		return nil
	}
}
const testAccCheckCherryServersProjectConfig_basic = `
resource "cherryservers_project" "foobar" {
	name       = "%s"
}`
