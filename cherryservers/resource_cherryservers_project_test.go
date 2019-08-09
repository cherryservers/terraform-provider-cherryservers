package cherryservers

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccCherryServersProject_Basic(t *testing.T) {
	var projectName = "test_project_" + acctest.RandString(5)
	projectTF := fmt.Sprintf(testAccCheckCherryServersProjectConfigBasic, projectName, teamID)
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCherryServersProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: projectTF,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCherryServersProjectExists("cherryservers_project.foobar"),
					//testAccCheckCherryServersProjectAttributes(projectName),
					resource.TestCheckResourceAttr(
						"cherryservers_project.foobar", "name", projectName),
					resource.TestCheckResourceAttr(
						"cherryservers_project.foobar", "team_id", teamID),
				),
			},
		},
	})
}

func testAccCheckCherryServersProjectDestroy(s *terraform.State) error {
	client, _ := testAccProvider.Meta().(*Config).Client()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cherryservers_project" {
			continue
		}

		id, converr := strconv.Atoi(rs.Primary.Attributes["team_id"])
		if converr != nil {
			return fmt.Errorf("Unable to convert Project ID")
		}
		// Try to find the domain
		res, _, err := client.client.Projects.List(id)

		if err != nil {
			return fmt.Errorf("Project listing error: %#v", err)
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

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		client, _ := testAccProvider.Meta().(*Config).Client()
		projectID, converr := strconv.Atoi(rs.Primary.ID)
		teamIDstr, converr2 := strconv.Atoi(teamID)
		if converr != nil {
			return fmt.Errorf("Unable to convert Project ID")
		}
		if converr2 != nil {
			return fmt.Errorf("Unable to convert Team ID")
		}
		// Try to find the project id
		foundDomain, _, err := client.client.Projects.List(teamIDstr)
		if err != nil {
			return err
		}
		for _, project := range foundDomain {
			if project.ID == projectID {
				return nil
			}
		}
		return fmt.Errorf("Project Record not found")
	}
}

const testAccCheckCherryServersProjectConfigBasic = `
resource "cherryservers_project" "foobar" {
  name = "%s"
  team_id = "%s"
}`
