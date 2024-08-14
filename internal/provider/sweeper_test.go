package provider

import (
	"fmt"
	"github.com/cherryservers/cherrygo/v3"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testProjectNamePrefix = "terraform_test_project_"

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

// sharedClientForRegion returns a common provider client configured for the specified region
func sharedClientForRegion(region string) (any, error) {
	apiKey := os.Getenv("CHERRY_AUTH_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("CHERRY_AUTH_TOKEN")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("CHERRY_AUTH_KEY or CHERRY_AUTH_TOKEN must be set for acceptance tests")
	}

	teamId := os.Getenv("CHERRY_TEST_TEAM_ID")
	if teamId == "" {
		return nil, fmt.Errorf("CHERRY_TEST_TEAM_ID must be set for acceptance tests")
	}

	//TODO
	//Make user agent version responsive.
	userAgent := fmt.Sprintf("terraform-provider/cherryservers/%s terraform/%s", "test", "1.0.0")
	args := []cherrygo.ClientOpt{cherrygo.WithAuthToken(apiKey), cherrygo.WithUserAgent(userAgent)}
	client, err := cherrygo.NewClient(args...)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// Deleting the projects is sufficient for sweeping, since all other
// resources are belong to projects.
func init() {
	resource.AddTestSweepers("cherryservers_projects", &resource.Sweeper{
		Name: "cherryservers_projects",
		F: func(region string) error {
			client, err := sharedClientForRegion(region)
			if err != nil {
				return fmt.Errorf("error getting client: %s", err)
			}

			conn := client.(*cherrygo.Client)
			teamId, err := strconv.Atoi(os.Getenv("CHERRY_TEST_TEAM_ID"))
			if err != nil {
				return fmt.Errorf("error parsing team id: %s", err)
			}

			projects, _, err := conn.Projects.List(teamId, nil)
			if err != nil {
				return fmt.Errorf("error listing projects: %s", err)
			}

			for _, project := range projects {
				if strings.HasPrefix(project.Name, testProjectNamePrefix) {
					_, err = conn.Projects.Delete(project.ID)
					if err != nil {
						return fmt.Errorf("error deleting project: %s", err)
					}
				}
			}
			return nil
		},
	})
}
