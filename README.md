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

Building the provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-cherryservers`

```bash
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:cherryservers/terraform-provider-cherryservers.git
```

You may want to get cherrygo library first:

```bash
go get github.com/cherryservers/cherrygo
```

Enter the provider directory and build provider

```bash
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-cherryservers
$ go build -o terraform-provider-cherryservers
```

Using the provider
------------------

The cherryservers provider will be installed on `terraform init` of a template using any of the cherryservers_* resources.

Usage
-----

The provider should be configured with proper credentials:

```
export CHERRY_AUTH_TOKEN="4bdc0acb8f7af4bdc0acb8f7afe78522e6dae9b7e03b0e78522e6dae9b7e03b0"
```

Examples
--------

#### Resource cherryservers_ssh

**cherryservers_ssh** module needed for adding public SSH keys to your account. After creation of this resource you may add one or more such keys to your newly ordered server resource by passing variable `${cherryservers_ssh.johny.id}` to `cherryservers_server` module`s resource.

You may update either `name` or `public_key` while working with your infrastructure file.

```
resource "cherryservers_ssh" "johny" {
    name = "johny"
    public_key = "${file("/path/to/public/key/johny.key")}"
}
```

##### Argument Reference

* **name** - label of new added SSH public key
* **public_key** - public key itself. You need to provide path to public key

During ssh addition process, some additional variable will be acquired from API:

* **fingerprint** - calculated fingerprint of added public ssh key
* **created** - the date public SSH key was added
* **updated** - the date when public key was updated

Those are needed internal module usage, but you may use them for other cases too.

#### Resource cherryservers_server

*c*herryservers_server** module needed for adding new bare metal server to your infrastructure.

```
# Create a server
resource "cherryservers_server" "my-dream-server-1" {
    project_id = "79813"
    region = "EU-East-1"
    hostname = "dream-server-1.example.com"
    image = "Ubuntu 16.04 64bit"
    plan_id = "86"
    ssh_keys_ids = ["95"]
}
```

##### Argument Reference

* **project_id** - (requered) ID of project of the servers
* **region** - region of the server. (EU-East-1 or EU-West-1) 
* **hostname** - define hostname of a server
* **image** - image to be installed on the server, e.g. ```Ubuntu 16.04 64bit```
* **plan_id** - plan for server creation
* **ssh_keys_ids** - SSH key`s ID for adding SSH key to server
* **ip_addresses_ids** - list of floating IP addresses UIDs to be added to a new server.

During server creation process, some additional variable will be acquired from API:

* **private_ip** - assigned private IP address of a server
* **primary_ip** - assigned primary (public) IP address of a server
* **power_state** - current power state of a server
* **state** - deployment state of a server

Those are needed internal module usage, but you may use them for other cases too.

#### Resource cherryservers_ip

**cherryservers_ip** - module needed for adding new floating IPs to your infrastructure. You may want to order new floating IP address and assign it to bare metal server which will be created just after you order you floating IP.

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

* **project_id** - ID of project
* **region** - region of the server. (EU-East-1 or EU-West-1) 
* **routed_to** - you need to specify UID of server`s IP address to which you with to route it to
* **routed_to_hostname** - on the other hand, you may specify hostname of server to which you want to route
* **routed_to_ip** - or you may want to specify IP address of the serfver to route to

During floating IP creation process, some additional variable will be acquired from API:

* **cidr** - subnet to which IP belongs
* **type** - type of IP address, it could be primary, private, floating, subnet etc.
* **gateway** - the gateway for newly created IP
* **address** - the address itself, you probably will need this for later use
* **a_record** - public A record which points to assigned IP address
* **ptr_record** - PTR record for IP address

Those are needed internal module usage, but you may use them for other cases too.