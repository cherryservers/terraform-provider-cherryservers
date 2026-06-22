terraform {
  required_providers {
    cherryservers = {
      source = "cherryservers/cherryservers"
    }
  }
}

# Set the variable value in *.tfvars file
# or using the -var="cherry_api_key=..." CLI option
variable "cherry_api_key" {
  sensitive = true
}

# Configure the Cherry Servers Provider.
provider "cherryservers" {
  api_token = var.cherry_api_key
}