package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRegionSingleDSByID(t *testing.T) {
	const dsName = "data.cherryservers_region.lt_region"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: regionByIdConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsName, "name", "Lithuania"),
					resource.TestCheckResourceAttr(dsName, "slug", "LT-Siauliai"),
					resource.TestCheckResourceAttr(dsName, "region_iso_2", "LT"),
					resource.TestCheckResourceAttrSet(dsName, "bgp.asn"),
					resource.TestMatchResourceAttr(dsName, "bgp.hosts.0", ipv4Regex),
				),
			},
		},
	})
}

const regionByIdConfig string = `

data "cherryservers_region" "lt_region" {
  id = 1
}
`

func TestAccRegionSingleDSBySlug(t *testing.T) {
	const dsName = "data.cherryservers_region.nl_region"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: regionBySlugConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsName, "name", "Amsterdam, NL"),
					resource.TestCheckResourceAttrSet(dsName, "id"),
					resource.TestCheckResourceAttr(dsName, "region_iso_2", "NL"),
					resource.TestCheckResourceAttrSet(dsName, "bgp.asn"),
					resource.TestMatchResourceAttr(dsName, "bgp.hosts.0", ipv4Regex),
				),
			},
		},
	})
}

const regionBySlugConfig string = `

data "cherryservers_region" "nl_region" {
  slug = "NL-Amsterdam"
}
`

func TestAccRegionsList(t *testing.T) {
	const dsName = "data.cherryservers_regions.all_regions"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: regionsListConfig,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(dsName, "regions.0.name", "Lithuania"),
					resource.TestCheckResourceAttr(dsName, "regions.0.slug", "LT-Siauliai"),
					resource.TestCheckResourceAttr(dsName, "regions.0.region_iso_2", "LT"),
					resource.TestCheckResourceAttrSet(dsName, "regions.0.bgp.asn"),
					resource.TestCheckResourceAttrSet(dsName, "regions.0.id"),
					resource.TestMatchResourceAttr(dsName, "regions.0.bgp.hosts.0", ipv4Regex),
				),
			},
		},
	})
}

const regionsListConfig string = `

data "cherryservers_regions" "all_regions" {
}
`
