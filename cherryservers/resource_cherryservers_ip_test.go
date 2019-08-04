package cherryservers

import (
	"fmt"
  //"github.com/cherryservers/cherrygo"
//	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
  "regexp"
  "sort"
  "os"
  //"context"
	//"strconv"
	"testing"
)
var teamID = "35587"
var floatingIP string

func TestAccCherryServersFloatingIP_Region(t *testing.T) {

	expectedURNRegEx, _ := regexp.Compile(`(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCherryServersFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCherryServersFloatingIPConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCherryServersFloatingIPExists("cherryservers_ip.foobar", floatingIP),
					resource.TestCheckResourceAttr(
						"cherryservers_ip.foobar", "region", "EU-East-1"),
					resource.TestMatchResourceAttr("cherryservers_ip.foobar", "address", expectedURNRegEx),
				),
			},
		},
	})
}
func testAccCheckCherryServersFloatingIPExists(n string,floatingIPlocal string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		/*if floatingIPlocal == "" {
			return fmt.Errorf("No Record ID is set")
		}*/

		client := testAccProvider.Meta().(*CombinedConfig).Client()

    project_id := rs.Primary.Attributes["project_id"]
		// Try to find the FloatingIP
		foundFloatingIP, _, err := client.IPAddresses.List(project_id)

    //return fmt.Errorf("ips%#v",foundFloatingIP)
		if err != nil {
			//return err
      fmt.Fprintln(os.Stdout,err)
		}
    if len(foundFloatingIP) > 0 {
      return nil
    }
    i := sort.Search(len(foundFloatingIP), func(k int) bool { return foundFloatingIP[k].ID == floatingIPlocal })
		if i < len(foundFloatingIP) && foundFloatingIP[i].ID == floatingIPlocal {
			return nil 
		}
    return fmt.Errorf("Record not found")
	}
}
func testAccCheckCherryServersFloatingIPDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).Client()

  var projectID string
	for _, rs := range s.RootModule().Resources {
    fmt.Fprintln(os.Stdout,rs.Type)
		if rs.Type != "cherryservers_ip" {
      fmt.Println("continuing")
			continue
		}
		if rs.Type == "cherryservers_project" {
      projectID = rs.Primary.Attributes["project_id"]
    //return fmt.Errorf(rs.Primary.Attributes["project_id"])
		// Try to find the key
		results , _, err := client.IPAddresses.List(projectID)

//    return fmt.Errorf("length: %d", results)
		// Try to find the key
		if err != nil {
			//return err
      fmt.Fprintln(os.Stdout,err)
		}
    if len(results) == 0 {
      return nil
      }
		if err == nil {
			return fmt.Errorf("Floating IP still exists")
		}
  }
	}

	return nil
}

func testAccCheckCherryServersFloatingIPConfig_basic() string {
	res := fmt.Sprintf(`
resource "cherryservers_project" "myproject" {
  team_id = "%v"
  name = "foobar-project"
}
resource "cherryservers_ip" "foobar" {
  project_id = "${cherryservers_project.myproject.id}"
  region = "EU-East-1"
}`,teamID)
	return res
}
