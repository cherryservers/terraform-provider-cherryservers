#Create a new server:
resource "cherryservers_server" "server" {
  plan       = "cloud_vps_1"
  project_id = 123456
  region     = "eu_nord_1"
}

#Create a new server with options:
resource "cherryservers_server" "server" {
  plan                   = "cloud_vps_1"
  hostname               = "sharing-wallaby"
  project_id             = 123456
  region                 = "eu_nord_1"
  image_slug             = "ubuntu_22_04"
  ssh_key_ids            = ["1", "2"]
  extra_ip_addresses_ids = ["8269de5d-9b89-af9a-8bcc-8efb4d9fa282"]
  spot_instance          = true
  tags = {
    Name        = "Example Instance"
    Environment = "Production"
  }
}
