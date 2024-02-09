terraform {
  required_providers {
    cherryservers = {
      source = "cherryservers/cherryservers"
    }
  }
}

# Set the variable value in *.tfvars file
# or using -var="cherry_api_token=..." CLI option
variable "cherry_api_token" {}

# Configure the Cherry Servers Provider
provider "cherryservers" {
  api_token = var.cherry_api_token // API key can be found in Cherry Servers client portal - https://portal.cherryservers.com/settings/api-keys
}
