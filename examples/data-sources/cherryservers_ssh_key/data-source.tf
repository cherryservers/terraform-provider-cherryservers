# Get SSH key info by ID.
data "cherryservers_ssh_key" "my_ssh_key" {
  ssh_key_id = "1234"
}

# Get SSH key info by label.
data "cherryservers_ssh_key" "by_label" {
  project_id = "123"
  label      = "ssh-key-label"
}
