// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServerDataSource_basic(t *testing.T) {
	projectId := os.Getenv("CHERRY_TEST_PROJECT_ID")
	resourceName := "cherryservers_server.test_server_data_source"
	dataSourceName := "data.cherryservers_server.test_server_data_source"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: testAccServerDataSourceConfig(projectId),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "plan", resourceName, "plan"),
					resource.TestCheckResourceAttrPair(dataSourceName, "project_id", resourceName, "project_id"),
					resource.TestCheckResourceAttrPair(dataSourceName, "region", resourceName, "region"),
					resource.TestCheckResourceAttrPair(dataSourceName, "hostname", resourceName, "hostname"),
					resource.TestCheckResourceAttrPair(dataSourceName, "name", resourceName, "name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "username", resourceName, "username"),
					resource.TestCheckResourceAttrPair(dataSourceName, "password", resourceName, "password"),
					resource.TestCheckResourceAttrPair(dataSourceName, "bmc", resourceName, "bmc"),
					resource.TestCheckResourceAttrPair(dataSourceName, "image", resourceName, "image"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ssh_key_ids", resourceName, "ssh_key_ids"),
					resource.TestCheckResourceAttrPair(dataSourceName, "extra_ip_addresses_ids", resourceName, "extra_ip_addresses_ids"),
					resource.TestCheckResourceAttrPair(dataSourceName, "user_data_file", resourceName, "user_data_file"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tags", resourceName, "tags"),
					resource.TestCheckResourceAttrPair(dataSourceName, "spot_instance", resourceName, "spot_instance"),
					resource.TestCheckResourceAttrPair(dataSourceName, "os_partition_size", resourceName, "os_partition_size"),
					resource.TestCheckResourceAttrPair(dataSourceName, "power_state", resourceName, "power_state"),
					resource.TestCheckResourceAttrPair(dataSourceName, "state", resourceName, "state"),
					resource.TestCheckResourceAttrPair(dataSourceName, "ip_addresses", resourceName, "ip_addresses"),
				),
			},
		},
	})
}

func testAccServerDataSourceConfig(projectID string) string {
	return fmt.Sprintf(`
resource "cherryservers_server" "test_server_data_source" {
  plan = "cloud_vps_1"
  region = "eu_nord_1"
  project_id = "%s"
}

data "cherryservers_server" "test_server_data_source" {
  id = cherryservers_server.test_server_data_source.id
}
`, projectID)
}
