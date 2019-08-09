---
layout: "cherryservers"
page_title: "CherryServers: cherryservers_server"
sidebar_current: "docs-cherryservers-resource-server"
description: |-
  Provides a CherryServers Server resource. This can be used to create, modify, and delete Servers. Servers also support provisioning.
---

# cherryservers\_server

Provides a CherryServers Server resource. This can be used to create,
modify, and delete Servers. Servers also support
[provisioning](/docs/provisioners/index.html).

## Example Usage

```hcl
variable "region" {
  default = "EU-East-1"
}
variable "image" {
  default = "Ubuntu 18.04 64bit"
}
variable "plan_id" {
  default = "86"
}
# Create a new Web Server in the EU-East-1 region
# Optional provisioning is specified as well.
resource "cherryservers_server" "my-server" {
  project_id = "${cherryservers_project.myproject.id}"
  region = "${var.region}"
  hostname = "production_server"
  image = "${var.image}"
  plan_id = "${var.plan_id}"
  ssh_keys_ids = [
    "${cherryservers_ssh.mykey.id}"]
  ip_addresses_ids = [
    "${cherryservers_ip.my-ip.id}"]

  # Upload your setup script
  provisioner "file" {
    source = "my_setup_script.sh"
    destination = "/tmp/my_setup_script.sh"

    connection {
      type = "ssh"
      user = "root"
      host = "${self.primary_ip}"
      private_key = "${file(var.private_key)}"
      timeout = "20m"
    }
  }
  # Make setup script executable and run it
  provisioner "remote-exec" {
    inline = [
      "chmod +x /tmp/my_setup_script.sh",
      "/tmp/my_setup_script.sh",
    ]
    connection {
      type = "ssh"
      user = "root"
      host = "${self.primary_ip}"
      private_key = "${file(var.private_key)}"
      timeout = "20m"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The Project ID to deploy the IP in.
* `region` - (Required) The region to deploy in. You can find these on [CherryServers API](https://api.cherryservers.com/doc/#tag/Regions).
* `image` - (Required) The Server image slug. You can find these on [CherryServers API](https://api.cherryservers.com/doc/#tag/Images/paths/~1v1~1plans~1{planId}~1images/get), after finding your plan_id.
* `hostname` - (Required) The Server name.
* `plan_id` - (Required) The unique slug that indentifies the type of Server. You can find a list of available slugs on [CherryServers API](https://api.cherryservers.com/doc/#tag/Plans/paths/~1v1~1teams~1{teamId}~1plans/get).
* `ssh_key_ids` - (Optional) A list of SSH IDs to enable in the format `[12345, 123456]`. 
* `ip_addresses_ids` - (Optional) A list of IP Address IDs to enable in the format `[12345, 123456]`. 

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Server
* `hostname` - The hostname of the Server
* `name`- The name of the Server
* `region` - The region of the Server
* `image` - The image of the Server
* `primary_ip` - The primary IPv4 address of the server. Servers will always have a primary ID in addition to any attached reserverd IPs
* `private_ip` - The private IPv4 address of the server
* `state` - The state of the server, such as "Pending"
* `power_state` - The power state of the server, such as "Powered off"
* `price` - The Server hourly price

## Import

Servers can be imported using the Server `id`, e.g.

```
terraform import cherryservers_server.myserver 100823
```
