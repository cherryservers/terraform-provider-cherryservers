# Create a new floating IP address
resource "cherryservers_ip" "floating-1" {
  project_id = 123456
  region     = "LT-Siauliai"
}

# Create a new floating IP address with optional parameters
resource "cherryservers_ip" "floating-1" {
  project_id      = 123
  region          = "LT-Siauliai"
  target_hostname = "gentle-turtle"
  ddos_scrubbing  = true
  tags = {
    Name        = "Example Instance"
    Environment = "Production"
  }
}
