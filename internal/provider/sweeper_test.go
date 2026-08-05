package provider

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/cherryservers/cherrygo/v3"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testProjectNamePrefix = "terraform_test_project_"

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

// sharedClient returns a common provider client.
func sharedClient() (any, error) {
	apiKey := os.Getenv(apiKeyVar)
	if apiKey == "" {
		return nil, fmt.Errorf("%s must be set for acceptance tests", apiKeyVar)
	}

	teamId := os.Getenv("CHERRY_TEST_TEAM_ID")
	if teamId == "" {
		return nil, fmt.Errorf("CHERRY_TEST_TEAM_ID must be set for acceptance tests")
	}

	// TODO
	// Make user agent version responsive.
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
			client, err := sharedClient()
			if err != nil {
				return fmt.Errorf("error getting client: %s", err)
			}

			conn, ok := client.(*cherrygo.Client)
			if !ok {
				return fmt.Errorf("expected cherrygo.Client, got %T", client)
			}
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
