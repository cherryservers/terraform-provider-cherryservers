---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cherryservers_project Data Source - cherryservers"
subcategory: ""
description: |-
  Provides a CherryServers Project data source. This can be used to read project data from CherryServers.
---

# cherryservers_project (Data Source)

Provides a CherryServers Project data source. This can be used to read project data from CherryServers.

## Example Usage

```terraform
# Create a Project data source.
data "cherryservers_project" "cool_project" {
  id = "123456"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `id` (Number) Project identifier.

### Read-Only

- `bgp` (Attributes) Project border gateway protocol(BGP) configuration. (see [below for nested schema](#nestedatt--bgp))
- `name` (String) The name of the project.

<a id="nestedatt--bgp"></a>
### Nested Schema for `bgp`

Read-Only:

- `enabled` (Boolean) BGP is enabled for the project.
- `local_asn` (Number) The local ASN of the project.
