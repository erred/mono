variable "googleapis" {
  # https://console.cloud.google.com/apis/library
  description = "GCP APIs"
  type        = list(string)
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
  # https://cloud.google.com/iam/docs/service-agents
  description = "GCP APIs with Service Agents"
  type        = list(string)
  default = [
    "artifactregistry",
    "cloudasset",
    "cloudbuild",
    "pubsub",
    "run",
  ]
}

locals {
  # Google APIs Service Agent
  # https://cloud.google.com/iam/docs/service-accounts#google-managed
  googleapis_service_agent = "${data.google_project.default.number}@cloudservices.gserviceaccount.com"

  # cloud build service account
  #
  cloudbuild_service_account = "${data.google_project.default.number}@cloudbuild.gserviceaccount.com"

  # appengine
  gae_service_service_account = "service-${data.google_project.default.number}@gcp-gae-service.iam.gserviceaccount.com"
  gae_api_service_account = "service-${data.google_project.default.number}@gae-api-prod.google.com.iam.gserviceaccount.com"
}

resource "google_project_iam_policy" "default" {
  project     = data.google_project.default.project_id
  policy_data = data.google_iam_policy.project.policy_data
}
data "google_iam_policy" "project" {
  audit_config {
    audit_log_configs {
      log_type = "ADMIN_READ"
    }
    audit_log_configs {
      log_type = "DATA_READ"
    }
    audit_log_configs {
      log_type = "DATA_WRITE"
    }
    service = "allServices"
  }

  binding {
    members = [
      "serviceAccount:${local.cloudbuild_service_account}",
    ]
    role = "roles/cloudbuild.builds.builder"
  }
  binding {
    members = [
      "serviceAccount:${local.googleapis_service_agent}",
    ]
    role = "roles/editor"
  }
  binding {
    members = [
      "user:seankhliao@gmail.com",
    ]
    role = "roles/owner"
  }
  binding {
    members = [
      "serviceAccount:${google_project_service_identity.svc["cloudbuild"].email}"
    ]
    role = "roles/run.admin"
  }
  binding {
    members = [
      "serviceAccount:${google_project_service_identity.svc["cloudbuild"].email}"
    ]
    role = "roles/artifactregistry.writer"
  }
  binding {
    members = [
      "user:seankhliao@gmail.com",
    ]
    role = "roles/storage.admin"
  }

  binding {
    members = [
      "serviceAccount:service-330311169810@compute-system.iam.gserviceaccount.com"
    ]
    role = "roles/compute.serviceAgent"
  }

  binding {
    members = [
      "serviceAccount:${local.gae_service_service_account}",
    ]
    role = "roles/appengine.serviceAgent"
  }

  dynamic "binding" {
    for_each = toset(var.googleapis_with_identity)
    iterator = svc

    content {
      members = [
        "serviceAccount:${google_project_service_identity.svc["${svc.value}"].email}",
      ]
      role = "roles/${svc.value}.serviceAgent"
    }
  }
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
