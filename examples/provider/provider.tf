terraform {
  required_providers {
    cherryservers = {
      source = "cherryservers/cherryservers"
    }
  }
}

# Set the variable value in variables.tf file.
# Or set the CHERRY_AUTH_KEY or CHERRY_AUTH_TOKEN environment variables.
variable "cherry_api_token" {
  description = "Cherry servers API token"
  type        = string
  default     = "my_api_token_goes_here"
}

# Configure the Cherry Servers Provider.
provider "cherryservers" {
  api_token = var.cherry_api_token // API token can be found in Cherry Servers client portal - https://portal.cherryservers.com/settings/api-keys
}