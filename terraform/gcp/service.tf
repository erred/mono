variable "googleapis" {
  default = [
    "containerregistry",
    "dns",
    "iamcredentials",
    "iam",
    "logging",
    "monitoring",
    "servicemanagement",
    "storage-api",
    "storage-component",
  ]
}
variable "googleapis_with_identity" {
  default = [
    "artifactregistry",
    "cloudasset",
    "cloudbuild",
    "pubsub",
    "run",
  ]
}


resource "google_project_service" "svc" {
  for_each = toset(concat(var.googleapis, var.googleapis_with_identity))

  project = data.google_project.default.number
  service = "${each.key}.googleapis.com"
}
resource "google_project_service_identity" "svc" {
  for_each = toset(var.googleapis_with_identity)

  project = data.google_project.default.project_id
  service = google_project_service.svc[each.key].service
}
