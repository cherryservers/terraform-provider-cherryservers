# Get server info by ID.
data "cherryservers_server" "my_server" {
  id = "123456"
}

# Get server info yby hostname.
data "cherryservers_server" "by_hostname" {
  project_id = "123456"
  hostname   = "foo-bar"
}
