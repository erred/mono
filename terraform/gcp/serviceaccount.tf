resource "google_project_default_service_accounts" "default" {
  project = data.google_project.default.project_id
  action  = "DEPRIVILEGE"
}

resource "google_service_account" "gtm-appengine" {
  account_id   = "gtm-appengine"
  display_name = "GTM AppEngine"
}
