package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCycleList(t *testing.T) {
	const dsName = "data.cherryservers_cycles.all"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: cyclesListAllConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(dsName, "cycles.0.id"),
					resource.TestCheckResourceAttrSet(dsName, "cycles.0.name"),
					resource.TestCheckResourceAttrSet(dsName, "cycles.0.slug"),
				),
			},
		},
	})
}

const cyclesListAllConfig string = `

data "cherryservers_cycles" "all" {
}
`
