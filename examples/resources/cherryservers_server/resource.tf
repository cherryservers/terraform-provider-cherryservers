#Create a new server:
resource "cherryservers_server" "server" {
  plan       = "B1-1-1gb-20s-shared"
  project_id = 123456
  region     = "LT-Siauliai"
}

#Create a new server with options:
resource "cherryservers_server" "server" {
  plan                   = "B1-1-1gb-20s-shared"
  hostname               = "sharing-wallaby"
  project_id             = 123456
  region                 = "LT-Siauliai"
  image                  = "ubuntu_22_04"
  ssh_key_ids            = ["1", "2"]
  extra_ip_addresses_ids = ["8269de5d-9b89-af9a-8bcc-8efb4d9fa282"]
  spot_instance          = true
  tags = {
    Name        = "Example Instance"
    Environment = "Production"
  }
}
