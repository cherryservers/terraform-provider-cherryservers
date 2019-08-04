Terraform Provider
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.12 (to build the provider plugin)

Building Terraform provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-cherryservers`

```bash
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone https://github.com/cherryservers/terraform-provider-cherryservers.git
```

You may want to get cherrygo library first:

```bash
go get github.com/cherryservers/cherrygo
```

Enter directory path for the Terraform provider and build it

```bash
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-cherryservers
$ go build -o terraform-provider-cherryservers
```

Using Terraform provider
------------------

You may download already builded binaries for your operating system from our mirror:

```
http://downloads.cherryservers.com/other/terraform/
```

There are build for Mac, Linux and Windows

The Terraform provider will be installed on `terraform init` command from a template by using any of the cherryservers_* resources.

Usage
-----

The Terraform provider should be configured with proper credentials:

```
export CHERRY_AUTH_TOKEN="4bdc0acb8f7af4bdc0acb8f7afe78522e6dae9b7e03b0e78522e6dae9b7e03b0"
```

or 

```
provider "cherryservers" {
  auth_token = "4bdc0acb8f7af4bdc0acb8f7afe78522e6dae9b7e03b0e78522e6dae9b7e03b0"
}
```

or even

```
variable "auth_token" {}

provider "cherryservers" {
    auth_token = "${var.auth_token}"
}
```

```
provider "cherryservers" {
  auth_token = "${var.auth_token}"
}
```

Examples
--------

#### Resource cherryservers_project

**cherryservers_project** module is needed to create a new project in your team. A project consists of the infrastructure you deployed. You may have several projects in a team with several servers within a project.

You may update your project with a new `name` while working with your infrastructure file.

```
resource "cherryservers_project" "DreamProject" {
    team_id = "28519"
    name = "DreamProject1"
}
```

##### Argument Reference
* **name** - the name of your project
* **team_id** - ID of the team your project will reside

#### Resource cherryservers_ssh

**cherryservers_ssh** module needed to add public SSH keys to your account. After creation of this resource you may assign one or more of such keys to your newly ordered server instance by passing variable `${cherryservers_ssh.johny.id}` to `cherryservers_server` module`s resource.

You may update either `name` or `public_key` while working with your infrastructure file.

```
resource "cherryservers_ssh" "johny" {
    name = "johny"
    public_key = "${file("/path/to/public/key/johny.key")}"
}
```

##### Argument Reference

* **name** - label of your newly added SSH public key
* **public_key** - public key itself. You need to provide path to public key

During ssh addition process, some additional variables will be acquired via API:

* **fingerprint** - calculated fingerprint of added public ssh key
* **created** - the date public SSH key was added
* **updated** - the date when public key was updated

These variables are needed for internal usage of the module, but you may use them for other purposes as well.

#### Resource cherryservers_server

*c*herryservers_server** module is needed for adding new bare metal servers to your infrastructure.

```
# Create a server
resource "cherryservers_server" "my-dream-server-1" {
    project_id = "79813"
    region = "EU-East-1"
    hostname = "dream-server-1.example.com"
    image = "Ubuntu 16.04 64bit"
    plan_id = "86"
    user_data = "I2Nsb3VkLWNvbmZpZwpwYWNrYWdlczoKICAtIGlmdG9wCiAgLSBubW9uCg=="
    ssh_keys_ids = ["95"]
}
```

##### Argument Reference

* **project_id** - your project ID
* **region** - server region ("EU-East-1" or "EU-West-1") 
* **hostname** - your defined server hostname
* **image** - your server image e.g. ```Ubuntu 16.04 64bit```
* **plan_id** - your server plan ID
* **ssh_keys_ids** - ID of your SSH key to be assigned to a new server
* **ip_addresses_ids** - UIDs of your floating IP addresses to be assigned to a new server
* **user_data** - base64 encoded User-Data blob. It should be either bash or cloud-config script.

During server creation process, some additional variables will be acquired via API:

* **private_ip** - assigned private IP address of a server
* **primary_ip** - assigned primary (public) IP address of a server
* **power_state** - current power state of a server
* **state** - deployment state of a server

These variables are needed for internal usage of the module, but you may use them for other purposes as well.

#### Resource cherryservers_ip

**cherryservers_ip** - module needed for adding new floating IPs to your infrastructure. For instance, you may want to order a new floating IP address and assign it to your server instance that will be created right after you order you floating IP.

```
# Create an IP address
resource "cherryservers_ip" "floating-ip1-server-1" {
    project_id = "79813"
    region = "EU-East-1"
    routed_to_ip = "188.214.132.48"
}
```

or 

```
# Create an IP address
resource "cherryservers_ip" "floating-ip1-server-1" {
    project_id = "79813"
    region = "EU-East-1"
    routed_to = "9268173d-f903-5b60-a19a-cd5f572c9377"
}
```

or

```
# Create an IP address
resource "cherryservers_ip" "floating-ip1-server-1" {
    project_id = "79813"
    region = "EU-East-1"
    routed_to_hostname = "server02.example.com"
}
```

##### Argument Reference

* **project_id** - project ID
* **region** - server region ("EU-East-1" or "EU-West-1") 
* **routed_to** - you need to specify a UID of your server`s primary IP address for routing your floating IP
* **routed_to_hostname** - you may also specify hostname of server for routing your floating IP
* **routed_to_ip** - alternativelly, you may specify a static IP address of the server for routing your floating IP

During floating IP creation process, some additional variables will be acquired via API:

* **cidr** - subnet to which IP belongs
* **type** - type of IP address (e.i. primary, private, floating, subnet etc.)
* **gateway** - the gateway for newly created IP
* **address** - the address itself, you will probably need this for later use
* **a_record** - public A record which points to the assigned IP address
* **ptr_record** - PTR record for IP address

These variables are needed for internal usage of the module, but you may use them for other purposes as well.

#### Real world example

Imagine you have started working on a new project and want to create an infrastructure for it. You hire a new developer and provide him with a new server, which has a floating IP assigned and should be reachable with developer`s ssh key.

In such case, you need to create the following Terraform scenario:

```
resource "cherryservers_project" "DreamProject99" {
    team_id = "28519"
    name = "DreamProject99"
}

resource "cherryservers_ssh" "johny-key-1" {
    name = "johny1"
    public_key = "${file("path/to/johny.key")}"
}

resource "cherryservers_ip" "floating-ip-server99" {
    project_id = "${cherryservers_project.DreamProject99.id}"
    region = "EU-East-1"
}

resource "cherryservers_server" "super-server99" {
    project_id = "${cherryservers_project.DreamProject99.id}"
    region = "EU-East-1"
    hostname = "virtual99.turbo.com"
    image = "Ubuntu 16.04 64bit"
    plan_id = "165"
    ssh_keys_ids = ["95", "${cherryservers_ssh.johny-key-1.id}"]
    ip_addresses_ids = ["${cherryservers_ip.floating-ip-server99.id}"]
}
```

As you can see, at first you create new **DreamProject99** project, than you add a **johny-key-1** key to portal, after that you add new floating IP address **floating-ip-server99** to your newly created project **DreamProject99** and then you create a new server **super-server99**, which will have **johny-key-1** and **floating-ip-server99** assigned to it.


## License

See the [LICENSE](LICENSE.md) file for license rights and limitations.