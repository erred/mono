locals {
  ar_run_url = "${google_artifact_registry_repository.run.location}-${lower(google_artifact_registry_repository.run.format)}.pkg.dev/${data.google_project.default.project_id}/${google_artifact_registry_repository.run.repository_id}"
}

resource "google_artifact_registry_repository" "run" {
  location      = "us-central1"
  repository_id = "run"
  format        = "DOCKER"
}
resource "google_artifact_registry_repository_iam_policy" "run" {
  location    = google_artifact_registry_repository.run.location
  repository  = google_artifact_registry_repository.run.name
  policy_data = data.google_iam_policy.ar_run.policy_data
}
data "google_iam_policy" "ar_run" {}
