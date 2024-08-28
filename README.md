# Cherry Servers Terraform Provider

![GitHub Release](https://img.shields.io/github/v/release/caliban0/terraform-provider-cherryservers?include_prereleases)
[![codecov](https://codecov.io/gh/caliban0/terraform-provider-cherryservers/graph/badge.svg?token=E0YQGYS8JH)](https://codecov.io/gh/caliban0/terraform-provider-cherryservers)

- Cherry Servers Website: https://www.cherryservers.com
- Terraform Website: https://www.terraform.io

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.8
- [Go](https://golang.org/doc/install) >= 1.21 (to build the provider plugin)

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Using the provider

See the documentation in [./docs/](/docs/) or [Cherry Servers Provider documentation](https://registry.terraform.io/providers/cherryservers/cherryservers/latest/docs) to get started using the Cherry Servers provider.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
