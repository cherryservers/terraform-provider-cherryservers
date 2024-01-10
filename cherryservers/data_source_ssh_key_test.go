package cherryservers

import (
	"fmt"
	"testing"

	acctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccCherryServersDataSourceSSHKey_Basic(t *testing.T) {
	rInt := acctest.RandInt()
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("cherryservers@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccCheckCherryServersSSHKeyDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCherryServersDataSourceSSHKeyConfigBasic(teamID, rInt, publicKeyMaterial),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.cherryservers_ssh_key.by_id", "id"),
					resource.TestCheckResourceAttrSet("data.cherryservers_ssh_key.by_id", "name"),
					resource.TestCheckResourceAttrSet("data.cherryservers_ssh_key.by_id", "public_key"),

					resource.TestCheckResourceAttrSet("data.cherryservers_ssh_key.by_name", "id"),
					resource.TestCheckResourceAttrSet("data.cherryservers_ssh_key.by_name", "name"),
					resource.TestCheckResourceAttrSet("data.cherryservers_ssh_key.by_name", "public_key"),
				),
			},
		},
	})
}

func testAccCheckCherryServersDataSourceSSHKeyConfigBasic(teamID string, rInt int, key string) string {
	res := fmt.Sprintf(`
resource "cherryservers_project" "project" {
	team_id = %s
  	name = "terraform-test-project-ssh"
}

resource "cherryservers_ssh_key" "foobar" {
    name = "foobar-%v"
    public_key = "%s"
}

data "cherryservers_ssh_key" "by_id" {
    ssh_key_id = "${cherryservers_ssh_key.foobar.id}"
}

data "cherryservers_ssh_key" "by_name" {
    name = "${cherryservers_ssh_key.foobar.name}"
	project_id = "${cherryservers_project.project.id}"
}`, teamID, rInt, key)
	return res
}
