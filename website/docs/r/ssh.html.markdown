---
layout: "cherryservers"
page_title: "CherryServers: cherryservers_ssh"
sidebar_current: "docs-cherryservers-resource-ssh"
description: |-
  Provides a CherryServers SSH key resource. This can be used to create, and delete SSH keys associated with your account.  Please note that you will not be able to add duplicate SSH keys to your account. 
---

# cherryservers\_ssh

Provides a CherryServers SSH Key resource. This can be used to create,
 and delete SSH Keys associated with your account. 

## Example Usage

```hcl
# Create a new SSH Key for your account

# (Optionally) specify the path to your key as a terraform variable 

variable "private_key" {
  default = "~/.ssh/cherry"
}

resource "cherryservers_ssh" "mysshkey" {
  name   = "mykey"
  public_key = "${file("${var.private_key}.pub")}" # The public key contents can also be stored specific here directly
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The Server name.
* `public_key` - (Required) The public key contents as a string to be added to your account

The following attributes are exported:

* `id` - The ID of the SSH Key. Use this attribute when associating an SSH Key to a cherryservers\_server resource
* `name`- The name of the SSH Key
* `fingerprint` - The fingerprint of your SSH Public key
* `created` - The date when this Key was added
* `updated` - The date when this Key was modified

## Import

Servers can be imported using the SSH Key `id`, e.g.

```
terraform import cherryservers_ssh.mysshkey 900
```
