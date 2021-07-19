########################################
# variables
########################################
variable "publicruns" {
  description = "cloud run services that are publicly available"
  type = map(object({
    image = string
    url   = string
  }))
  default = {
    "go-seankhliao-com" = {
      image = "vanity"
      url   = "go.seankhliao.com"
    }
    "seankhliao-com" = {
      image = "w16"
      url   = "seankhliao.com"
    }
  }
}

########################################
# shared settings
########################################
resource "google_logging_project_sink" "cloudrun_logs" {
  destination            = "storage.googleapis.com/${google_storage_bucket.cloudrun_logs.name}"
  filter                 = <<-EOT
    resource.type = "cloud_run_revision"
    resource.labels.location = "us-central1"
  EOT
  name                   = "cloudrun-logs"
  unique_writer_identity = true
}
resource "google_storage_bucket" "cloudrun_logs" {
  location                    = "US-CENTRAL1"
  name                        = "${data.google_project.default.project_id}-cloudrun-logs"
  storage_class               = "STANDARD"
  uniform_bucket_level_access = true
}
resource "google_storage_bucket_acl" "cloudrun_logs" {
  bucket = google_storage_bucket.cloudrun_logs.name
}
resource "google_storage_bucket_iam_policy" "cloudrun_logs" {
  bucket      = google_storage_bucket.cloudrun_logs.name
  policy_data = data.google_iam_policy.bucket_cloudrun_logs.policy_data
}
data "google_iam_policy" "bucket_cloudrun_logs" {
  binding {
    members = [
      "projectEditor:${data.google_project.default.project_id}",
      "projectOwner:${data.google_project.default.project_id}",
      google_logging_project_sink.cloudrun_logs.writer_identity,
    ]
    role = "roles/storage.legacyBucketOwner"
  }
}


########################################
# publicruns
########################################

resource "google_service_account" "publicruns" {
  for_each = var.publicruns

  account_id  = each.key
  description = "cloud run sa for ${replace(each.key, "-", ".")}"
}
resource "google_service_account_iam_policy" "publicruns" {
  for_each = var.publicruns

  service_account_id = google_service_account.publicruns[each.key].id
  policy_data        = data.google_iam_policy.publicruns_service_account[each.key].policy_data
}
data "google_iam_policy" "publicruns_service_account" {
  for_each = var.publicruns

  binding {
    members = [
      "serviceAccount:${local.cloudbuild_service_account}",
    ]
    role = "roles/iam.serviceAccountUser"
  }
}

resource "google_cloud_run_service" "publicruns" {
  for_each = var.publicruns

  location = "us-central1"
  name     = each.key
  template {
    spec {
      container_concurrency = 80
      service_account_name  = google_service_account.publicruns[each.key].email
      timeout_seconds       = 10

      containers {
        # image = "${google_artifact_registry_repository.run.location}-docker.pkg.dev/${data.google_project.default.project_id}/${google_artifact_registry_repository.run.repository_id}/${each.value.image}:latest"
        image = "${local.ar_run_url}/${each.value.image}:latest"

        ports {
          container_port = 8080
        }

        resources {
          limits = {
            "cpu"    = "1"
            "memory" = "128Mi"
          }
        }
      }
    }
  }
  traffic {
    latest_revision = true
    percent         = 100
  }
}
resource "google_cloud_run_service_iam_policy" "publicruns" {
  for_each = var.publicruns

  location = google_cloud_run_service.publicruns[each.key].location
  service  = google_cloud_run_service.publicruns[each.key].name

  policy_data = data.google_iam_policy.publicruns.policy_data
}
data "google_iam_policy" "publicruns" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}

resource "google_cloud_run_domain_mapping" "publicruns" {
  for_each = var.publicruns

  location = "us-central1"
  name     = each.value.url
  metadata {
    namespace = data.google_project.default.project_id
  }
  spec {
    route_name = each.key
  }
}
