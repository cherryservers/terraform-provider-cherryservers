# Get SSH key info by ID
data "cherryservers_ssh_key" "my_ssh_key" {
  ssh_key_id = "1234"
}
