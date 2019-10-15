---
layout: "cherryservers"
page_title: "Provider: CherryServers"
sidebar_current: "docs-cherryservers-index"
description: |-
  The CherryServers provider is used to interact with the resources supported by CherryServers. No terraform configuration is needed to configure the provider.  You must set a `CHERRY_AUTH_TOKEN="1238haufhiasufhaisf"` in  your environment with your appropriate API Key.  You will likely also need to find your `Team ID` for your account in order use a `cherryservers_ip` or `cherryserver_server` resource.

---

# CherryServers Provider

The CherryServers provider is used to interact with the
resources supported by CherryServers. The provider needs to be configured
with the proper credentials before it can be used. You must set a `CHERRY_AUTH_TOKEN="1238haufhiasufhaisf"` in  your environment with your appropriate API Key. 

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# No additional terraform configuration is necessary to use the CherryServers Provider
# Simply have your team_id ready in order to begin creating projects, IPs, and Servers.
# No team_id is necessary to create an SSH key for your account. 

# Create a web server
resource "cherryservers_server" "web" {
  # ...
}
```

## Argument Reference

- None

## Attributes Reference

* `auth_token` - (Optional) This is CherryServers access token. It must be provided, but it can also be specified from the `CHERRY_AUTH_TOKEN` environment variable.