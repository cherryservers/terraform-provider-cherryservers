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

func TestAccProjectResource_basic(t *testing.T) {
	teamId := os.Getenv("CHERRY_TEST_TEAM_ID")
	name := "terraform_test_project_" + acctest.RandString(5)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCherryServersProjectDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccProjectResourceConfig(name, teamId),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCherryServersProjectExists("cherryservers_project.test"),
					resource.TestMatchResourceAttr("cherryservers_project.test", "href", regexp.MustCompile("/projects/[0-9]+")),
					resource.TestCheckResourceAttr("cherryservers_project.test", "bgp.enabled", "false"),
					resource.TestCheckResourceAttrSet("cherryservers_project.test", "bgp.local_asn"),
					resource.TestMatchResourceAttr("cherryservers_project.test", "id", regexp.MustCompile("[0-9]+")),
				),
			},
			// ImportState testing
			{
				ResourceName:            "cherryservers_project.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"team_id"},
			},
			// Update and Read testing
			{
				Config: testAccProjectResourceConfig(name+"_update", teamId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cherryservers_project.test", "name", name+"_update"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccProjectResourceConfig(name string, teamId string) string {
	return fmt.Sprintf(`
resource "cherryservers_project" "test" {
  name = "%s"
  team_id = "%s"
}
`, name, teamId)
}

func testAccCheckCherryServersProjectExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("project ID is not set")
		}
		client := testCherryGoClient
		projectID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("unable to convert Project ID")
		}

		// Try to get the project id
		_, _, err = client.Projects.Get(projectID, nil)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckCherryServersProjectDestroy(s *terraform.State) error {
	client := testCherryGoClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cherryservers_project" {
			continue
		}

		projectID, converr := strconv.Atoi(rs.Primary.ID)
		if converr != nil {
			return fmt.Errorf("unable to convert Project ID")
		}
		// Try to get the project
		_, resp, err := client.Projects.Get(projectID, nil)

		if err != nil {
			if is404Error(resp) {
				continue
			}

			return fmt.Errorf("project listing error: %#v", err)
		}

		return fmt.Errorf("project still exists")
	}

	return nil
}
