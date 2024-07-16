# Create a new SSH Key for your account
# (Optionally) specify the path to your key as a terraform variable
variable "ssh_key_path" {
  type        = string
  description = "The file path to an ssh public key"
  default     = "~/.ssh/cherry.pub"
}

resource "cherryservers_ssh_key" "my_ssh_key" {
  label      = "mykey"
  public_key = file(var.ssh_key_path) # The public key contents can also be stored specific here directly
}
