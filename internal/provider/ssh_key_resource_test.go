package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"regexp"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSSHKeyResource_basic(t *testing.T) {
	name := "terraform_test_ssh_" + acctest.RandString(5)
	publicKey, _, err := acctest.RandSSHKeyPair("cherryservers@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}
	publicKeyUpdate, _, err := acctest.RandSSHKeyPair("cherryservers@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckCherryServersSSHKeyDestroy,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccSSHKeyConfig(name, publicKey),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckCherryServersSSHKeyExists("cherryservers_ssh_key.test_ssh_key_ssh_key"),
					resource.TestCheckResourceAttrSet("cherryservers_ssh_key.test_ssh_key_ssh_key", "created"),
					resource.TestCheckResourceAttrSet("cherryservers_ssh_key.test_ssh_key_ssh_key", "fingerprint"),
					resource.TestCheckResourceAttrSet("cherryservers_ssh_key.test_ssh_key_ssh_key", "updated"),
					resource.TestMatchResourceAttr("cherryservers_ssh_key.test_ssh_key_ssh_key", "id", regexp.MustCompile(`[0-9]+`)),
				),
			},
			// ImportState testing
			{
				ResourceName:      "cherryservers_ssh_key.test_ssh_key_ssh_key",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccSSHKeyConfig(name+"_update", publicKeyUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("cherryservers_ssh_key.test_ssh_key_ssh_key", "name", name+"_update"),
					resource.TestCheckResourceAttr("cherryservers_ssh_key.test_ssh_key_ssh_key", "public_key", publicKeyUpdate),
				),
			},

			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccCheckCherryServersSSHKeyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("SSH key ID is not set")
		}
		client := testCherryGoClient
		serverID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("unable to convert SSH key ID")
		}

		// Try to get the ssh key id
		_, _, err = client.SSHKeys.Get(serverID, nil)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckCherryServersSSHKeyDestroy(s *terraform.State) error {
	client := testCherryGoClient

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cherryservers_ssh_key" {
			continue
		}

		sshID, converr := strconv.Atoi(rs.Primary.ID)
		if converr != nil {
			return fmt.Errorf("unable to convert SSH key ID")
		}

		sshKey, resp, err := client.SSHKeys.Get(sshID, nil)

		if err != nil {
			if is404Error(resp) {
				continue
			}

			return fmt.Errorf("ssh key listing error: %#v", err)
		}

		return fmt.Errorf("ssh key still exists: %#v", sshKey.ID)
	}
	return nil
}

func testAccSSHKeyConfig(name string, publicKey string) string {
	return fmt.Sprintf(`
resource "cherryservers_ssh_key" "test_ssh_key_ssh_key" {
  name = "%s"
  public_key = "%s"
}
`, name, publicKey)
}
