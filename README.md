Cherry Servers Terraform Provider
==================
- Cherry Servers Website: https://www.cherryservers.com
- Terraform Website: https://www.terraform.io

<img src="https://raw.githubusercontent.com/hashicorp/terraform-website/master/public/img/logo-hashicorp.svg" width="600px">

Requirements
------------

-   [Terraform](https://www.terraform.io/downloads.html) >= 1.0
-   [Go](https://golang.org/doc/install) >= 1.20 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-cherryservers`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone https://github.com/cherryservers/terraform-provider-cherryservers.git
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-cherryservers
$ go build -o terraform-provider-cherryservers
```

Generate documentation

```sh
tfplugindocs generate
```

Using the provider
----------------------

See the [Cherry Servers Provider documentation](https://registry.terraform.io/providers/cherryservers/cherryservers/latest/docs) to get started using the Cherry Servers provider.

The Terraform provider should be configured with proper credentials:

```
export CHERRY_AUTH_TOKEN="4bdc0acb8f7af4bdc0acb8f7afe78522e6dae9b7e03b0e78522e6dae9b7e03b0"
```

or 

```
provider "cherryservers" {
  api_token = "4bdc0acb8f7af4bdc0acb8f7afe78522e6dae9b7e03b0e78522e6dae9b7e03b0"
}
```