resource "google_storage_bucket" "terraform" {
  location                    = "EUROPE-WEST4"
  name                        = "com-seankhliao-terraform"
  storage_class               = "STANDARD"
  uniform_bucket_level_access = true
}
resource "google_storage_bucket_acl" "terraform" {
  bucket = google_storage_bucket.terraform.name
}
resource "google_storage_bucket_iam_policy" "terraform" {
  bucket      = google_storage_bucket.terraform.name
  policy_data = data.google_iam_policy.bucket_terraform.policy_data
}
data "google_iam_policy" "bucket_terraform" {}
