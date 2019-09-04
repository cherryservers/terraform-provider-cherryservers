---
layout: "cherryservers"
page_title: "CherryServers: cherryservers_ip"
sidebar_current: "docs-cherryservers-resource-ip"
description: |-
  Provides a CherryServers IP resource. This can be used to create, modify, and delete Static IP addresses. These IPs can be added to your server and have the network interfaces configured at creation time.
---

# cherryservers\_ip

Provides a CherryServers IP resource. This can be used to create,
modify, and delete IP. 

## Example Usage

```hcl
# Optionally configure the Region to launch in as a variable or specify inline below
variable "region" {
  default = "EU-East-1"
}

# To see how to configure a project, see the Cherryservers_project documentation. 
# You will need to have an existing project in order to reserve an IP resource. You do not need to specify the Routing options.
resource "cherryservers_ip" "my_ip_address" {
  project_id = "${cherryservers_project.myproject.id}"
  region = "${var.region}"
  routed_to_ip = "127.0.0.1" # Optional
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The Project ID to deploy the IP in.
* `region` - (Required) The region to deply in.
* `routed_to_hostname` - (Optional) A hostname to route network traffic to
* `routed_to_ip` - (Optional) An IP address to route network traffic to

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the IP Address
* `address` - The address of the IP eg. 127.0.0.1
* `cidr`- The CIDR block of the IP
* `gateway` - The Gateway address for the IP
* `ptr` - The PTR record that resolves to the IP address
* `type` - The type of IP 
* `routed_to` - The hostname or IP that the IP routes to
* `region` - The region of the IP 

## Import

IPs can be imported using the IP `id`, e.g.

```
terraform import cherryservers_ip.my_ip_address 123
```
