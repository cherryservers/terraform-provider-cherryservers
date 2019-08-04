package cherryservers
import (
	"fmt"
  "strconv"
  //"github.com/cherryservers/cherrygo"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)
func TestAccCherryServersProject_Basic(t *testing.T){
  var teamID = "35587"
	projectName := fmt.Sprintf(testAccCheckCherryServersProjectConfig_basic, "test_project_"+acctest.RandString(5),teamID)
  fmt.Println(projectName)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCherryServersProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: projectName,
	/*			Check: resource.ComposeTestCheckFunc(
					testAccCheckCherryServersProjectExists("cherryservers_project.foobar"),
					//testAccCheckCherryServersProjectAttributes(projectName),
					resource.TestCheckResourceAttr(
						"cherryservers_project.foobar", "name", projectName),
					resource.TestCheckResourceAttr(
						"cherryservers_project.foobar", "href", "192.168.0.10"),
				),*/
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

    id, converr := strconv.Atoi(rs.Primary.Attributes["team_id"])
		if converr != nil {
			return fmt.Errorf("Unable to convert Project ID")
		}
		// Try to find the domain
		res, _, err := client.Projects.List(id)

		if err != nil {
			return fmt.Errorf("Project listing error: %#v",err)
		}
    if len(res) == 0 {
      return nil
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

    return fmt.Errorf("found %#v",rs)
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*CombinedConfig).Client()
    id, converr := strconv.Atoi(rs.Primary.ID)
		if converr != nil {
			return fmt.Errorf("Unable to convert Project ID")
		}
		// Try to find the domain

		foundDomain, _, err := client.Projects.List(id)

    return fmt.Errorf("found %#v",foundDomain)
		if err != nil {
			return err
		}
    fmt.Println(foundDomain)

		//if foundDomain.Name != rs.Primary.ID {
		//	return fmt.Errorf("Record not found")
		//}

		return nil
	}
}
const testAccCheckCherryServersProjectConfig_basic = `
resource "cherryservers_project" "foobar" {
  name = "%s"
  team_id = "%s"
}`
