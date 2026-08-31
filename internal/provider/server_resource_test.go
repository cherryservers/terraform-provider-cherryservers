package provider

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"testing"

	"github.com/cherryservers/cherrygo/v4"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServerResource_basic(t *testing.T) {
	ctx := t.Context()
	serverResourceName := "terraform_test_server_" + acctest.RandString(5)
	projectName := testProjectNamePrefix + acctest.RandString(5)
	testPlan := "B1-1-1gb-20s-shared"
	testRegion := "LT-Siauliai"
	teamID := os.Getenv("CHERRY_TEST_TEAM_ID")
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             checkWithContext(ctx, testAccCheckCherryServersServerDestroy),
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccServerResourceConfigOnlyReq(projectName, testPlan, testRegion, serverResourceName, teamID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCherryServersServerExists(ctx, "cherryservers_server."+serverResourceName),
					resource.TestMatchResourceAttr("cherryservers_server."+serverResourceName, "hostname", regexp.MustCompile("[a-z]+-[a-z]+")),
					resource.TestMatchResourceAttr("cherryservers_server."+serverResourceName, "id", regexp.MustCompile("[0-9]+")),
					resource.TestMatchTypeSetElemNestedAttrs("cherryservers_server."+serverResourceName, "ip_addresses.*", map[string]*regexp.Regexp{
						"id":             regexp.MustCompile("^.+$"),
						"type":           regexp.MustCompile("^primary-ip$"),
						"cidr":           regexp.MustCompile(`^.*\/\d\d$`),
						"address_family": regexp.MustCompile("^4$"),
						"address":        regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`),
					}),
					resource.TestMatchTypeSetElemNestedAttrs("cherryservers_server."+serverResourceName, "ip_addresses.*", map[string]*regexp.Regexp{
						"id":             regexp.MustCompile("^.+$"),
						"type":           regexp.MustCompile("^private-ip$"),
						"cidr":           regexp.MustCompile(`^.*\/\d\d$`),
						"address_family": regexp.MustCompile("^4$"),
						"address":        regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`),
						"vlan_id":        regexp.MustCompile(`^\d+$`),
					}),
					resource.TestCheckResourceAttrSet("cherryservers_server."+serverResourceName, "name"),
					resource.TestCheckResourceAttr("cherryservers_server."+serverResourceName, "power_state", "on"),
					resource.TestMatchResourceAttr("cherryservers_server."+serverResourceName, "project_id", regexp.MustCompile(`[0-9]+`)),
					resource.TestCheckResourceAttr("cherryservers_server."+serverResourceName, "spot_instance", "false"),
					resource.TestCheckResourceAttr("cherryservers_server."+serverResourceName, "state", "active"),
					resource.TestMatchResourceAttr("cherryservers_server."+serverResourceName, "pricing.price", regexp.MustCompile(`[+-]?([0-9]*[.])?[0-9]+`)),
					resource.TestCheckResourceAttr("cherryservers_server."+serverResourceName, "pricing.currency", "EUR"),
				),
			},
			// ImportState testing
			{
				ResourceName:            "cherryservers_server." + serverResourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"allow_reinstall"},
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
	ctx := t.Context()
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
		CheckDestroy:             checkWithContext(ctx, testAccCheckCherryServersServerDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccServerResourceFullConfig(projectName, teamID, label, publicKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCherryServersServerExists(ctx, "cherryservers_server.test_server_server"),
					resource.TestCheckResourceAttr("cherryservers_server.test_server_server", "image", "ubuntu_24_04_64bit"),
					resource.TestMatchResourceAttr("cherryservers_server.test_server_server", "id", regexp.MustCompile("[0-9]+")),
					resource.TestMatchResourceAttr("cherryservers_server.test_server_server", "ip_addresses.0.address", regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`)),
					resource.TestMatchResourceAttr("cherryservers_server.test_server_server", "ip_addresses.1.address", regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`)),
					resource.TestCheckResourceAttrSet("cherryservers_server.test_server_server", "name"),
					resource.TestCheckResourceAttr("cherryservers_server.test_server_server", "power_state", "on"),
					resource.TestMatchResourceAttr("cherryservers_server.test_server_server", "project_id", regexp.MustCompile(`[0-9]+`)),
					resource.TestCheckResourceAttr("cherryservers_server.test_server_server", "spot_instance", "false"),
					resource.TestCheckResourceAttr("cherryservers_server.test_server_server", "state", "active"),
				),
			},
			// Reinstall testing
			{
				Config: testAccServerResourceFullUpdateWithReinstall(projectName, teamID, label, publicKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCherryServersServerExists(ctx, "cherryservers_server.test_server_server"),
					resource.TestMatchResourceAttr("cherryservers_server.test_server_server", "id", regexp.MustCompile("[0-9]+")),
					resource.TestMatchResourceAttr("cherryservers_server.test_server_server", "ip_addresses.0.address", regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`)),
					resource.TestMatchResourceAttr("cherryservers_server.test_server_server", "ip_addresses.1.address", regexp.MustCompile(`^(((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})`)),
					resource.TestCheckResourceAttr("cherryservers_server.test_server_server", "power_state", "on"),
					resource.TestCheckResourceAttr("cherryservers_server.test_server_server", "spot_instance", "false"),
					resource.TestCheckResourceAttr("cherryservers_server.test_server_server", "state", "active"),
				),
			},
		},
	})
}

func TestAccServerIPXE(t *testing.T) {
	// We need an API key for setup, so skip early if it's not an acceptance test.
	if acc := os.Getenv(resource.EnvTfAcc); acc == "" {
		t.Skip("skipping ipxe acceptance test, since TF_ACC is not set")
	}

	const resourceName = "cherryservers_server.ipxe_test"
	project := testProjectNamePrefix + "ipxe"
	plan, region := ipxePlanRegion(t, testCherryGoClient, testTeam)
	ipxeCreate := ipxeScript(t, filepath.Join("testdata", "ubuntu.ipxe"))
	ipxeReinstall := ipxeScript(t, filepath.Join("testdata", "alma.ipxe"))
	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy: func(s *terraform.State) error {
			return testAccCheckCherryServersServerDestroy(t.Context(), s)
		},
		Steps: []resource.TestStep{
			{
				// Fail when iPXE image is configured with no script.
				Config:      ipxeOnlyImageConfig(project, region, plan, ipxeImage, testTeam, false),
				ExpectError: regexp.MustCompile("Missing Attribute Configuration"),
			},
			{
				// Fail when iPXE script is configured with user_data.
				Config:      ipxeInvalidWithUserData,
				ExpectError: regexp.MustCompile("Attribute .* cannot be specified"),
			},
			{
				// Fail when iPXE script is configured with os_partition_size.
				Config:      ipxeInvalidWithOSPartitionSize,
				ExpectError: regexp.MustCompile("Attribute .* cannot be specified"),
			},
			{
				// Fail when iPXE script is configured with ssh_key_ids.
				Config:      ipxeInvalidWithSSHKeyIDs,
				ExpectError: regexp.MustCompile("Attribute .* cannot be specified"),
			},
			{
				// Fail when persist_ipxe is configured without iPXE script.
				Config:      persistIPXEWithoutIPXE,
				ExpectError: regexp.MustCompile("Invalid Attribute Combination"),
			},
			{
				// Succeed when valid iPXE script is provided and image is left
				// for the provider to set.
				Config: ipxeConfig(region, plan, ipxeCreate, testTeam, false, false),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							resourceName,
							tfjsonpath.New("image"),
							knownvalue.StringExact(ipxeImage),
						),
					},
				},
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "image", ipxeImage),
					resource.TestMatchTypeSetElemNestedAttrs(resourceName, "ip_addresses.*", map[string]*regexp.Regexp{
						"address_family": regexp.MustCompile("^6$"),
					}),
				),
			},
			{
				// Fail when updating iPXE script, but allow_reinstall is not enabled.
				Config:      ipxeConfig(region, plan, ipxeReinstall, testTeam, false, false),
				ExpectError: regexp.MustCompile("allow_reinstall attribute not set"),
			},
			{
				// Fail when updating persist_ipxe, but allow_reinstall is not enabled.
				Config:      ipxeConfig(region, plan, ipxeReinstall, testTeam, false, true),
				ExpectError: regexp.MustCompile("allow_reinstall attribute not set"),
			},
			{
				// Succeed when reinstalling with a new iPXE script and persist_ipxe.
				Config: ipxeConfig(region, plan, ipxeReinstall, testTeam, true, true),
			},
			{
				// Fail when a standard image is configured along with an iPXE script.
				Config:      ipxeWithImageConfig(project, region, plan, defaultTestImage, ipxeReinstall, testTeam, true),
				ExpectError: regexp.MustCompile("Invalid Attribute Configuration"),
			},
			{
				// Succeed when reinstalling an iPXE server into one with a standard image.
				Config: ipxeOnlyImageConfig(project, region, plan, defaultTestImage, testTeam, true),
			},
			{
				// Succeed when reinstalling a server with a standard image into one with iPXE
				// and the image is configured by the user.
				Config: ipxeWithImageConfig(project, region, plan, ipxeImage, ipxeCreate, testTeam, true),
			},
			{
				// Imported server gets a new OS image planned on reinstall, when
				// current image is iPXE, but no script is configured.
				Config: ipxeOSPartitionConfig(project, region, plan, testTeam, true),

				ResourceName:    resourceName,
				ImportState:     true,
				ImportStateKind: resource.ImportBlockWithID,

				ImportPlanChecks: resource.ImportPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							resourceName,
							tfjsonpath.New("image"),
							knownvalue.StringRegexp(regexp.MustCompile("ubuntu")),
						),
					},
				},

				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func ipxePlanRegion(t *testing.T, client *cherrygo.Client, team int) (plan, region string) {
	t.Helper()

	// Provisioning bare-metal servers is expensive and slow, so we'll fold IPv6 tests
	// into iPXE tests.
	ip6Regions := []string{"SG-Singapore", "DE-Frankfurt", "SE-Stockholm", "US-Chicago"}

	plans, _, err := client.Plans.List(t.Context(), team, nil)
	if err != nil {
		t.Fatalf("failed to list plans: %s)", err.Error())
	}

	stock := 0
	for _, p := range plans {
		for _, r := range p.AvailableRegions {
			if r.StockQty > stock && slices.ContainsFunc(
				p.Softwares, func(s cherrygo.SoftwareImage) bool {
					return s.Image.Slug == ipxeImage
				},
			) && slices.Contains(ip6Regions, r.Slug) {
				stock = r.StockQty
				plan = p.Slug
				region = r.Slug
			}
		}
	}

	if plan == "" || region == "" {
		t.Fatalf("failed to find ipxe plan in ipv6 supported region with any stock %d", stock)
	}
	return
}

func ipxeScript(t *testing.T, path string) string {
	t.Helper()

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read ipxe script file: %s", err.Error())
	}

	return base64.StdEncoding.EncodeToString(b)
}

func testAccCheckCherryServersServerExists(ctx context.Context, resourceName string) resource.TestCheckFunc {
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
		_, _, err = client.Servers.Get(ctx, serverID, nil)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckCherryServersServerDestroy(ctx context.Context, s *terraform.State) error {
	client := testCherryGoClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cherryservers_server" {
			continue
		}

		serverID, converr := strconv.Atoi(rs.Primary.ID)
		if converr != nil {
			return fmt.Errorf("unable to convert Server ID")
		}

		server, resp, err := client.Servers.Get(ctx, serverID, nil)
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
  name = "%s"
  public_key = "%s"
}

resource "cherryservers_ip" "test_server_ip" {
  project_id = "${cherryservers_project.test_server_project.id}"
  region = "LT-Siauliai"
}

resource "cherryservers_server" "test_server_server" {
  region = "LT-Siauliai"
  plan = "B1-1-1gb-20s-shared"
  project_id = "${cherryservers_project.test_server_project.id}"
  name = "test"
  hostname = "server-fullconfig-test"
  image = "ubuntu_24_04_64bit"
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

func testAccServerResourceFullUpdateWithReinstall(projectName string, teamID string, sshKeyLabel string, sshKeyPublicKey string) string {
	return fmt.Sprintf(`
resource "cherryservers_project" "test_server_project" {
  name = "%s"
  team_id = "%s"
}

resource "cherryservers_ssh_key" "test_server_ssh_key" {
  name = "%s"
  public_key = "%s"
}

resource "cherryservers_ip" "test_server_ip" {
  project_id = "${cherryservers_project.test_server_project.id}"
  region = "LT-Siauliai"
}

resource "cherryservers_server" "test_server_server" {
  region = "LT-Siauliai"
  plan = "B1-1-1gb-20s-shared"
  project_id = "${cherryservers_project.test_server_project.id}"
  name = "test-reinstall"
  hostname = "server-reinstall-test"
  image = "debian_12_64bit"
  ssh_key_ids = []
  extra_ip_addresses_ids = []
  tags = {
    env = "reinstall"
  }
  spot_instance = "false"
  timeouts = {
    create = "20m"
	update = "10m"
  }
  allow_reinstall = true
}
`, projectName, teamID, sshKeyLabel, sshKeyPublicKey)
}

func ipxeConfig(region, plan, ipxe string, team int, allowReinstall, persistIPXE bool) string {
	return fmt.Sprintf(`
resource "cherryservers_project" "test_server_project" {
  name = "%s"
  team_id = "%d"
}

resource "cherryservers_server" "ipxe_test" {
  region = "%s" 
  plan = "%s"
  project_id = "${cherryservers_project.test_server_project.id}"
  ipxe = "%s"
  persist_ipxe = %t
  allow_reinstall = %t
  configure_ipv6 = true
}
`, testProjectNamePrefix+"ipxe", team, region, plan, ipxe, persistIPXE, allowReinstall)
}

func ipxeOnlyImageConfig(projectName, region, plan, image string, team int, allowReinstall bool) string {
	return fmt.Sprintf(`
resource "cherryservers_project" "test_server_project" {
  name = "%s"
  team_id = "%d"
}

resource "cherryservers_ssh_key" "test_server_ssh_key" {
  name = "test-key"
  public_key = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMSagPMjsdBnJ1CsF+ChdfnqZK3wc1n8m6MSRy7CK2Dz key-1787825911386"
}

resource "cherryservers_server" "ipxe_test" {
  region = "%s" 
  plan = "%s"
  project_id = "${cherryservers_project.test_server_project.id}"
  image = "%s"
  ssh_key_ids = ["${cherryservers_ssh_key.test_server_ssh_key.id}"]
  allow_reinstall = %t
  configure_ipv6 = true
}
`, projectName, team, region, plan, image, allowReinstall)
}

func ipxeWithImageConfig(projectName, region, plan, image, ipxe string, team int, allowReinstall bool) string {
	return fmt.Sprintf(`
resource "cherryservers_project" "test_server_project" {
  name = "%s"
  team_id = "%d"
}

resource "cherryservers_server" "ipxe_test" {
  region = "%s"
  plan = "%s"
  project_id = "${cherryservers_project.test_server_project.id}"
  ipxe = "%s"
  image = "%s"
  allow_reinstall = %t
  configure_ipv6 = true
}
`, projectName, team, region, plan, ipxe, image, allowReinstall)
}

func ipxeOSPartitionConfig(projectName, region, plan string, team int, allowReinstall bool) string {
	return fmt.Sprintf(`
resource "cherryservers_project" "test_server_project" {
  name = "%s"
  team_id = "%d"
}

resource "cherryservers_server" "ipxe_test" {
  region = "%s" 
  plan = "%s"
  os_partition_size = 100
  project_id = "${cherryservers_project.test_server_project.id}"
  allow_reinstall = %t
}
`, projectName, team, region, plan, allowReinstall)
}

const (
	ipxeInvalidWithUserData = `
resource "cherryservers_server" "ipxe_test" {
  region = "test"
  plan = "test"
  user_data = "test"
  project_id = 1
  ipxe = "test"
}
`
	ipxeInvalidWithOSPartitionSize = `
resource "cherryservers_server" "ipxe_test" {
  region = "test"
  plan = "test"
  os_partition_size = 1
  project_id = 1
  ipxe = "test"
}
`
	ipxeInvalidWithSSHKeyIDs = `
resource "cherryservers_server" "ipxe_test" {
  region = "test"
  plan = "test"
  ssh_key_ids = [ "1" ]
  project_id = 1
  ipxe = "test"
}
`
	persistIPXEWithoutIPXE = `
resource "cherryservers_server" "ipxe_test" {
  region = "test"
  plan = "test"
  project_id = 1
  persist_ipxe = true
}
`
)
