# Create a new Project for your team
# You must find your team_id for your account by login into the CherryServers portal: [https://portal.cherryservers.com/#/login](https://portal.cherryservers.com/#/login)
resource "cherryservers_project" "my_project" {
  team_id = "123456"
  name    = "My Cool New Project"
}

# Create a new Project with BGP enabled
resource "cherryservers_project" "project_with_bgp" {
  team_id = "123456"
  name    = "Cool project with BGP"
  bgp = {
    enabled = "true"
  }
}