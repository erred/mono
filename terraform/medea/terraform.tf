terraform {
  required_providers {
    vault = {
      source = "hashicorp/vault"
      version = "3.5.0"
    }
  }
}

provider "vault" {
  address = "http://192.168.100.1:29001"
}
