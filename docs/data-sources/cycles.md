---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "cherryservers_cycles Data Source - cherryservers"
subcategory: ""
description: |-
  Provides a CherryServers billing cycles data source. This can be used to read available billing cycle data.
---

# cherryservers_cycles (Data Source)

Provides a CherryServers billing cycles data source. This can be used to read available billing cycle data.

## Example Usage

```terraform
# Get available billing cycles.
data "cherryservers_cycles" "all" {

}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `cycles` (Attributes List) Available billing cycles. (see [below for nested schema](#nestedatt--cycles))

<a id="nestedatt--cycles"></a>
### Nested Schema for `cycles`

Read-Only:

- `id` (Number)
- `name` (String)
- `slug` (String) A more readable substitute for id.
