# Get IP address info by ID
data "cherryservers_ip" "by_id" {
  ip_id = "123"
}

# Get IP address info by IP address
data "cherryservers_ip" "by_address" {
  ip_address = "8269de5d-9b89-af9a-8bcc-8efb4d9fa282"
  project_id = "123"
}
