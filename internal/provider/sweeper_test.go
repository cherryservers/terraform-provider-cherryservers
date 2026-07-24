package provider

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Deleting the projects is sufficient for sweeping, since all other
// resources are belong to projects.
func init() {
	resource.AddTestSweepers("cherryservers_projects", &resource.Sweeper{
		Name: "cherryservers_projects",
		F: func(region string) error {
			client, err := setupClient()
			if err != nil {
				return fmt.Errorf("error getting client: %s", err)
			}

			teamId, err := strconv.Atoi(os.Getenv("CHERRY_TEST_TEAM_ID"))
			if err != nil {
				return fmt.Errorf("error parsing team id: %s", err)
			}

			projects, _, err := client.Projects.List(teamId, nil)
			if err != nil {
				return fmt.Errorf("error listing projects: %s", err)
			}

			for _, project := range projects {
				if strings.HasPrefix(project.Name, testProjectNamePrefix) {
					_, err = client.Projects.Delete(project.ID)
					if err != nil {
						return fmt.Errorf("error deleting project: %s", err)
					}
				}
			}
			return nil
		},
	})
}
