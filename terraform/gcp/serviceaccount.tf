resource "google_project_default_service_accounts" "default" {
  project = data.google_project.default.project_id
  action  = "DEPRIVILEGE"
}

resource "google_service_account" "cluster31_kaniko" {
  account_id = "cluster31-kaniko"
}
resource "google_service_account_iam_policy" "cluster31_kaniko" {
  service_account_id = google_service_account.cluster31_kaniko.name
  policy_data        = data.google_iam_policy.service_account_cluster31_kaniko.policy_data
}
data "google_iam_policy" "service_account_cluster31_kaniko" {}
resource "google_service_account_key" "cluster31_kaniko" {
  service_account_id = google_service_account.cluster31_kaniko.name
}


resource "google_service_account" "hetzner_medea" {
  account_id = "hetzner-medea"
}
resource "google_service_account_iam_policy" "hetzner_medea" {
  service_account_id = google_service_account.hetzner_medea.name
  policy_data        = data.google_iam_policy.service_account_hetzner_medea.policy_data
}
data "google_iam_policy" "service_account_hetzner_medea" {}
resource "google_service_account_key" "hetzner_medea" {
  service_account_id = google_service_account.hetzner_medea.name
}
