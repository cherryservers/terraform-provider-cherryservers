# Get region by ID.
data "cherryservers_region" "lt_region" {
  id = 1
}

# Get region by slug.
data "cherryservers_region" "lt_region_by_slug" {
  slug = "LT-Siauliai"
}