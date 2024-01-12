terraform {
  required_providers {
    cherryservers = {
      source = "cherryservers/cherryservers"
    }
  }
}

provider "cherryservers" {
  api_token = "CHEERRY_SERVERS_API_KEY" // API key can be found in Cherry Servers client portal - https://portal.cherryservers.com/settings/api-keys
}
