---
layout: "cherryservers"
page_title: "CherryServers: cherryservers_project"
sidebar_current: "docs-cherryservers-resource-project"
description: |-
  Provides a CherryServers Project resource. This can be used to create, modify, and delete Projects. These projects can contain segregated Server and IPs that can be destroyed all at once by deleting the Project.
---

# cherryservers\_project

Provides a CherryServers Project resource. This can be used to create,
modify, and delete Projects. 

## Example Usage

```hcl
# You must find your team_id for your account by login into the CherryServers portal: [https://portal.cherryservers.com/#/login](https://portal.cherryservers.com/#/login)
variable "team_id" {
  default = "12345"
}
# Specify a name for your project
variable "project_name" {
  default = "My Cool New Project"
}

# Create a new Project for your team 
resource "cherryservers_project" "myproject" {
  team_id = "${var.team_id}"
  name = "${var.project_name}"
}
```

## Argument Reference

The following arguments are supported:

* `team_id` - (Required) The Server image ID or slug.
* `name` - (Required) The Server name.

## Attributes Reference

The following attributes are exported:

* `team_id` - The ID of team that owns the Project
* `name`- The name of the Project
* `project_id` - The computed ID of the Project, to be used with other Resources 

## Import

Servers can be imported using the Project `id`, e.g.

```
terraform import cherryservers_project.myproject 8123
```
