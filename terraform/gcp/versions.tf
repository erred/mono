terraform {
  required_providers {
    google = {
      version = "~> 3.0"
    }
  }
}
provider "google" {
  project = "com-seankhliao"
}
