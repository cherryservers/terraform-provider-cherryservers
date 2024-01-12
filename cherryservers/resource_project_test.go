package cherryservers

import (
	"fmt"
	"strconv"
	"testing"

	acctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCherryServersProject_Basic(t *testing.T) {
	var projectName = "terraform_test_project_" + acctest.RandString(5)
	projectTF := fmt.Sprintf(testAccCheckCherryServersProjectConfigBasic, projectName, teamID)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCherryServersProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: projectTF,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCherryServersProjectExists("cherryservers_project.foobar"),
					resource.TestCheckResourceAttr("cherryservers_project.foobar", "name", projectName),
					resource.TestCheckResourceAttr("cherryservers_project.foobar", "team_id", teamID),
				),
			},
		},
	})
}

func testAccCheckCherryServersProjectDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client).cherrygoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cherryservers_project" {
			continue
		}

		projectID, converr := strconv.Atoi(rs.Primary.ID)
		if converr != nil {
			return fmt.Errorf("Unable to convert Project ID")
		}
		// Try to get the project
		_, resp, err := client.Projects.Get(projectID, nil)

		if err != nil {
			if is404Error(resp) {
				return nil
			}

			return fmt.Errorf("Project listing error: %#v", err)
		}
	}

	return nil
}

func testAccCheckCherryServersProjectExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Project ID is not set")
		}
		client := testAccProvider.Meta().(*Client).cherrygoClient()
		projectID, converr := strconv.Atoi(rs.Primary.ID)
		if converr != nil {
			return fmt.Errorf("Unable to convert Project ID")
		}

		// Try to get the project id
		_, _, err := client.Projects.Get(projectID, nil)
		if err != nil {
			return err
		}

		return nil
	}
}

const testAccCheckCherryServersProjectConfigBasic = `
resource "cherryservers_project" "foobar" {
  name = "%s"
  team_id = "%s"
}`
