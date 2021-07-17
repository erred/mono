resource "google_app_engine_application" "default" {
  auth_domain    = "gmail.com"
  database_type  = "CLOUD_FIRESTORE"
  location_id    = "us-central"
  serving_status = "USER_DISABLED"

  feature_settings {
    split_health_checks = true
  }
}
