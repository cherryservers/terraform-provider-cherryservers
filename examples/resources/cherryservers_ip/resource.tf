# Create a new floating IP addess
resource "cherryservers_ip" "floating-1" {
  project_id = 123
  region     = "eu_nord_1"
}

# Create a new floating IP addess with options
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
