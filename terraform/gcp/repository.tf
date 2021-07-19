resource "google_sourcerepo_repository" "mono" {
  name = "github_seankhliao_mono"
}
resource "google_sourcerepo_repository_iam_policy" "mono" {
  repository  = google_sourcerepo_repository.mono.name
  policy_data = data.google_iam_policy.repository_mono.policy_data
}
data "google_iam_policy" "repository_mono" {}
