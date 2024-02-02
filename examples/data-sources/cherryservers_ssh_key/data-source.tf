# Get SSH key info by ID
data "cherryservers_ssh_key" "by_id" {
  ssh_key_id = "123"
}

# Get SSH key info by name
data "cherryservers_ssh_key" "by_name" {
  project_id = "123"
  name       = "ssh-key-name"
}

