package cherryservers

import (
	"fmt"
	"github.com/cherryservers/cherrygo"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	//"os"
	"sort"
	"strconv"
	"testing"
)

func TestAccCherryServersSSHKey_Basic(t *testing.T) {
	var key cherrygo.SSHKey
	rInt := acctest.RandInt()
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("cherryservers@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCherryServersSSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCherryServersSSHKeyConfig_basic(rInt, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCherryServersSSHKeyExists("cherryservers_ssh.foobar", &key),
					resource.TestCheckResourceAttr(
						"cherryservers_ssh.foobar", "name", fmt.Sprintf("foobar-%d", rInt)),
					resource.TestCheckResourceAttr(
						"cherryservers_ssh.foobar", "public_key", publicKeyMaterial),
				),
			},
		},
	})
}
func testAccCheckCherryServersSSHKeyDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).Client()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cherryservers_ssh" {
			continue
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		// Try to find the key
		list, _, err := client.SSHKeys.List()
		i := sort.Search(len(list), func(k int) bool { return list[k].ID == id })
		if i < len(list) && list[i].ID == id {
			return fmt.Errorf("SSH key still exists")
		}

		if err == nil {
			return nil
		}
	}

	return nil
}
func testAccCheckCherryServersSSHKeyExists(n string, key *cherrygo.SSHKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*CombinedConfig).Client()

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		// Try to find the key
		list, _, err := client.SSHKeys.List()
		i := sort.Search(len(list), func(k int) bool { return list[k].ID == id })
		if i < len(list) && list[i].ID == id {
			return nil
		}

		if err != nil {
			return err
		}

		/*if strconv.Itoa(foundKey.ID) != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}
		*/
		return fmt.Errorf("Record not found")

		//*key = *foundKey

		return nil
	}
}
func testAccCheckCherryServersSSHKeyConfig_basic(rInt int, key string) string {
	res := fmt.Sprintf(`
resource "cherryservers_ssh" "foobar" {
    name = "foobar-%v"
    public_key = "%s"
}`, rInt, key)
	return res
}
