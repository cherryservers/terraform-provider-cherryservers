---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cherryservers_ssh_key Resource - cherryservers"
subcategory: ""
description: |-
  Provides a CherryServers SSH Key resource. This can be used to create, and delete SSH Keys associated with your Cherry account.
---

# cherryservers_ssh_key (Resource)

Provides a CherryServers SSH Key resource. This can be used to create, and delete SSH Keys associated with your Cherry account.

## Example Usage

```terraform
# Create a new SSH Key for your account
# (Optionally) specify the path to your key as a terraform variable
variable "ssh_key_path" {
  type        = string
  description = "The file path to an ssh public key"
  default     = "~/.ssh/cherry.pub"
}

resource "cherryservers_ssh_key" "my_ssh_key" {
  name       = "mykey"
  public_key = file(var.ssh_key_path) # The public key contents can also be stored specific here directly
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Label of the SSH key.
- `public_key` (String) Public SSH key.

### Read-Only

- `created` (String) Date when this Key was created.
- `fingerprint` (String) Fingerprint of the SSH public key.
- `id` (String) ID of the SSH Key.
- `updated` (String) Date when this Key was last modified.

## Import

Import is supported using the following syntax:

```shell
# Import existing SSH key via ID
terraform import cherryservers_ssh_key.main-ssh-key 1234
```