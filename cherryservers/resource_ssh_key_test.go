package cherryservers

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/cherryservers/cherrygo/v3"
	acctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccCherryServersSSHKey_Basic(t *testing.T) {
	var key cherrygo.SSHKey
	rInt := acctest.RandInt()
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("cherryservers@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCherryServersSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCherryServersSSHKeyConfigBasic(rInt, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCherryServersSSHKeyExists("cherryservers_ssh_key.ssh_key", &key),
					resource.TestCheckResourceAttr("cherryservers_ssh_key.ssh_key", "name", fmt.Sprintf("foobar-%d", rInt)),
					resource.TestCheckResourceAttr("cherryservers_ssh_key.ssh_key", "public_key", publicKeyMaterial),
					resource.TestCheckResourceAttrSet("cherryservers_ssh_key.ssh_key", "fingerprint"),
				),
			},
		},
	})
}
func testAccCheckCherryServersSSHKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Client).cherrygoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cherryservers_ssh_key" {
			continue
		}

		sshKeyID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		// Try to find the key
		_, resp, err := client.SSHKeys.Get(sshKeyID, nil)
		if err != nil {
			if is404Error(resp) {
				return nil
			}

			return fmt.Errorf("Error getting SSH key (%s)", rs.Primary.ID)
		}
	}

	return nil
}
func testAccCheckCherryServersSSHKeyExists(n string, key *cherrygo.SSHKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("SSH key ID is not set")
		}

		client := testAccProvider.Meta().(*Client).cherrygoClient()

		sshKeyID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		// Try to get the key
		_, _, err = client.SSHKeys.Get(sshKeyID, nil)
		if err != nil {
			return err
		}

		return nil
	}
}
func testAccCheckCherryServersSSHKeyConfigBasic(rInt int, key string) string {
	res := fmt.Sprintf(`
resource "cherryservers_ssh_key" "ssh_key" {
    name = "foobar-%v"
    public_key = "%s"
}`, rInt, key)
	return res
}
