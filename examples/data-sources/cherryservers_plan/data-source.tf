# Get server plan by ID.
data "cherryservers_plan" "vps" {
  id = 625
}

# Get server plan by slug.
data "cherryservers_plan" "vps_by_slug" {
  slug = "B1-1-1gb-20s-shared"
}