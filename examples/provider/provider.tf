terraform {
  required_providers {
    cherryservers = {
      source = "cherryservers/cherryservers"
    }
  }
}

# Set the variable value in variables.tf file.
# Or set the CHERRY_API_KEY environment variable.
variable "cherry_api_key" {
  description = "Cherry servers API key"
  type        = string
  default     = "my_api_key_goes_here"
}

# Configure the Cherry Servers Provider.
provider "cherryservers" {
  api_key = var.cherry_api_key // API keys can be found at the Cherry Servers client portal - https://portal.cherryservers.com/settings/api-keys
}