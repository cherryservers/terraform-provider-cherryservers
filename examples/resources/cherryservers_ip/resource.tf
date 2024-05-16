# Create a new floating IP address
resource "cherryservers_ip" "floating-1" {
  project_id = 123456
  region     = "eu_nord_1"
}

# Create a new floating IP address with optional parameters
resource "cherryservers_ip" "floating-1" {
  project_id      = 123
  region          = "eu_nord_1"
  target_hostname = "gentle-turtle"
  ddos_scrubbing  = true
  tags = {
    Name        = "Example Instance"
    Environment = "Production"
  }
}
