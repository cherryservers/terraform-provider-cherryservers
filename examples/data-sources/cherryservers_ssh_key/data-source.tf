# Get SSH key info by ID.
data "cherryservers_ssh_key" "my_ssh_key" {
  id = "1234"
}

# Get SSH key info by label.
data "cherryservers_ssh_key" "by_label" {
  project_id = "123"
  label      = "ssh-key-label"
}
