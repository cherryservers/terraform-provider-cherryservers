# Get SSH key info by ID.
data "cherryservers_ssh_key" "my_ssh_key" {
  id = "1234"
}

# Get SSH key info by name.
data "cherryservers_ssh_key" "by_name" {
  name = "ssh-key-label"
}
