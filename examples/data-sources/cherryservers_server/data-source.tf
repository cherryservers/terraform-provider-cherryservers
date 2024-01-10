# Get server info by hostname
data "cherryservers_server" "by_hostname" {
  project_id = "123"
  hostname   = "foo-bar"
}

# Get server info by ID
data "cherryservers_server" "by_id" {
  project_id = "123"
  server_id  = "321"
}

