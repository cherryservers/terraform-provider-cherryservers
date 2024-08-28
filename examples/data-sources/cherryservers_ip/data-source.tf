# Get IP address info by ID.
data "cherryservers_ip" "my_ip" {
  id = "123abc456def"
}

# Get IP address info by address and project ID.
data "cherryservers_ip" "by_address" {
  ip_address = "0.0.0.0"
  project_id = "123456"
}