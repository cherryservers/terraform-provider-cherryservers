# You must find your team_id for your account by login into the CherryServers portal: [https://portal.cherryservers.com/#/login](https://portal.cherryservers.com/#/login)
variable "team_id" {
  default = "12345"
}

# Create a new Project for your team 
resource "cherryservers_project" "myproject" {
  team_id = var.team_id
  name    = "My Cool New Project"
}
