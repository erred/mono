resource "google_project_default_service_accounts" "default" {
  project = data.google_project.default.project_id
  action  = "DEPRIVILEGE"
}
