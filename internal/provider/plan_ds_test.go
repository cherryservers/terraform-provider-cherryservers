package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testCheckPlanLT(dsName string, prefix string) resource.TestCheckFunc {
	return resource.ComposeAggregateTestCheckFunc(
		resource.TestCheckResourceAttrSet(dsName, prefix+"name"),
		resource.TestCheckResourceAttrSet(dsName, prefix+"slug"),
		resource.TestCheckResourceAttrSet(dsName, prefix+"id"),
		resource.TestCheckResourceAttr(dsName, prefix+"type", "vps"),

		// Images
		resource.TestCheckResourceAttrSet(dsName, prefix+"softwares.0.image.name"),
		resource.TestCheckResourceAttrSet(dsName, prefix+"softwares.0.image.slug"),

		// Specs
		resource.TestCheckResourceAttr(dsName, prefix+"specs.cpus.count", "1"),
		resource.TestCheckResourceAttrSet(dsName, prefix+"specs.cpus.name"),
		resource.TestCheckResourceAttr(dsName, prefix+"specs.cpus.cores", "1"),
		resource.TestCheckResourceAttrSet(dsName, prefix+"specs.cpus.frequency"),
		resource.TestCheckResourceAttrSet(dsName, prefix+"specs.cpus.unit"),

		resource.TestCheckResourceAttrSet(dsName, prefix+"specs.memory.count"),
		resource.TestCheckResourceAttrSet(dsName, prefix+"specs.memory.total"),
		resource.TestCheckResourceAttrSet(dsName, prefix+"specs.memory.unit"),
		resource.TestCheckResourceAttrSet(dsName, prefix+"specs.memory.name"),

		resource.TestCheckResourceAttrSet(dsName, prefix+"specs.storage.0.count"),
		resource.TestCheckResourceAttrSet(dsName, prefix+"specs.storage.0.name"),
		resource.TestCheckResourceAttrSet(dsName, prefix+"specs.storage.0.size"),
		resource.TestCheckResourceAttrSet(dsName, prefix+"specs.storage.0.unit"),

		resource.TestCheckResourceAttrSet(dsName, prefix+"specs.nics.name"),
		resource.TestCheckResourceAttrSet(dsName, prefix+"specs.bandwidth.name"),

		// Pricing
		resource.TestCheckResourceAttrSet(dsName, prefix+"pricing.0.unit"),
		resource.TestCheckResourceAttrSet(dsName, prefix+"pricing.0.price"),
		resource.TestCheckResourceAttrSet(dsName, prefix+"pricing.0.currency"),

		// Regions
		resource.TestCheckResourceAttr(dsName, prefix+"available_regions.0.id", "1"),
		resource.TestCheckResourceAttr(dsName, prefix+"available_regions.0.name", "Lithuania"),
		resource.TestCheckResourceAttr(dsName, prefix+"available_regions.0.region_iso_2", "LT"),
		resource.TestCheckResourceAttrSet(dsName, prefix+"available_regions.0.stock_qty"),
		resource.TestCheckResourceAttrSet(dsName, prefix+"available_regions.0.spot_qty"),
		resource.TestCheckResourceAttr(dsName, prefix+"available_regions.0.slug", "LT-Siauliai"),
		resource.TestCheckResourceAttr(dsName, prefix+"available_regions.0.location", "Lithuania, Šiauliai"),

		resource.TestMatchResourceAttr(dsName, prefix+"available_regions.0.bgp.hosts.0", ipv4Regex),
		resource.TestCheckResourceAttrSet(dsName, prefix+"available_regions.0.bgp.asn"),
	)
}

func TestAccPlanSingleDSByID(t *testing.T) {
	const dsName = "data.cherryservers_plan.by_id"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: planByIDConfig,
				Check:  testCheckPlanLT(dsName, ""),
			},
		},
	})
}

const planByIDConfig string = `

data "cherryservers_plan" "by_id" {
  id = 625
}
`

func TestAccPlanSingleDSBySlug(t *testing.T) {
	const dsName = "data.cherryservers_plan.by_slug"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: planBySlugConfig,
				Check:  testCheckPlanLT(dsName, ""),
			},
		},
	})
}

const planBySlugConfig string = `

data "cherryservers_plan" "by_slug" {
  slug = "B1-1-1gb-20s-shared"
}
`

func TestAccPlanList(t *testing.T) {
	const dsName = "data.cherryservers_plans.all_plans"
	const simple_vps_id = 625
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: planAllConfig,
				Check:  testAccCheckPlanListContains(simple_vps_id, dsName),
			},
		},
	})
}

const planAllConfig string = `

data "cherryservers_plans" "all_plans" {
}
`

func testAccCheckPlanListContains(id int, dsName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		i, err := findPlanIndex(id, testCherryGoClient)
		if err != nil {
			return err
		}

		return testCheckPlanLT(dsName, fmt.Sprintf("plans.%d.", i))(s)
	}
}
