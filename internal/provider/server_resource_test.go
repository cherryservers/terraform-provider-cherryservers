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

func TestAccServerResource_basic(t *testing.T) {
	serverResourceName := "terraform_test_server_" + acctest.RandString(5)
	projectName := testProjectNamePrefix + acctest.RandString(5)
	testPlan := "cloud_vps_1"
	testRegion := "eu_nord_1"
	teamID := os.Getenv("CHERRY_TEST_TEAM_ID")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCherryServersServerDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccServerResourceConfigOnlyReq(projectName, testPlan, testRegion, serverResourceName, teamID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCherryServersServerExists("cherryservers_server."+serverResourceName),
					resource.TestCheckResourceAttr("cherryservers_server."+serverResourceName, "bmc.password", ""),
					resource.TestCheckResourceAttr("cherryservers_server."+serverResourceName, "bmc.user", ""),
					resource.TestMatchResourceAttr("cherryservers_server."+serverResourceName, "hostname", regexp.MustCompile("[a-z]+-[a-z]+")),
					resource.TestMatchResourceAttr("cherryservers_server."+serverResourceName, "id", regexp.MustCompile("[0-9]+")),
					resource.TestMatchResourceAttr("cherryservers_server."+serverResourceName, "ip_addresses.0.address", regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`)),
					resource.TestMatchResourceAttr("cherryservers_server."+serverResourceName, "ip_addresses.1.address", regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`)),
					resource.TestCheckResourceAttrSet("cherryservers_server."+serverResourceName, "name"),
					resource.TestCheckResourceAttrSet("cherryservers_server."+serverResourceName, "password"),
					resource.TestCheckResourceAttr("cherryservers_server."+serverResourceName, "power_state", "on"),
					resource.TestMatchResourceAttr("cherryservers_server."+serverResourceName, "project_id", regexp.MustCompile(`[0-9]+`)),
					resource.TestCheckResourceAttr("cherryservers_server."+serverResourceName, "spot_instance", "false"),
					resource.TestCheckResourceAttr("cherryservers_server."+serverResourceName, "state", "active"),
					resource.TestCheckResourceAttr("cherryservers_server."+serverResourceName, "username", "root"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "cherryservers_server." + serverResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccServerResourceConfigUpdate(projectName, testPlan, testRegion, serverResourceName, teamID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cherryservers_server."+serverResourceName, "name", "update"),
					resource.TestCheckResourceAttr("cherryservers_server."+serverResourceName, "hostname", "server-update-test"),
					resource.TestCheckResourceAttr("cherryservers_server."+serverResourceName, "tags.env", "test"),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

func TestAccServerResource_fullConfig(t *testing.T) {
	projectName := testProjectNamePrefix + acctest.RandString(5)
	teamID := os.Getenv("CHERRY_TEST_TEAM_ID")
	label := "terraform_test_ssh_" + acctest.RandString(5)
	publicKey, _, err := acctest.RandSSHKeyPair("cherryservers@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCherryServersServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServerResourceFullConfig(projectName, teamID, label, publicKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCherryServersServerExists("cherryservers_server.test_server_server"),
					resource.TestCheckResourceAttr("cherryservers_server.test_server_server", "bmc.password", ""),
					resource.TestCheckResourceAttr("cherryservers_server.test_server_server", "bmc.user", ""),
					resource.TestCheckResourceAttr("cherryservers_server.test_server_server", "image", "ubuntu_22_04"),
					resource.TestMatchResourceAttr("cherryservers_server.test_server_server", "id", regexp.MustCompile("[0-9]+")),
					resource.TestMatchResourceAttr("cherryservers_server.test_server_server", "ip_addresses.0.address", regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`)),
					resource.TestMatchResourceAttr("cherryservers_server.test_server_server", "ip_addresses.1.address", regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`)),
					resource.TestCheckResourceAttrSet("cherryservers_server.test_server_server", "name"),
					resource.TestCheckResourceAttrSet("cherryservers_server.test_server_server", "password"),
					resource.TestCheckResourceAttr("cherryservers_server.test_server_server", "power_state", "on"),
					resource.TestMatchResourceAttr("cherryservers_server.test_server_server", "project_id", regexp.MustCompile(`[0-9]+`)),
					resource.TestCheckResourceAttr("cherryservers_server.test_server_server", "spot_instance", "false"),
					resource.TestCheckResourceAttr("cherryservers_server.test_server_server", "state", "active"),
					resource.TestCheckResourceAttr("cherryservers_server.test_server_server", "username", "root"),
				),
			},
		},
	})
}

func testAccCheckCherryServersServerExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("server ID is not set")
		}
		client := testCherryGoClient
		serverID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("unable to convert Server ID")
		}

		// Try to get the server id
		_, _, err = client.Servers.Get(serverID, nil)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckCherryServersServerDestroy(s *terraform.State) error {
	client := testCherryGoClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cherryservers_server" {
			continue
		}

		serverID, converr := strconv.Atoi(rs.Primary.ID)
		if converr != nil {
			return fmt.Errorf("unable to convert Server ID")
		}

		server, resp, err := client.Servers.Get(serverID, nil)

		if err != nil {
			if is404Error(resp) {
				continue
			}

			return fmt.Errorf("server listing error: %#v", err)
		}

		if server.State != "terminating" {
			return fmt.Errorf("server state is not terminating: %s", server.State)
		}
	}
	return nil
}

func testAccServerResourceConfigOnlyReq(projectName string, plan string, region string, serverResourceName string, teamID string) string {
	return fmt.Sprintf(`
resource "cherryservers_project" "test_server_project" {
  name = "%s"
  team_id = "%s"
}

resource "cherryservers_server" "%s" {
  region = "%s"
  plan = "%s"
  project_id = "${cherryservers_project.test_server_project.id}"
}
`, projectName, teamID, serverResourceName, region, plan)
}

func testAccServerResourceConfigUpdate(projectName string, plan string, region string, serverResourceName string, teamID string) string {
	return fmt.Sprintf(`
resource "cherryservers_project" "test_server_project" {
  name = "%s"
  team_id = "%s"
}

resource "cherryservers_server" "%s" {
  region = "%s"
  plan = "%s"
  project_id = "${cherryservers_project.test_server_project.id}"
  name = "update"
  hostname = "server-update-test"
  tags = {
    env = "test"
  }
}
`, projectName, teamID, serverResourceName, region, plan)
}

func testAccServerResourceFullConfig(projectName string, teamID string, sshKeyLabel string, sshKeyPublicKey string) string {
	return fmt.Sprintf(`
resource "cherryservers_project" "test_server_project" {
  name = "%s"
  team_id = "%s"
}

resource "cherryservers_ssh_key" "test_server_ssh_key" {
  label = "%s"
  public_key = "%s"
}

resource "cherryservers_ip" "test_server_ip" {
  project_id = "${cherryservers_project.test_server_project.id}"
  region = "eu_nord_1"
}

resource "cherryservers_server" "test_server_server" {
  region = "eu_nord_1"
  plan = "cloud_vps_1"
  project_id = "${cherryservers_project.test_server_project.id}"
  name = "test"
  hostname = "server-fullconfig-test"
  image = "ubuntu_22_04"
  ssh_key_ids = ["${cherryservers_ssh_key.test_server_ssh_key.id}"]
  extra_ip_addresses_ids = ["${cherryservers_ip.test_server_ip.id}"]
  tags = {
    env = "test"
  }
  spot_instance = "false"
  timeouts = {
    create = "20m"
  }
}
`, projectName, teamID, sshKeyLabel, sshKeyPublicKey)
}
