package cherryservers

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCherryServersDataSourceProject_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckCherryServersProjectDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "cherryservers_project" "project" {
					name = "tf-tests-terraform-account-data-project"
					team_id = "%s"
				}
				data cherryservers_project "by_id" {
						project_id = cherryservers_project.project.id
				}`, teamID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cherryservers_project.by_id", "name"),
					resource.TestCheckResourceAttrSet("data.cherryservers_project.by_id", "project_id"),
				),
			},
		},
	})
}
