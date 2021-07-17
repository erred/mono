resource "google_project_default_service_accounts" "default" {
  project = data.google_project.default.project_id
  action  = "DEPRIVILEGE"
}

resource "google_service_account" "cluster30_kaniko" {
  account_id = "cluster30-kaniko"
}
resource "google_service_account_iam_policy" "cluster30_kaniko" {
  service_account_id = google_service_account.cluster30_kaniko.name
  policy_data        = data.google_iam_policy.service_account_cluster30_kaniko.policy_data
}
data "google_iam_policy" "service_account_cluster30_kaniko" {}


resource "google_service_account" "cluster30_traefik" {
  account_id = "cluster30-traefik"
}
resource "google_service_account_iam_policy" "cluster30_traefik" {
  service_account_id = google_service_account.cluster30_traefik.name
  policy_data        = data.google_iam_policy.service_account_cluster30_traefik.policy_data
}
data "google_iam_policy" "service_account_cluster30_traefik" {}

resource "google_service_account" "hetzner_medea" {
  account_id = "hetzner-medea"
}
resource "google_service_account_iam_policy" "hetzner_medea" {
  service_account_id = google_service_account.hetzner_medea.name
  policy_data        = data.google_iam_policy.service_account_hetzner_medea.policy_data
}
data "google_iam_policy" "service_account_hetzner_medea" {}
