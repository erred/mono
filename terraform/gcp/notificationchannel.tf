locals {
  # not possible to manage with API
  notification_default_app = "projects/${data.google_project.default.project_id}/notificationChannels/3234509444054832325"
}

resource "google_monitoring_notification_channel" "default_email" {
  description  = "Email for error notifications"
  display_name = "seankhliao@gmail.com"
  labels = {
    "email_address" = "seankhliao@gmail.com"
  }
  type = "email"
}
