terraform {
  required_providers {
    google = {
      # use beta as default
      source  = "hashicorp/google-beta"
      version = "~> 3.0"
    }
  }
}

provider "google" {
  project = "com-seankhliao"
}

data "google_project" "default" {}
