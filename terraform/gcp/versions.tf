terraform {
  backend "gcs" {
    bucket = "com-seankhliao-terraform"
    prefix = "mono/terraform/gcp"
  }
  required_providers {
    google = {
      # use beta as default
      source  = "hashicorp/google-beta"
      version = "~> 4.0"
    }
  }
}

provider "google" {
  project = "com-seankhliao"
}

data "google_project" "default" {}
